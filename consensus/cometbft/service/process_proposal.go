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

	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/consensus/cometbft/service/encoding"
	consensusTypes "github.com/berachain/beacon-kit/consensus/types"
	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/math"
	cmtabci "github.com/cometbft/cometbft/abci/types"
)

const (
	// BeaconBlockTxIndex represents the index of the beacon block transaction.
	// It is the first transaction in the tx list.
	BeaconBlockTxIndex uint = iota
	// BlobSidecarsTxIndex represents the index of the blob sidecar transaction.
	// It follows the beacon block transaction in the tx list.
	BlobSidecarsTxIndex
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
	s.processProposalState = s.resetState()
	if req.Height > s.initialHeight {
		s.finalizeBlockState = s.resetState()
	}

	s.processProposalState.SetContext(
		s.getContextForProposal(
			s.processProposalState.Context(),
			req.Height,
		),
	)

	const BeaconBlockTxIndex = 0

	// Decode the beacon block.
	blk, err := encoding.
		UnmarshalBeaconBlockFromABCIRequest[*types.BeaconBlock](
		req,
		BeaconBlockTxIndex,
		s.chainSpec.ActiveForkVersionForSlot(math.U64(req.Height)),
	)
	if err != nil {
		return createProcessProposalResponse(errors.WrapNonFatal(err))
	}

	var consensusBlk consensusTypes.ConsensusBlock[*types.BeaconBlock]
	consensusBlk = *consensusBlk.New(
		blk,
		req.GetProposerAddress(),
		req.GetTime(),
	)

	// Decode the blob sidecars.
	sidecars, err := encoding.
		UnmarshalBlobSidecarsFromABCIRequest[*datypes.BlobSidecars](
		req,
		BlobSidecarsTxIndex,
	)
	if err != nil {
		return createProcessProposalResponse(errors.WrapNonFatal(err))
	}
	var consensusSidecars *consensusTypes.ConsensusSidecars[
		*datypes.BlobSidecars,
		*types.BeaconBlockHeader,
	]
	consensusSidecars = consensusSidecars.New(
		sidecars,
		blk.GetHeader(),
	)

	_, err = s.Blockchain.ProcessProposal(ctx, consensusBlk, consensusSidecars)
	if err != nil {
		return nil, err
	}

	return createProcessProposalResponse(nil)

}

// createResponse generates the appropriate ProcessProposalResponse based on the
// error.
func createProcessProposalResponse(
	err error,
) (*cmtabci.ProcessProposalResponse, error) {
	status := cmtabci.PROCESS_PROPOSAL_STATUS_REJECT
	if !errors.IsFatal(err) {
		status = cmtabci.PROCESS_PROPOSAL_STATUS_ACCEPT
		err = nil
	}
	return &cmtabci.ProcessProposalResponse{Status: status}, err
}
