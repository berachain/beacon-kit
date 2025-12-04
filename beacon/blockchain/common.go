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

package blockchain

import (
	"fmt"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/consensus/cometbft/service/encoding"
	"github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/primitives/math"
	cmtabci "github.com/cometbft/cometbft/abci/types"
)

func (s *Service) ParseProcessProposalRequest(req *cmtabci.ProcessProposalRequest) (
	*ctypes.SignedBeaconBlock,
	types.BlobSidecars,
	error,
) {
	blobConsensusEnabled := s.chainSpec.IsBlobConsensusEnabledAtHeight(req.Height)

	maxTxCount := MaxConsensusTxsCount
	if blobConsensusEnabled {
		maxTxCount = 1 // After BlobEnableHeight: only 1 tx expected (block), blobs in Blob field
	}

	if len(req.GetTxs()) > maxTxCount {
		return nil, nil, fmt.Errorf("max expected %d txs, got %d", maxTxCount, len(req.GetTxs()))
	}

	forkVersion := s.chainSpec.ActiveForkVersionForTimestamp(math.U64(req.GetTime().Unix())) //#nosec: G115
	signedBlk, err := encoding.UnmarshalBeaconBlockFromABCIRequest(req.GetTxs(), BeaconBlockTxIndex, forkVersion)
	if err != nil {
		return nil, nil, err
	}

	if signedBlk == nil {
		s.logger.Warn("Aborting block verification - beacon block not found in proposal")
		return nil, nil, ErrNilBlk
	}

	// Extract sidecars using the common helper
	sidecars, err := encoding.ExtractBlobSidecarsFromRequest(
		req.GetTxs(),
		req.GetBlob(),
		req.Height,
		s.chainSpec,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to extract blob sidecars at height %d: %w", req.Height, err)
	}

	return signedBlk, sidecars, nil
}
