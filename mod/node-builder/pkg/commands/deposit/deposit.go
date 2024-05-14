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
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/commands/utils/parser"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/components/signer"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/config/spec"
	"github.com/cosmos/cosmos-sdk/client"
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
func validateDepositMessage(
	_ *cobra.Command,
	args []string,
) error {
	pubkey, err := parser.ConvertPubkey(args[0])
	if err != nil {
		return err
	}

	credentials, err := parser.ConvertWithdrawalCredentials(args[1])
	if err != nil {
		return err
	}

	amount, err := parser.ConvertAmount(args[2])
	if err != nil {
		return err
	}

	signature, err := parser.ConvertSignature(args[3])
	if err != nil {
		return err
	}

	currentVersion, err := parser.ConvertVersion(args[4])
	if err != nil {
		return err
	}

	genesisValidatorRoot, err := parser.ConvertGenesisValidatorRoot(args[5])
	if err != nil {
		return err
	}

	depositMessage := types.DepositMessage{
		Pubkey:      pubkey,
		Credentials: credentials,
		Amount:      amount,
	}

	return depositMessage.VerifyCreateValidator(
		types.NewForkData(currentVersion, genesisValidatorRoot),
		signature,
		signer.BLSSigner{}.VerifySignature,
		// TODO: needs to be configurable.
		spec.LocalnetChainSpec().DomainTypeDeposit(),
	)
}
