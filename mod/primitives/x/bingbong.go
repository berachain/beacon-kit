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

package main

import (
	"fmt"

	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/math"
	"github.com/berachain/beacon-kit/mod/primitives/ssz"
	"github.com/berachain/beacon-kit/mod/primitives/x/bet"
)

func main() {
	// U64T U64[U64T], RootT ~[32]byte,
	// SpecT any, C Container[RootT]](value C) ([]byte, error) {
	x := bet.ItemA{
		Val1: math.U64(1),
		Val2: []math.U64{math.U64(2)},
	}
	y := bet.ItemB{
		Val1: math.U64(1),
		Val2: []uint64{2},
	}
	bz, err := ssz.SerializeContainer[math.U64, primitives.Root, any](x)
	if err != nil {
		panic(err)
	}

	bz2, err := y.MarshalSSZ()
	if err != nil {
		panic(err)
	}

	fmt.Println(bz, bz2)
}
