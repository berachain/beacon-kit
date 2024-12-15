// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
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

func (s *Service[LoggerT]) processProposal(
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

	// Since the application can get access to FinalizeBlock state and write to
	// it, we must be sure to reset it in case ProcessProposal timeouts and is
	// called
	// again in a subsequent round. However, we only want to do this after we've
	// processed the first block, as we want to avoid overwriting the
	// finalizeState
	// after state changes during InitChain.
	s.processProposalState = s.resetState(ctx)
	if req.Height > s.initialHeight {
		s.finalizeBlockState = s.resetState(ctx)
	}

	s.processProposalState.SetContext(
		s.getContextForProposal(
			s.processProposalState.Context(),
			req.Height,
		),
	)

	resp, err := s.Blockchain.ProcessProposal(
		s.processProposalState.Context(),
		req,
	)
	if err != nil {
		s.logger.Error(
			"failed to process proposal",
			"height",
			req.Height,
			"time",
			req.Time,
			"hash",
			fmt.Sprintf("%X", req.Hash),
			"err",
			err,
		)
		return &cmtabci.ProcessProposalResponse{
			Status: cmtabci.PROCESS_PROPOSAL_STATUS_REJECT,
		}, nil
	}

	return resp, nil
}
