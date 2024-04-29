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
	"github.com/berachain/beacon-kit/mod/node-builder/service"
)

// WithBaseService returns an Option that sets the BaseService for the Service.
func WithBaseService(base service.BaseService) service.Option[Service] {
	return func(s *Service) error {
		s.BaseService = base
		return nil
	}
}

// WithBlockValidator is a function that returns an Option.
// It sets the BlockValidator of the Service to the provided Service.
func WithBlockValidator(bv *core.BlockValidator) service.Option[Service] {
	return func(s *Service) error {
		s.bv = bv
		return nil
	}
}

// WithExecutionService is a function that returns an Option.
// It sets the ExecutionService of the Service to the provided Service.
func WithExecutionEngine(ee ExecutionEngine) service.Option[Service] {
	return func(s *Service) error {
		s.ee = ee
		return nil
	}
}

// WithLocalBuilder is a function that returns an Option.
// It sets the BuilderService of the Service to the provided Service.
func WithLocalBuilder(lb LocalBuilder) service.Option[Service] {
	return func(s *Service) error {
		s.lb = lb
		return nil
	}
}

// WithPayloadValidator is a function that returns an Option.
func WithPayloadValidator(pv *core.PayloadValidator) service.Option[Service] {
	return func(s *Service) error {
		s.pv = pv
		return nil
	}
}

// WithExecutionService is a function that returns an Option.
// It sets the ExecutionService of the Service to the provided Service.
func WithStakingService(sks StakingService) service.Option[Service] {
	return func(s *Service) error {
		s.sks = sks
		return nil
	}
}

// WithStateProcessor is a function that returns an Option.
func WithStateProcessor(sp *core.StateProcessor) service.Option[Service] {
	return func(s *Service) error {
		s.sp = sp
		return nil
	}
}
