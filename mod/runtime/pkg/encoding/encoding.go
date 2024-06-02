// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package encoding

import (
	"reflect"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
)

// ExtractBlobsAndBlockFromRequest extracts the blobs and block from an ABCI
// request.
func ExtractBlobsAndBlockFromRequest[
	BeaconBlockT interface {
		ssz.Marshallable
		NewFromSSZ([]byte, uint32) (BeaconBlockT, error)
	}, BlobsSidecarsT ssz.Marshallable,
](
	req ABCIRequest,
	beaconBlkIndex uint,
	blobSidecarsIndex uint,
	forkVersion uint32,
) (BeaconBlockT, BlobsSidecarsT, error) {
	var (
		blobs BlobsSidecarsT
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

	blobs, err = UnmarshalBlobSidecarsFromABCIRequest[BlobsSidecarsT](
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
func UnmarshalBeaconBlockFromABCIRequest[BeaconBlockT interface {
	ssz.Marshallable
	NewFromSSZ([]byte, uint32) (BeaconBlockT, error)
}](
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
	T interface{ UnmarshalSSZ([]byte) error },
](
	req ABCIRequest,
	bzIndex uint,
) (T, error) {
	var sidecars T

	sidecars, ok := reflect.New(reflect.TypeOf(sidecars).Elem()).Interface().(T)
	if !ok {
		return sidecars, ErrInvalidType
	}

	if req == nil {
		return sidecars, ErrNilABCIRequest
	}

	txs := req.GetTxs()
	lenTxs := uint(len(txs))

	// Ensure there are transactions in the request and that the request is
	// valid.
	if txs == nil || lenTxs == 0 {
		return sidecars, ErrNoBeaconBlockInRequest
	}
	if bzIndex >= lenTxs {
		return sidecars, ErrBzIndexOutOfBounds
	}

	// Extract the blob sidecars from the ABCI request.
	sidecarBz := txs[bzIndex]
	if sidecarBz == nil {
		return sidecars, ErrNilBeaconBlockInRequest
	}

	err := sidecars.UnmarshalSSZ(sidecarBz)
	return sidecars, err
}
