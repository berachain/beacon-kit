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
	"github.com/berachain/beacon-kit/mod/execution"
	"github.com/berachain/beacon-kit/mod/node-builder/service"
	stakingabi "github.com/berachain/beacon-kit/mod/runtime/services/staking/abi"
	"github.com/berachain/beacon-kit/mod/storage/deposit"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

// WithBaseService sets the BaseService for the Service.
func WithBaseService(
	base service.BaseService,
) service.Option[Service] {
	return func(s *Service) error {
		s.BaseService = base
		return nil
	}
}

// WithDepositABI returns an Option that sets the deposit
// contract's ABI for the Service.
func WithDepositABI(
	depositABI *abi.ABI,
) service.Option[Service] {
	return func(s *Service) error {
		s.abi = stakingabi.NewWrappedABI(depositABI)
		return nil
	}
}

func WithExecutionEngine(
	ee *execution.Engine,
) service.Option[Service] {
	return func(s *Service) error {
		s.ee = ee
		return nil
	}
}

func WithDepositStore(
	ds *deposit.KVStore,
) service.Option[Service] {
	return func(s *Service) error {
		s.ds = ds
		return nil
	}
}
