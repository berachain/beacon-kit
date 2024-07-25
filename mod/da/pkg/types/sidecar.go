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
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	typesv2 "github.com/berachain/beacon-kit/mod/consensus-types/pkg/types/v2"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/merkle"
	"github.com/karalabe/ssz"
)

// BlobSidecar as per the Ethereum 2.0 specification:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/deneb/p2p-interface.md?ref=bankless.ghost.io#blobsidecar
//
//nolint:lll
type BlobSidecar struct {
	// Index represents the index of the blob in the block.
	Index uint64
	// Blob represents the blob data.
	Blob eip4844.Blob
	// KzgCommitment is the KZG commitment of the blob.
	KzgCommitment eip4844.KZGCommitment
	// Kzg proof allows folr the verification of the KZG commitment.
	KzgProof eip4844.KZGProof
	// BeaconBlockHeader represents the beacon block header for which this blob
	// is being included.
	BeaconBlockHeader *typesv2.BeaconBlockHeader
	// InclusionProof is the inclusion proof of the blob in the beacon block
	// body.
	InclusionProof [][32]byte
}

// BuildBlobSidecar creates a blob sidecar from the given blobs and
// beacon block.
func BuildBlobSidecar(
	index math.U64,
	header *types.BeaconBlockHeader,
	blob *eip4844.Blob,
	commitment eip4844.KZGCommitment,
	proof eip4844.KZGProof,
	inclusionProof [][32]byte,
) *BlobSidecar {
	return &BlobSidecar{
		Index:         index.Unwrap(),
		Blob:          *blob,
		KzgCommitment: commitment,
		KzgProof:      proof,
		BeaconBlockHeader: &typesv2.BeaconBlockHeader{
			Slot:            math.U64(header.Slot),
			ProposerIndex:   math.U64(header.ProposerIndex),
			ParentBlockRoot: header.ParentBlockRoot,
			StateRoot:       header.StateRoot,
			BodyRoot:        header.BodyRoot,
		},
		InclusionProof: inclusionProof,
	}
}

// HasValidInclusionProof verifies the inclusion proof of the
// blob in the beacon body.
func (b *BlobSidecar) HasValidInclusionProof(
	kzgOffset uint64,
) bool {
	// Calculate the hash tree root of the KZG commitment.
	leaf, err := b.KzgCommitment.HashTreeRoot()
	if err != nil {
		return false
	}

	gIndex := kzgOffset + b.Index

	// Verify the inclusion proof.
	result := merkle.IsValidMerkleBranch(
		leaf,
		b.InclusionProof,
		//#nosec:G701 // safe.
		uint8(
			len(b.InclusionProof),
		), // TODO: use KZG_INCLUSION_PROOF_DEPTH calculation.
		gIndex,
		b.BeaconBlockHeader.BodyRoot,
	)
	return result
}

// DefineSSZ defines the SSZ encoding for the BlobSidecar object.
func (b *BlobSidecar) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineUint64(codec, &b.Index)
	ssz.DefineStaticBytes(codec, &b.Blob)
	ssz.DefineStaticBytes(codec, &b.KzgCommitment)
	ssz.DefineStaticBytes(codec, &b.KzgProof)
	ssz.DefineStaticObject(codec, &b.BeaconBlockHeader)
	//nolint:mnd // depth of 8
	ssz.DefineCheckedArrayOfStaticBytes(codec, &b.InclusionProof, 8)
}

// SizeSSZ returns the size of the BlobSidecar object in SSZ encoding.
func (b *BlobSidecar) SizeSSZ() uint32 {
	return 8 + // Index
		131072 + // Blob
		48 + // KzgCommitment
		48 + // KzgProof
		112 + // BeaconBlockHeader
		8*32 // InclusionProof
}

// MarshalSSZ marshals the BlobSidecar object to SSZ format.
func (b *BlobSidecar) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, b.SizeSSZ())
	return buf, ssz.EncodeToBytes(buf, b)
}

// UnmarshalSSZ unmarshals the BlobSidecar object from SSZ format.
func (b *BlobSidecar) UnmarshalSSZ(buf []byte) error {
	if b.BeaconBlockHeader == nil {
		b.BeaconBlockHeader = &typesv2.BeaconBlockHeader{}
	}
	return ssz.DecodeFromBytes(buf, b)
}

// MarshalSSZTo marshals the BlobSidecar object to the provided buffer in SSZ
// format.
func (b *BlobSidecar) MarshalSSZTo(buf []byte) ([]byte, error) {
	return buf, ssz.EncodeToBytes(buf, b)
}

// HashTreeRoot computes the SSZ hash tree root of the BlobSidecar object.
func (b *BlobSidecar) HashTreeRoot() ([32]byte, error) {
	return ssz.HashSequential(b), nil
}
