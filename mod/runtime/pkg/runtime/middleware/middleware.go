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

package middleware

import (
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
)

// FinalizeBlockMiddleware is a struct that encapsulates the necessary
// components to handle the proposal processes.
type FinalizeBlockMiddleware[
	BeaconBlockT ssz.Marshallable,
	BeaconStateT any,
	BlobsSidecarsT ssz.Marshallable,
] struct {
	// chainSpec is the chain specification.
	chainSpec primitives.ChainSpec
	// chainService represents the blockchain service.
	chainService BlockchainService[BeaconBlockT, BlobsSidecarsT]
	// metrics is the metrics for the middleware.
	metrics *finalizeMiddlewareMetrics
	// valUpdates caches the validator updates as they are produced.
	valUpdates []*transition.ValidatorUpdate
}

// NewFinalizeBlockMiddleware creates a new instance of the Middleware struct.
func NewFinalizeBlockMiddleware[
	BeaconBlockT ssz.Marshallable,
	BeaconStateT any, BlobsSidecarsT ssz.Marshallable,
](
	chainSpec primitives.ChainSpec,
	chainService BlockchainService[BeaconBlockT, BlobsSidecarsT],
	telemetrySink TelemetrySink,
) *FinalizeBlockMiddleware[BeaconBlockT, BeaconStateT, BlobsSidecarsT] {
	// This is just for nilaway, TODO: remove later.
	if chainService == nil {
		panic("chain service is nil")
	}

	return &FinalizeBlockMiddleware[BeaconBlockT, BeaconStateT, BlobsSidecarsT]{
		chainSpec:    chainSpec,
		chainService: chainService,
		metrics:      newFinalizeMiddlewareMetrics(telemetrySink),
	}
}
