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

package math

import (
	"fmt"

	"github.com/holiman/uint256"
)

const (
	// WeiPerEther is the number of Wei in an Eth.
	WeiPerEther = 1e18

	// GweiPerEther is the number of Gwei in an Eth.
	GweiPerEther = 1e9

	// WeiPerGwei is the number of Wei in a Gwei.
	WeiPerGwei = 1e9
)

type (
	// Wei is the smallest unit of Ether, represented as a pointer to a Uint256.
	Wei struct {
		*uint256.Int
	}

	// Gwei is a denomination of 1e9 Wei represented as an uint64.
	Gwei uint64
)

// ZeroWei returns a zero Wei.
func ZeroWei() Wei {
	return Wei{uint256.NewInt(0)}
}

// WeiFromBytes converts a Wei to a byte slice.
func WeiFromBytes(bz []byte) Wei {
	return Wei{uint256.NewInt(0).SetBytes(bz)}
}

// ToGwei converts Wei to uint64 gwei.
// It DOES not modify the underlying value.
func (w Wei) ToGwei() Gwei {
	if w.Int == nil {
		return 0
	}
	copied := new(uint256.Int).Set(w.Int)
	copied.Div(copied, uint256.NewInt(WeiPerGwei))
	return Gwei(copied.Uint64())
}

// WeiToEther returns the value of a Wei as an Ether.
// FOR DISPLAY PURPOSES ONLY. Do not use for actual
// blockchain things.
func (w Wei) ToEther() string {
	return fmt.Sprintf("%.4f", w.Int.Float64()/WeiPerEther)
}
