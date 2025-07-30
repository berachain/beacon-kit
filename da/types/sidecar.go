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
	fastssz "github.com/ferranbt/fastssz"
)

// Compile-time assertions to ensure BlobSidecar implements necessary interfaces.
var (
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
	commitmentRoot, err := b.KzgCommitment.HashTreeRoot()
	if err != nil {
		return false
	}
	return header != nil && merkle.IsValidMerkleBranch(
		common.Root(commitmentRoot),
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

// SizeSSZ returns the size of the BlobSidecar object in SSZ encoding.
func (b *BlobSidecar) SizeSSZ() int {
	return 8 + // Index
		131072 + // Blob
		48 + // KzgCommitment
		48 + // KzgProof
		208 + // SignedBeaconBlockHeader (112 + 96)
		ctypes.KZGInclusionProofDepth*32 // InclusionProof
}

// MarshalSSZ marshals the BlobSidecar object to SSZ format.
func (b *BlobSidecar) MarshalSSZ() ([]byte, error) {
	if len(b.InclusionProof) != ctypes.KZGInclusionProofDepth {
		return []byte{}, errors.New("invalid inclusion proof length")
	}
	return b.MarshalSSZTo(make([]byte, 0, b.SizeSSZ()))
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
func (b *BlobSidecar) MarshalSSZTo(dst []byte) ([]byte, error) {
	if len(b.InclusionProof) != ctypes.KZGInclusionProofDepth {
		return nil, errors.New("invalid inclusion proof length")
	}

	// Field (0) 'Index'
	dst = fastssz.MarshalUint64(dst, b.Index)

	// Field (1) 'Blob' (131072 bytes)
	dst = append(dst, b.Blob[:]...)

	// Field (2) 'KzgCommitment' (48 bytes)
	dst = append(dst, b.KzgCommitment[:]...)

	// Field (3) 'KzgProof' (48 bytes)
	dst = append(dst, b.KzgProof[:]...)

	// Field (4) 'SignedBeaconBlockHeader'
	if b.SignedBeaconBlockHeader == nil {
		return nil, errors.New("SignedBeaconBlockHeader is nil")
	}
	dst, err := b.SignedBeaconBlockHeader.MarshalSSZTo(dst)
	if err != nil {
		return nil, err
	}

	// Field (5) 'InclusionProof' (vector of 17 roots)
	for i := 0; i < ctypes.KZGInclusionProofDepth; i++ {
		dst = append(dst, b.InclusionProof[i][:]...)
	}

	return dst, nil
}

// HashTreeRoot computes the SSZ hash tree root of the BlobSidecar object.
func (b *BlobSidecar) HashTreeRoot() ([32]byte, error) {
	hh := fastssz.DefaultHasherPool.Get()
	defer fastssz.DefaultHasherPool.Put(hh)
	if err := b.HashTreeRootWith(hh); err != nil {
		return [32]byte{}, err
	}
	return hh.HashRoot()
}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// UnmarshalSSZ ssz unmarshals the BlobSidecar object.
func (b *BlobSidecar) UnmarshalSSZ(buf []byte) error {
	expectedSize := 8 + 131072 + 48 + 48 + 208 + ctypes.KZGInclusionProofDepth*32
	if len(buf) != expectedSize {
		return fastssz.ErrSize
	}

	var offset int

	// Field (0) 'Index'
	b.Index = fastssz.UnmarshallUint64(buf[offset : offset+8])
	offset += 8

	// Field (1) 'Blob' (131072 bytes)
	copy(b.Blob[:], buf[offset:offset+131072])
	offset += 131072

	// Field (2) 'KzgCommitment' (48 bytes)
	copy(b.KzgCommitment[:], buf[offset:offset+48])
	offset += 48

	// Field (3) 'KzgProof' (48 bytes)
	copy(b.KzgProof[:], buf[offset:offset+48])
	offset += 48

	// Field (4) 'SignedBeaconBlockHeader'
	if b.SignedBeaconBlockHeader == nil {
		b.SignedBeaconBlockHeader = &ctypes.SignedBeaconBlockHeader{}
	}
	if err := b.SignedBeaconBlockHeader.UnmarshalSSZ(buf[offset : offset+208]); err != nil {
		return err
	}
	offset += 208

	// Field (5) 'InclusionProof' (vector of 17 roots)
	b.InclusionProof = make([]common.Root, ctypes.KZGInclusionProofDepth)
	for i := 0; i < ctypes.KZGInclusionProofDepth; i++ {
		copy(b.InclusionProof[i][:], buf[offset:offset+32])
		offset += 32
	}

	return b.ValidateAfterDecodingSSZ()
}

// HashTreeRootWith ssz hashes the BlobSidecar object with a hasher.
func (b *BlobSidecar) HashTreeRootWith(hh fastssz.HashWalker) error {
	indx := hh.Index()

	// Field (0) 'Index'
	hh.PutUint64(b.Index)

	// Field (1) 'Blob' (131072 bytes)
	hh.PutBytes(b.Blob[:])

	// Field (2) 'KzgCommitment' (48 bytes)
	hh.PutBytes(b.KzgCommitment[:])

	// Field (3) 'KzgProof' (48 bytes)
	hh.PutBytes(b.KzgProof[:])

	// Field (4) 'SignedBeaconBlockHeader'
	if b.SignedBeaconBlockHeader == nil {
		return errors.New("SignedBeaconBlockHeader is nil")
	}
	if err := b.SignedBeaconBlockHeader.HashTreeRootWith(hh); err != nil {
		return err
	}

	// Field (5) 'InclusionProof' (vector of 17 roots)
	{
		if len(b.InclusionProof) != ctypes.KZGInclusionProofDepth {
			return fmt.Errorf("expected %d roots in InclusionProof, got %d", ctypes.KZGInclusionProofDepth, len(b.InclusionProof))
		}
		for i := 0; i < ctypes.KZGInclusionProofDepth; i++ {
			hh.PutBytes(b.InclusionProof[i][:])
		}
	}

	hh.Merkleize(indx)
	return nil
}

// GetTree ssz hashes the BlobSidecar object.
func (b *BlobSidecar) GetTree() (*fastssz.Node, error) {
	return fastssz.ProofTree(b)
}
