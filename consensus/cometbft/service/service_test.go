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

package cometbft_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/berachain/beacon-kit/beacon/blockchain/mocks"
	vmocks "github.com/berachain/beacon-kit/beacon/validator/mocks"
	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	"github.com/berachain/beacon-kit/log/phuslu"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/storage/db"
	cmtcfg "github.com/cometbft/cometbft/config"
	dbm "github.com/cosmos/cosmos-db"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/stretchr/testify/require"
)

func TestInitGenesisState(t *testing.T) {
	t.Parallel()

	logBuffer := new(bytes.Buffer)
	logger := phuslu.NewLogger(logBuffer, nil)

	appDB, err := db.OpenDB("app", dbm.MemDBBackend)
	require.NoError(t, err)

	var (
		chain      = mocks.NewBlockchainI(t)
		telem      = mocks.NewTelemetrySink(t)
		blkBuilder = vmocks.NewBlockBuilderI(t)
	)

	cmtCfg := setupTestConfig(t)

	dummyOpts := []func(*cometbft.Service){}

	var s *cometbft.Service
	require.NotPanics(t, func() {
		s = cometbft.NewService(logger, appDB, chain, blkBuilder, cmtCfg, telem, dummyOpts...)
	})

	// Can't mock consensus, so do not call Start nor Close
	// TODO: add dependency injection to mock consensus calls

	_, err = s.AppVersion(context.Background())
	require.NoError(t, err)
}

// TODO: duplicated from node-api/backend/validator_test.go. TO CONSOLIDATE
func setupTestConfig(t *testing.T) *cmtcfg.Config {
	t.Helper()

	// Create a temporary directory for CometBFT config
	tmpDir := t.TempDir()

	// Create CometBFT config with temporary directory
	cmtCfg := cometbft.DefaultConfig()
	cmtCfg.SetRoot(tmpDir)

	// Create config directory
	configDir := filepath.Join(tmpDir, "config")
	err := os.MkdirAll(configDir, 0o755)
	require.NoError(t, err)

	// Create app genesis
	appGenesis := genutiltypes.NewAppGenesisWithVersion("test-chain", []byte("{}"))
	appGenesis.Consensus.Params = cometbft.DefaultConsensusParams(crypto.CometBLSType)

	// Save genesis file to the config directory
	genesisFile := filepath.Join(configDir, "genesis.json")
	err = appGenesis.SaveAs(genesisFile)
	require.NoError(t, err)

	return cmtCfg
}
