// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
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

package primitives_test

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

func TestWeiToGwei(t *testing.T) {
	tests := []struct {
		name string
		v    *uint256.Int
		want primitives.Gwei
	}{
		{"just below 1 Gwei", uint256.NewInt(1e9 - 1), 0},
		{"exactly 1 Gwei", uint256.NewInt(1e9), 1},
		{"10 Gwei", uint256.NewInt(1e10), 10},
		{"large number", uint256.NewInt(239489233849348394), 239489233},
		{"1 Eth", uint256.NewInt(1e18), 1000000000},
		{"1.5 Eth", uint256.NewInt(15e17), 1500000000},
		{
			"edge case large number",
			uint256.NewInt(999999999999999999),
			999999999,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := primitives.Wei{tt.v}.ToGwei()
			if got != tt.want {
				t.Errorf("WeiToGwei() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWeiToGwei_CopyOk(t *testing.T) {
	v := uint256.NewInt(1e9)
	got := primitives.Wei{v}.ToGwei()

	require.Equal(
		t, primitives.Gwei(1), got, "conversion result mismatch")
	require.Equal(
		t, uint256.NewInt(1e9).Uint64(), v.Uint64(), "original value modified")
}

func TestWeiToEther(t *testing.T) {
	tests := []struct {
		name string
		v    *uint256.Int
		want string
	}{
		{"just below 1 Ether", uint256.NewInt(1e18 - 1e14), "0.9999"},
		{"exactly 1 Ether", uint256.NewInt(1e18), "1.0000"},
		{"1.5 Ether", uint256.NewInt(15e17), "1.5000"},
		{"10 Ether", uint256.NewInt(1e19), "10.0000"},
		{"large number", uint256.NewInt(1e18 - 1e17), "0.9000"},
		{"edge case large number", uint256.NewInt(9999900000e9), "9.9999"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := primitives.Wei{tt.v}.ToEther()
			if got != tt.want {
				t.Errorf("WeiToEther() = %v, want %v", got, tt.want)
			}
		})
	}
}
