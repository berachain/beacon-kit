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

package types

import (
	"fmt"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/merkle"
	"github.com/karalabe/ssz"
)

// Compile-time assertions to ensure BlobSidecar implements necessary interfaces.
var (
	_ ssz.StaticObject                    = (*BlobSidecar)(nil)
	_ constraints.SSZMarshallableRootable = (*BlobSidecar)(nil)
)

// BlobSidecar as per the Ethereum 2.0 specification:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/deneb/p2p-interface.md#blobsidecar
//
// NOTE: This struct is only ever (un)marshalled with SSZ and NOT with JSON.
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
	SignedBeaconBlockHeader *ctypes.SignedBeaconBlockHeader
	// InclusionProof is the inclusion proof of the blob in the beacon block
	// body.
	InclusionProof []common.Root
}

// BuildBlobSidecar creates a blob sidecar from the given blobs and
// beacon block.
func BuildBlobSidecar(
	index math.U64,
	header *ctypes.SignedBeaconBlockHeader,
	blob *eip4844.Blob,
	commitment eip4844.KZGCommitment,
	proof eip4844.KZGProof,
	inclusionProof []common.Root,
) *BlobSidecar {
	return &BlobSidecar{
		Index:                   index.Unwrap(),
		Blob:                    *blob,
		KzgCommitment:           commitment,
		KzgProof:                proof,
		SignedBeaconBlockHeader: header,
		InclusionProof:          inclusionProof,
	}
}

// HasValidInclusionProof verifies the inclusion proof of the
// blob in the beacon body.
func (b *BlobSidecar) HasValidInclusionProof() bool {
	header := b.GetBeaconBlockHeader()
	return header != nil && merkle.IsValidMerkleBranch(
		b.KzgCommitment.HashTreeRoot(),
		b.InclusionProof,
		ctypes.KZGInclusionProofDepth,
		ctypes.KZGOffset+b.Index,
		header.BodyRoot,
	)
}

/* -------------------------------------------------------------------------- */
/*                                   Getters                                  */
/* -------------------------------------------------------------------------- */

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

func (b *BlobSidecar) GetBeaconBlockHeader() *ctypes.BeaconBlockHeader {
	return b.SignedBeaconBlockHeader.Header
}

func (b *BlobSidecar) GetInclusionProof() []common.Root {
	return b.InclusionProof
}

func (b *BlobSidecar) GetSignature() crypto.BLSSignature {
	return b.SignedBeaconBlockHeader.Signature
}

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// DefineSSZ defines the SSZ encoding for the BlobSidecar object.
func (b *BlobSidecar) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineUint64(codec, &b.Index)
	ssz.DefineStaticBytes(codec, &b.Blob)
	ssz.DefineStaticBytes(codec, &b.KzgCommitment)
	ssz.DefineStaticBytes(codec, &b.KzgProof)
	ssz.DefineStaticObject(codec, &b.SignedBeaconBlockHeader)
	ssz.DefineCheckedArrayOfStaticBytes(codec, &b.InclusionProof, ctypes.KZGInclusionProofDepth)
}

// SizeSSZ returns the size of the BlobSidecar object in SSZ encoding.
// TODO: get from accessible chainspec field params.
func (b *BlobSidecar) SizeSSZ(sizer *ssz.Sizer) uint32 {
	ssize := (*ctypes.SignedBeaconBlockHeader)(nil).SizeSSZ(sizer)
	return 8 + // Index
		131072 + // Blob
		48 + // KzgCommitment
		48 + // KzgProof
		ssize + // SignedBeaconBlockHeader
		ctypes.KZGInclusionProofDepth*32 // InclusionProof
}

// MarshalSSZ marshals the BlobSidecar object to SSZ format.
func (b *BlobSidecar) MarshalSSZ() ([]byte, error) {
	if len(b.InclusionProof) != ctypes.KZGInclusionProofDepth {
		return []byte{}, errors.New("invalid inclusion proof length")
	}
	buf := make([]byte, ssz.Size(b))
	return buf, ssz.EncodeToBytes(buf, b)
}

func (b *BlobSidecar) ValidateAfterDecodingSSZ() error {
	// Verify inclusion proof length
	if len(b.InclusionProof) != ctypes.KZGInclusionProofDepth {
		return fmt.Errorf("invalid inclusion proof length, got %d, expect %d",
			b.InclusionProof,
			ctypes.KZGInclusionProofDepth,
		)
	}

	// Ensure SignedBeaconBlockHeader is not nil
	if b.SignedBeaconBlockHeader == nil {
		b.SignedBeaconBlockHeader = &ctypes.SignedBeaconBlockHeader{}
	}
	return nil
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
