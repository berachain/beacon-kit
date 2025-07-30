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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package types

import (
	"fmt"

	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/constraints"
	fastssz "github.com/ferranbt/fastssz"
	"github.com/sourcegraph/conc/iter"
)

// Compile-time check to ensure BlobSidecars implements the necessary interfaces.
var (
	_ constraints.SSZMarshallable = (*BlobSidecars)(nil)
	_ constraints.SSZRootable     = (*BlobSidecars)(nil)
)

// Sidecars is a slice of blob side cars to be included in the block.
type BlobSidecars []*BlobSidecar

// ValidateBlockRoots checks to make sure that
// all blobs in the sidecar are from the same block.
func (bs *BlobSidecars) ValidateBlockRoots() error {
	if bs == nil {
		return ErrAttemptedToVerifyNilSidecar
	}
	sidecars := *bs
	// We only need to check if there is more than
	// a single blob in the sidecar.
	if len(sidecars) > 1 {
		firstHtrBytes, err := sidecars[0].SignedBeaconBlockHeader.HashTreeRoot()
		if err != nil {
			return err
		}
		firstHtr := common.NewRootFromBytes(firstHtrBytes[:])
		for i := 1; i < len(sidecars); i++ {
			htrBytes, err := sidecars[i].SignedBeaconBlockHeader.HashTreeRoot()
			if err != nil {
				return err
			}
			if firstHtr != common.NewRootFromBytes(htrBytes[:]) {
				return ErrSidecarContainsDifferingBlockRoots
			}
		}
	}
	return nil
}

// VerifyInclusionProofs verifies the inclusion proofs for all sidecars.
func (bs *BlobSidecars) VerifyInclusionProofs() error {
	return errors.Join(iter.Map(
		*bs,
		func(sidecar **BlobSidecar) error {
			sc := *sidecar
			if sc == nil {
				return ErrAttemptedToVerifyNilSidecar
			}

			// Verify the KZG inclusion proof.
			if !sc.HasValidInclusionProof() {
				return ErrInvalidInclusionProof
			}
			return nil
		},
	)...)
}

// SizeSSZ returns the size of the BlobSidecars object in SSZ encoding.
func (bs *BlobSidecars) SizeSSZ() int {
	// BlobSidecar size: 8 + 131072 + 48 + 48 + 208 + 17*32 = 131720 bytes
	blobSidecarSize := 131720
	return 4 + len(*bs)*blobSidecarSize // offset + each blob sidecar
}

// MarshalSSZ marshals the BlobSidecars object to SSZ format.
func (bs *BlobSidecars) MarshalSSZ() ([]byte, error) {
	return bs.MarshalSSZTo(make([]byte, 0, bs.SizeSSZ()))
}

func (bs *BlobSidecars) ValidateAfterDecodingSSZ() error {
	if len(*bs) > constants.MaxBlobSidecarsPerBlock {
		return fmt.Errorf(
			"invalid number of blob sidecars, got %d max %d",
			len(*bs), constants.MaxBlobSidecarsPerBlock,
		)
	}
	return nil
}

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// MarshalSSZTo marshals the BlobSidecars object to a target array.
func (bs *BlobSidecars) MarshalSSZTo(dst []byte) ([]byte, error) {
	// Write offset
	offset := 4
	dst = fastssz.MarshalUint32(dst, uint32(offset))

	// Write sidecars
	for _, sidecar := range *bs {
		var err error
		dst, err = sidecar.MarshalSSZTo(dst)
		if err != nil {
			return nil, err
		}
	}

	return dst, nil
}

// UnmarshalSSZ ssz unmarshals the BlobSidecars object.
func (bs *BlobSidecars) UnmarshalSSZ(buf []byte) error {
	if len(buf) < 4 {
		return fastssz.ErrSize
	}

	// Read offset
	offset := fastssz.UnmarshallUint32(buf[0:4])
	if offset != 4 {
		return fastssz.ErrInvalidVariableOffset
	}

	// Calculate number of sidecars
	blobSidecarSize := 131720 // 8 + 131072 + 48 + 48 + 208 + 17*32
	remaining := len(buf) - 4
	if remaining%blobSidecarSize != 0 {
		return errors.New("invalid buffer size for blob sidecars")
	}

	count := remaining / blobSidecarSize
	if count > constants.MaxBlobSidecarsPerBlock {
		return fmt.Errorf("too many blob sidecars: %d > %d", count, constants.MaxBlobSidecarsPerBlock)
	}

	// Unmarshal each sidecar
	*bs = make(BlobSidecars, count)
	for i := 0; i < count; i++ {
		(*bs)[i] = &BlobSidecar{}
		start := 4 + i*blobSidecarSize
		end := start + blobSidecarSize
		if err := (*bs)[i].UnmarshalSSZ(buf[start:end]); err != nil {
			return err
		}
	}

	return bs.ValidateAfterDecodingSSZ()
}

// HashTreeRoot returns the hash tree root of the BlobSidecars.
func (bs *BlobSidecars) HashTreeRoot() ([32]byte, error) {
	hh := fastssz.DefaultHasherPool.Get()
	defer fastssz.DefaultHasherPool.Put(hh)
	if err := bs.HashTreeRootWith(hh); err != nil {
		return [32]byte{}, err
	}
	return hh.HashRoot()
}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// HashTreeRootWith ssz hashes the BlobSidecars object with a hasher.
func (bs *BlobSidecars) HashTreeRootWith(hh fastssz.HashWalker) error {
	indx := hh.Index()
	num := uint64(len(*bs))
	if num > constants.MaxBlobSidecarsPerBlock {
		return fastssz.ErrIncorrectListSize
	}
	for _, elem := range *bs {
		if err := elem.HashTreeRootWith(hh); err != nil {
			return err
		}
	}
	hh.MerkleizeWithMixin(indx, num, constants.MaxBlobSidecarsPerBlock)
	return nil
}

// GetTree ssz hashes the BlobSidecars object.
func (bs *BlobSidecars) GetTree() (*fastssz.Node, error) {
	return fastssz.ProofTree(bs)
}
