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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
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

func TestU256(t *testing.T) {
	var u math.U256
	require.Equal(t, constants.U256Size, u.SizeSSZ())
	require.True(t, u.IsFixed())
	require.Equal(t, schema.U256(), u.Type())
	require.Equal(t, uint64(1), u.ChunkCount())
}

func TestU256_MarshalSSZ(t *testing.T) {
	tests := []struct {
		input    *math.U256
		expected []byte
	}{
		{
			input:    math.NewU256FromUint64(1),
			expected: append([]byte{1}, make([]byte, constants.U256Size-1)...)},
		{
			input:    math.NewU256FromUint64(0),
			expected: make([]byte, constants.U256Size)},
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

func TestU256_HashTreeRoot(t *testing.T) {
	tests := []struct {
		input    *math.U256
		expected [32]byte
	}{
		{input: math.NewU256FromUint64(1), expected: [32]byte{1}},
		{input: math.NewU256FromUint64(0), expected: [32]byte{}},
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

func TestU256_NewFromSSZ(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected *math.U256
		err      bool
	}{
		{
			name:     "Valid1",
			input:    append([]byte{1}, make([]byte, constants.U256Size-1)...),
			expected: math.NewU256FromUint64(1), err: false},
		{
			name:     "Valid0",
			input:    make([]byte, constants.U256Size),
			expected: math.NewU256FromUint64(0), err: false},
		{
			name:     "InvalidLength",
			input:    []byte{1, 0, 0},
			expected: nil, err: true},
	}

	var u256 math.U256

	for _, tt := range tests {
		result, err := u256.NewFromSSZ(tt.input)
		if tt.err {
			require.Error(t, err, "Test name %s", tt.name)
		} else {
			require.NoError(t, err, "Test name %s", tt.name)
			require.Equal(t, tt.expected, result,
				"Test name %s", tt.name)
		}
	}
}
