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

package execution

import (
<<<<<<< HEAD
	"github.com/itsdevbear/bolaris/beacon/execution/logs"
	"github.com/itsdevbear/bolaris/beacon/staking"
=======
>>>>>>> 2ef1c7a568f885e74b3d5be498a841a391453bac
	"github.com/itsdevbear/bolaris/cache"
	"github.com/itsdevbear/bolaris/execution/engine"
)

type Option func(*Service) error

// WithEngineCaller is an option to set the Caller for the Service.
func WithEngineCaller(ec engine.Caller) Option {
	return func(s *Service) error {
		s.engine = ec
		return nil
	}
}

// WithPayloadCache is an option to set the PayloadIDCache for the Service.
func WithPayloadCache(pc *cache.PayloadIDCache) Option {
	return func(s *Service) error {
		s.payloadCache = pc
		return nil
	}
}

// WithLogProcessor is an option to set the LogProcessor for the Service.
func WithLogProcessor(lp *logs.Processor) Option {
	return func(s *Service) error {
		s.logProcessor = lp
		return nil
	}
}

func WithStakingService(st *staking.Service) Option {
	return func(s *Service) error {
		s.st = st
		return nil
	}
}
