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

package math

import (
	"github.com/holiman/uint256"
)

// GweiPerEth is the number of Gwei in an Eth.
const GweiPerEth = 1e9

type (
	// Wei is the smallest unit of Ether, represented as a pointer to a Uint256.
	Wei = *uint256.Int
	// Gwei is a denomination of 1e9 Wei represented as an uint64.
	Gwei uint64
)

// BytesToWei converts a byte slice to a Wei.
func BytesToWei(v []byte) Wei {
	return uint256.NewInt(0).SetBytes(v)
}

// WeiToBytes converts a Wei to a byte slice.
func WeiToBytes(v Wei) []byte {
	return v.Bytes()
}

// WeiToGwei converts Wei to uint64 gwei.
// The input `v` is copied before being modified.
func WeiToGwei(v Wei) Gwei {
	if v == nil {
		return 0
	}
	copied := new(uint256.Int).Set(v)
	copied.Div(copied, uint256.NewInt(GweiPerEth))
	return Gwei(copied.Uint64())
}
