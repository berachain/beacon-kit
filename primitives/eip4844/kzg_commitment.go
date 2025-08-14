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

package eip4844

import (
	"fmt"

	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/crypto/sha256"
	"github.com/prysmaticlabs/gohashtree"
)

// KZGCommitment is a KZG commitment.
type KZGCommitment bytes.B48

// ToVersionedHash converts this KZG commitment into a versioned hash
// as per the Ethereum 2.0 specification:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/deneb/beacon-chain.md#kzg_commitment_to_versioned_hash
func (c KZGCommitment) ToVersionedHash() common.Hash32 {
	sum := sha256.Hash(c[:])
	// Prefix the hash with the BlobCommitmentVersion
	// to create a versioned hash.
	sum[0] = constants.BlobCommitmentVersion
	return sum
}

// ToHashChunks converts this KZG commitment into a set of hash chunks.
func (c KZGCommitment) ToHashChunks() [][32]byte {
	chunks := make([][32]byte, 2) //nolint:mnd // 2 chunks.
	copy(chunks[0][:], c[:constants.RootLength])
	copy(chunks[1][:], c[constants.RootLength:])
	gohashtree.HashChunks(chunks, chunks)
	return chunks
}

// HashTreeRoot returns the hash tree root of the commitment.
func (c KZGCommitment) HashTreeRoot() ([32]byte, error) {
	// B48 already has a HashTreeRoot method that returns B32
	// which is aliased to Hash32
	return [32]byte(bytes.B48(c).HashTreeRoot()), nil
}

// UnmarshalJSON parses a commitment in hex syntax.
func (c *KZGCommitment) UnmarshalJSON(input []byte) error {
	var b48 bytes.B48
	err := b48.UnmarshalJSON(input)
	if err != nil {
		return err
	}
	*c = KZGCommitment(b48)
	return nil
}

// MarshalText returns the hex representation of c.
func (c KZGCommitment) MarshalText() ([]byte, error) {
	return bytes.B48(c).MarshalText()
}

// MarshalSSZ implements the SSZ marshaling for KZGCommitment.
func (c KZGCommitment) MarshalSSZ() ([]byte, error) {
	return bytes.B48(c).MarshalSSZ()
}

// UnmarshalSSZ implements the SSZ unmarshaling for KZGCommitment.
func (c *KZGCommitment) UnmarshalSSZ(buf []byte) error {
	if len(buf) != 48 {
		return fmt.Errorf("invalid buffer length for KZGCommitment: expected 48, got %d", len(buf))
	}
	copy((*c)[:], buf)
	return nil
}

// SizeSSZ returns the size of the KZGCommitment in bytes.
func (c KZGCommitment) SizeSSZ() int {
	return 48
}

// KZGCommitments represents a slice of KZG commitments.
type KZGCommitments[HashT ~[32]byte] []KZGCommitment

// ToVersionedHashes converts the commitments to a set of
// versioned hashes. It is simplify a convenience method
// for converting a slice of commitments to a slice of
// versioned hashes.
func (c KZGCommitments[HashT]) ToVersionedHashes() []HashT {
	hashes := make([]HashT, len(c))
	for i, bz := range c {
		hashes[i] = HashT(bz.ToVersionedHash())
	}
	return hashes
}

// Leafify converts the commitments to a slice of leaves. Each leaf is the
// hash tree root of each commitment.
func (c KZGCommitments[HashT]) Leafify() []common.Root {
	leaves := make([]common.Root, len(c))
	for i, commitment := range c {
		// KZGCommitment.HashTreeRoot() never returns an error
		root, _ := commitment.HashTreeRoot()
		leaves[i] = common.Root(root)
	}
	return leaves
}
