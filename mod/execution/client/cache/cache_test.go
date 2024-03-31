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
	"time"

	"github.com/berachain/beacon-kit/mod/execution/client/cache"
	"github.com/berachain/beacon-kit/mod/primitives"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestEth1HeaderCache(t *testing.T) {
	cacheConfig := cache.Config{
		HeaderSize: 10,
		HeaderTTL:  5 * time.Second,
	}

	cacheUnderTest := cache.NewEngineCache(cacheConfig)

	t.Run("Get from empty cache", func(t *testing.T) {
		h, ok := cacheUnderTest.HeaderByNumber(0)
		require.False(t, ok)
		require.Nil(t, h)
	})

	t.Run("Set and Get", func(t *testing.T) {
		number := uint64(1234)
		h := &primitives.ExecutionHeader{
			Number: new(big.Int).SetUint64(number),
		}
		hash := h.Hash()
		cacheUnderTest.AddHeader(h)

		h2, ok := cacheUnderTest.HeaderByNumber(number)
		require.True(t, ok)
		require.Equal(t, h, h2)

		h3, ok := cacheUnderTest.HeaderByHash(hash)
		require.True(t, ok)
		require.Equal(t, h, h3)
	})

	t.Run("Overwrite existing", func(t *testing.T) {
		number := uint64(1234)
		h1, ok := cacheUnderTest.HeaderByNumber(number)
		require.True(t, ok)
		require.NotNil(t, h1)
		require.Equal(t, ethcommon.HexToHash("0x0"), h1.ParentHash)

		oldHash := h1.Hash()

		parentHash := ethcommon.HexToHash("0x1234")
		newHeader := &primitives.ExecutionHeader{
			Number:     new(big.Int).SetUint64(number),
			ParentHash: parentHash,
		}
		newHash := newHeader.Hash()
		cacheUnderTest.AddHeader(newHeader)

		h2, ok := cacheUnderTest.HeaderByNumber(number)
		require.True(t, ok)
		require.Equal(t, newHeader, h2)
		require.Equal(t, parentHash, h2.ParentHash)

		h3, ok := cacheUnderTest.HeaderByHash(newHash)
		require.True(t, ok)
		require.Equal(t, newHeader, h3)
		require.Equal(t, parentHash, h3.ParentHash)

		h4, ok := cacheUnderTest.HeaderByHash(oldHash)
		require.False(t, ok)
		require.Nil(t, h4)
	})

	t.Run("Prune and verify deletion", func(t *testing.T) {
		for i := uint64(0); i < 20; i++ {
			h := &primitives.ExecutionHeader{
				Number: new(big.Int).SetUint64(i),
			}
			cacheUnderTest.AddHeader(h)
		}
		_, ok := cacheUnderTest.HeaderByNumber(uint64(1234))
		require.False(t, ok)
	})
}
