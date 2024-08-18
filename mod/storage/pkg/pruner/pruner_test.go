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

package pruner_test

import (
	"context"
	"testing"
	"time"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/events"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/storage/pkg/pruner"
	"github.com/berachain/beacon-kit/mod/storage/pkg/pruner/mocks"
	"github.com/stretchr/testify/mock"
)

func pruneRangeFn[BlockT pruner.BeaconBlock](
	event events.Event[BlockT],
) (uint64, uint64) {
	slot := event.Data().GetSlot().Unwrap()
	return slot, slot
}

func TestPruner(t *testing.T) {
	tests := []struct {
		name          string
		pruneIndexes  []uint64
		expectedCalls int
	}{
		{
			name:          "PruneSingleIndex",
			pruneIndexes:  []uint64{1},
			expectedCalls: 1,
		},
		{
			name:          "PruneMultipleIndexes",
			pruneIndexes:  []uint64{1, 2, 3, 4, 5},
			expectedCalls: 5,
		},
		{
			name:          "NoPruneIndexes",
			pruneIndexes:  []uint64{},
			expectedCalls: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := log.NewNopLogger()
			ch := make(chan events.Event[pruner.BeaconBlock])
			mockPrunable := new(mocks.Prunable)
			mockPrunable.On("Prune", mock.Anything, mock.Anything).
				Return(nil)

			// create Pruner with a Noop logger
			testPruner := pruner.NewPruner[
				pruner.BeaconBlock,
				pruner.Prunable,
			](logger, mockPrunable, "TestPruner", ch, pruneRangeFn)

			ctx, cancel := context.WithCancel(context.Background())
			// need to ensure goroutine is stopped
			defer cancel()

			testPruner.Start(ctx)

			// notify the pruner with the indexes to prune
			for _, index := range tt.pruneIndexes {
				block := mocks.BeaconBlock{}
				block.On("GetSlot").Return(math.U64(index))
				event := events.New[pruner.BeaconBlock](
					context.Background(),
					events.BeaconBlockFinalizedEvent,
					&block,
				)
				ch <- event
			}

			// some time for the goroutine to process the requests
			time.Sleep(100 * time.Millisecond)

			// assert that prune was called expected number of times
			mockPrunable.AssertNumberOfCalls(
				t,
				"Prune",
				tt.expectedCalls,
			)

			// assert that prune was called on correct indices
			for _, index := range tt.pruneIndexes {
				mockPrunable.AssertCalled(
					t,
					"Prune",
					index,
					mock.Anything,
				)
			}
		})
	}
}
