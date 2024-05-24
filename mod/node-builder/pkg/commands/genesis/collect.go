// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package genesis

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/genesis"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// CollectGenTxsCmd - return the cobra command to collect genesis transactions.
func CollectGenesisDepositsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collect-premined-deposits",
		Short: "adds a validator to the genesis file",
		RunE: func(cmd *cobra.Command, _ []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			appGenesis, err := genutiltypes.AppGenesisFromFile(
				config.GenesisFile(),
			)
			if err != nil {
				return errors.Wrap(err, "failed to read genesis doc from file")
			}

			// create the app state
			appGenesisState, err := genutiltypes.GenesisStateFromAppGenesis(
				appGenesis,
			)
			if err != nil {
				return err
			}

			var deposits []*types.Deposit
			if deposits, err = CollectValidatorJSONFiles(
				filepath.Join(config.RootDir, "config", "gentx"),
				appGenesis,
			); err != nil {
				return errors.Wrap(
					err,
					"failed to collect validator json files",
				)
			}

			genesisInfo := &genesis.Genesis[
				*types.Deposit,
				*types.ExecutionPayloadHeaderDeneb,
			]{}

			if err = json.Unmarshal(
				appGenesisState["beacon"], genesisInfo,
			); err != nil {
				return errors.Wrap(err, "failed to unmarshal beacon genesis")
			}

			for i, deposit := range deposits {
				//#nosec:G701 // won't realistically overflow.
				deposit.Index = uint64(i)
				genesisInfo.Deposits = append(genesisInfo.Deposits, deposit)
			}

			appGenesisState["beacon"], err = json.Marshal(genesisInfo)
			if err != nil {
				return errors.Wrap(err, "failed to marshal beacon genesis")
			}

			if appGenesis.AppState, err = json.MarshalIndent(
				appGenesisState, "", "  ",
			); err != nil {
				return err
			}

			return genutil.ExportGenesisFile(appGenesis, config.GenesisFile())
		},
	}

	return cmd
}

// CollectValidatorJSONFiles.
func CollectValidatorJSONFiles(
	genTxsDir string,
	genesis *genutiltypes.AppGenesis,
) ([]*types.Deposit, error) {
	// prepare a map of all balances in genesis state to then validate
	// against the validators addresses
	var appState map[string]json.RawMessage
	if err := json.Unmarshal(genesis.AppState, &appState); err != nil {
		return nil, err
	}

	// get the list of files in the genTxsDir
	fos, err := os.ReadDir(genTxsDir)
	if err != nil {
		return nil, err
	}

	// prepare the list of validators
	deposits := make([]*types.Deposit, 0)
	for _, fo := range fos {
		if fo.IsDir() {
			continue
		}
		if !strings.HasSuffix(fo.Name(), ".json") {
			continue
		}

		var bz []byte
		bz, err = afero.ReadFile(
			afero.NewOsFs(),
			filepath.Join(genTxsDir, fo.Name()),
		)
		if err != nil {
			return nil, err
		}

		val := &types.Deposit{}
		if err = json.Unmarshal(bz, val); err != nil {
			return nil, err
		}

		deposits = append(deposits, val)
	}

	return deposits, nil
}
