// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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

package node_test

import (
	"testing"

	"github.com/berachain/beacon-kit/node-api/handlers/node"
	"github.com/berachain/beacon-kit/node-api/handlers/node/mocks"
	"github.com/berachain/beacon-kit/node-api/handlers/node/types"
	"github.com/stretchr/testify/require"
)

func TestNodeSyncing(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                string
		setMockExpectations func(*mocks.Backend)
		check               func(t *testing.T, res any, err error)
	}{
		{
			name: "syncing",
			setMockExpectations: func(b *mocks.Backend) {
				latestHeight := int64(2025)
				syncToHeight := int64(2026)
				b.EXPECT().GetSyncData().Return(latestHeight, syncToHeight).Once()
			},
			check: func(t *testing.T, res any, err error) {
				t.Helper()

				require.NoError(t, err)
				require.NotNil(t, res)
				require.IsType(t, types.DataResponse{}, res)
				dr, _ := res.(types.DataResponse)
				require.IsType(t, &types.SyncingData{}, dr.Data)
				data, _ := dr.Data.(*types.SyncingData)

				latestHeight := int64(2025) // duplicated from setMockExpectations
				syncToHeight := int64(2026) // duplicated from setMockExpectations

				require.Equal(t, latestHeight, data.HeadSlot)
				require.Equal(t, syncToHeight-latestHeight, data.SyncDistance)
				require.True(t, data.IsSyncing)
				require.False(t, data.IsOptimistic)
				require.False(t, data.ELOffline)
			},
		},
		{
			name: "normal operations mode",
			setMockExpectations: func(b *mocks.Backend) {
				latestHeight := int64(1492)
				syncToHeight := int64(1492)
				b.EXPECT().GetSyncData().Return(latestHeight, syncToHeight).Once()
			},
			check: func(t *testing.T, res any, err error) {
				t.Helper()

				require.NoError(t, err)
				require.NotNil(t, res)
				require.IsType(t, types.DataResponse{}, res)
				dr, _ := res.(types.DataResponse)
				require.IsType(t, &types.SyncingData{}, dr.Data)
				data, _ := dr.Data.(*types.SyncingData)

				latestHeight := int64(1492) // duplicated from setMockExpectations
				syncToHeight := int64(1492) // duplicated from setMockExpectations

				require.Equal(t, latestHeight, data.HeadSlot)
				require.Equal(t, syncToHeight-latestHeight, data.SyncDistance)
				require.False(t, data.IsSyncing)
				require.False(t, data.IsOptimistic)
				require.False(t, data.ELOffline)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// setup test
			backend := mocks.NewBackend(t)
			h := node.NewHandler(backend)

			tc.setMockExpectations(backend)

			// test
			res, err := h.Syncing(nil) // input does not matter here

			tc.check(t, res, err)
		})
	}
}
