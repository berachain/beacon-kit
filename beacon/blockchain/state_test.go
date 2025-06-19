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

package blockchain_test

import (
	"testing"

	"cosmossdk.io/collections"
	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/node-core/components/metrics"
	"github.com/berachain/beacon-kit/node-core/components/storage"
	"github.com/berachain/beacon-kit/primitives/math"
	statetransition "github.com/berachain/beacon-kit/testing/state-transition"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestStateCopyIsIndependentFromOriginal(t *testing.T) {
	t.Parallel()

	// 1. Setup the test objects
	cs, err := spec.MainnetChainSpec()
	require.NoError(t, err)

	cms, kvStore, _, err := statetransition.BuildTestStores()
	require.NoError(t, err)
	sdkCtx := sdk.NewContext(cms.CacheMultiStore(), true, log.NewNopLogger())

	backend := storage.NewBackend(
		cs,
		nil,
		kvStore,
		nil,
		nil,
		log.NewNopLogger(),
		metrics.NewNoOpTelemetrySink(),
	)

	// 2. Query both state and its copy
	state := backend.StateFromContext(sdkCtx)
	copiedState := state.Copy(sdkCtx)

	_, err = state.GetSlot()
	require.ErrorIs(t, err, collections.ErrNotFound)

	_, err = copiedState.GetSlot()
	require.ErrorIs(t, err, collections.ErrNotFound)

	// 3. Modify state and check that the copied state is not affected
	targetSlot := math.Slot(12345)
	require.NoError(t, state.KVStore.SetSlot(targetSlot))

	retrieved1, err := state.GetSlot()
	require.NoError(t, err)
	require.Equal(t, targetSlot, retrieved1)

	_, err = copiedState.GetSlot()
	require.ErrorIs(t, err, collections.ErrNotFound)
}
