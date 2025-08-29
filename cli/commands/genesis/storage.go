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
	stdbytes "bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"path/filepath"

	"github.com/berachain/beacon-kit/cli/commands/genesis/types"
	clitypes "github.com/berachain/beacon-kit/cli/commands/server/types"
	"github.com/berachain/beacon-kit/cli/context"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	gethprimitives "github.com/berachain/beacon-kit/geth-primitives"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/crypto"
	cmtcfg "github.com/cometbft/cometbft/config"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/sha3"
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
			// Get the deposits from the beacon chain genesis appstate.
			config := context.GetConfigFromCmd(cmd)
			appOpts := context.GetViperFromCmd(cmd)
			chainSpec, err := chainSpecCreator(appOpts)
			if err != nil {
				return err
			}
			return SetDepositStorage(chainSpec, config, elGenesisFilePath)
		},
	}

	return cmd
}

func SetDepositStorage(
	chainSpec ChainSpec,
	config *cmtcfg.Config,
	elGenesisFilePath string,
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

	// Unmarshal the genesis file.
	elGenesis := &types.DefaultEthGenesisJSON{}
	allocsKey := types.DefaultAllocsKey
	if err = json.Unmarshal(elGenesisBz, elGenesis); err != nil {
		return errors.Wrap(err, "failed to unmarshal eth1 genesis")
	}

	depositAddr := common.Address(chainSpec.DepositContractAddress())
	allocs := writeDepositStorage(elGenesis, deposits, depositAddr)

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
	deposits ctypes.Deposits,
	depositAddr common.Address,
) gethprimitives.GenesisAlloc {
	slot0 := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
	slot1 := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001")

	depositsCount := big.NewInt(int64(len(deposits)))
	depositsRoot := deposits.HashTreeRoot()

	allocs := elGenesis.Alloc()
	if entry, ok := allocs[depositAddr]; ok {
		if entry.Storage == nil {
			entry.Storage = make(map[common.Hash]common.Hash)
		}
		entry.Storage[slot0] = common.BigToHash(depositsCount)
		entry.Storage[slot1] = common.BytesToHash(depositsRoot[:])
		allocs[depositAddr] = entry

		// Store operators keys for each validator, reusing their BLS key
		// TODO: this is good enough for testing over devnets, but we may
		// want to extend the command to be able to explicitly pass a list
		// of pre-arranged operator keys.
		for i, d := range deposits {
			storageKey := encodeSlot(d.Pubkey)
			operatorAddr, err := crypto.GetAddressFromPubKey(d.Pubkey) // reuse val BLS key for simplicity
			if err != nil {
				panic(fmt.Errorf("failed getting address from validator %d pub key: %w", i, err))
			}

			k := common.BytesToHash(storageKey)
			v := common.BytesToHash(operatorAddr)
			entry.Storage[k] = v
		}
	}
	return allocs
}

// encodeSlot mimics Solidity's keccak256(abi.encodePacked(...)) for:
// - pubKey: 48-byte public key
// - baseSlot: 32-byte storage slot
func encodeSlot(pubkey crypto.BLSPubkey) []byte {
	// Decode pubkey
	pubKeyStr := pubkey.String()
	packed, err := hex.DecodeString(pubKeyStr[2:])
	if err != nil {
		panic(err)
	}
	if len(packed) != bytes.B48Size {
		panic(fmt.Errorf("expected 48-byte pubkey, got %d", len(packed)))
	}

	// Convert mapping slot (uint256) to 32-byte left-padded value
	const mappintStorageBaseSlot = 2
	slotBytes := new(big.Int).SetUint64(mappintStorageBaseSlot).FillBytes(make([]byte, common.HashLength))

	// abi.encodePacked => direct byte concatenation
	packed = append(packed, slotBytes...)

	// keccak256 hash
	hash := sha3.NewLegacyKeccak256()
	hash.Write(packed)
	return hash.Sum(nil)
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

	// Unmarshal existing genesis using json.Number to preserve integer precision
	var existingGenesis map[string]interface{}
	decoder := json.NewDecoder(stdbytes.NewReader(existingBz))
	decoder.UseNumber()
	if err = decoder.Decode(&existingGenesis); err != nil {
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
