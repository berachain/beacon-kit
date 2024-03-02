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

	ethcommon "github.com/ethereum/go-ethereum/common"
	engineclient "github.com/itsdevbear/bolaris/engine/client"
	"github.com/pkg/errors"
)

type Processor struct {
	fcp ReadOnlyForkChoiceProvider

	// engine gives the access to the Engine API
	// of the execution client.
	engine engineclient.Caller

	// logFactory is the factory for creating
	// objects from Ethereum logs.
	factory LogFactory

	// sigToCache is a map of log signatures to their caches.
	sigToCache map[ethcommon.Hash]LogCache
}

// ProcessBlocksInBatch processes the blocks in batch,
// from the last processed block (exclusive)
// to the latest block (inclusive) for each cache.
// This function will be called in a goroutine
// to prediodically process the logs and backfill
// the caches in background.
func (p *Processor) ProcessBlocksInBatch(
	ctx context.Context,
) error {
	// Get the latest finalized block hash and block number.
	forkChoicer := p.fcp.ForkchoiceStore(ctx)
	finalizedBlockHash := forkChoicer.GetFinalizedEth1BlockHash()
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
	for _, cache := range p.sigToCache {
		lastFinalizedBlockInCache := cache.LastFinalizedBlock()
		if lastFinalizedBlockInCache < minLastFinalizedBlockInCache {
			minLastFinalizedBlockInCache = lastFinalizedBlockInCache
		}
	}

	// If all caches have processed the latest finalized block,
	// we don't need to process it again.
	if minLastFinalizedBlockInCache == finalizedBlockNumber {
		return nil
	}

	// Gather all the logs corresponding to
	// the addresses of interest in the range
	// from the last processed block to the latest block.
	// TODO: Can we assume that the logs are returned in order?
	batchedLogs, err := p.engine.GetLogs(
		ctx,
		minLastFinalizedBlockInCache+1,
		finalizedBlockNumber,
		p.factory.GetRegisteredAddresses(),
	)
	if err != nil {
		// TODO: Handle TooMuchDataRequestedError.
		return errors.Wrapf(err, "failed to get logs")
	}

	// TODO: Use MapErr
	for i := range batchedLogs {
		log := &batchedLogs[i]
		cache, ok := p.sigToCache[log.Topics[0]]
		// Skip the log if it is not registered.
		if !ok {
			continue
		}
		// Cache determine if the log should be processed,
		// based on its last processed block.
		// TODO: Should we also consider the last processed index?
		if cache.ShouldProcess(log) {
			var container LogValueContainer
			container, err = p.factory.ProcessLog(log)
			if err != nil {
				return errors.Wrapf(err, "failed to process log")
			}
			err = cache.Push(container)
			if err != nil {
				return errors.Wrapf(err, "failed to push container")
			}
		}
	}

	// Update the caches with the new finalized block.
	for _, cache := range p.sigToCache {
		cache.SetLastFinalizedBlock(finalizedBlockNumber)
	}

	return nil
}
