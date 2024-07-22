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
	"cosmossdk.io/depinject"
	cmdlib "github.com/berachain/beacon-kit/mod/cli/pkg/commands"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	cmtcfg "github.com/cometbft/cometbft/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cobra"
)

// Opt is a type that defines a function that modifies CLIBuilder.
type Opt[
	T types.Node,
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
] func(*CLIBuilder[T, ExecutionPayloadT])

// WithName sets the name for the CLIBuilder.
func WithName[
	T types.Node,
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
](name string) Opt[T, ExecutionPayloadT] {
	return func(cb *CLIBuilder[T, ExecutionPayloadT]) {
		cb.name = name
	}
}

// WithDescription sets the description for the CLIBuilder.
func WithDescription[
	T types.Node,
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
](description string) Opt[T, ExecutionPayloadT] {
	return func(cb *CLIBuilder[T, ExecutionPayloadT]) {
		cb.description = description
	}
}

// WithDepInjectConfig sets the depinject config for the CLIBuilder.
func WithDepInjectConfig[
	T types.Node,
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
](
	cfg depinject.Config) Opt[T, ExecutionPayloadT] {
	return func(cb *CLIBuilder[T, ExecutionPayloadT]) {
		cb.depInjectCfg = cfg
	}
}

// WithComponents sets the components for the CLIBuilder.
func WithComponents[
	T types.Node,
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
](components []any) Opt[T, ExecutionPayloadT] {
	return func(cb *CLIBuilder[T, ExecutionPayloadT]) {
		cb.components = components
	}
}

// SupplyModuleDeps populates the slice of direct module dependencies to be
// supplied to depinject.
func SupplyModuleDeps[
	T types.Node,
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
](deps []any) Opt[T, ExecutionPayloadT] {
	return func(cb *CLIBuilder[T, ExecutionPayloadT]) {
		cb.suppliers = append(cb.suppliers, deps...)
	}
}

// WithRunHandler sets the run handler for the CLIBuilder.
func WithRunHandler[
	T types.Node,
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
](
	runHandler func(cmd *cobra.Command,
		customAppConfigTemplate string,
		customAppConfig interface{},
		cmtConfig *cmtcfg.Config,
	) error,
) Opt[T, ExecutionPayloadT] {
	return func(cb *CLIBuilder[T, ExecutionPayloadT]) {
		cb.runHandler = runHandler
	}
}

// WithDefaultRootCommandSetup sets the root command setup func to the default.
func WithDefaultRootCommandSetup[
	T types.Node,
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
]() Opt[T, ExecutionPayloadT] {
	return func(cb *CLIBuilder[T, ExecutionPayloadT]) {
		cb.rootCmdSetup = cmdlib.DefaultRootCommandSetup[T, ExecutionPayloadT]
	}
}

// WithNodeBuilderFunc sets the cosmos app creator for the CLIBuilder.
func WithNodeBuilderFunc[
	T types.Node,
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
](
	nodeBuilderFunc servertypes.AppCreator[T],
) Opt[T, ExecutionPayloadT] {
	return func(cb *CLIBuilder[T, ExecutionPayloadT]) {
		cb.nodeBuilderFunc = nodeBuilderFunc
	}
}
