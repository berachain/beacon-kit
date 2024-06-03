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

	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/commands/utils/parser"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/components"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/components/signer"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// CollectGenTxsCmd - return the cobra command to collect genesis transactions.
func AddGenesisDepositCmd(cs primitives.ChainSpec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-premined-deposit",
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

			var (
				depositAmountString string
				depositAmount       math.Gwei
			)

			// Get the BLS signer.
			blsSigner, err := getBLSSigner()
			if err != nil {
				return err
			}

			// Get the deposit amount.
			depositAmountString, err = cmd.Flags().GetString(depositAmountFlag)
			if err != nil {
				return err
			}

			depositAmount, err = parser.ConvertAmount(depositAmountString)
			if err != nil {
				return err
			}

			// TODO: configurable.
			currentVersion := version.FromUint32[primitives.Version](
				version.Deneb,
			)

			depositMsg, signature, err := types.CreateAndSignDepositMessage(
				types.NewForkData(currentVersion, common.Root{}),
				cs.DomainTypeDeposit(),
				blsSigner,
				// TODO: configurable.
				types.NewCredentialsFromExecutionAddress(
					common.ExecutionAddress{},
				),
				depositAmount,
			)
			if err != nil {
				return err
			}

			// Verify the deposit message.
			if err = depositMsg.VerifyCreateValidator(
				types.NewForkData(currentVersion, common.Root{}),
				signature,
				cs.DomainTypeDeposit(),
				signer.BLSSigner{}.VerifySignature,
			); err != nil {
				return err
			}

			deposit := types.Deposit{
				Pubkey:      depositMsg.Pubkey,
				Amount:      depositMsg.Amount,
				Signature:   signature,
				Credentials: depositMsg.Credentials,
			}

			//#nosec:G703 // Ignore errors on this line.
			outputDocument, _ := cmd.Flags().GetString(flags.FlagOutputDocument)
			if outputDocument == "" {
				outputDocument, err = makeOutputFilepath(config.RootDir,
					crypto.BLSPubkey(valPubKey.Bytes()).String())
				if err != nil {
					return errors.Wrap(err, "failed to create output file path")
				}
			}

			if err = writeDepositToFile(outputDocument, &deposit); err != nil {
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
	writePath := filepath.Join(rootDir, "config", "premined-deposits")
	if err := afero.NewOsFs().MkdirAll(writePath, os.ModePerm); err != nil {
		return "", errors.Newf(
			"could not create directory %q: %w",
			writePath,
			err,
		)
	}

	return filepath.Join(
		writePath,
		fmt.Sprintf("premined-deposit-%v.json", pubkey),
	), nil
}

func writeDepositToFile(
	outputDocument string,
	depositMessage *types.Deposit,
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

	bz, err := json.Marshal(depositMessage)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(outputFile, "%s\n", bz)

	return err
}

// getBLSSigner returns a BLS signer based on the override node key flag.
func getBLSSigner() (crypto.BLSSigner, error) {
	var blsSigner crypto.BLSSigner
	if err := depinject.Inject(
		depinject.Configs(
			depinject.Supply(
				viper.GetViper(),
			),
			depinject.Provide(
				components.ProvideBlsSigner,
			),
		),
		&blsSigner,
	); err != nil {
		return nil, err
	}

	return blsSigner, nil
}
