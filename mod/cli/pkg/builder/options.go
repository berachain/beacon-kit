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
	"cosmossdk.io/core/transaction"
	"cosmossdk.io/depinject"
	serverv2 "cosmossdk.io/server/v2"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/spf13/cobra"
)

// Opt is a type that defines a function that modifies CLIBuilder.
type Opt[
	
	NodeT types.Node[T], T transaction.Tx any,
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
	ValidatorUpdateT any,
] func(*CLIBuilder[NodeT, T, ExecutionPayloadT, ValidatorUpdateT])

// WithName sets the name for the CLIBuilder.
func WithName[
	
	NodeT types.Node[T], T transaction.Tx any,
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
	ValidatorUpdateT any,
](name string) Opt[NodeT, T, ExecutionPayloadT, ValidatorUpdateT] {
	return func(cb *CLIBuilder[NodeT, T, ExecutionPayloadT, ValidatorUpdateT]) {
		cb.name = name
	}
}

// WithDescription sets the description for the CLIBuilder.
func WithDescription[
	
	NodeT types.Node[T], T transaction.Tx any,
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
	ValidatorUpdateT any,
](description string) Opt[NodeT, T, ExecutionPayloadT, ValidatorUpdateT] {
	return func(cb *CLIBuilder[NodeT, T, ExecutionPayloadT, ValidatorUpdateT]) {
		cb.description = description
	}
}

// WithDepInjectConfig sets the depinject config for the CLIBuilder.
func WithDepInjectConfig[
	
	NodeT types.Node[T], T transaction.Tx any,
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
	ValidatorUpdateT any,
](
	cfg depinject.Config) Opt[NodeT, T, ExecutionPayloadT, ValidatorUpdateT] {
	return func(cb *CLIBuilder[NodeT, T, ExecutionPayloadT, ValidatorUpdateT]) {
		cb.depInjectCfg = cfg
	}
}

// WithComponents sets the components for the CLIBuilder.
func WithComponents[
	
	NodeT types.Node[T], T transaction.Tx any,
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
	ValidatorUpdateT any,
](components ...any) Opt[NodeT, T, ExecutionPayloadT, ValidatorUpdateT] {
	return func(cb *CLIBuilder[NodeT, T, ExecutionPayloadT, ValidatorUpdateT]) {
		cb.components = components
	}
}

// SupplyModuleDeps populates the slice of direct module dependencies to be
// supplied to depinject.
func SupplyModuleDeps[
	
	NodeT types.Node[T], T transaction.Tx any,
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
	ValidatorUpdateT any,
](deps []any) Opt[NodeT, T, ExecutionPayloadT, ValidatorUpdateT] {
	return func(cb *CLIBuilder[NodeT, T, ExecutionPayloadT, ValidatorUpdateT]) {
		cb.suppliers = append(cb.suppliers, deps...)
	}
}

// WithRunHandler sets the run handler for the CLIBuilder.
func WithRunHandler[
	NodeT types.Node[T], T transaction.Tx any,
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
	ValidatorUpdateT any,
](
	runHandler func(cmd *cobra.Command,
		customAppConfigTemplate string,
		customAppConfig interface{},
		cmtConfig *cmtcfg.Config,
	) error,
) Opt[NodeT, T, ExecutionPayloadT, ValidatorUpdateT] {
	return func(cb *CLIBuilder[NodeT, T, ExecutionPayloadT, ValidatorUpdateT]) {
		cb.runHandler = runHandler
	}
}

// WithNodeBuilderFunc sets the cosmos app creator for the CLIBuilder.
func WithNodeBuilderFunc[
	NodeT types.Node[T], T transaction.Tx any,
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
	ValidatorUpdateT any,
](
	nodeBuilderFunc serverv2.AppCreator[NodeT, T],
) Opt[NodeT, T, ExecutionPayloadT, ValidatorUpdateT] {
	return func(cb *CLIBuilder[NodeT, T, ExecutionPayloadT, ValidatorUpdateT]) {
		cb.nodeBuilderFunc = nodeBuilderFunc
	}
}

// WithServer sets the server for the CLIBuilder.
func WithServer[NodeT types.Node[T], T transaction.Tx, ValidatorUpdateT any](
	server *serverv2.Server[NodeT, T],
) Opt[NodeT, T, ValidatorUpdateT] {
	return func(cb *CLIBuilder[NodeT, T, ValidatorUpdateT]) {
		cb.server = server
	}
}
