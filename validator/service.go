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

package validator

import (
	"github.com/itsdevbear/bolaris/beacon/execution/engine"
	"github.com/itsdevbear/bolaris/runtime/service"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/cache"
)

type BlockBuilder interface {
}

// TODO: Decouple from ABCI and have this validator run on a seperate thread
// have it configured itself and not be a service persay.
type Service struct {
	service.BaseService
	en           engine.Caller
	payloadCache *cache.PayloadIDCache
}

func NewService(
	base service.BaseService,
	opts ...Option,
) *Service {
	s := &Service{
		BaseService: base,
	}

	for _, opt := range opts {
		if err := opt(s); err != nil {
			panic(err)
		}
	}
	return s
}

func (s *Service) Start() {
}

func (s *Service) Stop() error {
	return nil
}

func (s *Service) Status() error {
	return nil
}
