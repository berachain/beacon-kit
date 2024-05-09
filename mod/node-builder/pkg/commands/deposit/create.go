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

package deposit

import (
	"encoding/hex"
	"os"

	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/commands/utils/parser"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/components"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/components/signer"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/config/spec"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewValidateDeposit creates a new command for validating a deposit message.
//
//nolint:gomnd // lots of magic numbers
func NewCreateValidator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-validator",
		Short: "Creates a validator deposit",
		Long: `Creates a validator deposit with the necessary credentials. The 
		arguments are expected in the order of withdrawal credentials, deposit
		amount, current version, and genesis validator root. If the broadcast
		flag is set to true, a private key must be provided to sign the transaction.`,
		Args: cobra.ExactArgs(4), //nolint:mnd // The number of arguments.
		RunE: createValidatorCmd(),
	}

	cmd.Flags().BoolP(
		broadcastDeposit, broadcastDepositShorthand,
		defaultBroadcastDeposit, broadcastDepositMsg,
	)
	cmd.Flags().String(privateKey, defaultPrivateKey, privateKeyMsg)
	cmd.Flags().BoolP(
		overrideNodeKey, overrideNodeKeyShorthand,
		defaultOverrideNodeKey, overrideNodeKeyMsg,
	)
	cmd.Flags().
		String(valPrivateKey, defaultValidatorPrivateKey, valPrivateKeyMsg)
	cmd.Flags().String(jwtSecretPath, defaultJWTSecretPath, jwtSecretPathMsg)
	cmd.Flags().String(engineRPCURL, defaultEngineRPCURL, engineRPCURLMsg)

	return cmd
}

// createValidatorCmd returns a command that builds a create validator request.
//
// TODO: Implement broadcast functionality. Currently, the implementation works
// for the geth client but something about the Deposit binding is not handling
// other execution layers correctly. Peep the commit history for what we had.
// 🤷‍♂️.
func createValidatorCmd() func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		var (
			logger = log.NewLogger(os.Stdout)
		)

		// Get the BLS signer.
		blsSigner, err := getBLSSigner(cmd)
		if err != nil {
			return err
		}

		credentials, err := parser.ConvertWithdrawalCredentials(args[0])
		if err != nil {
			return err
		}

		amount, err := parser.ConvertAmount(args[1])
		if err != nil {
			return err
		}

		currentVersion, err := parser.ConvertVersion(args[2])
		if err != nil {
			return err
		}

		genesisValidatorRoot, err := parser.ConvertGenesisValidatorRoot(args[3])
		if err != nil {
			return err
		}

		// Create and sign the deposit message.
		depositMsg, signature, err := types.CreateAndSignDepositMessage(
			types.NewForkData(currentVersion, genesisValidatorRoot),
			spec.LocalnetChainSpec().DomainTypeDeposit(),
			blsSigner,
			credentials,
			amount,
		)
		if err != nil {
			return err
		}

		// Verify the deposit message.
		if err = depositMsg.VerifyCreateValidator(
			types.NewForkData(currentVersion, genesisValidatorRoot),
			signature,
			signer.BLSSigner{}.VerifySignature,
			spec.LocalnetChainSpec().DomainTypeDeposit(),
		); err != nil {
			return err
		}

		// If the broadcast flag is not set, output the deposit message and
		// signature and return early.
		logger.Info(
			"Deposit Message CallData",
			"pubkey", parser.EncodeTo0xPrefixedString(depositMsg.Pubkey[:]),
			"withdrawal credentials",
			parser.EncodeTo0xPrefixedString(depositMsg.Credentials[:]),
			"amount", depositMsg.Amount,
			"signature", parser.EncodeTo0xPrefixedString(signature[:]),
		)

		// TODO: once broadcast is fixed, remove this.
		logger.Info("Send the above calldata to the deposit contract 🫡")

		return nil
	}
}

// getBLSSigner returns a BLS signer based on the override node key flag.
func getBLSSigner(
	cmd *cobra.Command,
) (crypto.BLSSigner, error) {
	var blsSigner crypto.BLSSigner
	// If the override node key flag is set, a validator private key must be
	// provided.
	overrideFlag, err := cmd.Flags().GetBool(overrideNodeKey)
	if err != nil {
		return nil, err
	}

	// Build the BLS signer.
	//nolint:nestif // complexity comes from parsing values
	if overrideFlag {
		var (
			validatorPrivKey   string
			validatorPrivKeyBz []byte
		)
		validatorPrivKey, err = cmd.Flags().GetString(valPrivateKey)
		if err != nil {
			return nil, err
		}
		if validatorPrivKey == "" {
			return nil, ErrValidatorPrivateKeyRequired
		}

		validatorPrivKeyBz, err = hex.DecodeString(validatorPrivKey)
		if err != nil {
			return nil, err
		}
		if len(validatorPrivKeyBz) != constants.BLSSecretKeyLength {
			return nil, ErrInvalidValidatorPrivateKeyLength
		}

		blsSigner, err = signer.NewBLSSigner(
			[constants.BLSSecretKeyLength]byte(validatorPrivKeyBz),
		)
		if err != nil {
			return nil, err
		}

		return blsSigner, nil
	}

	if err = depinject.Inject(
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
