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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package lv_test

import (
	"testing"

	"github.com/berachain/beacon-kit/observability/metrics/lv"
	"github.com/stretchr/testify/require"
)

func TestLabelValues_With(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		initial  lv.LabelValues
		input    []string
		expected lv.LabelValues
	}{
		{
			name:     "empty to single pair",
			initial:  lv.LabelValues{},
			input:    []string{"key1", "value1"},
			expected: lv.LabelValues{"key1", "value1"},
		},
		{
			name:     "append to existing",
			initial:  lv.LabelValues{"key1", "value1"},
			input:    []string{"key2", "value2"},
			expected: lv.LabelValues{"key1", "value1", "key2", "value2"},
		},
		{
			name:     "multiple pairs at once",
			initial:  lv.LabelValues{},
			input:    []string{"key1", "value1", "key2", "value2"},
			expected: lv.LabelValues{"key1", "value1", "key2", "value2"},
		},
		{
			name:     "odd number of values",
			initial:  lv.LabelValues{},
			input:    []string{"key1", "value1", "key2"},
			expected: lv.LabelValues{"key1", "value1", "key2", "unknown"},
		},
		{
			name:     "empty input",
			initial:  lv.LabelValues{"key1", "value1"},
			input:    []string{},
			expected: lv.LabelValues{"key1", "value1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial.With(tt.input...)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestLabelValues_Immutability(t *testing.T) {
	t.Parallel()
	// Test that With doesn't modify the original
	original := lv.LabelValues{"key1", "value1"}
	modified := original.With("key2", "value2")

	require.Equal(t, lv.LabelValues{"key1", "value1"}, original)
	require.Equal(t, lv.LabelValues{"key1", "value1", "key2", "value2"}, modified)
}

func TestLabelValues_Chaining(t *testing.T) {
	t.Parallel()
	// Test chaining multiple With calls
	var lvs lv.LabelValues
	lvs = lvs.With("method", "GET")
	lvs = lvs.With("status", "200")
	lvs = lvs.With("endpoint", "/api/v1/blocks")

	expected := lv.LabelValues{
		"method", "GET",
		"status", "200",
		"endpoint", "/api/v1/blocks",
	}
	require.Equal(t, expected, lvs)
}
