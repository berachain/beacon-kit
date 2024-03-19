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

package ssz_test

// import (
// 	"math/rand"
// 	"testing"

// 	"github.com/berachain/beacon-kit/lib/ssz"
// 	"github.com/berachain/beacon-kit/lib/ssz/mocks"
// 	"github.com/berachain/beacon-kit/lib/ssz/types"
// 	"github.com/stretchr/testify/require"
// )

// func Test_Vector4SSZ(t *testing.T) {
// 	sszVec := make([]uint64, 4)
// 	for i := range sszVec {
// 		sszVec[i] = rand.Uint64()
// 	}
// 	sszContainer := &mocks.Vector4Container{
// 		VectorField: sszVec,
// 	}

// 	vec := types.Vector[types.Uint64]{
// 		Length:   4,
// 		ElemType: types.TypeUint,
// 		Elements: make([]types.Uint64, 4),
// 	}

// 	for i := range sszVec {
// 		vec.Elements[i] = types.Uint64(sszVec[i])
// 	}
// 	container := &types.Container{
// 		Fields: []types.SSZObject{&vec},
// 	}

// 	h1, err := sszContainer.HashTreeRoot()
// 	require.NoError(t, err)

// 	h2, err := ssz.MerkleizeContainerSSZ(container)
// 	require.NoError(t, err)

// 	require.Equal(t, h1, h2)
// }

// func Test_Vector5SSZ(t *testing.T) {
// 	vec := make([]uint64, 5)
// 	for i := range vec {
// 		vec[i] = uint64(i) + 1
// 	}
// 	container := &mocks.Vector5Container{
// 		VectorField: vec,
// 	}

// 	mockVec := make([]mocks.Uint64, 5)
// 	for i := range mockVec {
// 		mockVec[i] = mocks.Uint64(vec[i])
// 	}
// 	mockContainer := &mocks.MockSingleFieldContainer[mocks.Vector[mocks.Uint64]]{
// 		Field: mockVec,
// 	}

// 	h1, err := container.HashTreeRoot()
// 	require.NoError(t, err)

// 	h2, err := ssz.MerkleizeContainerSSZ(mockContainer)
// 	require.NoError(t, err)

// 	require.Equal(t, h1, h2)
// }

// func Test_Vector6SSZ(t *testing.T) {
// 	vec := make([]uint64, 6)
// 	for i := range vec {
// 		vec[i] = uint64(i) + 1
// 	}
// 	container := &mocks.Vector6Container{
// 		VectorField: vec,
// 	}

// 	mockVec := make([]mocks.Uint64, 6)
// 	for i := range mockVec {
// 		mockVec[i] = mocks.Uint64(vec[i])
// 	}
// 	mockContainer := &mocks.MockSingleFieldContainer[mocks.Vector[mocks.Uint64]]{
// 		Field: mockVec,
// 	}

// 	// Per https://simpleserialize.com, the root of the vector [1, 2, 3, 4, 5, 6] is:
// 	// 0xac136edda3bdd2e949a19a945b1ac554e4b607d339a43c540c336098fff97f2b

// 	h1, err := container.HashTreeRoot()
// 	require.NoError(t, err)

// 	h2, err := ssz.MerkleizeContainerSSZ(mockContainer)
// 	require.NoError(t, err)

// 	require.Equal(t, h1, h2)
// }
