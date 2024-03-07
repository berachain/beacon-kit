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

package cache_test

import (
	"math/big"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethcoretypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/itsdevbear/bolaris/engine/client/cache"
	"github.com/stretchr/testify/require"
)

func TestEth1HeaderCache(t *testing.T) {
	cacheUnderTest := cache.NewHeaderCache()

	t.Run("Get from empty cache", func(t *testing.T) {
		h, ok := cacheUnderTest.GetByNumber(0)
		require.False(t, ok)
		require.Nil(t, h)
	})

	t.Run("Set and Get", func(t *testing.T) {
		number := uint64(1234)
		h := &ethcoretypes.Header{
			Number: new(big.Int).SetUint64(number),
		}
		hash := h.Hash()
		cacheUnderTest.Add(h)

		h2, ok := cacheUnderTest.GetByNumber(number)
		require.True(t, ok)
		require.Equal(t, h, h2)

		h3, ok := cacheUnderTest.GetByHash(hash)
		require.True(t, ok)
		require.Equal(t, h, h3)
	})

	t.Run("Overwrite existing", func(t *testing.T) {
		number := uint64(1234)
		h1, ok := cacheUnderTest.GetByNumber(number)
		require.True(t, ok)
		require.NotNil(t, h1)
		require.Equal(t, ethcommon.HexToHash("0x0"), h1.ParentHash)

		oldHash := h1.Hash()

		parentHash := ethcommon.HexToHash("0x1234")
		newHeader := &ethcoretypes.Header{
			Number:     new(big.Int).SetUint64(number),
			ParentHash: parentHash,
		}
		newHash := newHeader.Hash()
		cacheUnderTest.Add(newHeader)

		h2, ok := cacheUnderTest.GetByNumber(number)
		require.True(t, ok)
		require.Equal(t, newHeader, h2)
		require.Equal(t, parentHash, h2.ParentHash)

		h3, ok := cacheUnderTest.GetByHash(newHash)
		require.True(t, ok)
		require.Equal(t, newHeader, h3)
		require.Equal(t, parentHash, h3.ParentHash)

		h4, ok := cacheUnderTest.GetByHash(oldHash)
		require.False(t, ok)
		require.Nil(t, h4)
	})

	t.Run("Prune and verify deletion", func(t *testing.T) {
		for i := uint64(0); i < 20; i++ {
			h := &ethcoretypes.Header{
				Number: new(big.Int).SetUint64(i),
			}
			cacheUnderTest.Add(h)
		}
		_, ok := cacheUnderTest.GetByNumber(uint64(1234))
		require.False(t, ok)
	})
}
