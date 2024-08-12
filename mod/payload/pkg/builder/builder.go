// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package builder

import (
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/payload/pkg/cache"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// PayloadBuilder is used to build payloads on the
// execution client.
type PayloadBuilder[
	AttributesFactoryT AttributesFactory[BeaconStateT, PayloadAttributesT],
	BeaconStateT BeaconState[ExecutionPayloadHeaderT, WithdrawalT],
	ExecutionEngineT ExecutionEngine[ExecutionPayloadT, PayloadAttributesT, PayloadIDT],
	ExecutionPayloadT ExecutionPayload[ExecutionPayloadT],
	ExecutionPayloadHeaderT ExecutionPayloadHeader,
	LoggerT log.Logger[any],
	PayloadAttributesT PayloadAttributes[PayloadAttributesT, WithdrawalT],
	PayloadIDT ~[8]byte,
	WithdrawalT any,
] struct {
	// cfg holds the configuration settings for the PayloadBuilder.
	cfg *Config
	// chainSpec holds the chain specifications for the PayloadBuilder.
	chainSpec common.ChainSpec
	// logger is used for logging within the PayloadBuilder.
	logger LoggerT
	// ee is the execution engine.
	ee ExecutionEngineT
	// pc is the payload ID cache, it is used to store
	// "in-flight" payloads that are being built on
	// the execution client.
	pc *cache.PayloadIDCache[
		PayloadIDT, [32]byte, math.Slot,
	]
	// attributesFactory is used to create attributes for the
	attributesFactory AttributesFactoryT
}

// New creates a new service.
func New[
	AttributesFactoryT AttributesFactory[BeaconStateT, PayloadAttributesT],
	BeaconStateT BeaconState[ExecutionPayloadHeaderT, WithdrawalT],
	ExecutionEngineT ExecutionEngine[ExecutionPayloadT, PayloadAttributesT, PayloadIDT],
	ExecutionPayloadT ExecutionPayload[ExecutionPayloadT],
	ExecutionPayloadHeaderT ExecutionPayloadHeader,
	LoggerT log.Logger[any],
	PayloadAttributesT PayloadAttributes[PayloadAttributesT, WithdrawalT],
	PayloadIDT ~[8]byte,
	WithdrawalT any,
](
	cfg *Config,
	chainSpec common.ChainSpec,
	logger LoggerT,
	ee ExecutionEngineT,
	pc *cache.PayloadIDCache[
		PayloadIDT, [32]byte, math.Slot,
	],
	af AttributesFactoryT,
) *PayloadBuilder[
	AttributesFactoryT, BeaconStateT, ExecutionEngineT,
	ExecutionPayloadT, ExecutionPayloadHeaderT, LoggerT,
	PayloadAttributesT, PayloadIDT, WithdrawalT,
] {
	return &PayloadBuilder[
		AttributesFactoryT, BeaconStateT, ExecutionEngineT,
		ExecutionPayloadT, ExecutionPayloadHeaderT, LoggerT,
		PayloadAttributesT, PayloadIDT, WithdrawalT,
	]{
		cfg:               cfg,
		chainSpec:         chainSpec,
		logger:            logger,
		ee:                ee,
		pc:                pc,
		attributesFactory: af,
	}
}

// Enabled returns true if the payload builder is enabled.
func (pb *PayloadBuilder[
	_, _, _, _, _, _, _, _, _,
]) Enabled() bool {
	return pb.cfg.Enabled
}
