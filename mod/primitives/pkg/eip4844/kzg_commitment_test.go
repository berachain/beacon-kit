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
	"encoding/hex"
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	"github.com/stretchr/testify/require"
)

func TestKzgCommitmentToVersionedHash(t *testing.T) {
	commitment := newTestCommitment("test commitment")
	expectedPrefix := constants.BlobCommitmentVersion

	hash := commitment.ToVersionedHash()
	require.Equal(t, expectedPrefix, hash[0],
		"First byte of hash should match BlobCommitmentVersion")
	require.Len(t, hash, 32, "Hash length should be 32 bytes")
}

func TestKzgCommitmentsToVersionedHashHashes(t *testing.T) {
	commitments := []eip4844.KZGCommitment{
		newTestCommitment("commitment 1"),
		newTestCommitment("commitment 2"),
	}

	hashes := eip4844.KZGCommitments[[32]byte](commitments).ToVersionedHashes()
	require.Len(t, hashes, len(commitments),
		"Number of hashes should match number of commitments")

	for i, hash := range hashes {
		require.Equal(t, constants.BlobCommitmentVersion, hash[0],
			"First byte of hash %d should match BlobCommitmentVersion", i)
	}
}

func TestKZGCommitmentToHashChunks(t *testing.T) {
	tests := []struct {
		name     string
		input    eip4844.KZGCommitment
		expected int
	}{
		{"Valid input",
			newTestCommitment("example commitment data that " +
				"exceeds root length to test chunking"),
			2},
		{"Short input", newTestCommitment("short"), 2},
		{"Empty input", eip4844.KZGCommitment{}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chunks := tt.input.ToHashChunks()
			require.Len(t, chunks, tt.expected,
				"Incorrect number of chunks for test: "+tt.name)
		})
	}
}

func TestKZGCommitmentHashTreeRoot(t *testing.T) {
	tests := []struct {
		name     string
		input    eip4844.KZGCommitment
		expected [32]byte
	}{
		{"Simple input", newTestCommitment("example commitment"),
			[32]byte{138, 20, 122, 217, 77, 116, 246, 111, 195, 118, 240,
				67, 111, 145, 176, 117, 67, 82, 153, 245, 152, 25, 235, 239, 171,
				54, 148, 169, 30, 169, 167, 229}},
		{"Empty input", eip4844.KZGCommitment{},
			[32]byte{245, 165, 253, 66, 209, 106, 32, 48, 39, 152, 239,
				110, 211, 9, 151, 155, 67, 0, 61, 35, 32, 217, 240, 232, 234, 152,
				49, 169, 39, 89, 251, 75}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashTreeRoot, err := tt.input.HashTreeRoot()
			require.NoError(t, err)
			require.Equal(
				t,
				tt.expected,
				hashTreeRoot,
				"Hash tree root does not "+
					"match expected for test: "+tt.name,
			)
		})
	}
}

func TestKZGCommitmentUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    eip4844.KZGCommitment
		shouldError bool
	}{
		{
			name: "Valid hex input",
			input: `"0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789` +
				`abcdef0123456789abcdef0123456789abcdef"`,
			expected: func() eip4844.KZGCommitment {
				var c eip4844.KZGCommitment
				data, _ := hex.DecodeString(
					"0123456789abcdef0123456789abcdef0" +
						"123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
				)
				copy(c[:], data)
				return c
			}(),
			shouldError: false,
		},
		{
			name:        "Invalid hex input",
			input:       `"0xG123456789abcdef"`,
			expected:    eip4844.KZGCommitment{},
			shouldError: true,
		},
		{
			name:        "Empty input",
			input:       `""`,
			expected:    eip4844.KZGCommitment{},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var commitment eip4844.KZGCommitment
			err := commitment.UnmarshalJSON([]byte(tt.input))
			if tt.shouldError {
				require.Error(t, err, "Expected an error for test: "+tt.name)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, commitment, "Unmarshaled commitment does "+
					"not match expected for test: "+tt.name)
			}
		})
	}
}

func TestKZGCommitment_MarshalText(t *testing.T) {
	testCases := []struct {
		name     string
		input    eip4844.KZGCommitment
		expected string
	}{
		{
			name:  "Empty Commitment",
			input: eip4844.KZGCommitment{},
			expected: "3078303030303030303030303030303030303030303030303030303030" +
				"3030303030303030303030303030303030303030303030303030303030303030" +
				"3030303030303030303030303030303030303030303030303030303030303030" +
				"3030303030",
		},
		{
			name: "Non-Empty Commitment",
			input: func() eip4844.KZGCommitment {
				var c eip4844.KZGCommitment
				for i := range c {
					c[i] = byte(i % 256)
				}
				return c
			}(),
			expected: "30783030303130323033303430353036303730383039306130623063306" +
				"43065306631303131313231333134313531363137313831393161316231633164" +
				"31653166323032313232323332343235323632373238323932613262326332643" +
				"2653266",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := tc.input.MarshalText()
			require.NoError(t, err)
			require.Equal(t, tc.expected, hex.EncodeToString(output),
				"Test case: %s", tc.name)
		})
	}
}

func TestKZGCommitments_Leafify(t *testing.T) {
	tests := []struct {
		name  string
		input []eip4844.KZGCommitment
	}{
		{
			name: "Single Commitment",
			input: []eip4844.KZGCommitment{
				newTestCommitment("single commitment"),
			},
		},
		{
			name: "Multiple Commitments",
			input: []eip4844.KZGCommitment{
				newTestCommitment("commitment one"),
				newTestCommitment("commitment two"),
			},
		},
		{
			name:  "Empty Commitments",
			input: []eip4844.KZGCommitment{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Dynamically compute expected values based on input
			expected := make([][32]byte, len(tt.input))
			for i, commitment := range tt.input {
				expected[i] = commitment.ToHashChunks()[0]
			}

			commitments := eip4844.KZGCommitments[[32]byte](tt.input)
			leaves, err := commitments.Leafify()
			require.NoError(t, err)
			require.Equal(t, expected, leaves,
				"Leaves do not match expected for test: "+tt.name)
		})
	}
}

func newTestCommitment(data string) eip4844.KZGCommitment {
	var c eip4844.KZGCommitment
	copy(c[:], data)
	return c
}
