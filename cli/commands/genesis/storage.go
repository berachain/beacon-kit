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

package genesis

import (
	"math/big"
	"path/filepath"

	"github.com/berachain/beacon-kit/cli/commands/genesis/types"
	clitypes "github.com/berachain/beacon-kit/cli/commands/server/types"
	"github.com/berachain/beacon-kit/cli/context"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	gethprimitives "github.com/berachain/beacon-kit/geth-primitives"
	libcommon "github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/encoding/json"
	cmtcfg "github.com/cometbft/cometbft/config"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// SetDepositStorageCmd sets deposit contract storage in genesis alloc file.
//
//nolint:lll // reads better if long description is one line
func SetDepositStorageCmd(chainSpecCreator clitypes.ChainSpecCreator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-deposit-storage [eth/genesis/file.json]",
		Short: "sets deposit contract storage in eth genesis",
		Long:  `Updates the deposit contract storage in the passed in eth genesis file. Creates a new EL genesis file with the changes in the BEACOND_HOME directory.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Read the EL genesis file.
			elGenesisFilePath := args[0]
			isNethermind, err := cmd.Flags().GetBool(nethermindGenesis)
			if err != nil {
				return err
			}
			// Get the deposits from the beacon chain genesis appstate.
			config := context.GetConfigFromCmd(cmd)
			appOpts := context.GetViperFromCmd(cmd)
			chainSpec, err := chainSpecCreator(appOpts)
			if err != nil {
				return err
			}
			return SetDepositStorage(chainSpec, config, elGenesisFilePath, isNethermind)
		},
	}

	cmd.Flags().BoolP(
		nethermindGenesis, nethermindGenesisShorthand,
		nethermindGenesisDefault, nethermindGenesisMsg,
	)
	return cmd
}

func SetDepositStorage(
	chainSpec ChainSpec,
	config *cmtcfg.Config,
	elGenesisFilePath string,
	isNethermind bool,
) error {
	elGenesisBz, err := afero.ReadFile(afero.NewOsFs(), elGenesisFilePath)
	if err != nil {
		return errors.Wrap(err, "failed to read eth1 genesis file")
	}

	clGenesis, err := genutiltypes.AppGenesisFromFile(
		config.GenesisFile(),
	)
	if err != nil {
		return errors.Wrap(err, "failed to read genesis doc from file")
	}

	genesisState, err := genutiltypes.GenesisStateFromAppGenesis(
		clGenesis,
	)
	if err != nil {
		return errors.Wrap(err, "failed to read appstate from genesis")
	}

	beaconStateRaw := genesisState["beacon"]
	var beaconState struct {
		Deposits ctypes.Deposits `json:"deposits"`
	}
	if err = json.Unmarshal(beaconStateRaw, &beaconState); err != nil {
		return errors.Wrap(err, "failed to unmarshal beacon state")
	}
	deposits := beaconState.Deposits

	// Set the storage of the deposit contract with deposits count and root.
	count := big.NewInt(int64(len(deposits)))
	root := deposits.HashTreeRoot()

	var allocsKey string

	// Unmarshal the genesis file.
	var elGenesis types.EthGenesis
	if isNethermind {
		elGenesis = &types.NethermindEthGenesisJSON{}
		allocsKey = types.NethermindAllocsKey
	} else {
		elGenesis = &types.DefaultEthGenesisJSON{}
		allocsKey = types.DefaultAllocsKey
	}
	if err = json.Unmarshal(elGenesisBz, elGenesis); err != nil {
		return errors.Wrap(err, "failed to unmarshal eth1 genesis")
	}

	depositAddr := common.Address(chainSpec.DepositContractAddress())
	allocs := writeDepositStorage(elGenesis, depositAddr, count, root)

	// Get just the filename from the path
	filename := filepath.Base(elGenesisFilePath)
	outputPath := filepath.Join(config.RootDir, filename)

	// Write to file.
	err = writeGenesisAllocToFile(depositAddr, outputPath, elGenesisFilePath, allocs, allocsKey)
	if err != nil {
		return errors.Wrap(err, "failed to write genesis alloc to file")
	}
	return nil
}

func writeDepositStorage(
	elGenesis types.EthGenesis,
	depositAddr common.Address,
	depositsCount *big.Int,
	depositsRoot libcommon.Root,
) gethprimitives.GenesisAlloc {
	slot0 := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
	slot1 := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001")

	allocs := elGenesis.Alloc()
	if entry, ok := allocs[depositAddr]; ok {
		if entry.Storage == nil {
			entry.Storage = make(map[common.Hash]common.Hash)
		}
		entry.Storage[slot0] = common.BigToHash(depositsCount)
		entry.Storage[slot1] = common.BytesToHash(depositsRoot[:])
		allocs[depositAddr] = entry
	}
	return allocs
}

func writeGenesisAllocToFile(
	depositAddr common.Address,
	outputDocument string,
	inputDocument string,
	genesisAlloc gethprimitives.GenesisAlloc,
	allocsKey string,
) error {
	// Read existing el genesis file
	existingBz, err := afero.ReadFile(afero.NewOsFs(), inputDocument)
	if err != nil {
		return err
	}

	// Unmarshal existing genesis
	var existingGenesis map[string]interface{}
	if err = json.Unmarshal(existingBz, &existingGenesis); err != nil {
		return err
	}

	// Get existing alloc.
	alloc, ok := existingGenesis[allocsKey].(map[string]interface{})
	if !ok {
		return errors.New("invalid alloc format in genesis file")
	}

	// Update only the deposit contract entry
	if account, exists := genesisAlloc[depositAddr]; exists {
		alloc[depositAddr.Hex()] = account
	}

	// Marshal back to JSON
	bz, err := json.MarshalIndent(existingGenesis, "", "  ")
	if err != nil {
		return err
	}

	// Write back to file
	return afero.WriteFile(
		afero.NewOsFs(),
		outputDocument,
		bz,
		0o644, //nolint:mnd // file permissions.
	)
}
