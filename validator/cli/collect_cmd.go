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
	"path/filepath"

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/cobra"
)

const flagGenTxDir = "gentx-dir"

// CollectGenTxsCmd - return the cobra command to collect genesis transactions.
func CollectGenTxsCmd(
	genBalIterator types.GenesisBalancesIterator,
	validator types.MessageValidator,
	valAddrCodec runtime.ValidatorAddressCodec,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "customcollect-gentxs",
		Short: "Collect genesis txs and output a genesis.json file",
		RunE: func(cmd *cobra.Command, _ []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			clientCtx := client.GetClientContextFromCmd(cmd)
			cdc := clientCtx.Codec

			nodeID, valPubKey, err := genutil.InitializeNodeValidatorFiles(
				config,
			)
			if err != nil {
				return errors.Wrap(
					err,
					"failed to initialize node validator files",
				)
			}

			appGenesis, err := types.AppGenesisFromFile(config.GenesisFile())
			if err != nil {
				return errors.Wrap(err, "failed to read genesis doc from file")
			}

			// Read the genesis transactions from the default directory
			// or from the one provided by the flag.
			var genTxDir string
			genTxDir, err = cmd.Flags().GetString(flagGenTxDir)
			if err != nil {
				return errors.Wrap(err, "failed to get gentx dir")
			}
			genTxsDir := genTxDir
			if genTxsDir == "" {
				genTxsDir = filepath.Join(config.RootDir, "config", "gentx")
			}

			initCfg := types.NewInitConfig(
				appGenesis.ChainID,
				genTxsDir,
				nodeID,
				valPubKey,
			)

			_, err = GenAppStateFromConfig(
				cdc,
				clientCtx.TxConfig,
				config,
				initCfg,
				appGenesis,
				genBalIterator,
				validator,
				valAddrCodec,
			)
			if err != nil {
				return errors.Wrap(
					err,
					"failed to get genesis app state from config",
				)
			}
			return nil
		},
	}

	cmd.Flags().
		String(
			flagGenTxDir, "",
			"override default \"gentx\" directory from which "+
				"collect and execute genesis transactions; "+
				"default [--home]/config/gentx/")
	return cmd
}
