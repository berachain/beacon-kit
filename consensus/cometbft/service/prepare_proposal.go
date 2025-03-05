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
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package cometbft

import (
	"context"
	"fmt"
	"time"

	"github.com/berachain/beacon-kit/consensus/types"
	"github.com/berachain/beacon-kit/primitives/math"
	cmtabci "github.com/cometbft/cometbft/abci/types"
)

func (s *Service) prepareProposal(
	ctx context.Context,
	req *cmtabci.PrepareProposalRequest,
) (*cmtabci.PrepareProposalResponse, error) {
	// Check if ctx is still good. CometBFT does not check this.
	if ctx.Err() != nil {
		// If the context is getting cancelled, we are shutting down.
		// It is ok returning an empty proposal.
		return &cmtabci.PrepareProposalResponse{Txs: req.Txs}, nil
	}
	startTime := time.Now()
	defer s.telemetrySink.MeasureSince(
		"beacon_kit.runtime.prepare_proposal_duration", startTime)

	// CometBFT must never call PrepareProposal with a height of 0.
	if req.Height < 1 {
		return nil, fmt.Errorf(
			"prepareProposal at height %v: %w",
			req.Height,
			errInvalidHeight,
		)
	}

	// Always reset state given that PrepareProposal can timeout
	// and be called again in a subsequent round.
	s.prepareProposalState = s.resetState(ctx)
	s.prepareProposalState.SetContext(
		s.getContextForProposal(
			s.prepareProposalState.Context(),
			req.Height,
			// Preserve the CosmosSDK context while using the correct base ctx.
		).WithContext(ctx),
	)

	slotData := types.NewSlotData(
		math.Slot(req.GetHeight()), // #nosec G115
		nil,                        // no attestations
		nil,                        // no slashings
		req.GetProposerAddress(),
		req.GetTime(),
	)

	//nolint:contextcheck // ctx already passed via resetState
	blkBz, sidecarsBz, err := s.BlockBuilder.BuildBlockAndSidecars(
		s.prepareProposalState.Context(),
		slotData,
	)
	if err != nil {
		s.logger.Error(
			"failed to prepare proposal",
			"height", req.Height,
			"time", req.Time,
			"err", err,
		)
		return &cmtabci.PrepareProposalResponse{Txs: [][]byte{}}, nil
	}

	return &cmtabci.PrepareProposalResponse{
		Txs: [][]byte{blkBz, sidecarsBz},
	}, nil
}
