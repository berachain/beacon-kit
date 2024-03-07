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

	"github.com/ethereum/go-ethereum/common"
	"github.com/sourcegraph/conc/iter"
)

// ConvertCommitmentToVersionedHash computes a SHA-256 hash of the given
// commitment and prefixes it with the BlobCommitmentVersion. This function is
// used to generate
// a versioned hash for KZG commitments.
//
// The restulting hash is intended for use in contexts where a versioned
// identifier
// for the commitment is required.
func ConvertCommitmentToVersionedHash(commitment [48]byte) common.Hash {
	hash := sha256.Sum256(commitment[:])
	// Prefix the hash with the BlobCommitmentVersion to create a versioned
	// hash.
	hash[0] = BlobCommitmentVersion
	return hash
}

// ConvertCommitmentsToVersionedHashes converts a slice of commitments to a
// slice of versioned hashes. This function is used to generate versioned hashes
// for KZG commitments.
//
// The resulting hashes are intended for use in contexts where a versioned
// identifier
// for the commitments is required.
func ConvertCommitmentsToVersionedHashes(commitments [][48]byte) []common.Hash {
	return iter.Map(commitments, func(bz *[48]byte) common.Hash {
		return ConvertCommitmentToVersionedHash(*bz)
	})
}
