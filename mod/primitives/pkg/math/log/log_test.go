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

package log_test

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/math/log"
)

func TestILog2Ceil(t *testing.T) {
	tests := []struct {
		input    uint64
		expected uint8
	}{
		{0, 0},
		{1, 0},
		{2, 1},
		{3, 2},
		{4, 2},
		{5, 3},
		{8, 3},
		{9, 4},
		{16, 4},
		{17, 5},
	}

	for _, test := range tests {
		result := log.ILog2Ceil(test.input)
		if result != test.expected {
			t.Errorf(
				"ILog2Ceil(%d) = %d; want %d",
				test.input, result, test.expected,
			)
		}
	}
}

func TestILog2Floor(t *testing.T) {
	tests := []struct {
		input    uint64
		expected uint8
	}{
		{0, 0},
		{1, 0},
		{2, 1},
		{3, 1},
		{4, 2},
		{5, 2},
		{8, 3},
		{9, 3},
		{16, 4},
		{17, 4},
	}

	for _, test := range tests {
		result := log.ILog2Floor(test.input)
		if result != test.expected {
			t.Errorf(
				"ILog2Floor(%d) = %d; want %d",
				test.input, result, test.expected,
			)
		}
	}
}
