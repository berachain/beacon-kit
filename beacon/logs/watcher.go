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

package logs

import (
	"context"
	"math/big"

	"cosmossdk.io/log"
	beacontypes "github.com/berachain/beacon-kit/beacon/core/types"
	engineclient "github.com/berachain/beacon-kit/engine/client"
	"github.com/berachain/beacon-kit/lib/skiplist"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	coretypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

// Watcher is responsible for ingestion logs from the execution client.
type Watcher struct {
	logger                    log.Logger
	lastProcessedDepositIndex uint64
	depositQueue              *skiplist.Skiplist[*beacontypes.Deposit]
	depositContractAddress    common.Address
	ec                        *engineclient.EngineClient
	ch                        chan struct{}
}

// Start spawns any goroutines required by the service.
func (w *Watcher) Start(ctx context.Context) {
	w.ch = make(chan struct{})
	go w.mainLoop(ctx)
}

// UpdateLastProcessedDepositIndex updates the last processed deposit index.
func (w *Watcher) UpdateLastProcessedDepositIndex(
	index uint64,
) {
	w.lastProcessedDepositIndex = index

	// We notify the mainLoop() that we have updated the last processed
	// block, and thus we need to sync logs to the head.
	w.ch <- struct{}{}
}

// mainLoop is the main loop of the watcher.
func (w *Watcher) mainLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-w.ch:
			// TODO: these should be pushed into a serial queue using
			// GCD, right now if we have a lot of logs to backfill
			// we are going to start firing off go-rountines like
			// no tomorrow.
			if err := w.syncLogsToHead(ctx); err != nil {
				w.logger.Error("failed to sync logs to head", "err", err)
			}
		}
	}
}

// syncLogsToHead syncs logs to the head of the chain, it also has backfilling
// capabilities.
func (w *Watcher) syncLogsToHead(
	ctx context.Context,
) error {
	var (
		logs       []coretypes.Log
		head       *beacontypes.Deposit
		finalBlock *coretypes.Block
		err        error
	)

	// Get the latest finalized block.
	finalBlock, err = w.ec.Client.BlockByNumber(
		ctx,
		big.NewInt(int64(rpc.FinalizedBlockNumber)),
	)
	if err != nil {
		return err
	}
	finalBlockNumber := finalBlock.NumberU64()

	for {
		if head, err = w.depositQueue.Front(); err != nil {
			return err
		}

		// If this is true, we are caught up and can exit.
		if head.Index <= w.lastProcessedDepositIndex {
			return nil
		}

		logs, err = w.ec.Client.FilterLogs(ctx, ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(finalBlockNumber)),
			ToBlock:   big.NewInt(int64(finalBlockNumber)),
			Addresses: []common.Address{w.depositContractAddress},
		})
		if err != nil {
			return err
		}

		_ = logs

		// Keep going back in time until we have caught up. This gives us
		// a built-in backfilling mechanism along the main loop.
		if finalBlockNumber == 0 {
			return nil
		}
		finalBlockNumber--
	}
}
