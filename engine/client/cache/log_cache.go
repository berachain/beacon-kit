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

package cache

import (
	"sync"

	"github.com/berachain/beacon-kit/lib/skiplist"
	"github.com/ethereum/go-ethereum/common"
	coretypes "github.com/ethereum/go-ethereum/core/types"
)

type LogCache struct {
	skiplist     *skiplist.Skiplist[coretypes.Log]
	blocknumMap  map[uint64][]coretypes.Log
	blockhashMap map[common.Hash][]coretypes.Log

	rwmutex *sync.RWMutex
}

func NewLogCache() *LogCache {
	return &LogCache{
		skiplist:     skiplist.NewSkiplist(logCompare{}),
		blocknumMap:  make(map[uint64][]coretypes.Log),
		blockhashMap: make(map[common.Hash][]coretypes.Log),
		rwmutex:      new(sync.RWMutex),
	}
}

type logCompare struct{}

func (l logCompare) Compare(lhs, rhs coretypes.Log) int {
	// Compare block number
	if lhs.BlockNumber < rhs.BlockNumber {
		return -1
	} else if lhs.BlockNumber > rhs.BlockNumber {
		return 1
	}

	// Compare transaction index
	if lhs.TxIndex < rhs.TxIndex {
		return -1
	} else if lhs.TxIndex > rhs.TxIndex {
		return 1
	}

	// Compare log index
	if lhs.Index < rhs.Index {
		return -1
	} else if lhs.Index > rhs.Index {
		return 1
	}

	return 0
}

// Add adds logs to the cache.
func (c *LogCache) Add(logs []coretypes.Log) {
	c.rwmutex.Lock()
	defer c.rwmutex.Unlock()
	if len(logs) > 0 {
		c.blocknumMap[logs[0].BlockNumber] = logs
		c.blockhashMap[logs[0].BlockHash] = logs
		for _, log := range logs {
			c.skiplist.Insert(log)
		}
	}
}

// GetByBlockNumber returns all logs with the given block number.
func (c *LogCache) GetByBlockNumber(blocknum uint64) ([]coretypes.Log, bool) {
	c.rwmutex.RLock()
	defer c.rwmutex.RUnlock()
	logs, ok := c.blocknumMap[blocknum]
	return logs, ok
}

// GetByBlockHash returns all logs with the given block hash.
func (c *LogCache) GetByBlockHash(
	blockhash common.Hash,
) ([]coretypes.Log, bool) {
	c.rwmutex.RLock()
	defer c.rwmutex.RUnlock()
	logs, ok := c.blockhashMap[blockhash]
	return logs, ok
}

// Remove removes a log from the cache.
func (c *LogCache) Remove(log coretypes.Log) {
	c.rwmutex.Lock()
	defer c.rwmutex.Unlock()
	c.skiplist.Remove(log)
	// Prune the individual log from the maps
	c.removeLog(log)
}

// RemoveByBlockNumber removes all logs with the given block number from the
// cache.
func (c *LogCache) RemoveByBlockNumber(blocknum uint64) {
	c.rwmutex.Lock()
	defer c.rwmutex.Unlock()
	logs, ok := c.blocknumMap[blocknum]
	if ok && len(logs) > 0 {
		// Remove logs from skiplist
		for _, log := range logs {
			c.skiplist.Remove(log)
		}
		// Remove logs from blockhash map
		delete(c.blockhashMap, logs[0].BlockHash)
		// Remove logs from blocknum map
		delete(c.blocknumMap, blocknum)
	}
}

// RemoveByBlockHash removes all logs with the given block hash from the cache.
func (c *LogCache) RemoveByBlockHash(blockhash common.Hash) {
	c.rwmutex.Lock()
	defer c.rwmutex.Unlock()
	logs, ok := c.blockhashMap[blockhash]
	if ok && len(logs) > 0 {
		// Remove logs from skiplist
		for _, log := range logs {
			c.skiplist.Remove(log)
		}
		// Remove logs from blocknum map
		delete(c.blocknumMap, logs[0].BlockNumber)
		// Remove logs from blockhash map
		delete(c.blockhashMap, blockhash)
	}
}

// ContainsBlockNumber returns true if the cache contains logs with the given
// block number.
func (c *LogCache) ContainsBlockNumber(blocknum uint64) bool {
	c.rwmutex.RLock()
	defer c.rwmutex.RUnlock()
	_, ok := c.blocknumMap[blocknum]
	return ok
}

// ContainsBlockHash returns true if the cache contains logs with the given
// block hash.
func (c *LogCache) ContainsBlockHash(blockhash common.Hash) bool {
	c.rwmutex.RLock()
	defer c.rwmutex.RUnlock()
	_, ok := c.blockhashMap[blockhash]
	return ok
}

// Front returns the first log in the skiplist.
func (c *LogCache) Front() (coretypes.Log, error) {
	c.rwmutex.RLock()
	defer c.rwmutex.RUnlock()
	log, err := c.skiplist.Front()
	return log, err
}

// RemoveFront removes the first log in the skiplist.
func (c *LogCache) RemoveFront() (coretypes.Log, error) {
	c.rwmutex.Lock()
	defer c.rwmutex.Unlock()
	log, err := c.skiplist.RemoveFront()
	if err == nil {
		c.removeLog(log)
	}
	return log, err
}

// Back returns the last log in the skiplist.
func (c *LogCache) Back() (coretypes.Log, error) {
	c.rwmutex.RLock()
	defer c.rwmutex.RUnlock()
	log, err := c.skiplist.Back()
	return log, err
}

// RemoveBack removes the last log in the skiplist.
func (c *LogCache) RemoveBack() (coretypes.Log, error) {
	c.rwmutex.Lock()
	defer c.rwmutex.Unlock()
	log, err := c.skiplist.RemoveBack()
	if err == nil {
		c.removeLog(log)
	}
	return log, err
}

// removeLog prunes an individual log from the maps.
func (c *LogCache) removeLog(log coretypes.Log) {
	// Lock is already held by the caller
	logs := c.blocknumMap[log.BlockNumber]
	newLogs := make([]coretypes.Log, 0, len(logs))
	// Remove the log from the blocknum map, if the log index is equal
	for _, l := range logs {
		if l.Index != log.Index {
			newLogs = append(newLogs, l)
		}
	}
	if len(newLogs) == 0 {
		delete(c.blocknumMap, log.BlockNumber)
		delete(c.blockhashMap, log.BlockHash)
	} else {
		c.blocknumMap[log.BlockNumber] = newLogs
		c.blockhashMap[log.BlockHash] = newLogs
	}
}
