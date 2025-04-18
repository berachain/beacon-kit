// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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

package deposit

import (
	clitypes "github.com/berachain/beacon-kit/cli/commands/server/types"
	"github.com/berachain/beacon-kit/cli/context"
	"github.com/berachain/beacon-kit/cli/utils/parser"
	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/node-core/components/signer"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/spf13/cobra"
)

const (
	validatePubKey0 = iota
	validateCreds1  = iota
	validateAmt2    = iota
	validateSign3   = iota

	minArgsValidateDeposit = 4
	maxArgsValidateDeposit = 5
)

// GetValidateDepositCmd creates a new command for validating a deposit message.
//
//nolint:lll // Reads better if long description is one line.
func GetValidateDepositCmd(chainSpecCreator clitypes.ChainSpecCreator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate [pubkey] [withdrawal-credentials] [amount] [signature] ?[beacond/genesis.json]",
		Short: "Validates a deposit message for creating a new validator",
		Long:  `Validates a deposit message (public key, withdrawal credentials, deposit amount) for creating a new validator. The args taken are in the order of the public key, withdrawal credentials, deposit amount, signature, and optionally the beacond genesis file. If the genesis validator root flag is NOT set, the beacond genesis file MUST be provided as the last argument.`,
		Args:  cobra.RangeArgs(minArgsValidateDeposit, maxArgsValidateDeposit),
		RunE:  validateDepositMessage(chainSpecCreator),
	}

	cmd.Flags().StringP(
		useGenesisValidatorRoot,
		useGenesisValidatorRootShorthand,
		defaultGenesisValidatorRoot,
		"Use the provided genesis validator root. If this is not set, the beacond genesis file must be provided manually as the last argument.",
	)

	return cmd
}

// validateDepositMessage validates a deposit message for creating a new validator.
func validateDepositMessage(chainSpecCreator clitypes.ChainSpecCreator) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		v := context.GetViperFromCmd(cmd)
		chainSpec, err := chainSpecCreator(v)
		if err != nil {
			return err
		}
		pubKeyStr := args[validatePubKey0]
		pubkey, err := parser.ConvertPubkey(pubKeyStr)
		if err != nil {
			return err
		}

		credsStr := args[validateCreds1]
		credentials, err := parser.ConvertWithdrawalCredentials(credsStr)
		if err != nil {
			return err
		}

		amountStr := args[validateAmt2]
		amount, err := parser.ConvertAmount(amountStr)
		if err != nil {
			return err
		}

		sigStr := args[validateSign3]
		signature, err := parser.ConvertSignature(sigStr)
		if err != nil {
			return err
		}

		genesisValidatorRoot, err := getGenesisValidatorRoot(
			cmd, chainSpec, args, maxArgsValidateDeposit,
		)
		if err != nil {
			return err
		}

		if err = ValidateDeposit(
			chainSpec, pubkey, credentials, amount, genesisValidatorRoot, signature,
		); err != nil {
			return err
		}

		cmd.Println("✅ Deposit message is valid!")
		return nil
	}
}

func ValidateDeposit(
	cs ChainSpec,
	pubkey crypto.BLSPubkey,
	creds types.WithdrawalCredentials,
	amount math.Gwei,
	genValRoot common.Root,
	signature crypto.BLSSignature,
) error {
	depositMessage := types.DepositMessage{
		Pubkey:      pubkey,
		Credentials: creds,
		Amount:      amount,
	}

	// All deposits are signed with the genesis version.
	return depositMessage.VerifyCreateValidator(
		types.NewForkData(cs.GenesisForkVersion(), genValRoot),
		signature,
		cs.DomainTypeDeposit(),
		signer.BLSSigner{}.VerifySignature,
	)
}
