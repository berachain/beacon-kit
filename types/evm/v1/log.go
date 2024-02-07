// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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

package v1

import (
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// Logs is a slice of pointers to Log.
type Logs []*Log

// ToEthLogs converts Logs to a slice of pointers to Geth Log.
func (logs Logs) ToEthLogs() []*ethtypes.Log {
	ethLogs := make([]*ethtypes.Log, len(logs))
	for i := range logs {
		ethLogs[i] = logs[i].ToEthLog()
	}
	return ethLogs
}

// LogsFromEthLogs converts a slice of pointers to Geth Log to Logs.
func LogsFromEthLogs(ethlogs []ethtypes.Log) Logs {
	logs := make(Logs, len(ethlogs))
	for i, ethlog := range ethlogs {
		logs[i] = NewLogFromGethLog(ethlog)
	}
	return logs
}

// ToEthLog converts a protobuf Log to a Geth Log.
func (log *Log) ToEthLog() *ethtypes.Log {
	topics := make([]common.Hash, len(log.GetTopics()))
	for i, topic := range log.GetTopics() {
		topics[i] = common.BytesToHash(topic)
	}

	return &ethtypes.Log{
		Address:     common.BytesToAddress(log.GetAddress()),
		Topics:      topics,
		Data:        log.GetData(),
		BlockNumber: log.GetBlockNumber(),
		TxHash:      common.BytesToHash(log.GetTxHash()),
		TxIndex:     uint(log.GetTxIndex()),
		Index:       uint(log.GetIndex()),
		BlockHash:   common.BytesToHash(log.GetBlockHash()),
		Removed:     log.GetRemoved(),
	}
}

// NewLogFromGethLog creates a new Log from a Geth Log.
func NewLogFromGethLog(log ethtypes.Log) *Log {
	topics := make([][]byte, len(log.Topics))
	for i, topic := range log.Topics {
		topics[i] = topic.Bytes()
	}

	return &Log{
		Address:     log.Address.Bytes(),
		Topics:      topics,
		Data:        log.Data,
		BlockNumber: log.BlockNumber,
		TxHash:      log.TxHash.Bytes(),
		TxIndex:     uint64(log.TxIndex),
		Index:       uint64(log.Index),
		BlockHash:   log.BlockHash.Bytes(),
		Removed:     log.Removed,
	}
}
