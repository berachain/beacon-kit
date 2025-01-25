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

package mock_consensus_test

import (
	"io"
	"os"
	"testing"

	"github.com/berachain/beacon-kit/cli/flags"
	"github.com/berachain/beacon-kit/config"
	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	"github.com/berachain/beacon-kit/da/kzg"
	"github.com/berachain/beacon-kit/log/phuslu"
	"github.com/berachain/beacon-kit/node-api/engines/echo"
	nodebuilder "github.com/berachain/beacon-kit/node-core/builder"
	"github.com/berachain/beacon-kit/node-core/components"
	nodetypes "github.com/berachain/beacon-kit/node-core/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

// DefaultComponents requires testing.T to avoid accidental misuse
func DefaultComponents(_ *testing.T) []any {
	c := []any{
		components.ProvideAttributesFactory[*phuslu.Logger],
		components.ProvideAvailabilityStore[*phuslu.Logger],
		components.ProvideDepositContract,
		components.ProvideBlockStore[*phuslu.Logger],
		components.ProvideBlsSigner,
		components.ProvideBlobProcessor[*phuslu.Logger],
		components.ProvideBlobProofVerifier,
		components.ProvideChainService[*phuslu.Logger],
		components.ProvideNode,
		components.ProvideChainSpec,
		components.ProvideConfig,
		components.ProvideServerConfig,
		// Using in-memory Deposit Store
		components.ProvideDepositStoreInMemory[*phuslu.Logger],
		components.ProvideEngineClient[*phuslu.Logger],
		components.ProvideExecutionEngine[*phuslu.Logger],
		components.ProvideJWTSecret,
		components.ProvideLocalBuilder[*phuslu.Logger],
		components.ProvideReportingService[*phuslu.Logger],
		components.ProvideCometBFTService[*phuslu.Logger],
		components.ProvideServiceRegistry[*phuslu.Logger],
		components.ProvideSidecarFactory,
		components.ProvideStateProcessor[*phuslu.Logger],
		components.ProvideKVStore,
		components.ProvideStorageBackend,
		components.ProvideTelemetrySink,
		components.ProvideTelemetryService,
		components.ProvideTrustedSetup,
		components.ProvideValidatorService[*phuslu.Logger],
		components.ProvideShutDownService[*phuslu.Logger],
	}
	c = append(c,
		components.ProvideNodeAPIServer[*phuslu.Logger, echo.Context],
		components.ProvideNodeAPIEngine,
		components.ProvideNodeAPIBackend,
	)
	//
	c = append(c, components.ProvideNodeAPIHandlers[echo.Context],
		components.ProvideNodeAPIBeaconHandler[echo.Context],
		components.ProvideNodeAPIBuilderHandler[echo.Context],
		components.ProvideNodeAPIConfigHandler[echo.Context],
		components.ProvideNodeAPIDebugHandler[echo.Context],
		components.ProvideNodeAPIEventsHandler[echo.Context],
		components.ProvideNodeAPINodeHandler[echo.Context],
		components.ProvideNodeAPIProofHandler[echo.Context],
	)

	return c
}

// newTestNode creates a node that is similar to production except keeps storage in-memory
// and does not start the CometBFT loop, allowing us to inject our own calls
func newTestNode(t *testing.T) (nodetypes.Node, *cometbft.Service[*phuslu.Logger]) {
	// 1. Build a node builder with your default or custom test components.
	nb := nodebuilder.New(
		nodebuilder.WithComponents[
			nodetypes.Node,
			*phuslu.Logger,
			*phuslu.Config,
		](DefaultComponents(t)),
	)

	// Create minimal parameters to pass into Build.
	logger := phuslu.NewLogger(os.Stdout, nil)

	// Use an in-memory DB
	db := dbm.NewMemDB()
	cmtCfg := cometbft.DefaultConfig()
	beaconCfg := config.DefaultConfig()

	appOpts := viper.New()

	appOpts.Set(flags.JWTSecretPath, "../files/jwt.hex")
	appOpts.Set(flags.RPCDialURL, "http://localhost:8551")
	appOpts.Set(flags.PrivValidatorKeyFile, "./test_priv_validator_key.json")
	appOpts.Set(flags.PrivValidatorStateFile, "./test_priv_validator_state.json")

	appOpts.Set(flags.BlockStoreServiceAvailabilityWindow, beaconCfg.BlockStoreService.AvailabilityWindow)
	appOpts.Set(flags.BlockStoreServiceEnabled, beaconCfg.BlockStoreService.Enabled)
	appOpts.Set(flags.KZGTrustedSetupPath, "../files/kzg-trusted-setup.json")
	appOpts.Set(flags.KZGImplementation, kzg.DefaultConfig().Implementation)

	// TODO: Cleanup this Set
	appOpts.Set("pruning", "default")

	node := nb.Build(
		logger,
		db,
		io.Discard, // or some other writer
		cmtCfg,
		appOpts,
	)

	var cometService *cometbft.Service[*phuslu.Logger]
	err := node.FetchService(&cometService)
	require.NoError(t, err)
	require.NotNil(t, cometService)
	return node, cometService
}
