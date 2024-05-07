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
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/payload/pkg/cache"
	"github.com/berachain/beacon-kit/mod/primitives"
	engineprimitves "github.com/berachain/beacon-kit/mod/primitives-engine"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// TODO: Decouple from ABCI and have this validator run on a separate thread
// have it configured itself and not be a service persay.
type PayloadBuilder struct {
	cfg       *Config
	chainSpec primitives.ChainSpec
	logger    log.Logger[any]

	// ee is the execution engine.
	ee ExecutionEngine

	// pc is the payload ID cache, it is used to store
	// "in-flight" payloads that are being built on
	// the execution client.
	pc *cache.PayloadIDCache[
		engineprimitves.PayloadID, [32]byte, math.Slot,
	]
}

// NewService creates a new service.
func New(opts ...Option) (*PayloadBuilder, error) {
	pb := &PayloadBuilder{}
	for _, opt := range opts {
		if err := opt(pb); err != nil {
			return nil, err
		}
	}
	return pb, nil
}
