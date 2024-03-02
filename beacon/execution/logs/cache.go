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
	"slices"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/maps"
)

type BlockchainService interface {
	IsFinalized(blockHash ethcommon.Hash, blockNumber uint64) bool
}

type BlockContainer struct {
	blockHash   ethcommon.Hash
	blockNumber uint64
	logValues   []LogValueContainer
}

type Cache struct {
	bs BlockchainService

	lastValidBlockInCache     uint64
	lastFinalizedBlockInCache uint64
	latestFinalizedBlock      uint64

	// The stores is ordered by block number in ascending order.
	// The pending store contains the blocks that are not finalized yet,
	// from (latestFinalizedBlock+1) to lastValidBlockInCache.
	pendingStore map[uint64]BlockContainer
	// The finalized store contains the finalized blocks,
	// which are before (lastFinalizedBlockInCache+1).
	//
	// Blocks in the range from (lastFinalizedBlockInCache+1) to
	// latestFinalizedBlock are not in the cache yet, but they
	// are already finalized and should be processed to add to the
	// cache.
	finalizedStore map[uint64]BlockContainer
}

func (c *Cache) Update(finalizedBlockNumber uint64) uint64 {
	c.latestFinalizedBlock = finalizedBlockNumber
	for _, blockNum := range OrderedMapKeys(c.pendingStore) {
		if blockNum > finalizedBlockNumber {
			// The block is not finalized yet, so it
			// is still valid wrt the finalized block.
			break
		}
		block := c.pendingStore[blockNum]
		if c.bs.IsFinalized(block.blockHash, block.blockNumber) {
			// The block is finalized, move it to the finalized store.
			c.finalizedStore[block.blockNumber] = block
			c.lastFinalizedBlockInCache = block.blockNumber
			delete(c.pendingStore, block.blockNumber)
		} else {
			// Discard all the remaining blocks in the pending store,
			// because they are built on top of an invalid block.
			// The last valid block in the cache is the last finalized block.
			c.pendingStore = make(map[uint64]BlockContainer)
			c.lastValidBlockInCache = c.lastFinalizedBlockInCache
			break
		}
	}
	return c.lastValidBlockInCache
}

func (c *Cache) ShouldProcess(log *ethtypes.Log) bool {
	// Logs from a finalized block should be processed.
	if log.BlockNumber <= c.latestFinalizedBlock {
		if log.BlockNumber < c.lastFinalizedBlockInCache {
			// The log is from a block that is completely processed.
			// Reorg cannot happen since the block is already finalized.
			return false
		}
		if log.BlockNumber == c.lastFinalizedBlockInCache {
			lastIndex, err := c.LastIndexInBlock(log.BlockNumber)
			if err != nil {
				return false
			}
			return log.Index > lastIndex
		}
		return true
	}

	// Logs before or at the last valid block should
	// be processed if reorg happens. The if condition
	// also covers the case log.BlockNumber > c.lastValidBlockInCache,
	// when the log is from a block that is not in the cache yet.
	if log.BlockHash != c.pendingStore[log.BlockNumber].blockHash {
		return true
	}

	// Logs from the last valid block is still being processed.
	if log.BlockNumber == c.lastValidBlockInCache {
		lastIndex, err := c.LastIndexInBlock(log.BlockNumber)
		if err != nil {
			return false
		}
		return log.Index > lastIndex
	}

	return false
}

func (c *Cache) Push(container LogValueContainer) error {
	blockNumber := container.BlockNumber()
	if blockNumber <= c.latestFinalizedBlock {
		// Block is finalized, add it to the finalized store.
		if block, ok := c.finalizedStore[blockNumber]; ok {
			block.logValues = append(block.logValues, container)
		} else {
			c.finalizedStore[blockNumber] = BlockContainer{
				blockHash:   container.BlockHash(),
				blockNumber: blockNumber,
				logValues:   []LogValueContainer{container},
			}
		}
		c.lastFinalizedBlockInCache = blockNumber
		// A finalized block is still being processed.
		// A non-empty pending store should be invalid.
		if len(c.pendingStore) > 0 {
			c.pendingStore = make(map[uint64]BlockContainer)
			c.lastValidBlockInCache = blockNumber
		}
		return nil
	}
	// Reorg or new block.
	if container.BlockHash() != c.pendingStore[blockNumber].blockHash {
		// Reorg happened if blockNumber < lastValidBlockInCache,
		// the block and its descendants are not valid anymore.
		for i := blockNumber; i <= c.lastValidBlockInCache; i++ {
			delete(c.pendingStore, i)
		}
		// Create a new block with the container.
		c.pendingStore[blockNumber] = BlockContainer{
			blockHash:   container.BlockHash(),
			blockNumber: blockNumber,
			logValues:   []LogValueContainer{container},
		}
	} else {
		// The last valid block is still being processed.
		// Add the container to the block.
		// ShouldProcess already checked the index.
		block := c.pendingStore[blockNumber]
		block.logValues = append(block.logValues, container)
	}
	c.lastValidBlockInCache = blockNumber
	return nil
}

func (c *Cache) LastIndexInBlock(blockNumber uint64) (uint, error) {
	var (
		block BlockContainer
		ok    bool
	)

	if blockNumber <= c.lastFinalizedBlockInCache {
		block, ok = c.finalizedStore[blockNumber]
	} else {
		block, ok = c.pendingStore[blockNumber]
	}
	if !ok {
		return 0, errors.Errorf("block %d not found in the cache", blockNumber)
	}
	return block.logValues[len(block.logValues)-1].LogIndex(), nil
}

// OrderedMapKeys returns the map keys in ascending order.
func OrderedMapKeys[K constraints.Ordered, V any](m map[K]V) []K {
	keys := maps.Keys(m)
	slices.Sort(keys)
	return keys
}
