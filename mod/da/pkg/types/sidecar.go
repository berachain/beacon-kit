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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
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
	BeaconBlockHeader *types.BeaconBlockHeader
	// InclusionProof is the inclusion proof of the blob in the beacon block
	// body.
	InclusionProof []common.Root
}

// BuildBlobSidecar creates a blob sidecar from the given blobs and
// beacon block.
func BuildBlobSidecar[
	BeaconBlockHeaderT any,
](
	index math.U64,
	header BeaconBlockHeaderT,
	blob *eip4844.Blob,
	commitment eip4844.KZGCommitment,
	proof eip4844.KZGProof,
	inclusionProof []common.Root,
) *BlobSidecar {
	return &BlobSidecar{
		Index:             index.Unwrap(),
		Blob:              *blob,
		KzgCommitment:     commitment,
		KzgProof:          proof,
		BeaconBlockHeader: any(header).(*types.BeaconBlockHeader),
		InclusionProof:    inclusionProof,
	}
}

// HasValidInclusionProof verifies the inclusion proof of the
// blob in the beacon body.
func (b *BlobSidecar) HasValidInclusionProof(
	kzgOffset uint64,
) bool {
	// Verify the inclusion proof.
	return merkle.IsValidMerkleBranch(
		b.KzgCommitment.HashTreeRoot(),
		b.InclusionProof,
		//#nosec:G701 // safe.
		uint8(
			len(b.InclusionProof),
		), // TODO: use KZG_INCLUSION_PROOF_DEPTH calculation.
		kzgOffset+b.Index,
		b.BeaconBlockHeader.BodyRoot,
	)
}

func (b *BlobSidecar) GetIndex() uint64 {
	return b.Index
}

func (b *BlobSidecar) GetBlob() eip4844.Blob {
	return b.Blob
}

func (b *BlobSidecar) GetKzgProof() eip4844.KZGProof {
	return b.KzgProof
}

func (b *BlobSidecar) GetKzgCommitment() eip4844.KZGCommitment {
	return b.KzgCommitment
}

func (b *BlobSidecar) GetBeaconBlockHeader() *types.BeaconBlockHeader {
	return b.BeaconBlockHeader
}

func (b *BlobSidecar) GetInclusionProof() []common.Root {
	return b.InclusionProof
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
func (b *BlobSidecar) SizeSSZ(*ssz.Sizer) uint32 {
	return 8 + // Index
		131072 + // Blob
		48 + // KzgCommitment
		48 + // KzgProof
		112 + // BeaconBlockHeader
		8*32 // InclusionProof
}

// MarshalSSZ marshals the BlobSidecar object to SSZ format.
func (b *BlobSidecar) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(b))
	return buf, ssz.EncodeToBytes(buf, b)
}

// UnmarshalSSZ unmarshals the BlobSidecar object from SSZ format.
func (b *BlobSidecar) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, b)
}

// MarshalSSZTo marshals the BlobSidecar object to the provided buffer in SSZ
// format.
func (b *BlobSidecar) MarshalSSZTo(buf []byte) ([]byte, error) {
	return buf, ssz.EncodeToBytes(buf, b)
}

// HashTreeRoot computes the SSZ hash tree root of the BlobSidecar object.
func (b *BlobSidecar) HashTreeRoot() common.Root {
	return ssz.HashSequential(b)
}
