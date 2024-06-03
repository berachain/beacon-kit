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
	engineprimitves "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/payload/pkg/cache"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// PayloadBuilder is used to build payloads on the
// execution client.
type PayloadBuilder[
	BeaconStateT BeaconState[ExecutionPayloadHeaderT],
	ExecutionPayloadT interface {
		IsNil() bool
		Empty(uint32) ExecutionPayloadT
		GetBlockHash() common.ExecutionHash
		GetFeeRecipient() common.ExecutionAddress
		GetParentHash() common.ExecutionHash
	},
	ExecutionPayloadHeaderT interface {
		GetBlockHash() common.ExecutionHash
		GetParentHash() common.ExecutionHash
	},
] struct {
	// cfg holds the configuration settings for the PayloadBuilder.
	cfg *Config
	// chainSpec holds the chain specifications for the PayloadBuilder.
	chainSpec primitives.ChainSpec
	// logger is used for logging within the PayloadBuilder.
	logger log.Logger[any]
	// ee is the execution engine.
	ee ExecutionEngine[ExecutionPayloadT]
	// pc is the payload ID cache, it is used to store
	// "in-flight" payloads that are being built on
	// the execution client.
	pc *cache.PayloadIDCache[
		engineprimitves.PayloadID, [32]byte, math.Slot,
	]
}

// NewService creates a new service.
func New[
	BeaconStateT BeaconState[ExecutionPayloadHeaderT],
	ExecutionPayloadT interface {
		IsNil() bool
		Empty(uint32) ExecutionPayloadT
		GetBlockHash() common.ExecutionHash
		GetParentHash() common.ExecutionHash
		GetFeeRecipient() common.ExecutionAddress
	},
	ExecutionPayloadHeaderT interface {
		GetBlockHash() common.ExecutionHash
		GetParentHash() common.ExecutionHash
	},
](
	cfg *Config,
	chainSpec primitives.ChainSpec,
	logger log.Logger[any],
	ee ExecutionEngine[ExecutionPayloadT],
	pc *cache.PayloadIDCache[
		engineprimitves.PayloadID, [32]byte, math.Slot,
	],
) *PayloadBuilder[
	BeaconStateT, ExecutionPayloadT, ExecutionPayloadHeaderT,
] {
	return &PayloadBuilder[
		BeaconStateT, ExecutionPayloadT, ExecutionPayloadHeaderT,
	]{
		cfg:       cfg,
		chainSpec: chainSpec,
		logger:    logger,
		ee:        ee,
		pc:        pc,
	}
}

// Enabled returns true if the payload builder is enabled.
func (pb *PayloadBuilder[
	BeaconStateT, ExecutionPayloadT, ExecutionPayloadHeaderT,
]) Enabled() bool {
	return pb.cfg.Enabled
}
