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
	"math/big"

	"github.com/berachain/beacon-kit/mod/node-builder/config/spec"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/itsdevbear/comet-bls12-381/bls/blst"
	"github.com/spf13/cobra"
)

// Commands creates a new command for deposit related actions.
func Commands() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "deposit",
		Short:                      "deposit subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2, //nolint:mnd // from sdk.
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewValidateDeposit(),
		NewCreateValidator(),
	)

	return cmd
}

// NewValidateDeposit creates a new command for validating a deposit message.
//
//nolint:mnd // lots of magic numbers
func NewValidateDeposit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validates a deposit message for creating a new validator",
		Long: `Validates a deposit message for creating a new validator. The
		deposit message includes the public key, withdrawal credentials,
		and deposit amount. The args taken are in the order of the public key,
		withdrawal credentials, deposit amount, signature, current version,
		and genesis validator root.`,
		Args: cobra.ExactArgs(6),
		RunE: validateDepositMessage,
	}

	return cmd
}

// validateDepositMessage validates a deposit message for creating a new
// validator.
func validateDepositMessage(_ *cobra.Command, args []string) error {
	pubkey, err := convertPubkey(args[0])
	if err != nil {
		return err
	}

	credentials, err := convertWithdrawalCredentials(args[1])
	if err != nil {
		return err
	}

	amount, err := convertAmount(args[2])
	if err != nil {
		return err
	}

	signature, err := convertSignature(args[3])
	if err != nil {
		return err
	}

	currentVersion, err := convertVersion(args[4])
	if err != nil {
		return err
	}

	genesisValidatorRoot, err := convertGenesisValidatorRoot(args[5])
	if err != nil {
		return err
	}

	depositMessage := primitives.DepositMessage{
		Pubkey:      pubkey,
		Credentials: credentials,
		Amount:      amount,
	}

	return depositMessage.VerifyCreateValidator(
		primitives.NewForkData(currentVersion, genesisValidatorRoot),
		signature,
		blst.VerifySignaturePubkeyBytes,
		// TODO: needs to be configurable.
		spec.LocalnetChainSpec().DomainTypeDeposit(),
	)
}

// convertPubkey converts a string to a public key.
func convertPubkey(pubkey string) (primitives.BLSPubkey, error) {
	// convert the public key to a BLSPubkey.
	pubkeyBytes, err := hex.DecodeString(pubkey)
	if err != nil {
		return primitives.BLSPubkey{}, err
	}
	if len(pubkeyBytes) != constants.BLSPubkeyLength {
		return primitives.BLSPubkey{}, ErrInvalidPubKeyLength
	}

	return primitives.BLSPubkey(pubkeyBytes), nil
}

// convertWithdrawalCredentials converts a string to a withdrawal credentials.
func convertWithdrawalCredentials(credentials string) (
	primitives.WithdrawalCredentials,
	error,
) {
	// convert the credentials to a WithdrawalCredentials.
	credentialsBytes, err := hex.DecodeString(credentials)
	if err != nil {
		return primitives.WithdrawalCredentials{}, err
	}
	if len(credentialsBytes) != constants.RootLength {
		return primitives.WithdrawalCredentials{},
			ErrInvalidWithdrawalCredentialsLength
	}
	return primitives.WithdrawalCredentials(credentialsBytes), nil
}

// convertAmount converts a string to a deposit amount.
//
//nolint:mnd // lots of magic numbers
func convertAmount(amount string) (math.Gwei, error) {
	// Convert the amount to a Gwei.
	amountBigInt, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		return 0, ErrInvalidAmount
	}
	return math.Gwei(amountBigInt.Uint64()), nil
}

// convertSignature converts a string to a signature.
func convertSignature(signature string) (primitives.BLSSignature, error) {
	// convert the signature to a BLSSignature.
	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		return primitives.BLSSignature{}, err
	}
	if len(signatureBytes) != constants.BLSSignatureLength {
		return primitives.BLSSignature{}, ErrInvalidSignatureLength
	}
	return primitives.BLSSignature(signatureBytes), nil
}

// convertVersion converts a string to a version.
//

func convertVersion(version string) (primitives.Version, error) {
	versionBytes, err := hex.DecodeString(version)
	if err != nil {
		return primitives.Version{}, err
	}
	if len(versionBytes) != constants.DomainTypeLength {
		return primitives.Version{}, ErrInvalidVersionLength
	}
	return primitives.Version(versionBytes), nil
}

// convertGenesisValidatorRoot converts a string to a genesis validator root.
//

func convertGenesisValidatorRoot(root string) (primitives.Root, error) {
	rootBytes, err := hex.DecodeString(root)
	if err != nil {
		return primitives.Root{}, err
	}
	if len(rootBytes) != constants.RootLength {
		return primitives.Root{}, ErrInvalidRootLength
	}
	return primitives.Root(rootBytes), nil
}
