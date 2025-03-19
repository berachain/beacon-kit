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

//nolint:mnd // todo fix.
package types

import (
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/karalabe/ssz"
	"github.com/sourcegraph/conc/iter"
)

// Compile-time check to ensure BlobSidecars implements the necessary interfaces.
var (
	_ ssz.DynamicObject                          = (*BlobSidecars)(nil)
	_ constraints.SSZMarshallable[*BlobSidecars] = (*BlobSidecars)(nil)
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
		firstHtr := sidecars[0].SignedBeaconBlockHeader.HashTreeRoot()
		for i := 1; i < len(sidecars); i++ {
			if firstHtr != sidecars[i].SignedBeaconBlockHeader.HashTreeRoot() {
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

// DefineSSZ defines the SSZ encoding for the BlobSidecars object.
// TODO: get from accessible chainspec field params.
func (bs *BlobSidecars) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineSliceOfStaticObjectsOffset(codec, (*[]*BlobSidecar)(bs), 6)
	ssz.DefineSliceOfStaticObjectsContent(codec, (*[]*BlobSidecar)(bs), 6)
}

// SizeSSZ returns the size of the BlobSidecars object in SSZ encoding.
func (bs *BlobSidecars) SizeSSZ(siz *ssz.Sizer, fixed bool) uint32 {
	if fixed {
		return 4
	}
	return 4 + ssz.SizeSliceOfStaticObjects(siz, *bs)
}

// MarshalSSZ marshals the BlobSidecars object to SSZ format.
func (bs *BlobSidecars) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(bs))
	return bs.MarshalSSZTo(buf)
}

// MarshalSSZTo marshals the BlobSidecars object to the provided buffer in SSZ
// format.
func (bs *BlobSidecars) MarshalSSZTo(buf []byte) ([]byte, error) {
	return buf, ssz.EncodeToBytes(buf, bs)
}

// NewFromSSZ unmarshals the BlobSidecars object from SSZ format.
func (bs *BlobSidecars) NewFromSSZ(buf []byte) (*BlobSidecars, error) {
	if bs == nil {
		bs = &BlobSidecars{}
	}
	return bs, ssz.DecodeFromBytes(buf, bs)
}
