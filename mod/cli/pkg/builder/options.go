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
	"github.com/berachain/beacon-kit/mod/node-core/pkg/types"
	cmtcfg "github.com/cometbft/cometbft/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cobra"
)

// Opt is a type that defines a function that modifies CLIBuilder.
type Opt[T types.Node] func(*CLIBuilder[T])

// WithName sets the name for the CLIBuilder.
func WithName[T types.Node](name string) Opt[T] {
	return func(cb *CLIBuilder[T]) {
		cb.name = name
	}
}

// WithDescription sets the description for the CLIBuilder.
func WithDescription[T types.Node](description string) Opt[T] {
	return func(cb *CLIBuilder[T]) {
		cb.description = description
	}
}

// WithDepInjectConfig sets the depinject config for the CLIBuilder.
func WithDepInjectConfig[T types.Node](
	cfg depinject.Config) Opt[T] {
	return func(cb *CLIBuilder[T]) {
		cb.depInjectCfg = cfg
	}
}

// WithComponents sets the components for the CLIBuilder.
func WithComponents[T types.Node](components []any) Opt[T] {
	return func(cb *CLIBuilder[T]) {
		cb.components = components
	}
}

// SupplyModuleDeps populates the slice of direct module dependencies to be
// supplied to depinject.
func SupplyModuleDeps[T types.Node](deps []any) Opt[T] {
	return func(cb *CLIBuilder[T]) {
		cb.suppliers = append(cb.suppliers, deps...)
	}
}

// WithRunHandler sets the run handler for the CLIBuilder.
func WithRunHandler[T types.Node](
	runHandler func(cmd *cobra.Command,
		customAppConfigTemplate string,
		customAppConfig interface{},
		cmtConfig *cmtcfg.Config,
	) error,
) Opt[T] {
	return func(cb *CLIBuilder[T]) {
		cb.runHandler = runHandler
	}
}

// WithRootCommandSetup sets the root command setup func to the default.
func WithRootCommandSetup[T types.Node](rootCmdSetup rootCmdSetup[T]) Opt[T] {
	return func(cb *CLIBuilder[T]) {
		cb.rootCmdSetup = rootCmdSetup
	}
}

// WithNodeBuilderFunc sets the cosmos app creator for the CLIBuilder.
func WithNodeBuilderFunc[T types.Node](
	nodeBuilderFunc servertypes.AppCreator[T],
) Opt[T] {
	return func(cb *CLIBuilder[T]) {
		cb.nodeBuilderFunc = nodeBuilderFunc
	}
}
