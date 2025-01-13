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

package genesis

import (
	"math/big"
	"path/filepath"

	"github.com/berachain/beacon-kit/cli/commands/genesis/types"
	"github.com/berachain/beacon-kit/cli/context"
	"github.com/berachain/beacon-kit/config/spec"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	gethprimitives "github.com/berachain/beacon-kit/geth-primitives"
	"github.com/berachain/beacon-kit/primitives/encoding/json"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// Set deposit contract storage in genesis alloc file.
func SetDepositStorageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-deposit-storage [eth/genesis/file.json]",
		Short: "sets deposit contract storage in eth genesis",
		Long: `Updates the deposit contract storage in the passed in eth genesis file. 
		Creates a new EL genesis file with the changes in the BEACOND_HOME directory.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get the deposits from the beacon chain genesis appstate.
			config := context.GetConfigFromCmd(cmd)

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

			// Read the EL genesis file.
			elGenesisBz, err := afero.ReadFile(afero.NewOsFs(), args[0])
			if err != nil {
				return errors.Wrap(err, "failed to read eth1 genesis file")
			}

			isNethermind, err := cmd.Flags().GetBool(nethermindGenesis)
			if err != nil {
				return err
			}
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

			// Create a map to store the storage of the deposit contract.
			storage := make(map[common.Hash]common.Hash)
			slot0 := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
			slot1 := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001")
			storage[slot0] = common.BigToHash(count)
			storage[slot1] = common.BytesToHash(root[:])

			// Assign storage to the deposit contract address.
			addr := common.HexToAddress(spec.DefaultDepositContractAddress)
			allocs := elGenesis.Alloc()
			if entry, ok := allocs[addr]; ok {
				entry.Storage = storage
				allocs[addr] = entry
			}

			// Get just the filename from the path
			filename := filepath.Base(args[0])
			outputPath := filepath.Join(config.RootDir, filename)

			// Write to file.
			err = writeGenesisAllocToFile(outputPath, args[0], allocs, allocsKey)
			if err != nil {
				return errors.Wrap(err, "failed to write genesis alloc to file")
			}

			return nil
		},
	}

	cmd.Flags().BoolP(
		nethermindGenesis, nethermindGenesisShorthand,
		nethermindGenesisDefault, nethermindGenesisMsg,
	)
	return cmd
}
func writeGenesisAllocToFile(
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
	depositAddr := spec.DefaultDepositContractAddress
	if account, exists := genesisAlloc[common.HexToAddress(depositAddr)]; exists {
		alloc[depositAddr] = account
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
