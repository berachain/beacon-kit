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

package logs_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/berachain/beacon-kit/beacon/execution/logs"
	logsMocks "github.com/berachain/beacon-kit/beacon/execution/logs/mocks"
	"github.com/berachain/beacon-kit/primitives"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethcoretypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

func TestLogProcessor(t *testing.T) {
	ctx := context.Background()

	blockNum := new(big.Int).SetUint64(100)
	header := &ethcoretypes.Header{
		Number: blockNum,
	}
	blockHash := header.Hash()
	logSignature := ethcommon.HexToHash("0x1234")
	contractAddress := ethcommon.HexToAddress("0x5678")

	log := ethcoretypes.Log{
		Address:     contractAddress,
		BlockNumber: blockNum.Uint64(),
		Topics:      []ethcommon.Hash{logSignature},
	}

	blockNumToHash := map[uint64]ethcommon.Hash{
		blockNum.Uint64(): blockHash,
	}

	container := &logsMocks.LogContainer{}
	container.EXPECT().Signature().
		Return(logSignature)

	client := &logsMocks.LogEngineClient{}
	client.EXPECT().HeaderByHash(ctx, blockHash).
		Return(header, nil)
	client.EXPECT().HeaderByNumber(ctx, blockNum).
		Return(header, nil)
	client.EXPECT().GetLogsInRange(
		ctx,
		blockNum.Uint64(),
		blockNum.Uint64(),
		[]primitives.ExecutionAddress{contractAddress},
	).Return([]ethcoretypes.Log{log}, nil)

	factory := &logsMocks.LogFactory{}
	factory.EXPECT().GetRegisteredSignatures().
		Return([]ethcommon.Hash{logSignature})
	factory.EXPECT().GetRegisteredAddresses().
		Return([]primitives.ExecutionAddress{contractAddress})
	factory.EXPECT().ProcessLogs(
		[]ethcoretypes.Log{log}, blockNumToHash,
	).Return([]logs.LogContainer{container}, nil)

	cache := &SimpleCache{}

	// Create a new log processor.
	processor, err := logs.NewProcessor(
		logs.WithEngineClient(client),
		logs.WithLogFactory(factory),
		logs.WithLogCache(logSignature, cache),
	)
	require.NoError(t, err)

	err = processor.ProcessEth1Block(
		ctx,
		blockHash,
	)
	require.NoError(t, err)
	require.Len(t, cache.store, 1)
	require.Equal(t, container, cache.store[0])
}

type SimpleCache struct {
	store []logs.LogContainer
}

func (c *SimpleCache) Insert(log logs.LogContainer) error {
	c.store = append(c.store, log)
	return nil
}

func (c *SimpleCache) RemoveMulti(
	index uint64, n uint64,
) ([]logs.LogContainer, error) {
	if len(c.store) == 0 || c.store[0].LogIndex() != index {
		return nil, errors.New("index not found")
	}
	logs := c.store[:n]
	c.store = c.store[n:]
	return logs, nil
}
