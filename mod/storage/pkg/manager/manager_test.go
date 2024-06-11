// // SPDX-License-Identifier: MIT
// //
// // Copyright (c) 2024 Berachain Foundation
// //
// // Permission is hereby granted, free of charge, to any person
// // obtaining a copy of this software and associated documentation
// // files (the "Software"), to deal in the Software without
// // restriction, including without limitation the rights to use,
// // copy, modify, merge, publish, distribute, sublicense, and/or sell
// // copies of the Software, and to permit persons to whom the
// // Software is furnished to do so, subject to the following
// // conditions:
// //
// // The above copyright notice and this permission notice shall be
// // included in all copies or substantial portions of the Software.
// //
// // THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// // EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// // OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// // NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// // HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// // WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// // FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// // OTHER DEALINGS IN THE SOFTWARE.

package manager_test

import (
	"context"
	"testing"
	"time"

	"cosmossdk.io/log"
	interfaceMocks "github.com/berachain/beacon-kit/mod/storage/pkg/interfaces/mocks"
	"github.com/berachain/beacon-kit/mod/storage/pkg/manager"
	"github.com/berachain/beacon-kit/mod/storage/pkg/manager/mocks"
	"github.com/berachain/beacon-kit/mod/storage/pkg/pruner"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDBManager_Start(t *testing.T) {
	mockPrunable := new(interfaceMocks.Prunable)
	feed := mocks.BlockFeed[
		manager.BeaconBlock,
		manager.BlockEvent[manager.BeaconBlock],
		manager.Subscription,
	]{}

	sub := mocks.Subscription{}
	sub.EXPECT().Unsubscribe().Return()
	feed.EXPECT().Subscribe(mock.Anything).Return(&sub)
	pruneParamsFn :=
		func(_ manager.BlockEvent[manager.BeaconBlock]) (uint64, uint64) {
			return 0, 0
		}

	logger := log.NewNopLogger()
	p1 := pruner.NewPruner[
		manager.BeaconBlock,
		manager.BlockEvent[manager.BeaconBlock],
		*interfaceMocks.Prunable,
		manager.Subscription,
	](logger, mockPrunable, "pruner1", &feed, pruneParamsFn)
	p2 := pruner.NewPruner[
		manager.BeaconBlock,
		manager.BlockEvent[manager.BeaconBlock],
		*interfaceMocks.Prunable,
		manager.Subscription,
	](logger, mockPrunable, "pruner2", &feed, pruneParamsFn)

	m, err := manager.NewDBManager[
		manager.BeaconBlock,
		manager.BlockEvent[manager.BeaconBlock],
		manager.Subscription,
	](logger, p1, p2)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = m.Start(ctx)
	require.NoError(t, err)
	time.Sleep(100 * time.Millisecond)
	feed.AssertCalled(t, "Subscribe", mock.Anything)
	feed.AssertNumberOfCalls(t, "Subscribe", 2)
	mockPrunable.AssertNotCalled(t, "PruneFromInclusive")
}
