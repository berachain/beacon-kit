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

	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
)

func TestFromUint32(t *testing.T) {
	tests := []struct {
		name     string
		input    uint32
		expected [4]byte
	}{
		{
			name:     "Phase0",
			input:    version.Phase0,
			expected: [4]byte{0, 0, 0, 0},
		},
		{
			name:     "Altair",
			input:    version.Altair,
			expected: [4]byte{1, 0, 0, 0},
		},
		{
			name:     "Bellatrix",
			input:    version.Bellatrix,
			expected: [4]byte{2, 0, 0, 0},
		},
		{
			name:     "Capella",
			input:    version.Capella,
			expected: [4]byte{3, 0, 0, 0},
		},
		{
			name:     "Deneb",
			input:    version.Deneb,
			expected: [4]byte{4, 0, 0, 0},
		},
		{
			name:     "Electra",
			input:    version.Electra,
			expected: [4]byte{5, 0, 0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := version.FromUint32[[4]byte](tt.input)
			if result != tt.expected {
				t.Errorf(
					"FromUint32(%d) = %v, expected %v",
					tt.input,
					result,
					tt.expected,
				)
			}
		})
	}
}

func TestFromUint32_CustomType(t *testing.T) {
	type CustomVersion [4]byte

	input := uint32(123456789)
	expected := CustomVersion{}
	binary.LittleEndian.PutUint32(expected[:], input)

	result := version.FromUint32[CustomVersion](input)
	if result != expected {
		t.Errorf("FromUint32(%d) = %v, expected %v", input, result, expected)
	}
}
