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
	"fmt"
	"os"
	"path/filepath"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/commands/utils/parser"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/afero"
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
				config, crypto.CometBLSType,
			)
			if err != nil {
				return errors.Wrap(
					err,
					"failed to initialize node validator files",
				)
			}

			genesis, err := genutiltypes.AppGenesisFromFile(
				config.GenesisFile(),
			)
			if err != nil {
				return errors.Wrap(err, "failed to read genesis doc from file")
			}

			// create the app state
			_, err = genutiltypes.GenesisStateFromAppGenesis(genesis)
			if err != nil {
				return err
			}

			var (
				depositAmountString string
				depositAmount       math.Gwei
			)
			// Get the deposit amount.
			depositAmountString, err = cmd.Flags().GetString(depositAmountFlag)
			if err != nil {
				return err
			}
			depositAmount, err = parser.ConvertAmount(depositAmountString)
			if err != nil {
				return err
			}

			// TODO: Should we do deposits here?
			validator := types.NewValidatorFromDeposit(
				crypto.BLSPubkey(valPubKey.Bytes()),
				types.NewCredentialsFromExecutionAddress(
					common.Address{},
				),
				depositAmount,
				depositAmount,
				32e9, //nolint:mnd // temp.
			)

			//#nosec:G703 // Ignore errors on this line.
			outputDocument, _ := cmd.Flags().GetString(flags.FlagOutputDocument)
			if outputDocument == "" {
				outputDocument, err = makeOutputFilepath(config.RootDir,
					crypto.BLSPubkey(valPubKey.Bytes()).String())
				if err != nil {
					return errors.Wrap(err, "failed to create output file path")
				}
			}
			if err = writeValidatorStruct(outputDocument, validator); err != nil {
				return errors.Wrap(err, "failed to write signed gen tx")
			}

			return nil
		},
	}

	cmd.Flags().
		String(depositAmountFlag, defaultDepositAmount, depositAmountFlagMsg)

	return cmd
}

func makeOutputFilepath(rootDir, pubkey string) (string, error) {
	writePath := filepath.Join(rootDir, "config", "gentx")
	if err := afero.NewOsFs().MkdirAll(writePath, os.ModePerm); err != nil {
		return "", errors.Newf(
			"could not create directory %q: %w",
			writePath,
			err,
		)
	}

	return filepath.Join(writePath, fmt.Sprintf("gentx-%v.json", pubkey)), nil
}

func writeValidatorStruct(
	outputDocument string,
	validator *types.Validator,
) error {
	//#nosec:G302,G304 // Ignore errors on this line.
	outputFile, err := afero.NewOsFs().OpenFile(
		outputDocument,
		os.O_CREATE|os.O_EXCL|os.O_WRONLY,
		0o644, //nolint:mnd // file permissions.
	)
	if err != nil {
		return err
	}

	//#nosec:G307 // Ignore errors on this line.
	defer outputFile.Close()

	bz, err := json.Marshal(validator)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(outputFile, "%s\n", bz)

	return err
}
