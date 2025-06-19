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
)

func (s *Service) ParseBeaconBlock(req encoding.ABCIRequest) (
	*ctypes.SignedBeaconBlock,
	types.BlobSidecars,
	error,
) {
	if countTx := len(req.GetTxs()); countTx > MaxConsensusTxsCount {
		return nil, nil, fmt.Errorf("max expected %d, got %d: %w",
			MaxConsensusTxsCount, countTx,
			ErrTooManyConsensusTxs,
		)
	}

	forkVersion := s.chainSpec.ActiveForkVersionForTimestamp(math.U64(req.GetTime().Unix())) //#nosec: G115
	// Decode signed block and sidecars.
	signedBlk, sidecars, err := encoding.ExtractBlobsAndBlockFromRequest(
		req,
		BeaconBlockTxIndex,
		BlobSidecarsTxIndex,
		forkVersion,
	)
	if err != nil {
		return nil, nil, err
	}
	if signedBlk == nil {
		s.logger.Warn(
			"Aborting block verification - beacon block not found in proposal",
		)
		return nil, nil, ErrNilBlk
	}
	if sidecars == nil {
		s.logger.Warn(
			"Aborting block verification - blob sidecars not found in proposal",
		)
		return nil, nil, ErrNilBlob
	}

	return signedBlk, sidecars, nil
}
