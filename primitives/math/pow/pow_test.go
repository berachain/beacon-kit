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

package pow_test

import (
	"testing"

	"github.com/berachain/beacon-kit/primitives/math/pow"
)

func TestPrevPowerOfTwo(t *testing.T) {
	tests := []struct {
		input    uint64
		expected uint64
	}{
		{0, 1},
		{1, 1},
		{2, 2},
		{3, 2},
		{4, 4},
		{5, 4},
		{6, 4},
		{7, 4},
		{8, 8},
		{9, 8},
		{16, 16},
		{17, 16},
		{31, 16},
		{32, 32},
		{33, 32},
	}

	for _, test := range tests {
		result := pow.PrevPowerOfTwo(test.input)
		if result != test.expected {
			t.Errorf(
				"PrevPowerOfTwo(%d) = %d; want %d",
				test.input, result, test.expected,
			)
		}
	}
}

func TestNextPowerOfTwo(t *testing.T) {
	tests := []struct {
		input    uint64
		expected uint64
	}{
		{0, 1},
		{1, 1},
		{2, 2},
		{3, 4},
		{4, 4},
		{5, 8},
		{6, 8},
		{7, 8},
		{8, 8},
		{9, 16},
		{16, 16},
		{17, 32},
		{31, 32},
		{32, 32},
		{33, 64},
	}

	for _, test := range tests {
		result := pow.NextPowerOfTwo(test.input)
		if result != test.expected {
			t.Errorf(
				"NextPowerOfTwo(%d) = %d; want %d",
				test.input, result, test.expected,
			)
		}
	}
}
