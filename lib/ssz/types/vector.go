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

package types

import (
	"github.com/berachain/beacon-kit/lib/ssz"
	"github.com/berachain/beacon-kit/lib/ssz/common"
)

type Vector[T common.SSZObject] struct {
	Size     uint64
	ElemType common.Type
	Elems    []T
}

func (v *Vector[T]) Marshal() ([]byte, error) {
	return ssz.MarshalComposite(v)
}

func (v *Vector[T]) HashTreeRoot() ([32]byte, error) {
	if common.IsBasicType(v.ElemType) {
		return ssz.MerkleizeBasic(v, true)
	}
	return ssz.MerkleizeComposite(v, true)
}

func (v *Vector[T]) Elements() []common.SSZObject {
	elements := make([]common.SSZObject, len(v.Elems))
	for i, elem := range v.Elems {
		elements[i] = elem
	}
	return elements
}

func (v *Vector[T]) Type() common.Type {
	return common.TypeVector
}
