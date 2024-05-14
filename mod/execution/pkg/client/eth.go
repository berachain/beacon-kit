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

	engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
)

// HeaderByNumber retrieves the block header by its number.
func (s *EngineClient[ExecutionPayloadDenebT]) HeaderByNumber(
	ctx context.Context,
	number *big.Int,
) (*engineprimitives.Header, error) {
	// Infer the latest height if the number is nil.
	if number == nil {
		latest, err := s.BlockNumber(ctx)
		if err != nil {
			return nil, err
		}

		number = new(big.Int).SetUint64(latest)
	}

	// Check the cache for the header.
	header, ok := s.engineCache.HeaderByNumber(number.Uint64())
	if ok {
		return header, nil
	}

	header, err := s.Client.HeaderByNumber(ctx, number)
	if err != nil {
		return nil, err
	}

	defer s.engineCache.AddHeader(header)

	return header, nil
}

// HeaderByHash retrieves the block header by its hash.
func (s *EngineClient[ExecutionPayloadDenebT]) HeaderByHash(
	ctx context.Context,
	hash common.ExecutionHash,
) (*engineprimitives.Header, error) {
	// Check the cache for the header.
	header, ok := s.engineCache.HeaderByHash(hash)
	if ok {
		return header, nil
	}
	header, err := s.Client.HeaderByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	s.engineCache.AddHeader(header)
	return header, nil
}
