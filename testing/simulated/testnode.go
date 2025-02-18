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

package simulated

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/beacon/blockchain"
	"github.com/berachain/beacon-kit/chain"
	servertypes "github.com/berachain/beacon-kit/cli/commands/server/types"
	"github.com/berachain/beacon-kit/cli/flags"
	"github.com/berachain/beacon-kit/config"
	"github.com/berachain/beacon-kit/da/kzg"
	"github.com/berachain/beacon-kit/log/phuslu"
	nodecomponents "github.com/berachain/beacon-kit/node-core/components"
	nodetypes "github.com/berachain/beacon-kit/node-core/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/net/url"
	"github.com/berachain/beacon-kit/storage/db"
	cmtcfg "github.com/cometbft/cometbft/config"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

// TestNodeInput takes the input for building and starting a node
type TestNodeInput struct {
	TempHomeDir string
	CometConfig *cmtcfg.Config
	AuthRPC     *url.ConnectionURL
	Logger      *phuslu.Logger
	AppOpts     *viper.Viper
	Components  []any
}

type TestNode struct {
	nodetypes.Node
	StorageBackend blockchain.StorageBackend
	ChainSpec      chain.Spec
	APIBackend     nodecomponents.NodeAPIBackend
}

// NewTestNode Uses the testnet chainspec.
func NewTestNode(
	t *testing.T,
	input TestNodeInput,
) TestNode {
	t.Helper()

	beaconKitConfig := createBeaconKitConfig(t)
	beaconKitConfig.Engine.RPCDialURL = input.AuthRPC
	appOpts := getAppOptions(t, input.AppOpts, beaconKitConfig, input.TempHomeDir)

	// Create a database
	database, err := db.OpenDB(input.TempHomeDir, dbm.PebbleDBBackend)
	require.NoError(t, err)

	// Build a node
	node := buildNode(
		input.Logger,
		database,
		os.Stdout, // or some other writer
		input.CometConfig,
		appOpts,
		input.Components,
	)
	return node
}

// buildNode run the same logic as primary build, but it returns the components allowing us to query them.
func buildNode(
	logger *phuslu.Logger,
	db dbm.DB,
	_ io.Writer,
	cmtCfg *cmtcfg.Config,
	appOpts servertypes.AppOptions,
	components []any,
) TestNode {
	// variables to hold the components needed to set up BeaconApp
	var (
		apiBackend     nodecomponents.NodeAPIBackend
		beaconNode     nodetypes.Node
		cmtService     nodetypes.ConsensusService
		config         *config.Config
		storageBackend blockchain.StorageBackend
		chainSpec      chain.Spec
	)

	// build all node components using depinject
	if err := depinject.Inject(
		depinject.Configs(
			depinject.Provide(
				components...,
			),
			depinject.Supply(
				appOpts,
				logger,
				db,
				cmtCfg,
			),
		),
		&apiBackend,
		&beaconNode,
		&cmtService,
		&config,
		&storageBackend,
		&chainSpec,
	); err != nil {
		panic(err)
	}
	if config == nil {
		panic("config is nil")
	}
	if apiBackend == nil {
		panic("node or api backend is nil")
	}

	logger.WithConfig(config.GetLogger())
	apiBackend.AttachQueryBackend(cmtService)
	return TestNode{
		Node:           beaconNode,
		StorageBackend: storageBackend,
		ChainSpec:      chainSpec,
		APIBackend:     apiBackend,
	}
}

// getAppOptions returns the Application Options we need to set for the Node Builder.
// Ideally we can avoid having to set the flags like this and just directly modify a config type.
func getAppOptions(t *testing.T, appOpts *viper.Viper, beaconKitConfig *config.Config, tempHomeDir string) *viper.Viper {
	t.Helper()
	// Execution Client Config
	relativePathJwt := "../files/jwt.hex"
	jwtPath, err := filepath.Abs(relativePathJwt)
	require.NoError(t, err)
	appOpts.Set(flags.JWTSecretPath, jwtPath)
	appOpts.Set(flags.RPCJWTRefreshInterval, beaconKitConfig.GetEngine().RPCJWTRefreshInterval.String())
	appOpts.Set(flags.RPCStartupCheckInterval, beaconKitConfig.GetEngine().RPCStartupCheckInterval.String())
	appOpts.Set(flags.RPCDialURL, beaconKitConfig.GetEngine().RPCDialURL.String())
	appOpts.Set(flags.RPCTimeout, beaconKitConfig.GetEngine().RPCTimeout.String())

	appOpts.Set(flags.LogLevel, "debug")

	// BLS Config
	appOpts.Set(flags.PrivValidatorKeyFile, "./config/priv_validator_key.json")
	appOpts.Set(flags.PrivValidatorStateFile, "./data/priv_validator_state.json")

	// Beacon Config
	appOpts.Set(flags.BlockStoreServiceAvailabilityWindow, beaconKitConfig.GetBlockStoreService().AvailabilityWindow)
	appOpts.Set(flags.BlockStoreServiceEnabled, beaconKitConfig.GetBlockStoreService().Enabled)
	appOpts.Set(flags.KZGTrustedSetupPath, "../files/kzg-trusted-setup.json")
	appOpts.Set(flags.KZGImplementation, kzg.DefaultConfig().Implementation)

	// Payload Builder Config
	beaconKitConfig.GetPayloadBuilder().SuggestedFeeRecipient = common.NewExecutionAddressFromHex("0x981114102592310C347E61368342DDA67017bf84")
	appOpts.Set(flags.BuilderEnabled, beaconKitConfig.GetPayloadBuilder().Enabled)
	appOpts.Set(flags.BuildPayloadTimeout, beaconKitConfig.GetPayloadBuilder().PayloadTimeout)
	appOpts.Set(flags.SuggestedFeeRecipient, beaconKitConfig.GetPayloadBuilder().SuggestedFeeRecipient)

	// TODO: Cleanup this Set
	appOpts.Set("pruning", "default")
	appOpts.Set("home", tempHomeDir)
	return appOpts
}

func createBeaconKitConfig(_ *testing.T) *config.Config {
	return config.DefaultConfig()
}
