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

package blockchain

import (
	"github.com/berachain/beacon-kit/mod/core"
	datypes "github.com/berachain/beacon-kit/mod/da/pkg/types"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives"
)

// Option is a function that sets a field on the Service.
type Option func(*Service) error

func WithBeaconStorageBackend(bsb BeaconStorageBackend) Option {
	return func(s *Service) error {
		s.bsb = bsb
		return nil
	}
}

// It sets the BlockVerifier of the Service to the provided Service.
func WithBlockVerifier(bv *core.BlockVerifier) Option {
	return func(s *Service) error {
		s.bv = bv
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

// WithExecutionService is a function that returns an Option.
// It sets the ExecutionService of the Service to the provided Service.
func WithExecutionEngine(ee ExecutionEngine) Option {
	return func(s *Service) error {
		s.ee = ee
		return nil
	}
}

// WithLocalBuilder is a function that returns an Option.
// It sets the BuilderService of the Service to the provided Service.
func WithLocalBuilder(lb LocalBuilder) Option {
	return func(s *Service) error {
		s.lb = lb
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

// WithPayloadVerifier is a function that returns an Option.
func WithPayloadVerifier(pv *core.PayloadVerifier) Option {
	return func(s *Service) error {
		s.pv = pv
		return nil
	}
}

// WithExecutionService is a function that returns an Option.
// It sets the ExecutionService of the Service to the provided Service.
func WithStakingService(sks StakingService) Option {
	return func(s *Service) error {
		s.sks = sks
		return nil
	}
}

// WithStateProcessor is a function that returns an Option.
func WithStateProcessor(
	sp *core.StateProcessor[*datypes.BlobSidecars],
) Option {
	return func(s *Service) error {
		s.sp = sp
		return nil
	}
}
