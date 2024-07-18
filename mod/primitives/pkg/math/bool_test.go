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

func TestBool(t *testing.T) {
	var b math.Bool
	require.Equal(t, constants.BoolSize, b.SizeSSZ())
	require.True(t, b.IsFixed())
	require.Equal(t, schema.Bool(), b.Type())
	require.Equal(t, uint64(0x1), b.ChunkCount())
}

func TestBool_MarshalSSZ(t *testing.T) {
	tests := []struct {
		input    math.Bool
		expected []byte
	}{
		{input: true, expected: []byte{1}},
		{input: false, expected: []byte{0}},
	}

	for _, tt := range tests {
		result, err := tt.input.MarshalSSZ()
		require.NoError(t, err)
		require.Equal(t, tt.expected, result)
	}
}

func TestBool_HashTreeRoot(t *testing.T) {
	tests := []struct {
		input    math.Bool
		expected [32]byte
	}{
		{input: true, expected: [32]byte{1}},
		{input: false, expected: [32]byte{}},
	}

	for _, tt := range tests {
		result, err := tt.input.HashTreeRoot()
		require.NoError(t, err)
		require.Equal(t, tt.expected, result)
	}
}

func TestBool_NewFromSSZ(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected math.Bool
		err      bool
	}{
		{name: "ValidTrue", input: []byte{1}, expected: true, err: false},
		{name: "ValidFalse", input: []byte{0}, expected: false, err: false},
		{name: "ValidTrue 2", input: []byte{2}, expected: true, err: false},
		{name: "InvalidLength", input: []byte{}, expected: false, err: true},
	}
	var inputBool math.Bool

	for _, tt := range tests {
		result, err := inputBool.NewFromSSZ(tt.input)
		if tt.err {
			require.Error(t, err, "Test name %s", tt.name)
		} else {
			require.NoError(t, err, "Test name %s", tt.name)
			require.Equal(t, tt.expected, result, "Test name %s", tt.name)
		}
	}
}
