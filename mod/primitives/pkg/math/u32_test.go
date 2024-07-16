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

package math_test

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/schema"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/stretchr/testify/require"
)

func TestU32(t *testing.T) {
	var u math.U32
	require.Equal(t, constants.U32Size, u.SizeSSZ())
	require.True(t, u.IsFixed())
	require.Equal(t, schema.U32(), u.Type())
	require.Equal(t, uint64(1), u.ChunkCount())
}

func TestU32_MarshalSSZ(t *testing.T) {
	tests := []struct {
		input    math.U32
		expected []byte
	}{
		{input: 1, expected: []byte{1, 0, 0, 0}},
		{input: 0, expected: []byte{0, 0, 0, 0}},
	}

	for _, tt := range tests {
		result, err := tt.input.MarshalSSZ()
		if err != nil {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			require.Equal(t, tt.expected, result)
		}
	}
}

func TestU32_HashTreeRoot(t *testing.T) {
	tests := []struct {
		input    math.U32
		expected [32]byte
	}{
		{input: 1, expected: [32]byte{1}},
		{input: 0, expected: [32]byte{}},
	}

	for _, tt := range tests {
		result, err := tt.input.HashTreeRoot()
		if err != nil {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			require.Equal(t, tt.expected, result)
		}
	}
}

func TestU32_NewFromSSZ(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected math.U32
		err      bool
	}{
		{name: "Valid1", input: []byte{1, 0, 0, 0}, expected: 1, err: false},
		{name: "Valid0", input: []byte{0, 0, 0, 0}, expected: 0, err: false},
		{name: "InvalidLength", input: []byte{1, 0, 0}, expected: 0, err: true},
	}

	var u32 math.U32 = 1

	for _, tt := range tests {
		result, err := u32.NewFromSSZ(tt.input)
		if tt.err {
			require.Error(t, err, "Test name %s", tt.name)
		} else {
			require.NoError(t, err, "Test name %s", tt.name)
			require.Equal(t, tt.expected, result,
				"Test name %s", tt.name)
		}
	}
}
