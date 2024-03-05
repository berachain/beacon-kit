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
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// Cache implements the LogCache interface.
var _ LogCache = (*Cache)(nil)

type Cache struct {
	// The final store contains processed
	// logs from complete blocks.
	finalStore []LogValueContainer

	// lastProcessedBlock records the block number
	// of the latest block up to which
	// the logs are processed and stored in cache.
	// This is in-memory only and will be reset to 0
	// when the node restarts.
	lastProcessedBlock uint64

	// The processing store contains logs from
	// the current block being processed.
	// They will be moved to the final store
	// when the block is set as the last finalized block.
	processingStore []LogValueContainer

	// currProcessingBlock records the block number
	// of the current block being processed.
	// It is always greater than lastProcessedBlock.
	currProcessingBlock uint64

	// currProcessingIndex records the index of the
	// current log being processed in the current block.
	currProcessingIndex uint64
}

// ShouldProcess returns true if the cache determines
// that the log should be processed and added to it.
func (c *Cache) ShouldProcess(log *ethtypes.Log) bool {
	return log.BlockNumber > c.currProcessingBlock ||
		(log.BlockNumber == c.currProcessingBlock &&
			(c.currProcessingIndex == 0 ||
				uint64(log.Index) > c.currProcessingIndex))
}

// Push pushes the log value container into the cache.
// In this implementation, the cache keeps the item
// in its temporary processing store until it finishes
// processing all the logs in the block and sets the
// block as the last finalized block.
func (c *Cache) Push(container LogValueContainer) error {
	// ShouldProcess should be called before Push
	// to avoid unnecessary processing.
	c.processingStore = append(c.processingStore, container)
	c.currProcessingBlock = container.BlockNumber()
	c.currProcessingIndex = container.LogIndex()
	return nil
}

// LastProcessedBlock returns the block number of
// the last completely processed block in cache.
func (c *Cache) LastProcessedBlock() uint64 {
	return c.lastProcessedBlock
}

// SetLastProcessedBlock sets the block number of
// the last completely processed block in cache.
// The cache will move the logs from the processing
// store to the final store at this time.
func (c *Cache) SetLastProcessedBlock(blockNumber uint64) {
	c.lastProcessedBlock = blockNumber
	c.finalStore = append(c.finalStore, c.processingStore...)
	c.processingStore = nil
	// We are expecting the block after blockNumber to be processed.
	c.currProcessingBlock = blockNumber + 1
	c.currProcessingIndex = 0
}

// SetCurrentProcessingBlock sets the block number of
// the current block being processed.
func (c *Cache) SetCurrentProcessingBlock(blockNumber uint64) {
	c.currProcessingBlock = blockNumber
}

// SetCurrentProcessingIndex sets the index of the
// current log being processed in the current block.
func (c *Cache) SetCurrentProcessingIndex(index uint64) {
	c.currProcessingIndex = index
}

// Rollback rolls back the cache to the last completely processed block
// if there is any error during processing the logs
// in the current block.
func (c *Cache) Rollback() {
	c.processingStore = nil
	c.currProcessingBlock = c.lastProcessedBlock + 1
	c.currProcessingIndex = 0
}
