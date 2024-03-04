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
	"time"

	"cosmossdk.io/log"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	engineclient "github.com/itsdevbear/bolaris/engine/client"
	"github.com/pkg/errors"
)

const (
	// DefaultBatchSize is the default size of the batch
	// for processing the logs in the background.
	DefaultBatchSize = 1000
)

type Processor struct {
	fcs ReadOnlyForkChoicer
	fls FinalizedLogsStore

	logger log.Logger

	// engine gives the access to the Engine API
	// of the execution client.
	engine engineclient.Caller

	// logFactory is the factory for creating
	// objects from Ethereum logs.
	factory LogFactory

	// sigToCache is a map of log signatures to their caches.
	sigToCache map[ethcommon.Hash]LogCache
}

func NewProcessor(opts ...Option[Processor]) (*Processor, error) {
	p := &Processor{}
	for _, opt := range opts {
		if err := opt(p); err != nil {
			return nil, err
		}
	}
	return p, nil
}

// ProcessPastLogs processes the blocks in batch,
// from the last processed block (exclusive)
// to the latest block (inclusive) for each cache.
// This function will be called in a goroutine
// to prediodically process the logs and backfill
// the caches in background.
func (p *Processor) ProcessPastLogs(
	ctx context.Context,
) error {
	// Get the latest finalized block hash and block number.
	finalizedBlockHash := p.fcs.FinalizedCheckpoint()
	finalizedHeader, err := p.engine.HeaderByHash(ctx, finalizedBlockHash)
	if err != nil {
		return errors.Wrapf(err, "failed to get finalized header")
	}
	finalizedBlockNumber := finalizedHeader.Number.Uint64()

	// Determine the earliest block to process
	// by checking the last finalized blocks among caches.
	// By doing so, we can avoid processing the same block
	// multiple times for different types of logs.
	minLastFinalizedBlockInCache := finalizedBlockNumber
	for sig, cache := range p.sigToCache {
		lastFinalizedBlockInCache := cache.LastFinalizedBlock()
		lastProcessedBlock := p.fls.GetLastProcessedBlockNumber(sig)
		// Update the block number from which we should start processing
		// logs to insert into the cache.
		if lastFinalizedBlockInCache < minLastFinalizedBlockInCache {
			minLastFinalizedBlockInCache = lastFinalizedBlockInCache
		}
		if lastProcessedBlock < minLastFinalizedBlockInCache {
			minLastFinalizedBlockInCache = lastProcessedBlock
		}
	}

	// If all caches have processed the latest finalized block,
	// we don't need to process it again.
	if minLastFinalizedBlockInCache == finalizedBlockNumber {
		return nil
	}

	// Get the registered addresses for the logs.
	registeredAddresses := p.factory.GetRegisteredAddresses()

	currBlock := minLastFinalizedBlockInCache
	for currBlock < finalizedBlockNumber {
		// Process the logs in batch.
		currBlock, err = p.processBlocksInBatch(
			ctx,
			currBlock+1,
			DefaultBatchSize,
			finalizedBlockNumber,
			registeredAddresses,
		)
		if err != nil {
			return errors.Wrapf(err, "failed to process logs in batch")
		}
	}

	return nil
}

// processBlocksInBatch processes the logs in the range
// from fromBlock (inclusive)
// to min(fromBlock + batchSize - 1, latestFinalizedBlock) (inclusive).
func (p *Processor) processBlocksInBatch(
	ctx context.Context,
	fromBlock uint64,
	batchSize uint64,
	latestFinalizedBlock uint64,
	registeredAddresses []ethcommon.Address,
) (uint64, error) {
	// Gather all the logs corresponding to
	// the addresses of interest in the range.
	// TODO: Can we assume that the logs are returned in order?
	toBlock := fromBlock + batchSize - 1
	if toBlock > latestFinalizedBlock {
		toBlock = latestFinalizedBlock
	}
	batchedLogs, err := p.engine.GetLogs(
		ctx,
		fromBlock,
		toBlock,
		registeredAddresses,
	)
	if err != nil {
		// TODO: Handle TooMuchDataRequestedError.
		return 0, errors.Wrapf(err, "failed to get logs")
	}

	blockToLogs := make(map[uint64][]ethtypes.Log)
	for _, log := range batchedLogs {
		blockToLogs[log.BlockNumber] = append(blockToLogs[log.BlockNumber], log)
	}

	defer func() {
		// If there are any erros, we need to rollback
		// the caches to the last finalized block.
		if err != nil {
			p.rollbackCaches()
		}
	}()
	for blockNum := fromBlock; blockNum <= toBlock; blockNum++ {
		logs, ok := blockToLogs[blockNum]
		if !ok {
			continue
		}
		var containers []LogValueContainer
		// Process the logs (in parallel) and push them into the caches.
		containers, err = p.factory.ProcessLogs(logs, blockNum)
		if err != nil {
			return 0, errors.Wrapf(err, "failed to process logs")
		}
		for _, container := range containers {
			sig := container.Signature()
			var cache LogCache
			cache, ok = p.sigToCache[sig]
			if !ok {
				continue
			}
			err = cache.Push(container)
			if err != nil {
				return 0, errors.Wrapf(err, "failed to push container")
			}
		}
		// We start processing the logs from a new block.
		// Notify the caches to update the last finalized block.
		p.setLastFinalizedBlockAllCaches(blockNum)
	}

	return toBlock, nil
}

// rollbackCaches rolls back all the caches to the last finalized block.
func (p *Processor) rollbackCaches() {
	for _, cache := range p.sigToCache {
		cache.Rollback()
	}
}

// setLastFinalizedBlockAllCaches sets the last finalized block
// to the given block number for all the caches.
func (p *Processor) setLastFinalizedBlockAllCaches(blockNumber uint64) {
	for _, cache := range p.sigToCache {
		cache.SetLastFinalizedBlock(blockNumber)
	}
}

// RunLoop processes the past logs in background.
func (p *Processor) RunLoop(ctx context.Context) {
	// TODO: Make this configurable?
	logPeriod := 1 * time.Minute
	logTicker := time.NewTicker(logPeriod)
	defer logTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-logTicker.C:
			err := p.ProcessPastLogs(ctx)
			if err != nil {
				p.logger.Error("failed to process past logs", "error", err)
				// TODO: Should we return error here or
				// continue to retry in the next tick?
			}
			continue
		}
	}
}
