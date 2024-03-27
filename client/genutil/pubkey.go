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
	"fmt"
	"os"
	"path/filepath"

	beacontypes "github.com/berachain/beacon-kit/beacon/core/types"
	"github.com/berachain/beacon-kit/primitives"
	"github.com/cockroachdb/errors"
	"github.com/cosmos/cosmos-sdk/client/flags"
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
			_, err = types.GenesisStateFromAppGenesis(genesis)
			if err != nil {
				return err
			}

			// TODO: Should we do deposits here?
			validator := beacontypes.NewValidatorFromDeposit(
				primitives.BLSPubkey(valPubKey.Bytes()),
				beacontypes.NewCredentialsFromExecutionAddress(
					common.Address{},
				),
				1e9,  //nolint:gomnd // temp.
				1e9,  //nolint:gomnd // temp.
				32e9, //nolint:gomnd // temp.
			)

			//#nosec:G703 // Ignore errors on this line.
			outputDocument, _ := cmd.Flags().GetString(flags.FlagOutputDocument)
			if outputDocument == "" {
				outputDocument, err = makeOutputFilepath(config.RootDir,
					primitives.BLSPubkey(valPubKey.Bytes()).String())
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

	return cmd
}

func makeOutputFilepath(rootDir, pubkey string) (string, error) {
	writePath := filepath.Join(rootDir, "config", "gentx")
	if err := os.MkdirAll(writePath, 0o700); err != nil {
		return "", fmt.Errorf(
			"could not create directory %q: %w",
			writePath,
			err,
		)
	}

	return filepath.Join(writePath, fmt.Sprintf("gentx-%v.json", pubkey)), nil
}

func writeValidatorStruct(
	outputDocument string,
	validator *beacontypes.Validator,
) error {
	//#nosec:G302,G304 // Ignore errors on this line.
	outputFile, err := os.OpenFile(
		outputDocument,
		os.O_CREATE|os.O_EXCL|os.O_WRONLY,
		0o644,
	)
	if err != nil {
		return err
	}

	//#nosec:G307 // Ignore errors on this line.
	defer outputFile.Close()

	bz, err := validator.MarshalJSON()
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(outputFile, "%s\n", bz)

	return err
}
