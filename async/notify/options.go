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

package notify

import (
	"github.com/berachain/beacon-kit/runtime/service"
	"github.com/ethereum/go-ethereum/event"
)

// WithBaseService is an option to set the BaseService for the Service.
func WithBaseService(bs service.BaseService) service.Option[Service] {
	return func(s *Service) error {
		s.BaseService = bs

		// We piggyback initialization of the feeds and handlers maps here.
		s.feeds = make(map[string]*event.Feed)
		s.handlers = make(map[string][]eventHandlerQueuePair)
		return nil
	}
}

// WithGCD is an option to set the GrandCentralDispatch for the Service.
func WithGCD(gcd GrandCentralDispatch) service.Option[Service] {
	return func(s *Service) error {
		s.gcd = gcd

		return nil
	}
}
