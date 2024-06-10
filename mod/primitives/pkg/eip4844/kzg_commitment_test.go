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
