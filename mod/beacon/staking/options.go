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
	"github.com/berachain/beacon-kit/mod/log"
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

// WithDepositContract returns an Option that sets the deposit
// contract's ABI for the Service.
func WithDepositContract(
	depositContract DepositContract,
) Option {
	return func(s *Service) error {
		s.depositContract = depositContract
		return nil
	}
}

// WithDepositStore returns an Option that sets the deposit.
func WithDepositStore(
	ds DepositStore,
) Option {
	return func(s *Service) error {
		s.ds = ds
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
