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

package blockchain

import (
	"context"

	"github.com/berachain/beacon-kit/primitives/math"
)

// HandleRoundChange is called when CometBFT enters a new consensus round
// (round > 0) at the same height. On the sequencer, this triggers a fresh
// payload build so that the new round's proposer gets an up-to-date payload
// with the latest transactions, and RPC nodes receive fresh flashblocks.
//
// The proposerAddress is the CometBFT address of the new round's proposer.
// The ctx should be an SDK context backed by a fresh committed multistore.
func (s *Service) HandleRoundChange(
	ctx context.Context,
	height int64,
	round int32,
	proposerAddress []byte,
) {
	if !s.preconfCfg.IsSequencer() || !s.localBuilder.Enabled() {
		s.logger.Error("Round change invoked on non-sequencer", "height", height, "round", round)
		return
	}

	s.metrics.markRoundChange()

	st := s.storageBackend.StateFromContext(ctx)

	slot := math.Slot(height) //#nosec:G115 // CometBFT heights are always positive.

	proposerPubkey, err := s.getNextProposerPubkey(st, proposerAddress)
	if err != nil {
		// Clear the tracker so the stale proposer can no longer fetch payloads.
		s.preconfProposerTracker.ClearExpectedProposer(slot)

		s.logger.Warn("Round change: failed to resolve proposer pubkey",
			"height", height,
			"round", round,
			"error", err,
		)
		return
	}

	// Update expected proposer immediately so the preconf server rejects
	// requests from the previous round's proposer, even if the new proposer
	// is not whitelisted or prefetch fails below.
	s.preconfProposerTracker.SetExpectedProposer(slot, proposerPubkey)

	if !s.preconfWhitelist.IsWhitelisted(proposerPubkey) {
		s.logger.Debug("Round change: new proposer not whitelisted, skipping rebuild",
			"height", height,
			"round", round,
			"proposer", proposerPubkey.String(),
		)
		return
	}

	s.logger.Info("Round change detected: rebuilding payload for whitelisted proposer",
		"height", height,
		"round", round,
		"proposer", proposerPubkey.String(),
	)

	lph, err := st.GetLatestExecutionPayloadHeader()
	if err != nil {
		s.logger.Error("Round change: failed to read latest execution payload header",
			"height", height,
			"round", round,
			"error", err,
		)
		return
	}
	currentTime := lph.GetTimestamp()
	ephemeralState := st.Protect(ctx)
	buildData, err := s.preFetchBuildData(ephemeralState, currentTime)
	if err != nil {
		s.logger.Error("Round change: failed to prefetch build data",
			"height", height,
			"round", round,
			"error", err,
		)
		return
	}

	// Clear stale cache so RequestPayloadAsync sends a fresh FCU.
	// Safe from races: CometBFT fires EventNewRound before enterPropose.
	s.localBuilder.InvalidatePayload(buildData.Slot, buildData.ParentBlockRoot)

	go s.handleOptimisticPayloadBuild(ctx, buildData)
}
