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
	"github.com/berachain/beacon-kit/mod/storage/pkg/manager"
	"github.com/berachain/beacon-kit/mod/storage/pkg/pruner"
	"github.com/berachain/beacon-kit/mod/storage/pkg/pruner/mocks"
	"github.com/stretchr/testify/require"
)

func TestDBManager_Start(t *testing.T) {
	mockPrunable := new(mocks.Prunable)

	ch := make(chan manager.BlockEvent[manager.BeaconBlock], 1)
	pruneParamsFn :=
		func(_ manager.BlockEvent[manager.BeaconBlock]) (uint64, uint64) {
			return 0, 0
		}

	logger := log.NewNopLogger()
	p1 := pruner.NewPruner[
		manager.BeaconBlock,
		manager.BlockEvent[manager.BeaconBlock],
		*mocks.Prunable,
	](logger, mockPrunable, "pruner1", ch, pruneParamsFn)
	p2 := pruner.NewPruner[
		manager.BeaconBlock,
		manager.BlockEvent[manager.BeaconBlock],
		*mocks.Prunable,
	](logger, mockPrunable, "pruner2", ch, pruneParamsFn)

	m, err := manager.NewDBManager(logger, p1, p2)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = m.Start(ctx)
	require.NoError(t, err)
	time.Sleep(100 * time.Millisecond)
	mockPrunable.AssertNotCalled(t, "PruneFromInclusive")
}
