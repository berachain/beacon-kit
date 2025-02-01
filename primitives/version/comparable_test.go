// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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
	"slices"
	"testing"

	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/stretchr/testify/require"
)

func TestIsBefore(t *testing.T) {
	tests := []struct {
		name     string
		a, b     common.Version
		expected bool
	}{
		{
			name:     "equal versions",
			a:        common.Version{1, 0, 0, 0},
			b:        common.Version{1, 0, 0, 0},
			expected: false,
		},
		{
			name:     "clearly before",
			a:        common.Version{1, 0, 0, 0},
			b:        common.Version{2, 0, 0, 0},
			expected: true,
		},
		{
			name:     "clearly after",
			a:        common.Version{2, 0, 0, 0},
			b:        common.Version{1, 0, 0, 0},
			expected: false,
		},
		{
			name:     "minor version before",
			a:        common.Version{1, 1, 0, 0},
			b:        common.Version{1, 2, 0, 0},
			expected: true,
		},
		{
			name:     "patch version before",
			a:        common.Version{1, 1, 1, 0},
			b:        common.Version{1, 1, 2, 0},
			expected: true,
		},
		{
			name:     "beacon version before",
			a:        version.Phase0(),
			b:        version.Altair(),
			expected: true,
		},
		{
			name:     "beacon version before minor",
			a:        version.Deneb(),
			b:        version.Deneb1(),
			expected: true,
		},
		{
			name:     "beacon version after minor",
			a:        common.Version{0x05, 0x00, 0x00, 0x00},
			b:        common.Version{0x04, 0x01, 0x00, 0x00},
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Expecto fail
			result := version.IsAfter(tc.a, tc.b)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestEquals(t *testing.T) {
	tests := []struct {
		name     string
		a, b     common.Version
		expected bool
	}{
		{
			name:     "empty versions",
			a:        common.Version{},
			b:        common.Version{},
			expected: true,
		},
		{
			name:     "equal versions",
			a:        common.Version{1, 0, 0, 0},
			b:        common.Version{1, 0, 0, 0},
			expected: true,
		},
		{
			name:     "different major version",
			a:        common.Version{1, 0, 0, 0},
			b:        common.Version{2, 0, 0, 0},
			expected: false,
		},
		{
			name:     "different minor version",
			a:        common.Version{1, 1, 0, 0},
			b:        common.Version{1, 2, 0, 0},
			expected: false,
		},
		{
			name:     "different patch version",
			a:        common.Version{1, 1, 1, 0},
			b:        common.Version{1, 1, 2, 0},
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := version.Equals(tc.a, tc.b)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestIsAfter(t *testing.T) {
	tests := []struct {
		name     string
		a, b     common.Version
		expected bool
	}{
		{
			name:     "equal versions",
			a:        common.Version{1, 0, 0, 0},
			b:        common.Version{1, 0, 0, 0},
			expected: false,
		},
		{
			name:     "clearly after",
			a:        common.Version{2, 0, 0, 0},
			b:        common.Version{1, 0, 0, 0},
			expected: true,
		},
		{
			name:     "clearly before",
			a:        common.Version{1, 0, 0, 0},
			b:        common.Version{2, 0, 0, 0},
			expected: false,
		},
		{
			name:     "minor version after",
			a:        common.Version{1, 2, 0, 0},
			b:        common.Version{1, 1, 0, 0},
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := version.IsAfter(tc.a, tc.b)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestSort(t *testing.T) {
	// Recommended function implementing cmp from the comments.
	cmp := func(a, b common.Version) int {
		if version.IsBefore(a, b) {
			return -1
		}
		if version.Equals(a, b) {
			return 0
		}
		return 1
	}

	tests := []struct {
		name     string
		input    []common.Version
		expected []common.Version
	}{
		{
			name: "already sorted",
			input: []common.Version{
				version.Phase0(),
				version.Altair(),
				version.Bellatrix(),
				version.Capella(),
			},
			expected: []common.Version{
				version.Phase0(),
				version.Altair(),
				version.Bellatrix(),
				version.Capella(),
			},
		},
		{
			name: "reverse sorted",
			input: []common.Version{
				version.Capella(),
				version.Bellatrix(),
				version.Altair(),
				version.Phase0(),
			},
			expected: []common.Version{
				version.Phase0(),
				version.Altair(),
				version.Bellatrix(),
				version.Capella(),
			},
		},
		{
			name: "mixed versions with minors",
			input: []common.Version{
				version.Deneb1(),
				version.Deneb(),
				version.Phase0(),
				version.Electra(),
				version.Capella(),
			},
			expected: []common.Version{
				version.Phase0(),
				version.Capella(),
				version.Deneb(),
				version.Deneb1(),
				version.Electra(),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Sort the input slice using the cmp function
			slices.SortFunc(tc.input, cmp)
			require.Equal(t, tc.expected, tc.input)
		})
	}
}
