// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package cometbft

import (
	"context"
	"time"

	cmttypes "github.com/cometbft/cometbft/types"
)

const (
	roundChangeSubscriber = "beacon-kit-round-change"

	// roundChangeBufferCapacity is kept at 1 so any overflow trips the
	// resubscribe path below instead of silently absorbing bursts.
	roundChangeBufferCapacity = 1

	// roundChangeResubscribeBackoff is the delay between resubscribe attempts
	// after the underlying subscription is cancelled (e.g. ErrOutOfCapacity).
	roundChangeResubscribeBackoff = 1 * time.Second
)

// startRoundChangeListener subscribes to CometBFT's EventNewRound on the local
// event bus and runs a goroutine that calls Blockchain.HandleRoundChange for
// rounds > 0. This allows the sequencer to detect consensus round timeouts and
// trigger a fresh payload build for the new proposer.
//
// If the subscription is cancelled while ctx is still alive (e.g. slow consumer
// overflow), the goroutine resubscribes with a fixed backoff so a one-off drop
// doesn't permanently disable the listener for the rest of the node's lifetime.
func (s *Service) startRoundChangeListener(ctx context.Context) {
	subscribe := func() cmttypes.Subscription {
		sub, err := s.node.EventBus().Subscribe(
			ctx, roundChangeSubscriber, cmttypes.EventQueryNewRound, roundChangeBufferCapacity,
		)
		if err != nil {
			s.logger.Error("NewRound subscribe failed", "error", err)
			return nil
		}
		return sub
	}

	sub := subscribe()

	go func() {
		for {
			// If we don't have a live subscription, wait out the backoff and
			// try again. Bail out early if ctx is cancelled during the wait.
			if sub == nil {
				timer := time.NewTimer(roundChangeResubscribeBackoff)
				select {
				case <-ctx.Done():
					timer.Stop()
					return
				case <-timer.C:
				}
				sub = subscribe()
				continue
			}

			select {
			case <-ctx.Done():
				return
			case <-sub.Canceled():
				s.logger.Warn("NewRound subscription cancelled, will resubscribe", "error", sub.Err())
				sub = nil
			case msg, open := <-sub.Out():
				if !open {
					sub = nil
					continue
				}
				s.dispatchRoundChangeEvent(ctx, msg.Data())
			}
		}
	}()
}

func (s *Service) dispatchRoundChangeEvent(ctx context.Context, data any) {
	event, ok := data.(cmttypes.EventDataNewRound)
	if !ok || event.Round < 1 {
		return
	}
	s.logger.Info("CometBFT round change detected",
		"height", event.Height,
		"round", event.Round,
		"proposer_address", event.Proposer.Address,
	)

	freshState := s.resetState(ctx)
	//nolint:contextcheck // ctx already passed via resetState
	s.Blockchain.HandleRoundChange(
		freshState.Context(),
		event.Height,
		event.Round,
		event.Proposer.Address,
	)
}
