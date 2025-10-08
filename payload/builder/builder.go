// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/primitives/common"
)

// PayloadBuilder is used to build payloads on the
// execution client.
type PayloadBuilder struct {
	// cfg holds the configuration settings for the PayloadBuilder.
	cfg *Config
	// chainSpec holds the chain specifications for the PayloadBuilder.
	chainSpec ChainSpec
	// logger is used for logging within the PayloadBuilder.
	logger log.Logger
	// suggestedFeeRecipient is the suggested fee recipient sent to
	// the execution client for the payload build.
	suggestedFeeRecipient common.ExecutionAddress
	// ee is the execution engine.
	ee ExecutionEngine
	// pc is the payload ID cache, it is used to store
	// "in-flight" payloads that are being built on
	// the execution client.
	pc PayloadCache
}

// New creates a new service.
func New(
	cfg *Config,
	chainSpec ChainSpec,
	logger log.Logger,
	ee ExecutionEngine,
	pc PayloadCache,
) *PayloadBuilder {
	return &PayloadBuilder{
		cfg:                   cfg,
		chainSpec:             chainSpec,
		logger:                logger,
		suggestedFeeRecipient: cfg.SuggestedFeeRecipient,
		ee:                    ee,
		pc:                    pc,
	}
}

// Enabled returns true if the payload builder is enabled.
func (pb *PayloadBuilder) Enabled() bool {
	return pb.cfg.Enabled
}
