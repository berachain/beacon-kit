// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package eip4844

import (
	"crypto/sha256"
	"reflect"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/prysmaticlabs/gohashtree"
)

// KZGCommitment is a KZG commitment.
type KZGCommitment [48]byte

// ToVersionedHash converts this KZG commitment into a versioned hash
// as per the Ethereum 2.0 specification:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/deneb/beacon-chain.md#kzg_commitment_to_versioned_hash
//
//nolint:lll
func (c KZGCommitment) ToVersionedHash() [32]byte {
	hash := sha256.Sum256(c[:])
	// Prefix the hash with the BlobCommitmentVersion
	// to create a versioned hash.
	hash[0] = constants.BlobCommitmentVersion
	return hash
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
		reflect.TypeOf(KZGCommitment{}),
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
