// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
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

package nodebuilder

import (
	"os"

	"cosmossdk.io/client/v2/autocli"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	consensustypes "github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/da/pkg/kzg/noop"
	dastore "github.com/berachain/beacon-kit/mod/da/pkg/store"
	engineclient "github.com/berachain/beacon-kit/mod/execution/pkg/client"
	cmdlib "github.com/berachain/beacon-kit/mod/node-core/pkg/commands"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/signer"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/node"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	depositdb "github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/types/module"
	gokzg4844 "github.com/crate-crypto/go-kzg-4844"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type NodeBuilder[NodeT types.NodeI] struct {
	node NodeT

	name         string
	description  string
	depInjectCfg depinject.Config
	chainSpec    primitives.ChainSpec
}

// New returns a new NodeBuilder.
func New[NodeT types.NodeI](opts ...Opt[NodeT]) *NodeBuilder[NodeT] {
	nb := &NodeBuilder[NodeT]{
		node: node.New[NodeT](),
	}
	for _, opt := range opts {
		opt(nb)
	}
	return nb
}

// Build builds the application.
func (nb *NodeBuilder[NodeT]) Build() (NodeT, error) {
	rootCmd, err := nb.buildRootCmd()
	if err != nil {
		return nb.node, err
	}

	nb.node.SetRootCmd(rootCmd)
	return nb.node, nil
}

// buildRootCmd builds the root command for the application.
func (nb *NodeBuilder[NodeT]) buildRootCmd() (*cobra.Command, error) {
	var (
		autoCliOpts autocli.AppOptions
		mm          *module.Manager
		clientCtx   client.Context
	)
	if err := depinject.Inject(
		depinject.Configs(
			nb.depInjectCfg,
			depinject.Supply(
				log.NewLogger(os.Stdout),
				viper.GetViper(),
				nb.chainSpec,
				&depositdb.KVStore[*consensustypes.Deposit]{},
				&engineclient.EngineClient[*consensustypes.ExecutionPayload]{},
				&gokzg4844.JSONTrustedSetup{},
				&noop.Verifier{},
				&dastore.Store[consensustypes.BeaconBlockBody]{},
				&signer.BLSSigner{},
			),
			depinject.Provide(
				components.ProvideNoopTxConfig,
				components.ProvideClientContext,
				components.ProvideKeyring,
				components.ProvideConfig,
				components.ProvideTelemetrySink,
			),
		),
		&autoCliOpts,
		&mm,
		&clientCtx,
	); err != nil {
		return nil, err
	}

	cmd := &cobra.Command{
		Use:   nb.name,
		Short: nb.description,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// set the default command outputs
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			var err error
			clientCtx, err = client.ReadPersistentCommandFlags(
				clientCtx,
				cmd.Flags(),
			)
			if err != nil {
				return err
			}

			customClientTemplate, customClientConfig := components.InitClientConfig()
			clientCtx, err = config.CreateClientConfig(
				clientCtx,
				customClientTemplate,
				customClientConfig,
			)
			if err != nil {
				return err
			}

			if err = client.SetCmdClientContextHandler(
				clientCtx, cmd,
			); err != nil {
				return err
			}

			return server.InterceptConfigsPreRunHandler(
				cmd,
				DefaultAppConfigTemplate(),
				DefaultAppConfig(),
				DefaultCometConfig(),
			)
		},
	}

	cmdlib.DefaultRootCommandSetup(
		cmd,
		mm,
		nb.AppCreator,
		nb.chainSpec,
	)

	if err := autoCliOpts.EnhanceRootCommand(cmd); err != nil {
		return nil, err
	}

	return cmd, nil
}
