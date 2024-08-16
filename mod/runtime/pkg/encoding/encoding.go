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

package encoding

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
)

// ExtractBlobsAndBlockFromRequest extracts the blobs and block from an ABCI
// request.
func ExtractBlobsAndBlockFromRequest[
	BeaconBlockT BeaconBlock[BeaconBlockT],
	BlobSidecarsT interface {
		constraints.SSZUnmarshaler
		Empty() BlobSidecarsT
	},
](
	req ABCIRequest,
	beaconBlkIndex uint,
	blobSidecarsIndex uint,
	forkVersion uint32,
) (BeaconBlockT, BlobSidecarsT, error) {
	var (
		blobs BlobSidecarsT
		blk   BeaconBlockT
	)

	if req == nil {
		return blk, blobs, ErrNilABCIRequest
	}

	blk, err := UnmarshalBeaconBlockFromABCIRequest[BeaconBlockT](
		req,
		beaconBlkIndex,
		forkVersion,
	)
	if err != nil {
		return blk, blobs, err
	}

	blobs, err = UnmarshalBlobSidecarsFromABCIRequest[BlobSidecarsT](
		req,
		blobSidecarsIndex,
	)
	if err != nil {
		return blk, blobs, err
	}

	return blk, blobs, nil
}

// UnmarshalBeaconBlockFromABCIRequest extracts a beacon block from an ABCI
// request.
func UnmarshalBeaconBlockFromABCIRequest[
	BeaconBlockT BeaconBlock[BeaconBlockT],
](
	req ABCIRequest,
	bzIndex uint,
	forkVersion uint32,
) (BeaconBlockT, error) {
	var blk BeaconBlockT
	if req == nil {
		return blk, ErrNilABCIRequest
	}

	txs := req.GetTxs()
	lenTxs := uint(len(txs))

	// Ensure there are transactions in the request and that the request is
	// valid.
	if txs == nil || lenTxs == 0 {
		return blk, ErrNoBeaconBlockInRequest
	}
	if bzIndex >= lenTxs {
		return blk, ErrBzIndexOutOfBounds
	}

	// Extract the beacon block from the ABCI request.
	blkBz := txs[bzIndex]
	if blkBz == nil {
		return blk, ErrNilBeaconBlockInRequest
	}

	return blk.NewFromSSZ(blkBz, forkVersion)
}

// UnmarshalBlobSidecarsFromABCIRequest extracts blob sidecars from an ABCI
// request.
func UnmarshalBlobSidecarsFromABCIRequest[
	BlobSidecarsT interface {
		constraints.SSZUnmarshaler
		Empty() BlobSidecarsT
	},
](
	req ABCIRequest,
	bzIndex uint,
) (BlobSidecarsT, error) {
	var sidecars BlobSidecarsT
	if req == nil {
		return sidecars, ErrNilABCIRequest
	}

	txs := req.GetTxs()
	if len(txs) == 0 || bzIndex >= uint(len(txs)) {
		return sidecars, ErrNoBeaconBlockInRequest
	}

	sidecarBz := txs[bzIndex]
	if sidecarBz == nil {
		return sidecars, ErrNilBeaconBlockInRequest
	}

	// TODO: Do some research to figure out how to make this more
	// elegant.
	sidecars = sidecars.Empty()
	return sidecars, sidecars.UnmarshalSSZ(sidecarBz)
}
