// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package math_test

import (
	"testing"

	"github.com/holiman/uint256"
	"github.com/itsdevbear/bolaris/math"
	"github.com/stretchr/testify/require"
)

func TestWeiToGwei(t *testing.T) {
	tests := []struct {
		name string
		v    *uint256.Int
		want math.Gwei
	}{
		{"just below 1 Gwei", uint256.NewInt(1e9 - 1), 0},
		{"exactly 1 Gwei", uint256.NewInt(1e9), 1},
		{"10 Gwei", uint256.NewInt(1e10), 10},
		{"large number", uint256.NewInt(239489233849348394), 239489233},
		{"1 Eth", uint256.NewInt(1e18), 1000000000},
		{"1.5 Eth", uint256.NewInt(15e17), 1500000000},
		{"edge case large number", uint256.NewInt(999999999999999999), 999999999},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := math.WeiToGwei(tt.v); got != tt.want {
				t.Errorf("WeiToGwei() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWeiToGwei_CopyOk(t *testing.T) {
	v := uint256.NewInt(1e9)
	got := math.WeiToGwei(v)

	require.Equal(t, math.Gwei(1), got, "conversion result mismatch")
	require.Equal(t, uint256.NewInt(1e9).Uint64(), v.Uint64(), "original value modified")
}
