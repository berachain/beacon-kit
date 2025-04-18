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

// implementation taken from github.com/cosmos/cosmos-sdk/blob/main/x/genutil/client/cli/init.go
// and modified to circumvent using default cometbft config which sets the timeout_commit to 0
package initialize

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math/unsafe"
	"github.com/berachain/beacon-kit/chain"
	clitypes "github.com/berachain/beacon-kit/cli/commands/server/types"
	"github.com/berachain/beacon-kit/cli/context"
	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/encoding/json"
	cfg "github.com/cometbft/cometbft/config"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/go-bip39"
	"github.com/spf13/cobra"
)

const (
	// FlagOverwrite defines a flag to overwrite an existing genesis JSON file.
	FlagOverwrite = "overwrite"

	// FlagSeed defines a flag to initialize the private validator key from a specific seed.
	FlagRecover = "recover"

	// FlagDefaultBondDenom defines the default denom to use in the genesis file.
	FlagDefaultBondDenom = "default-denom"

	// In BeaconKit we use crypto.CometBLSType only so we don't allow to specify
	// any consensus key.
	consensusKeyAlgo = crypto.CometBLSType
)

type printInfo struct {
	Moniker    string          `json:"moniker"     yaml:"moniker"`
	ChainID    string          `json:"chain_id"    yaml:"chain_id"`
	NodeID     string          `json:"node_id"     yaml:"node_id"`
	GenTxsDir  string          `json:"gentxs_dir"  yaml:"gentxs_dir"`
	AppMessage json.RawMessage `json:"app_message" yaml:"app_message"`
}

func newPrintInfo(moniker, chainID, nodeID, genTxsDir string, appMessage json.RawMessage) printInfo {
	return printInfo{
		Moniker:    moniker,
		ChainID:    chainID,
		NodeID:     nodeID,
		GenTxsDir:  genTxsDir,
		AppMessage: appMessage,
	}
}

func displayInfo(dst io.Writer, info printInfo) error {
	out, err := json.MarshalIndent(info, "", " ")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(dst, "%s\n", out)

	return err
}

//nolint:funlen,gocognit,mnd // based on cosmossdk implementation
func InitCmd(creator clitypes.ChainSpecCreator, mm interface {
	DefaultGenesis(chain.Spec) map[string]json.RawMessage
	ValidateGenesis(genesisData map[string]json.RawMessage) error
}) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init <moniker>",
		Short: "Initialize private validator, p2p, genesis, and application configuration files",
		Long:  `Initialize validators's and node's configuration files.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			// it's very important here to use `GetConfigFromCmd` as redefined
			// in BeaconKit instead of cosmos sdk original implementation. Failure
			// to do so, would result in cosmos sdk default configs being picked
			// instead of BeaconKit ones
			config := context.GetConfigFromCmd(cmd)
			chainID, err := cmd.Flags().GetString(flags.FlagChainID)
			if err != nil {
				return errors.New("failed to parse FlagChainID")
			}

			switch {
			case chainID != "":
			case clientCtx.ChainID != "":
				chainID = clientCtx.ChainID
			default:
				chainID = fmt.Sprintf("test-chain-%v", unsafe.Str(6))
			}
			if config.RootDir == "" {
				config.RootDir = clientCtx.HomeDir
			}

			// Get bip39 mnemonic
			var mnemonic string
			shouldRecover, err := cmd.Flags().GetBool(FlagRecover)
			if err != nil {
				return errors.New("failed to parse FlagRecover")
			}
			if shouldRecover {
				inBuf := bufio.NewReader(cmd.InOrStdin())
				var value string
				value, err = input.GetString("Enter your bip39 mnemonic", inBuf)
				if err != nil {
					return err
				}

				mnemonic = value
				if !bip39.IsMnemonicValid(mnemonic) {
					return errors.New("invalid mnemonic")
				}
			}

			// Get initial height
			initHeight, err := cmd.Flags().GetInt64(flags.FlagInitHeight)
			if err != nil {
				return errors.New("failed to parse FlagInitHeight")
			}
			if initHeight < 1 {
				initHeight = 1
			}

			nodeID, _, err := genutil.InitializeNodeValidatorFilesFromMnemonic(config, mnemonic, consensusKeyAlgo)
			if err != nil {
				return err
			}

			config.Moniker = args[0]

			genFile := config.GenesisFile()
			overwrite, err := cmd.Flags().GetBool(FlagOverwrite)
			if err != nil {
				return errors.New("failed to parse FlagOverwrite")
			}
			defaultDenom, err := cmd.Flags().GetString(FlagDefaultBondDenom)
			if err != nil {
				return errors.New("failed to parse FlagDefaultBondDenom")
			}

			// use os.Stat to check if the file exists
			_, err = os.Stat(genFile)
			if !overwrite && !os.IsNotExist(err) {
				return fmt.Errorf("genesis.json file already exists: %v", genFile)
			}

			// Overwrites the SDK default denom for side-effects
			if defaultDenom != "" {
				sdk.DefaultBondDenom = defaultDenom
			}

			chainSpec, err := creator(context.GetViperFromCmd(cmd))
			if err != nil {
				return fmt.Errorf("faile to create chain spec: %w", err)
			}
			appGenState := mm.DefaultGenesis(chainSpec)

			appState, err := json.MarshalIndent(appGenState, "", " ")
			if err != nil {
				return errorsmod.Wrap(err, "Failed to marshal default genesis state")
			}

			appGenesis := &types.AppGenesis{}
			if _, err = os.Stat(genFile); err != nil {
				if !os.IsNotExist(err) {
					return err
				}
			} else {
				appGenesis, err = types.AppGenesisFromFile(genFile)
				if err != nil {
					return errorsmod.Wrap(err, "Failed to read genesis doc from file")
				}
			}

			appGenesis.AppName = version.AppName
			appGenesis.AppVersion = version.Version
			appGenesis.ChainID = chainID
			appGenesis.AppState = appState
			appGenesis.InitialHeight = initHeight
			appGenesis.Consensus = &types.ConsensusGenesis{
				Validators: nil,
				Params:     cometbft.DefaultConsensusParams(consensusKeyAlgo),
			}

			if err = genutil.ExportGenesisFile(appGenesis, genFile); err != nil {
				return errorsmod.Wrap(err, "Failed to export genesis file")
			}

			toPrint := newPrintInfo(config.Moniker, chainID, nodeID, "", appState)

			// Note: the config file was already creating before execution this command
			// by [SetupCommand], and it is being overwritten here. The only difference,
			// post default values cleanups, should be in the moniker, which is only setup
			// correctly here
			cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)
			return displayInfo(cmd.ErrOrStderr(), toPrint)
		},
	}

	cmd.Flags().BoolP(FlagOverwrite, "o", false, "overwrite the genesis.json file")
	cmd.Flags().Bool(FlagRecover, false, "provide seed phrase to recover existing key instead of creating")
	cmd.Flags().String(flags.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().String(FlagDefaultBondDenom, "", "genesis file default denomination, if left blank default value is 'stake'")
	cmd.Flags().Int64(flags.FlagInitHeight, 1, "specify the initial block height at genesis")
	return cmd
}
