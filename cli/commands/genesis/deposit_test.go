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

package genesis_test

import (
	"path"
	"testing"

	"github.com/berachain/beacon-kit/cli/commands/genesis"
	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/node-core/components/signer"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cometbft/cometbft/crypto/bls12381"
	"github.com/cometbft/cometbft/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
)

func TestGenesisDeposit(t *testing.T) {
	t.Parallel()
	homeDir := t.TempDir()
	t.Log("Home folder:", homeDir)

	chainSpec, err := spec.MainnetChainSpec()
	require.NoError(t, err)

	cometConfig := cmtcfg.DefaultConfig()

	cometConfig.SetRoot(homeDir)

	// Forces Comet to Create it
	cometConfig.NodeKey = "nodekey.json"

	depositAmount := math.Gwei(250_000 * params.GWei)
	withdrawalAdress := common.NewExecutionAddressFromHex("0x981114102592310C347E61368342DDA67017bf84")
	outputDocument := ""

	blsSigner := signer.BLSSigner{PrivValidator: types.NewMockPVWithKeyType(bls12381.KeyType)}

	err = genesis.AddGenesisDeposit(chainSpec, cometConfig, blsSigner, depositAmount, withdrawalAdress, outputDocument)
	require.NoError(t, err)

	require.FileExists(t, path.Join(homeDir, "nodekey.json"))
	require.FileExists(t, path.Join(homeDir, "data", "priv_validator_state.json"))
	require.FileExists(t, path.Join(homeDir, "config", "priv_validator_key.json"))
	require.DirExists(t, path.Join(homeDir, "config", "premined-deposits"))
	// TODO: Extend tests to assert on the contents of the files
}
