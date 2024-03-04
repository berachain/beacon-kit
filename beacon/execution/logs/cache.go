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
	// logs from finalized blocks.
	finalStore []LogValueContainer
	// lastFinalizedBlock records the block number
	// of the latest finalized block up to which
	// the logs are processed and stored in cache.
	// This is in-memory only and will be reset to 0
	// when the node restarts.
	lastFinalizedBlock uint64

	// The processing store contains logs from
	// the current block being processed.
	// They will be moved to the final store
	// when the block is set as the last finalized block.
	processingStore []LogValueContainer
}

func (c *Cache) ShouldProcess(log *ethtypes.Log) bool {
	return log.BlockNumber > c.lastFinalizedBlock
}

func (c *Cache) Push(container LogValueContainer) error {
	// ShouldProcess should be called before Push
	// to avoid unnecessary processing.
	c.processingStore = append(c.processingStore, container)
	return nil
}

func (c *Cache) LastFinalizedBlock() uint64 {
	return c.lastFinalizedBlock
}

func (c *Cache) SetLastFinalizedBlock(blockNumber uint64) {
	c.lastFinalizedBlock = blockNumber
	c.finalStore = append(c.finalStore, c.processingStore...)
	c.processingStore = nil
}

func (c *Cache) Rollback() {
	c.processingStore = nil
}
