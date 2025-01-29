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
	"testing"

	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/stretchr/testify/require"
)

func TestEquals(t *testing.T) {
	tests := []struct {
		name     string
		a, b     func() common.Version
		expected bool
	}{
		{
			name:     "same version (deneb)",
			a:        version.Deneb,
			b:        version.Deneb,
			expected: true,
		},
		{
			name:     "different versions (deneb vs electra)",
			a:        version.Deneb,
			b:        version.Electra,
			expected: false,
		},
		{
			name:     "phase0 vs phase0",
			a:        version.Phase0,
			b:        version.Phase0,
			expected: true,
		},
		{
			name:     "deneb vs deneb1",
			a:        version.Deneb,
			b:        version.Deneb1,
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := version.Equals(tc.a(), tc.b())
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestPrecedes(t *testing.T) {
	tests := []struct {
		name     string
		a, b     func() common.Version
		expected bool
	}{
		{
			name:     "same version (deneb)",
			a:        version.Deneb,
			b:        version.Deneb,
			expected: true,
		},
		{
			name:     "phase0 precedes altair",
			a:        version.Phase0,
			b:        version.Altair,
			expected: true,
		},
		{
			name:     "altair precedes bellatrix",
			a:        version.Altair,
			b:        version.Bellatrix,
			expected: true,
		},
		{
			name:     "deneb does not precede altair",
			a:        version.Deneb,
			b:        version.Altair,
			expected: false,
		},
		{
			name:     "deneb precedes electra",
			a:        version.Deneb,
			b:        version.Electra,
			expected: true,
		},
		{
			name:     "deneb precedes deneb1",
			a:        version.Deneb,
			b:        version.Deneb1,
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := version.Precedes(tc.a(), tc.b())
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestFollows(t *testing.T) {
	tests := []struct {
		name     string
		a, b     func() common.Version
		expected bool
	}{
		{
			name:     "same version (deneb)",
			a:        version.Deneb,
			b:        version.Deneb,
			expected: true,
		},
		{
			name:     "altair follows phase0",
			a:        version.Altair,
			b:        version.Phase0,
			expected: true,
		},
		{
			name:     "bellatrix follows altair",
			a:        version.Bellatrix,
			b:        version.Altair,
			expected: true,
		},
		{
			name:     "altair does not follow deneb",
			a:        version.Altair,
			b:        version.Deneb,
			expected: false,
		},
		{
			name:     "electra follows deneb",
			a:        version.Electra,
			b:        version.Deneb,
			expected: true,
		},
		{
			name:     "deneb1 follows deneb",
			a:        version.Deneb1,
			b:        version.Deneb,
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := version.Follows(tc.a(), tc.b())
			require.Equal(t, tc.expected, result)
		})
	}
}
