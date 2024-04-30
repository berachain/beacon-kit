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

package eip4844_test

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	"github.com/stretchr/testify/require"
)

func TestKzgCommitmentToVersionedHash(t *testing.T) {
	commitment := eip4844.KZGCommitment{}
	copy(commitment[:], []byte("test commitment"))
	// Assuming BlobCommitmentVersion is a byte value
	expectedPrefix := constants.BlobCommitmentVersion

	hash := commitment.ToVersionedHash()
	if hash[0] != expectedPrefix {
		t.Errorf(
			"expected first byte of hash to be %v, got %v",
			expectedPrefix,
			hash[0],
		)
	}

	require.Len(t, hash, 32)
}

func TestKzgCommitmentsToVersionedHashHashes(t *testing.T) {
	commitments := make([]eip4844.KZGCommitment, 2)
	copy(commitments[0][:], "commitment 1")
	copy(commitments[1][:], "commitment 2")

	hashes := eip4844.KZGCommitments[[32]byte](commitments).ToVersionedHashes()

	if len(hashes) != len(commitments) {
		t.Errorf("expected %d hashes, got %d", len(commitments), len(hashes))
	}

	for i, hash := range hashes {
		if hash[0] != constants.BlobCommitmentVersion {
			t.Errorf(
				"expected first byte of hash %d to be %v, got %v",
				i,
				constants.BlobCommitmentVersion,
				hash[0],
			)
		}
	}
}
