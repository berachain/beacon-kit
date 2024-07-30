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
	"context"
	"testing"
	"time"

	"github.com/berachain/beacon-kit/mod/node-api/backend/cache"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/stretchr/testify/require"
)

func TestQueryCache(t *testing.T) {
	cacheConfig := cache.Config{
		QueryContextSize: 10,
		QueryContextTTL:  5 * time.Second,
	}

	cacheUnderTest := cache.NewQueryCache(cacheConfig)
	//nolint:revive,staticcheck // ok for test.
	ctx := context.WithValue(context.Background(), "testK", "testV")

	t.Run("Get from empty cache", func(t *testing.T) {
		emtpyCtx, ok := cacheUnderTest.GetQueryContext(0)
		require.False(t, ok)
		require.Nil(t, emtpyCtx)
	})

	t.Run("Set and Get", func(t *testing.T) {
		slot := math.Slot(1234)
		cacheUnderTest.AddQueryContext(slot, ctx)

		ctx2, ok := cacheUnderTest.GetQueryContext(slot)
		require.True(t, ok)
		require.Equal(t, ctx.Value("testK"), ctx2.Value("testK"))
	})

	t.Run("Overwrite existing", func(t *testing.T) {
		slot := math.Slot(12345)
		cacheUnderTest.AddQueryContext(slot, ctx)

		ctx2, ok := cacheUnderTest.GetQueryContext(slot)
		require.True(t, ok)
		require.Equal(t, ctx.Value("testK"), ctx2.Value("testK"))

		//nolint:revive,staticcheck // ok for test.
		newCtx := context.WithValue(context.Background(), "testK", "testV_new")
		//nolint:contextcheck // ok for test.
		cacheUnderTest.AddQueryContext(slot, newCtx)

		ctx3, ok := cacheUnderTest.GetQueryContext(slot)
		require.True(t, ok)
		require.Equal(t, newCtx.Value("testK"), ctx3.Value("testK"))
	})

	t.Run("Prune and verify deletion", func(t *testing.T) {
		cacheUnderTest.AddQueryContext(math.Slot(1234), ctx)
		for i := range math.Slot(cacheConfig.QueryContextSize) {
			cacheUnderTest.AddQueryContext(i, ctx)
		}
		_, found := cacheUnderTest.GetQueryContext(math.Slot(1234))
		require.False(t, found)
	})
}
