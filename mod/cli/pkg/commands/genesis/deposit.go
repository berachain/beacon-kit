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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/mod/cli/pkg/utils/context"
	"github.com/berachain/beacon-kit/mod/cli/pkg/utils/parser"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/errors"
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/signer"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// AddGenesisDepositCmd - returns the cobra command to
// add a premined deposit to the genesis file.
func AddGenesisDepositCmd(cs common.ChainSpec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-premined-deposit",
		Short: "adds a validator to the genesis file",
		RunE: func(cmd *cobra.Command, _ []string) error {
			serverCtx := context.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			_, valPubKey, err := genutil.InitializeNodeValidatorFiles(
				config, crypto.CometBLSType,
			)
			if err != nil {
				return errors.Wrap(
					err,
					"failed to initialize commands validator files",
				)
			}

			var (
				depositAmountString string
				depositAmount       math.Gwei
			)

			// Get the BLS signer.
			blsSigner, err := getBLSSigner(client.GetViperFromCmd(cmd))
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
			currentVersion := version.FromUint32[common.Version](
				version.Deneb,
			)

			depositMsg, signature, err := types.CreateAndSignDepositMessage(
				types.NewForkData(currentVersion, common.Root{}),
				cs.DomainTypeDeposit(),
				blsSigner,
				// TODO: configurable.
				types.NewCredentialsFromExecutionAddress(
					gethprimitives.ExecutionAddress{},
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
func getBLSSigner(v *viper.Viper) (crypto.BLSSigner, error) {
	var blsSigner crypto.BLSSigner
	if err := depinject.Inject(
		depinject.Configs(
			depinject.Supply(
				v,
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
