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
	"fmt"
	"os"
	"path/filepath"

	"github.com/berachain/beacon-kit/chain-spec/chain"
	"github.com/berachain/beacon-kit/cli/context"
	"github.com/berachain/beacon-kit/cli/utils/parser"
	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/node-core/components"
	"github.com/berachain/beacon-kit/node-core/components/signer"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/encoding/json"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// AddGenesisDepositCmd - returns the cobra command to
// add a premined deposit to the genesis file.
//

func AddGenesisDepositCmd(cs chain.ChainSpec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-premined-deposit",
		Short: "adds a validator to the genesis file",
		Long: `Adds a validator to the genesis file with the necessary
		credentials. The arguments are expected in the order of the deposit
		amount and withdrawal address.`,
		Args: cobra.ExactArgs(2), //nolint:mnd // The number of arguments.
		RunE: func(cmd *cobra.Command, args []string) error {
			config := context.GetConfigFromCmd(cmd)

			_, valPubKey, err := genutil.InitializeNodeValidatorFiles(
				config, crypto.CometBLSType,
			)
			if err != nil {
				return errors.Wrap(
					err,
					"failed to initialize commands validator files",
				)
			}

			// Get the BLS signer.
			blsSigner, err := components.ProvideBlsSigner(
				components.BlsSignerInput{
					AppOpts: context.GetViperFromCmd(cmd),
				},
			)
			if err != nil {
				return err
			}

			// Get the deposit amount.
			var depositAmount math.Gwei
			depositAmount, err = parser.ConvertAmount(args[0])
			if err != nil {
				return err
			}

			// TODO: configurable.
			currentVersion := version.FromUint32[common.Version](
				version.Deneb,
			)

			// Get the withdrawal address.
			withdrawalAddress := common.NewExecutionAddressFromHex(args[1])

			depositMsg, signature, err := types.CreateAndSignDepositMessage(
				types.NewForkData(currentVersion, common.Root{}),
				cs.DomainTypeDeposit(),
				blsSigner,
				types.NewCredentialsFromExecutionAddress(withdrawalAddress),
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

	return cmd
}

func makeOutputFilepath(rootDir, pubkey string) (string, error) {
	writePath := filepath.Join(rootDir, "config", "premined-deposits")
	if err := afero.NewOsFs().MkdirAll(writePath, os.ModePerm); err != nil {
		return "", errors.Wrapf(
			errors.New("could not create directory"), "%q: %w",
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
