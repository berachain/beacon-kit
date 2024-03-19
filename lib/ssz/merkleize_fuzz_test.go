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

// func FuzzUint8SSZ(f *testing.F) {
// 	f.Add(uint8(0))
// 	f.Fuzz(func(t *testing.T, value uint8) {
// 		container := &mocks.Uint8Container{
// 			Uint8Field: value,
// 		}

// 		mockContainer := &mocks.MockSingleFieldContainer[mocks.Uint8]{
// 			Field: mocks.Uint8(value),
// 		}

// 		h1, err := container.HashTreeRoot()
// 		require.NoError(t, err)

// 		h2, err := ssz.MerkleizeContainerSSZ(mockContainer)
// 		require.NoError(t, err)

// 		require.Equal(t, h1, h2)
// 	})
// }

// func FuzzUint16SSZ(f *testing.F) {
// 	f.Add(uint16(0))
// 	f.Fuzz(func(t *testing.T, value uint16) {
// 		container := &mocks.Uint16Container{
// 			Uint16Field: value,
// 		}

// 		mockContainer := &mocks.MockSingleFieldContainer[mocks.Uint16]{
// 			Field: mocks.Uint16(value),
// 		}

// 		h1, err := container.HashTreeRoot()
// 		require.NoError(t, err)

// 		h2, err := ssz.MerkleizeContainerSSZ(mockContainer)
// 		require.NoError(t, err)

// 		require.Equal(t, h1, h2)
// 	})
// }

// func FuzzUint32SSZ(f *testing.F) {
// 	f.Add(uint32(0))
// 	f.Fuzz(func(t *testing.T, value uint32) {
// 		container := &mocks.Uint32Container{
// 			Uint32Field: value,
// 		}

// 		mockContainer := &mocks.MockSingleFieldContainer[mocks.Uint32]{
// 			Field: mocks.Uint32(value),
// 		}

// 		h1, err := container.HashTreeRoot()
// 		require.NoError(t, err)

// 		h2, err := ssz.MerkleizeContainerSSZ(mockContainer)
// 		require.NoError(t, err)

// 		require.Equal(t, h1, h2)
// 	})
// }

// func FuzzUint64SSZ(f *testing.F) {
// 	f.Add(uint64(0))
// 	f.Fuzz(func(t *testing.T, value uint64) {
// 		container := &mocks.Uint64Container{
// 			Uint64Field: value,
// 		}

// 		mockContainer := &mocks.MockSingleFieldContainer[mocks.Uint64]{
// 			Field: mocks.Uint64(value),
// 		}

// 		h1, err := container.HashTreeRoot()
// 		require.NoError(t, err)

// 		h2, err := ssz.MerkleizeContainerSSZ(mockContainer)
// 		require.NoError(t, err)

// 		require.Equal(t, h1, h2)
// 	})
// }

// func FuzzByteSSZ(f *testing.F) {
// 	f.Add(byte(0))
// 	f.Fuzz(func(t *testing.T, value byte) {
// 		container := &mocks.ByteContainer{
// 			ByteField: value,
// 		}

// 		mockContainer := &mocks.MockSingleFieldContainer[mocks.Byte]{
// 			Field: mocks.Byte(value),
// 		}

// 		h1, err := container.HashTreeRoot()
// 		require.NoError(t, err)

// 		h2, err := ssz.MerkleizeContainerSSZ(mockContainer)
// 		require.NoError(t, err)

// 		require.Equal(t, h1, h2)
// 	})
// }

// func FuzzBoolSSZ(f *testing.F) {
// 	f.Add(false)
// 	f.Fuzz(func(t *testing.T, value bool) {
// 		container := &mocks.BoolContainer{
// 			BoolField: value,
// 		}

// 		mockContainer := &mocks.MockSingleFieldContainer[mocks.Bool]{
// 			Field: mocks.Bool(value),
// 		}

// 		h1, err := container.HashTreeRoot()
// 		require.NoError(t, err)

// 		h2, err := ssz.MerkleizeContainerSSZ(mockContainer)
// 		require.NoError(t, err)

// 		require.Equal(t, h1, h2)
// 	})
// }

// func FuzzVectorSSZ(f *testing.F) {
// 	f.Add(uint64(0))
// 	f.Fuzz(func(t *testing.T, value uint64) {
// 		vec := make([]uint64, 20)
// 		for i := range vec {
// 			vec[i] = value + uint64(i)
// 		}
// 		container := &mocks.Vector4Container{
// 			VectorField: vec,
// 		}

// 		mockVec := make([]mocks.Uint64, 20)
// 		for i := range mockVec {
// 			mockVec[i] = mocks.Uint64(vec[i])
// 		}
// 		mockContainer := &mocks.MockSingleFieldContainer[mocks.Vector[mocks.Uint64]]{
// 			Field: mockVec,
// 		}

// 		h1, err := container.HashTreeRoot()
// 		require.NoError(t, err)

// 		h2, err := ssz.MerkleizeContainerSSZ(mockContainer)
// 		require.NoError(t, err)

// 		require.Equal(t, h1, h2)
// 	})
// }
