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

package initialsync

import (
	"context"
	"math/big"

	"cosmossdk.io/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Option is a function that modifies the Service.
type Option func(*Service) error

// ethClient is an interface that wraps the ChainSyncReader from the go-ethereum package.
type ethClient interface {
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error)
}

// WithEthClient is an Option that sets the ethClient of the Service.
func WithEthClient(ethClient ethClient) Option {
	return func(s *Service) error {
		s.ethClient = ethClient
		return nil
	}
}

// WithLogger is an Option that sets the logger of the Service.
func WithLogger(logger log.Logger) Option {
	return func(r *Service) error {
		r.logger = logger.With("module", "beacon-kit-sync")
		return nil
	}
}

// WithForkChoiceStoreProvider is an Option that sets the ForkChoiceStoreProvider
// of the Service.
func WithForkChoiceStoreProvider(fcsp ForkChoiceStoreProvider) Option {
	return func(r *Service) error {
		r.fcsp = fcsp
		return nil
	}
}
