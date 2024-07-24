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

package types

import (
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/sourcegraph/conc/iter"
)

// BlobSidecars is a slice of blob side cars to be included in the block.
type BlobSidecars struct {
	// Sidecars is a slice of blob side cars to be included in the block.
	Sidecars []*BlobSidecar `ssz-max:"6"`
}

// IsNil checks to see if blobs are nil.
func (bs *BlobSidecars) IsNil() bool {
	return bs == nil || bs.Sidecars == nil
}

// ValidateBlockRoots checks to make sure that
// all blobs in the sidecar are from the same block.
func (bs *BlobSidecars) ValidateBlockRoots() error {
	// We only need to check if there is more than
	// a single blob in the sidecar.
	if sc := bs.Sidecars; len(sc) > 1 {
		firstHtr, err := sc[0].BeaconBlockHeader.HashTreeRoot()
		if err != nil {
			return err
		}

		var nextHtr [32]byte
		for i := 1; i < len(sc); i++ {
			nextHtr, err = sc[i].BeaconBlockHeader.HashTreeRoot()
			if err != nil {
				return err
			}
			if firstHtr != nextHtr {
				return ErrSidecarContainsDifferingBlockRoots
			}
		}
	}
	return nil
}

// VerifyInclusionProofs verifies the inclusion proofs for all sidecars.
func (bs *BlobSidecars) VerifyInclusionProofs(
	kzgOffset uint64,
) error {
	return errors.Join(iter.Map(
		bs.Sidecars,
		func(sidecar **BlobSidecar) error {
			sc := *sidecar
			if sc == nil {
				return ErrAttemptedToVerifyNilSidecar
			}

			// Verify the KZG inclusion proof.
			if !sc.HasValidInclusionProof(kzgOffset) {
				return ErrInvalidInclusionProof
			}
			return nil
		},
	)...)
}

// Len returns the number of sidecars in the sidecar.
func (bs *BlobSidecars) Len() int {
	return len(bs.Sidecars)
}
