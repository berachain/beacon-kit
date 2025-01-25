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
	"fmt"
	"github.com/berachain/beacon-kit/cli/flags"
	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"testing"

	cmtcfg "github.com/cometbft/cometbft/config"

	// Cosmos server types (for AppOptions)

	// Your logger
	"github.com/berachain/beacon-kit/log/phuslu"

	// Node-core references
	nodebuilder "github.com/berachain/beacon-kit/node-core/builder"
	nodetypes "github.com/berachain/beacon-kit/node-core/types"
	// If your ABCI service type is in a different package, import it:
	// "github.com/berachain/your-repo/cometbft"
)

// TestABCIFlow shows how to instantiate the Node in-memory, retrieve
// the CometBFTService, and call its ABCI methods.
func TestABCIFlow(t *testing.T) {
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
	db := dbm.NewMemDB()
	cmtCfg := cmtcfg.DefaultConfig()

	// Step 1: Create a pointer to viper.Viper:
	appOpts := viper.New()

	appOpts.Set(flags.JWTSecretPath, "../files/jwt.hex")
	appOpts.Set(flags.RPCDialURL, "http://localhost:8551")
	appOpts.Set(flags.PrivValidatorKeyFile, "./test_priv_validator_key.json")
	appOpts.Set(flags.PrivValidatorStateFile, "./test_priv_validator_state.json")

	//appOpts.Set
	// 2. Build the node in-memory.
	fmt.Println("REZ: Reached A")

	node := nb.Build(
		logger,
		db,
		io.Discard, // or some other writer
		cmtCfg,
		appOpts,
	)
	fmt.Println("REZ: Reached B")
	// 3. Extract your CometBFTService from the node.
	//    This depends on how your Node implements GetService(...).
	var cometService *cometbft.Service[*phuslu.Logger]
	err := node.FetchService(cometService)
	require.NoError(t, err)
	require.NotNil(t, cometService)

	//ctx, cancelFunc := context.WithCancel(context.Background())
	//defer cancelFunc()
	//
	//request := &comettypes.InitChainRequest{
	//	ChainId: "80090",
	//}
	//_, err = cometService.InitChain(ctx, request)
	//require.NoError(t, err)

	//// Cast it to the correct concrete type:
	//service, ok := cmtService.(*Service[*phuslu.Logger]) // Adjust package/type as needed
	//require.True(t, ok, "CometBFTService is not the expected type")
	//
	//// 4. Now call the ABCI methods in memory:
	//initRes, err := service.InitChain(context.Background(), &cmtabci.InitChainRequest{
	//	ChainId: "test-chain",
	//	// Fill in any other fields you want
	//})
	//require.NoError(t, err)
	//t.Logf("InitChain response: %+v", initRes)
	//
	//prepResp, err := service.PrepareProposal(context.Background(), &cmtabci.PrepareProposalRequest{
	//	Height: 1,
	//	// ...
	//})
	//require.NoError(t, err)
	//t.Logf("PrepareProposal response: %+v", prepResp)
	//
	//procResp, err := service.ProcessProposal(context.Background(), &cmtabci.ProcessProposalRequest{
	//	// ...
	//})
	//require.NoError(t, err)
	//t.Logf("ProcessProposal response: %+v", procResp)
	//
	//finResp, err := service.FinalizeBlock(context.Background(), &cmtabci.FinalizeBlockRequest{
	//	Height: 1,
	//	// ...
	//})
	//require.NoError(t, err)
	//t.Logf("FinalizeBlock response: %+v", finResp)

	// 5. Optionally test Commit or other ABCI methods:
	// commitResp, err := service.Commit(context.Background(), &cmtabci.CommitRequest{})
	// require.NoError(t, err)
	// t.Logf("Commit response: %+v", commitResp)

	// Add any assertions about the resulting app state, logs, etc.
}
