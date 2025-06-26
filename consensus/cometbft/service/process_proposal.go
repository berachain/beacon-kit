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

	cmtabci "github.com/cometbft/cometbft/abci/types"
)

func (s *Service) processProposal(
	ctx context.Context,
	req *cmtabci.ProcessProposalRequest,
) (*cmtabci.ProcessProposalResponse, error) {
	startTime := time.Now()
	defer s.telemetrySink.MeasureSince(
		"beacon_kit.runtime.process_proposal_duration", startTime)

	// CometBFT must never call ProcessProposal with a height of 0.
	if req.Height < 1 {
		return nil, fmt.Errorf(
			"processProposal at height %v: %w",
			req.Height,
			errInvalidHeight,
		)
	}

	// processProposalState is used for ProcessProposal, which is set based on
	// the previous block's state. This state is never committed. In case of
	// multiple consensus rounds, the state is always reset to the previous
	// block's state.
	processProposalState := s.resetState(ctx)

	//nolint:contextcheck // ctx already passed via resetState
	processProposalState.SetContext(
		s.getContextForProposal(
			processProposalState.Context(),
			req.Height,
		),
	)

	// errors to consensus indicate that the node was not able to understand
	// whether the block was valid or not. Viceversa, we signal that a block
	// is invalid by its status, but we do return nil error in such a case.
	valUpdates, err := s.Blockchain.ProcessProposal(
		processProposalState.Context(),
		req,
	)
	if err != nil {
		status := cmtabci.PROCESS_PROPOSAL_STATUS_REJECT
		s.logger.Error(
			"failed to process proposal",
			"height", req.Height,
			"time", req.Time,
			"hash", fmt.Sprintf("%X", req.Hash),
			"err", err,
		)
		return &cmtabci.ProcessProposalResponse{Status: status}, nil
	}

	// We must not cache execution of the first block post initialHeight
	// because its state must be handled in a different way
	if req.Height > s.initialHeight {
		stateHash := string(req.Hash)
		toCache := &CacheElement{
			state:      processProposalState,
			valUpdates: valUpdates,
		}
		s.cachedStates.cache(stateHash, toCache)
	}
	status := cmtabci.PROCESS_PROPOSAL_STATUS_ACCEPT
	return &cmtabci.ProcessProposalResponse{Status: status}, nil
}
