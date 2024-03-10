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

package client

import (
	"context"
	"math/big"

	"github.com/berachain/beacon-kit/primitives"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	coretypes "github.com/ethereum/go-ethereum/core/types"
)

// HeaderByNumber retrieves the block header by its number.
func (s *EngineClient) HeaderByNumber(
	ctx context.Context,
	number *big.Int,
) (*coretypes.Header, error) {
	if header, ok := s.engineCache.HeaderByNumber(number.Uint64()); ok {
		return header, nil
	}

	header, err := s.Client.HeaderByNumber(ctx, number)
	if err != nil {
		return nil, err
	}

	s.engineCache.AddHeader(header)
	return header, nil
}

// GetLogs retrieves the logs from the Ethereum execution client.
// It calls the eth_getLogs method via JSON-RPC.
func (s *EngineClient) GetLogs(
	ctx context.Context,
	blockHash common.Hash,
	addresses []primitives.ExecutionAddress,
) ([]coretypes.Log, error) {
	// Gather all the logs according to the query.
	logs, err := s.FilterLogs(ctx, ethereum.FilterQuery{
		Addresses: addresses,
		BlockHash: &blockHash,
	})
	if err != nil {
		return []coretypes.Log{}, err
	}

	// Add logs to cache.
	return logs, nil
}
