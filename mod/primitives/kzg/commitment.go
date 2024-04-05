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

package kzg

import (
	"crypto/sha256"
	"reflect"

	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/constants"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/prysmaticlabs/gohashtree"
)

// Commitments represents a slice of KZG commitments.
type Commitments []Commitment

// CommitmentsFromBz converts byte slices to commitments.
func CommitmentsFromBz[T ~[]byte](slices []T) Commitments {
	commitments := make([]Commitment, len(slices))
	for i, slice := range slices {
		copy(commitments[i][:], slice)
	}
	return commitments
}

// ToVersionedHashes converts the commitments to a set of
// versioned hashes. It is simplify a convenience method
// for converting a slice of commitments to a slice of
// versioned hashes.
func (c Commitments) ToVersionedHashes() []primitives.ExecutionHash {
	hashes := make([]primitives.ExecutionHash, len(c))
	for i, bz := range c {
		hashes[i] = bz.ToVersionedHash()
	}
	return hashes
}

// Commitment is a KZG commitment.
type Commitment [48]byte

// ToVersionedHash converts this KZG commitment into a versioned hash
// as per the Ethereum 2.0 specification:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/deneb/beacon-chain.md#kzg_commitment_to_versioned_hash
//
//nolint:lll
func (c Commitment) ToVersionedHash() primitives.ExecutionHash {
	hash := sha256.Sum256(c[:])
	// Prefix the hash with the BlobCommitmentVersion
	// to create a versioned hash.
	hash[0] = constants.BlobCommitmentVersion
	return hash
}

// ToHashChunks converts this KZG commitment into a set of hash chunks.
func (c Commitment) ToHashChunks() [][32]byte {
	chunks := make([][32]byte, 2) //nolint:gomnd // 2 chunks.
	copy(chunks[0][:], c[:])
	copy(chunks[1][:], c[constants.RootLength:])
	gohashtree.HashChunks(chunks, chunks)
	return chunks
}

// UnmarshalJSON parses a commitment in hex syntax.
func (c *Commitment) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(reflect.TypeOf(Commitment{}), input, c[:])
}

// MarshalText returns the hex representation of c.
func (c Commitment) MarshalText() ([]byte, error) {
	return hexutil.Bytes(c[:]).MarshalText()
}
