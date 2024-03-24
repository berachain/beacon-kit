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

package cli

import (
	"encoding/json"

	"github.com/berachain/beacon-kit/beacon/core/state"
	beacontypes "github.com/berachain/beacon-kit/beacon/core/types"
	"github.com/cockroachdb/errors"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

// CollectGenTxsCmd - return the cobra command to collect genesis transactions.
func AddPubkeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-validator",
		Short: "adds a validator to the genesis file",
		RunE: func(cmd *cobra.Command, _ []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			_, valPubKey, err := genutil.InitializeNodeValidatorFiles(
				config,
			)
			if err != nil {
				return errors.Wrap(
					err,
					"failed to initialize node validator files",
				)
			}

			genesis, err := types.AppGenesisFromFile(config.GenesisFile())
			if err != nil {
				return errors.Wrap(err, "failed to read genesis doc from file")
			}

			// create the app state
			appGenesisState, err := types.GenesisStateFromAppGenesis(genesis)
			if err != nil {
				return err
			}

			// Create a new validator and add it to the app state.
			validator := &beacontypes.Validator{
				Pubkey: [48]byte(valPubKey.Bytes()),
				// TODO: credentials.
				Credentials: beacontypes.NewCredentialsFromExecutionAddress(
					common.Address{},
				),
				EffectiveBalance: 1,
				Slashed:          false,
			}

			beaconState := &state.BeaconStateDeneb{}
			if err = json.Unmarshal(
				appGenesisState["beacon"], beaconState,
			); err != nil {
				return errors.Wrap(err, "failed to unmarshal beacon state")
			}

			beaconState.Validators = append(beaconState.Validators, validator)

			appGenesisState["beacon"], err = json.Marshal(beaconState)
			if err != nil {
				return errors.Wrap(err, "failed to marshal beacon state")
			}

			if genesis.AppState, err = json.MarshalIndent(
				appGenesisState, "", "  ",
			); err != nil {
				return err
			}

			return genutil.ExportGenesisFile(genesis, config.GenesisFile())
		},
	}

	return cmd
}
