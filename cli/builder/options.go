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
	servertypes "github.com/berachain/beacon-kit/cli/commands/server/types"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/node-core/types"
	"github.com/berachain/beacon-kit/primitives/constraints"
)

// Opt is a type that defines a function that modifies CLIBuilder.
type Opt[
	T types.Node,
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
	LoggerT log.AdvancedLogger[LoggerT],
	ChainServiceT any,
] func(*CLIBuilder[T, ExecutionPayloadT, LoggerT, ChainServiceT])

// WithName sets the name for the CLIBuilder.
func WithName[
	T types.Node,
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
	LoggerT log.AdvancedLogger[LoggerT],
	ChainServiceT any,
](name string) Opt[T, ExecutionPayloadT, LoggerT, ChainServiceT] {
	return func(cb *CLIBuilder[T, ExecutionPayloadT, LoggerT, ChainServiceT]) {
		cb.name = name
	}
}

// WithDescription sets the description for the CLIBuilder.
func WithDescription[
	T types.Node,
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
	LoggerT log.AdvancedLogger[LoggerT],
	ChainServiceT any,
](description string) Opt[T, ExecutionPayloadT, LoggerT, ChainServiceT] {
	return func(cb *CLIBuilder[T, ExecutionPayloadT, LoggerT, ChainServiceT]) {
		cb.description = description
	}
}

// WithComponents sets the components for the CLIBuilder.
func WithComponents[
	T types.Node,
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
	LoggerT log.AdvancedLogger[LoggerT],
	ChainServiceT any,
](components []any) Opt[T, ExecutionPayloadT, LoggerT, ChainServiceT] {
	return func(cb *CLIBuilder[T, ExecutionPayloadT, LoggerT, ChainServiceT]) {
		cb.components = components
	}
}

// WithNodeBuilderFunc sets the cosmos app creator for the CLIBuilder.
func WithNodeBuilderFunc[
	T types.Node,
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
	LoggerT log.AdvancedLogger[LoggerT],
	ChainServiceT any,
](
	nodeBuilderFunc servertypes.AppCreator[T, LoggerT, ChainServiceT],
) Opt[T, ExecutionPayloadT, LoggerT, ChainServiceT] {
	return func(cb *CLIBuilder[T, ExecutionPayloadT, LoggerT, ChainServiceT]) {
		cb.nodeBuilderFunc = nodeBuilderFunc
	}
}
