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

package eip4844

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto/sha256"
	"github.com/prysmaticlabs/gohashtree"
)

// KZGCommitment is a KZG commitment.
type KZGCommitment [48]byte

// ToVersionedHash converts this KZG commitment into a versioned hash
// as per the Ethereum 2.0 specification:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/deneb/beacon-chain.md#kzg_commitment_to_versioned_hash
//
//nolint:lll // link.
func (c KZGCommitment) ToVersionedHash() [32]byte {
	sum := sha256.Hash(c[:])
	// Prefix the hash with the BlobCommitmentVersion
	// to create a versioned hash.
	sum[0] = constants.BlobCommitmentVersion
	return sum
}

// ToHashChunks converts this KZG commitment into a set of hash chunks.
func (c KZGCommitment) ToHashChunks() [][32]byte {
	chunks := make([][32]byte, 2) //nolint:mnd // 2 chunks.
	copy(chunks[0][:], c[:])
	copy(chunks[1][:], c[constants.RootLength:])
	gohashtree.HashChunks(chunks, chunks)
	return chunks
}

// HashTreeRoot returns the hash tree root of the commitment.
func (c KZGCommitment) HashTreeRoot() ([32]byte, error) {
	chunks := c.ToHashChunks()
	return chunks[0], nil
}

// UnmarshalJSON parses a commitment in hex syntax.
func (c *KZGCommitment) UnmarshalJSON(input []byte) error {
	return bytes.UnmarshalFixedJSON(
		input,
		c[:],
	)
}

// MarshalText returns the hex representation of c.
func (c KZGCommitment) MarshalText() ([]byte, error) {
	return bytes.Bytes(c[:]).MarshalText()
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
		hashes[i] = bz.ToVersionedHash()
	}
	return hashes
}

// Leafify converts the commitments to a slice of leaves.
func (c KZGCommitments[HashT]) Leafify() [][32]byte {
	leaves := make([][32]byte, len(c))
	for i, commitment := range c {
		leaves[i] = commitment.ToHashChunks()[0]
	}
	return leaves
}
