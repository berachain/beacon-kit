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

package types_test

import (
	"strconv"
	"testing"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/errors"
	byteslib "github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/karalabe/ssz"
	"github.com/sourcegraph/conc/iter"
	"github.com/stretchr/testify/require"
)

// BlobSidecarsStruct is a slice of blob side cars to be included in the block.
type BlobSidecarsStruct struct {
	// Sidecars is a slice of blob side cars to be included in the block.
	Sidecars []*types.BlobSidecar
}

// NewBlobSidecars creates a new BlobSidecars object.
func (bs *BlobSidecarsStruct) Empty() *BlobSidecarsStruct {
	return &BlobSidecarsStruct{}
}

func (bs *BlobSidecarsStruct) Len() int {
	return len(bs.Sidecars)
}

func (bs *BlobSidecarsStruct) GetSidecars() []*types.BlobSidecar {
	return bs.Sidecars
}

func (bs *BlobSidecarsStruct) Get(index int) *types.BlobSidecar {
	return bs.Sidecars[index]
}

// IsNil checks to see if blobs are nil.
func (bs *BlobSidecarsStruct) IsNil() bool {
	return bs == nil || bs.Sidecars == nil
}

// ValidateBlockRoots checks to make sure that
// all blobs in the sidecar are from the same block.
func (bs *BlobSidecarsStruct) ValidateBlockRoots() error {
	// We only need to check if there is more than
	// a single blob in the sidecar.
	if sc := bs.Sidecars; len(sc) > 1 {
		firstHtr := sc[0].SignedBeaconBlockHeader.HashTreeRoot()
		for i := 1; i < len(sc); i++ {
			if firstHtr != sc[i].SignedBeaconBlockHeader.HashTreeRoot() {
				return types.ErrSidecarContainsDifferingBlockRoots
			}
		}
	}
	return nil
}

// VerifyInclusionProofs verifies the inclusion proofs for all sidecars.
func (bs *BlobSidecarsStruct) VerifyInclusionProofs(
	kzgOffset uint64,
) error {
	return errors.Join(iter.Map(
		bs.Sidecars,
		func(sidecar **types.BlobSidecar) error {
			sc := *sidecar
			if sc == nil {
				return types.ErrAttemptedToVerifyNilSidecar
			}

			// Verify the KZG inclusion proof.
			if !sc.HasValidInclusionProof(kzgOffset) {
				return types.ErrInvalidInclusionProof
			}
			return nil
		},
	)...)
}

// DefineSSZ defines the SSZ encoding for the BlobSidecarsStruct object.
func (bs *BlobSidecarsStruct) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineSliceOfStaticObjectsOffset(codec, &bs.Sidecars, 6)
	ssz.DefineSliceOfStaticObjectsContent(codec, &bs.Sidecars, 6)
}

// SizeSSZ returns the size of the BlobSidecarsStruct object in SSZ encoding.
func (bs *BlobSidecarsStruct) SizeSSZ(siz *ssz.Sizer, fixed bool) uint32 {
	if fixed {
		return 4
	}
	return 4 + ssz.SizeSliceOfStaticObjects(siz, bs.Sidecars)
}

// MarshalSSZ marshals the BlobSidecarsStruct object to SSZ format.
func (bs *BlobSidecarsStruct) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(bs))
	return bs.MarshalSSZTo(buf)
}

// MarshalSSZTo marshals the BlobSidecarsStruct object to the provided buffer in SSZ
// format.
func (bs *BlobSidecarsStruct) MarshalSSZTo(buf []byte) ([]byte, error) {
	return buf, ssz.EncodeToBytes(buf, bs)
}

// UnmarshalSSZ unmarshals the BlobSidecarsStruct object from SSZ format.
func (bs *BlobSidecarsStruct) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, bs)
}

func TestSidecarsMarshalling(t *testing.T) {
	sidecar := types.BuildBlobSidecar(
		math.U64(1),
		&ctypes.SignedBeaconBlockHeader{
			Header: &ctypes.BeaconBlockHeader{
				Slot:            math.Slot(1),
				ProposerIndex:   math.ValidatorIndex(3),
				ParentBlockRoot: common.Root{},
				StateRoot:       common.Root{},
			},
			Signature: crypto.BLSSignature{},
		},
		&eip4844.Blob{},
		eip4844.KZGCommitment{},
		[48]byte{},
		[]common.Root{},
	)

	// marshal BlobSidecar (defined as []*types.BlobSidecars)
	var sidecars = types.BlobSidecars{sidecar}
	bytes, err := sidecars.MarshalSSZ()
	require.NoError(t, err, "Marshalling BlobSidecar should not produce an error")

	// marshal BlobSidecarsStruct
	var sidecars2 = BlobSidecarsStruct{Sidecars: []*types.BlobSidecar{sidecar}}
	bytes2, err := sidecars2.MarshalSSZ()
	require.NoError(t, err, "Marshalling BlobSidecarsStruct should not produce an error")

	// the marshalled bytes should be equal between types.BlobSidecars and BlobSidecarsStruct
	require.Equal(t, bytes, bytes2, "The marshalled bytes should be equal between types.BlobSidecars and BlobSidecarsStruct")
}

func TestEmptySidecarMarshalling(t *testing.T) {
	// Create an empty BlobSidecar
	inclusionProof := make([]common.Root, 0)
	for i := int(1); i <= 8; i++ {
		it := byteslib.ExtendToSize([]byte(strconv.Itoa(i)), byteslib.B32Size)
		proof, err := byteslib.ToBytes32(it)
		require.NoError(t, err)
		inclusionProof = append(inclusionProof, common.Root(proof))
	}

	sidecar := types.BuildBlobSidecar(
		math.U64(0),
		&ctypes.SignedBeaconBlockHeader{
			Header:    &ctypes.BeaconBlockHeader{},
			Signature: crypto.BLSSignature{},
		},
		&eip4844.Blob{},
		eip4844.KZGCommitment{},
		[48]byte{},
		inclusionProof,
	)

	// Marshal the empty sidecar
	marshalled, err := sidecar.MarshalSSZ()
	require.NoError(
		t,
		err,
		"Marshalling empty sidecar should not produce an error",
	)
	require.NotNil(
		t,
		marshalled,
		"Marshalling empty sidecar should produce a result",
	)

	// Unmarshal the empty sidecar
	unmarshalled := &types.BlobSidecar{}
	err = unmarshalled.UnmarshalSSZ(marshalled)
	require.NoError(
		t,
		err,
		"Unmarshalling empty sidecar should not produce an error",
	)

	// Compare the original and unmarshalled empty sidecars
	require.Equal(
		t,
		sidecar,
		unmarshalled,
		"The original and unmarshalled empty sidecars should be equal",
	)
}

func TestValidateBlockRoots(t *testing.T) {
	// Create a sample BlobSidecar with valid roots
	inclusionProof := make([]common.Root, 0)
	for i := int(1); i <= 8; i++ {
		it := byteslib.ExtendToSize([]byte(strconv.Itoa(i)), byteslib.B32Size)
		proof, err := byteslib.ToBytes32(it)
		require.NoError(t, err)
		inclusionProof = append(inclusionProof, common.Root(proof))
	}

	validSidecar := types.BuildBlobSidecar(
		math.U64(0),
		&ctypes.SignedBeaconBlockHeader{
			Header: &ctypes.BeaconBlockHeader{
				StateRoot: [32]byte{1},
				BodyRoot:  [32]byte{2},
			},
			Signature: crypto.BLSSignature{},
		},
		&eip4844.Blob{},
		[48]byte{},
		[48]byte{},
		inclusionProof,
	)

	// Validate the sidecar with valid roots
	sidecars := types.BlobSidecars{
		validSidecar,
	}
	err := sidecars.ValidateBlockRoots()
	require.NoError(
		t,
		err,
		"Validating sidecar with valid roots should not produce an error",
	)

	// Create a sample BlobSidecar with invalid roots
	differentBlockRootSidecar := types.BuildBlobSidecar(
		math.U64(0),
		&ctypes.SignedBeaconBlockHeader{
			Header: &ctypes.BeaconBlockHeader{
				StateRoot: [32]byte{1},
				BodyRoot:  [32]byte{3},
			},
			Signature: crypto.BLSSignature{},
		},
		&eip4844.Blob{},
		eip4844.KZGCommitment{},
		eip4844.KZGProof{},
		inclusionProof,
	)
	// Validate the sidecar with invalid roots
	sidecarsInvalid := types.BlobSidecars{
		validSidecar,
		differentBlockRootSidecar,
	}
	err = sidecarsInvalid.ValidateBlockRoots()
	require.Error(
		t,
		err,
		"Validating sidecar with invalid roots should produce an error",
	)
}
