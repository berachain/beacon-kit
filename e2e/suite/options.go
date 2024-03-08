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

package suite

import (
	"context"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/kurtosis"
)

// Type Option is a function that sets a field on the KurtosisE2ESuite.
type Option func(*KurtosisE2ESuite) error

// WithConfig sets the E2ETestConfig for the test suite.
func WithConfig(cfg *kurtosis.E2ETestConfig) Option {
	return func(s *KurtosisE2ESuite) error {
		s.cfg = cfg
		return nil
	}
}

// WithContext sets the context for the test suite.
func WithContext(ctx context.Context) Option {
	return func(s *KurtosisE2ESuite) error {
		s.ctx = ctx
		return nil
	}
}

// WithLogger sets the logger for the test suite.
func WithLogger(logger log.Logger) Option {
	return func(s *KurtosisE2ESuite) error {
		s.logger = logger
		return nil
	}
}
