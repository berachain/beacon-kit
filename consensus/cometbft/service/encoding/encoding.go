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

package encoding

import (
	"fmt"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/consensus/cometbft/service/blobreactor"
	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/encoding/ssz"
	cmtabci "github.com/cometbft/cometbft/abci/types"
)

// UnmarshalBeaconBlockFromABCIRequest extracts a beacon block from an ABCI
// request.
func UnmarshalBeaconBlockFromABCIRequest(
	txs [][]byte,
	bzIndex uint,
	forkVersion common.Version,
) (*ctypes.SignedBeaconBlock, error) {
	var signedBlk *ctypes.SignedBeaconBlock
	lenTxs := uint(len(txs))

	// Ensure there are transactions in the request and that the request is
	// valid.
	if txs == nil || lenTxs == 0 {
		return signedBlk, ErrNoBeaconBlockInRequest
	}
	if bzIndex >= lenTxs {
		return signedBlk, ErrBzIndexOutOfBounds
	}

	// Extract the beacon block from the ABCI request.
	blkBz := txs[bzIndex]
	if blkBz == nil {
		return signedBlk, ErrNilBeaconBlockInRequest
	}

	block, err := ctypes.NewEmptySignedBeaconBlockWithVersion(forkVersion)
	if err != nil {
		return nil, fmt.Errorf("attempt at building block with wrong version %s: %w", forkVersion, err)
	}
	if err = ssz.Unmarshal(blkBz, block); err != nil {
		return nil, err
	}
	return block, nil
}

// UnmarshalBlobSidecarsFromABCIRequest extracts blob sidecars from an ABCI FinalizeBlockRequest based on blob consensus parameters.
func UnmarshalBlobSidecarsFromABCIRequest(req *cmtabci.FinalizeBlockRequest, cfg blobreactor.ConfigGetter) (datypes.BlobSidecars, error) {
	if req == nil {
		return nil, ErrNilABCIRequest
	}

	// If we are at or after BlobEnableHeight, then blobs must be retrieved from cache or via blobreactor p2p
	if cfg.IsBlobConsensusEnabledAtHeight(req.Height) {
		return datypes.BlobSidecars{}, nil
	}

	// Otherwise, blobs are in Txs[1] if present in this block
	txs := req.GetTxs()
	if len(txs) <= 1 {
		return datypes.BlobSidecars{}, nil
	}

	sidecarBz := txs[1]
	if len(sidecarBz) == 0 {
		return datypes.BlobSidecars{}, nil
	}

	var sidecars datypes.BlobSidecars
	if err := ssz.Unmarshal(sidecarBz, &sidecars); err != nil {
		return nil, fmt.Errorf("failed to unmarshal blobs from Txs[1]: %w", err)
	}

	return sidecars, nil
}

// ExtractBlobSidecarsFromRequest is a generic helper that extracts blob sidecars from either
// ProcessProposal or FinalizeBlock requests based on the blob consensus configuration.
// It handles the transition from storing blobs in Txs to using the BlobReactor.
func ExtractBlobSidecarsFromRequest(
	txs [][]byte,
	blobData []byte,
	height int64,
	cfg blobreactor.ConfigGetter,
) (datypes.BlobSidecars, error) {
	// If we are at or after BlobEnableHeight, then we use blobData for blobs if present
	if cfg.IsBlobConsensusEnabledAtHeight(height) {
		if len(blobData) == 0 {
			return datypes.BlobSidecars{}, nil
		}

		var sidecars datypes.BlobSidecars
		if err := ssz.Unmarshal(blobData, &sidecars); err != nil {
			return nil, fmt.Errorf("failed to unmarshal blobs from Blob field at height %d: %w", height, err)
		}

		return sidecars, nil
	}

	// Otherwise, blobs are in Txs[1] if present
	if len(txs) <= 1 {
		return datypes.BlobSidecars{}, nil
	}

	sidecarBz := txs[1]
	if len(sidecarBz) == 0 {
		return datypes.BlobSidecars{}, nil
	}

	var sidecars datypes.BlobSidecars
	if err := ssz.Unmarshal(sidecarBz, &sidecars); err != nil {
		return nil, fmt.Errorf("failed to unmarshal blobs from Txs[1] at height %d: %w", height, err)
	}

	return sidecars, nil
}
