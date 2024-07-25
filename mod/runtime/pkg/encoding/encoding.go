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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package encoding

import (
	"fmt"
	"reflect"
	"runtime/debug"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
)

// ExtractBlobsAndBlockFromRequest extracts the blobs and block from an ABCI
// request.
func ExtractBlobsAndBlockFromRequest[
	BeaconBlockT BeaconBlock[BeaconBlockT],
	BlobSidecarsT constraints.SSZMarshallableDynamic,
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

	fmt.Println("ABC")

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
	b, err := blk.NewFromSSZ(blkBz, forkVersion)
	return b, err

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
		fmt.Printf("Debug: Invalid type conversion for sidecars\n")
		return sidecars, ErrInvalidType
	}

	if req == nil {
		fmt.Printf("Debug: Received nil ABCI request\n")
		return sidecars, ErrNilABCIRequest
	}

	txs := req.GetTxs()
	lenTxs := uint(len(txs))
	fmt.Printf("Debug: Number of transactions in request: %d\n", lenTxs)

	// Ensure there are transactions in the request and that the request is
	// valid.
	if txs == nil || lenTxs == 0 {
		fmt.Printf("Debug: No transactions found in request\n")
		return sidecars, ErrNoBeaconBlockInRequest
	}
	if bzIndex >= lenTxs {
		fmt.Printf("Debug: bzIndex (%d) out of bounds. Max index: %d\n", bzIndex, lenTxs-1)
		return sidecars, ErrBzIndexOutOfBounds
	}

	// Extract the blob sidecars from the ABCI request.
	sidecarBz := txs[bzIndex]
	if sidecarBz == nil {
		fmt.Printf("Debug: Nil beacon block found at index %d\n", bzIndex)
		return sidecars, ErrNilBeaconBlockInRequest
	}

	var err error
	fmt.Printf("Debug: Attempting to unmarshal sidecar of length %d\n", len(sidecarBz))
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Debug: Panic occurred during UnmarshalSSZ: %v\n", r)
				fmt.Printf("Debug: Call stack:\n%s\n", debug.Stack())
				err = fmt.Errorf("panic during UnmarshalSSZ: %v", r)
			}
		}()
		err = sidecars.UnmarshalSSZ(sidecarBz)
	}()
	if err != nil {
		fmt.Printf("Debug: Error unmarshalling sidecars: %v\n", err)
	} else {
		fmt.Printf("Debug: Successfully unmarshalled sidecars\n")
	}
	return sidecars, err
}
