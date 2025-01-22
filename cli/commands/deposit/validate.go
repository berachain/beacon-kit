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
	"github.com/berachain/beacon-kit/chain"
	"github.com/berachain/beacon-kit/cli/utils/parser"
	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/node-core/components/signer"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/spf13/cobra"
)

const (
	validateDepPubKeyIdx     = iota
	validateDepCredsIdx      = iota
	validateDepAmtIdx        = iota
	validateDepSignIdx       = iota
	validateDepGenValRootIdx = iota
	validateDepArgsCount     = iota
)

// NewValidateDeposit creates a new command for validating a deposit message.
//
//nolint:lll // reads better if long description is one line
func NewValidateDeposit(chainSpec chain.Spec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validates a deposit message for creating a new validator",
		Long:  `Validates a deposit message for creating a new validator. The deposit message includes the public key, withdrawal credentials, and deposit amount. The args taken are in the order of the public key, withdrawal credentials, deposit amount, signature, and genesis validator root.`,
		Args:  cobra.ExactArgs(validateDepArgsCount),
		RunE:  validateDepositMessage(chainSpec),
	}

	return cmd
}

// validateDepositMessage validates a deposit message for creating a new
// validator.
func validateDepositMessage(chainSpec chain.Spec) func(
	_ *cobra.Command,
	args []string,
) error {
	return func(_ *cobra.Command, args []string) error {
		pubKeyStr := args[validateDepPubKeyIdx]
		pubkey, err := parser.ConvertPubkey(pubKeyStr)
		if err != nil {
			return err
		}

		credsStr := args[validateDepCredsIdx]
		credentials, err := parser.ConvertWithdrawalCredentials(credsStr)
		if err != nil {
			return err
		}

		amountStr := args[validateDepAmtIdx]
		amount, err := parser.ConvertAmount(amountStr)
		if err != nil {
			return err
		}

		sigStr := args[validateDepSignIdx]
		signature, err := parser.ConvertSignature(sigStr)
		if err != nil {
			return err
		}

		genValRootStr := args[validateDepGenValRootIdx]
		genesisValidatorRoot, err := parser.ConvertGenesisValidatorRoot(genValRootStr)
		if err != nil {
			return err
		}

		return ValidateDeposit(chainSpec, pubkey, genesisValidatorRoot, credentials, amount, signature)
	}
}

func ValidateDeposit(
	cs chain.Spec,
	pubkey crypto.BLSPubkey,
	genValRoot common.Root,
	creds types.WithdrawalCredentials,
	amount math.Gwei,
	signature crypto.BLSSignature,
) error {
	depositMessage := types.DepositMessage{
		Pubkey:      pubkey,
		Credentials: creds,
		Amount:      amount,
	}

	// All deposits are signed with the genesis version.
	genesisVersion := version.FromUint32[common.Version](constants.GenesisVersion)

	return depositMessage.VerifyCreateValidator(
		types.NewForkData(genesisVersion, genValRoot),
		signature,
		cs.DomainTypeDeposit(),
		signer.BLSSigner{}.VerifySignature,
	)
}
