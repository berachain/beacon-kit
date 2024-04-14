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

package uint256_test

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/uint256"
	huint256 "github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

func TestLittleEndian_UInt256(t *testing.T) {
	le := uint256.NewLittleEndian([]byte{1, 2, 3, 4, 5})
	expected := new(huint256.Int).SetBytes([]byte{5, 4, 3, 2, 1})
	require.Equal(t, expected, le.UInt256())
}

func TestLittleEndian_Big(t *testing.T) {
	le := uint256.NewLittleEndian([]byte{1, 2, 3, 4, 5})
	expected := new(huint256.Int).SetBytes([]byte{5, 4, 3, 2, 1})
	require.Equal(t, expected.ToBig(), le.ToBig())
}

func TestLittleEndian_MarshalJSON(t *testing.T) {
	le := uint256.NewLittleEndian([]byte{1, 2, 3, 4, 5})
	expected := []byte("\"0x504030201\"")
	result, err := le.MarshalJSON()
	require.NoError(t, err)
	require.JSONEq(t, string(expected), string(result))
}

func TestLittleEndian_UnmarshalJSON(t *testing.T) {
	le := new(uint256.LittleEndian)
	err := le.UnmarshalJSON([]byte("\"0x504030201\""))
	require.NoError(t, err)
	expected := uint256.NewLittleEndian([]byte{1, 2, 3, 4, 5})
	require.Equal(t, expected, *le)
}
