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

package types_test

import (
	"testing"

	"github.com/berachain/beacon-kit/lib/ssz/common"
	"github.com/berachain/beacon-kit/lib/ssz/mocks"
	"github.com/berachain/beacon-kit/lib/ssz/types"
	"github.com/stretchr/testify/require"
)

func Test_VectorBasicObjects(t *testing.T) {
	size := 6
	elems := make([]uint64, size)
	for i := range elems {
		elems[i] = uint64(i) + 1
	}

	mockContainer := &mocks.Vector6Container{
		VectorField: elems,
	}

	sszElems := make([]types.Uint64, size)
	for i := range elems {
		sszElems[i] = types.Uint64(elems[i])
	}
	sszVec := &types.Vector[types.Uint64]{
		Size:     uint64(size),
		ElemType: common.TypeUint,
		Elems:    sszElems,
	}

	h1, err := mockContainer.HashTreeRoot()
	require.NoError(t, err)

	h2, err := sszVec.HashTreeRoot()
	require.NoError(t, err)

	require.Equal(t, h1, h2)
}
