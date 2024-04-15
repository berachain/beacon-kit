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

	"github.com/berachain/beacon-kit/mod/primitives"
	consensusprimitives "github.com/berachain/beacon-kit/mod/primitives-consensus"
	"github.com/berachain/beacon-kit/mod/primitives/constants"
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
		SuggestionsMinimumDistance: 2, //nolint:gomnd // from sdk.
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewValidateDeposit(),
	)

	return cmd
}

// NewValidateDeposit creates a new command for validating a deposit message.
//
//nolint:gomnd // lots of magic numbers
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
	pubkey, err := ConvertPubkey(args[0])
	if err != nil {
		return err
	}

	credentials, err := ConvertWithdrawalCredentials(args[1])
	if err != nil {
		return err
	}

	amount, err := ConvertAmount(args[2])
	if err != nil {
		return err
	}

	signature, err := ConvertSignature(args[3])
	if err != nil {
		return err
	}

	currentVersion, err := ConvertVersion(args[4])
	if err != nil {
		return err
	}

	genesisValidatorRoot, err := ConvertGenesisValidatorRoot(args[5])
	if err != nil {
		return err
	}

	depositMessage := consensusprimitives.DepositMessage{
		Pubkey:      pubkey,
		Credentials: credentials,
		Amount:      amount,
	}

	return depositMessage.VerifyCreateValidator(
		consensusprimitives.NewForkData(currentVersion, genesisValidatorRoot),
		signature,
		blst.VerifySignaturePubkeyBytes,
	)
}

// ConvertPubkey converts a string to a public key.
func ConvertPubkey(pubkey string) (primitives.BLSPubkey, error) {
	// Convert the public key to a BLSPubkey.
	pubkeyBytes, err := hex.DecodeString(pubkey)
	if err != nil {
		return primitives.BLSPubkey{}, err
	}
	if len(pubkeyBytes) != constants.BLSPubkeyLength {
		return primitives.BLSPubkey{}, ErrInvalidPubKeyLength
	}

	return primitives.BLSPubkey(pubkeyBytes), nil
}

// ConvertWithdrawalCredentials converts a string to a withdrawal credentials.
func ConvertWithdrawalCredentials(credentials string) (
	consensusprimitives.WithdrawalCredentials,
	error,
) {
	// Convert the credentials to a WithdrawalCredentials.
	credentialsBytes, err := hex.DecodeString(credentials)
	if err != nil {
		return consensusprimitives.WithdrawalCredentials{}, err
	}
	if len(credentialsBytes) != constants.RootLength {
		return consensusprimitives.WithdrawalCredentials{},
			ErrInvalidWithdrawalCredentialsLength
	}
	return consensusprimitives.WithdrawalCredentials(credentialsBytes), nil
}

// ConvertAmount converts a string to a deposit amount.
//
//nolint:gomnd // lots of magic numbers
func ConvertAmount(amount string) (primitives.Gwei, error) {
	// Convert the amount to a Gwei.
	amountBigInt, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		return 0, ErrInvalidAmount
	}
	return primitives.Gwei(amountBigInt.Uint64()), nil
}

// ConvertSignature converts a string to a signature.
func ConvertSignature(signature string) (primitives.BLSSignature, error) {
	// Convert the signature to a BLSSignature.
	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		return primitives.BLSSignature{}, err
	}
	if len(signatureBytes) != constants.BLSSignatureLength {
		return primitives.BLSSignature{}, ErrInvalidSignatureLength
	}
	return primitives.BLSSignature(signatureBytes), nil
}

// ConvertVersion converts a string to a version.
//

func ConvertVersion(version string) (primitives.Version, error) {
	versionBytes, err := hex.DecodeString(version)
	if err != nil {
		return primitives.Version{}, err
	}
	if len(versionBytes) != constants.DomainTypeLength {
		return primitives.Version{}, ErrInvalidVersionLength
	}
	return primitives.Version(versionBytes), nil
}

// ConvertGenesisValidatorRoot converts a string to a genesis validator root.
//

func ConvertGenesisValidatorRoot(root string) (primitives.Root, error) {
	rootBytes, err := hex.DecodeString(root)
	if err != nil {
		return primitives.Root{}, err
	}
	if len(rootBytes) != constants.RootLength {
		return primitives.Root{}, ErrInvalidRootLength
	}
	return primitives.Root(rootBytes), nil
}
