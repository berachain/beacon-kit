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

	cmttypes "github.com/cometbft/cometbft/types"
)

const roundChangeSubscriber = "beacon-kit-round-change"

// subscribeToRoundChanges subscribes to CometBFT's EventNewRound on the
// local event bus and calls Blockchain.HandleRoundChange for rounds > 0.
// This allows the sequencer to detect consensus round timeouts and trigger
// a fresh payload build for the new proposer.
func (s *Service) subscribeToRoundChanges(ctx context.Context) {
	sub, err := s.node.EventBus().Subscribe(ctx, roundChangeSubscriber, cmttypes.EventQueryNewRound)
	if err != nil {
		// Should never fail after a successful node start. Degrades gracefully.
		s.logger.Error("Failed to subscribe to NewRound events", "error", err)
		return
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-sub.Canceled():
				s.logger.Warn("NewRound subscription cancelled", "error", sub.Err())
				return
			case msg, open := <-sub.Out():
				if !open {
					return
				}
				event, ok := msg.Data().(cmttypes.EventDataNewRound)
				if !ok || event.Round < 1 {
					continue
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
		}
	}()
}
