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

package version_test

import (
	"encoding/binary"
	"testing"

	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/stretchr/testify/require"
)

func TestFromUint32(t *testing.T) {
	tests := []struct {
		name     string
		input    uint32
		expected common.Version
	}{
		{
			name:     "Phase0",
			input:    version.Phase0,
			expected: common.Version{0, 0, 0, 0},
		},
		{
			name:     "Altair",
			input:    version.Altair,
			expected: common.Version{1, 0, 0, 0},
		},
		{
			name:     "Bellatrix",
			input:    version.Bellatrix,
			expected: common.Version{2, 0, 0, 0},
		},
		{
			name:     "Capella",
			input:    version.Capella,
			expected: common.Version{3, 0, 0, 0},
		},
		{
			name:     "Deneb",
			input:    version.Deneb,
			expected: common.Version{4, 0, 0, 0},
		},
		{
			name:     "Deneb1",
			input:    version.Deneb1,
			expected: common.Version{4, 1, 0, 0},
		},
		{
			name:     "Electra",
			input:    version.Electra,
			expected: common.Version{5, 0, 0, 0},
		},
		{
			name:     "Electra1",
			input:    version.Electra1,
			expected: common.Version{5, 1, 0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := version.FromUint32(tt.input)
			require.Equal(t, tt.expected, result, "Test case: %s", tt.name)
		})
	}
}
func TestFromUint32_CustomType(t *testing.T) {
	input := uint32(123456789)
	expected := common.Version{}
	binary.LittleEndian.PutUint32(expected[:], input)

	result := version.FromUint32(input)
	require.Equal(t, expected, result)
}

func TestToUint32(t *testing.T) {
	tests := []struct {
		name     string
		input    common.Version
		expected uint32
	}{
		{
			name:     "Phase0",
			input:    common.Version{0, 0, 0, 0},
			expected: version.Phase0,
		},
		{
			name:     "Altair",
			input:    common.Version{1, 0, 0, 0},
			expected: version.Altair,
		},
		{
			name:     "Bellatrix",
			input:    common.Version{2, 0, 0, 0},
			expected: version.Bellatrix,
		},
		{
			name:     "Capella",
			input:    common.Version{3, 0, 0, 0},
			expected: version.Capella,
		},
		{
			name:     "Deneb",
			input:    common.Version{4, 0, 0, 0},
			expected: version.Deneb,
		},
		{
			name:     "Deneb1",
			input:    common.Version{4, 1, 0, 0},
			expected: version.Deneb1,
		},
		{
			name:     "Electra",
			input:    common.Version{5, 0, 0, 0},
			expected: version.Electra,
		},
		{
			name:     "Electra1",
			input:    common.Version{5, 1, 0, 0},
			expected: version.Electra1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := version.ToUint32(tt.input)
			require.Equal(t, tt.expected, result, "Test case: %s", tt.name)
		})
	}
}

func TestToUint32_CustomType(t *testing.T) {
	input := common.Version{0x15, 0xCD, 0x5B, 0x07}
	expected := uint32(123456789)

	result := version.ToUint32(input)
	require.Equal(t, expected, result)
}
