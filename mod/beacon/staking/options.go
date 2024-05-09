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

package staking

import (
	stakingabi "github.com/berachain/beacon-kit/mod/beacon/staking/abi"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/execution/pkg/engine"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

// Option is a function that sets a field on the Service.
type Option func(*Service) error

// WithBeaconStorageBackend is a function that returns an Option.
func WithBeaconStorageBackend(bsb BeaconStorageBackend) Option {
	return func(s *Service) error {
		s.bsb = bsb
		return nil
	}
}

// WithChainSpec is a function that returns an Option.
// It sets the ChainSpec of the Service to the provided ChainSpec.
func WithChainSpec(cs primitives.ChainSpec) Option {
	return func(s *Service) error {
		s.cs = cs
		return nil
	}
}

// WithDepositABI returns an Option that sets the deposit
// contract's ABI for the Service.
func WithDepositABI(
	depositABI *abi.ABI,
) Option {
	return func(s *Service) error {
		s.abi = stakingabi.NewWrappedABI(depositABI)
		return nil
	}
}

// WithDepositStore returns an Option that sets the deposit.
func WithDepositStore(
	ds *deposit.KVStore,
) Option {
	return func(s *Service) error {
		s.ds = ds
		return nil
	}
}

// WithExecutionEngine is a function that returns an Option.
func WithExecutionEngine(
	ee *engine.Engine[
		types.ExecutionPayload, *types.ExecutableDataDeneb,
	],
) Option {
	return func(s *Service) error {
		s.ee = ee
		return nil
	}
}

// WithLogger is a function that returns an Option.
// It sets the Logger of the Service to the provided Logger.
func WithLogger(logger log.Logger[any]) Option {
	return func(s *Service) error {
		s.logger = logger
		return nil
	}
}
