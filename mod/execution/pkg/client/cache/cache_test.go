// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package cache_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/berachain/beacon-kit/mod/execution/pkg/client/cache"
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
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
		h := &gethprimitives.Header{
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
		newHeader := &gethprimitives.Header{
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
		for i := range uint64(20) {
			h := &gethprimitives.Header{
				Number: new(big.Int).SetUint64(i),
			}
			cacheUnderTest.AddHeader(h)
		}
		_, ok := cacheUnderTest.HeaderByNumber(uint64(1234))
		require.False(t, ok)
	})
}
