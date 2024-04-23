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

package builder

import (
	"time"

	"github.com/berachain/beacon-kit/mod/primitives"
)

const (
	// defaultLocalBuilderEnabled is the default value for local builder.
	defaultLocalBuilderEnabled = true
	// defaultPayloadTimeout is the default value for local build
	// payload timeout.
	defaultPayloadTimeout = 2500 * time.Millisecond
)

// Config is the configuration for the payload builder.
//
//nolint:lll // struct tags.
type Config struct {
	// Enabled determines if the local builder is enabled.
	Enabled bool `mapstructure:"enabled"`

	// PayloadTimeout is the timeout for the payload build.
	PayloadTimeout time.Duration `mapstructure:"payload-timeout"`

	// Suggested FeeRecipient is the address that will receive the transaction
	// fees
	// produced by any blocks from this node.
	SuggestedFeeRecipient primitives.ExecutionAddress `mapstructure:"suggested-fee-recipient"`
}

// DefaultConfig returns the default fork configuration.
func DefaultConfig() Config {
	return Config{
		Enabled:               defaultLocalBuilderEnabled,
		SuggestedFeeRecipient: primitives.ExecutionAddress{},
		PayloadTimeout:        defaultPayloadTimeout,
	}
}
