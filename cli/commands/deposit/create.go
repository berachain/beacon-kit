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
	"encoding/json"
	"fmt"

	"github.com/berachain/beacon-kit/chain"
	clicontext "github.com/berachain/beacon-kit/cli/context"
	"github.com/berachain/beacon-kit/cli/utils/parser"
	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/node-core/components"
	"github.com/berachain/beacon-kit/node-core/components/signer"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/spf13/cobra"
)

const (
	createAddr0     = iota
	createAmt1      = iota
	createRoot2     = iota
	createArgsCount = iota

	overrideNodeKey = "override-node-key"
	valPrivateKey   = "validator-private-key"
)

// NewCreateValidator creates a new command to create a validator deposit.
//
//nolint:lll // reads better if long description is one line
func NewCreateValidator(
	chainSpec chain.Spec,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-validator",
		Short: "Creates a validator deposit",
		Long:  `Creates a validator deposit with the necessary credentials. The arguments are expected in the order of withdrawal address, deposit amount, and genesis validator root. If the broadcast flag is set to true, a private key must be provided to sign the transaction.`,
		Args:  cobra.ExactArgs(createArgsCount),
		RunE:  createValidatorCmd(chainSpec),
	}

	cmd.Flags().BoolP(
		overrideNodeKey,
		"o",
		false, // no override by default
		"override the node private key",
	)
	cmd.Flags().String(
		valPrivateKey,
		"", // no default private key
		"validator private key. This is required if the override-node-key flag is set.",
	)

	return cmd
}

// createValidatorCmd returns a command that builds a create validator request.
func createValidatorCmd(
	chainSpec chain.Spec,
) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// Get the BLS signer.
		blsSigner, err := getBLSSigner(cmd)
		if err != nil {
			return err
		}

		withdrawalAddressStr := args[createAddr0]
		withdrawalAddress, err := parser.ConvertWithdrawalAddress(withdrawalAddressStr)
		if err != nil {
			return err
		}
		credentials := types.NewCredentialsFromExecutionAddress(withdrawalAddress)

		amountStr := args[createAmt1]
		amount, err := parser.ConvertAmount(amountStr)
		if err != nil {
			return err
		}

		genValRootStr := args[createRoot2]
		genesisValidatorRoot, err := parser.ConvertGenesisValidatorRoot(genValRootStr)
		if err != nil {
			return err
		}

		depositMsg, signature, err := CreateDepositMessage(chainSpec, blsSigner, genesisValidatorRoot, credentials, amount)
		if err != nil {
			return err
		}

		val, err := json.Marshal(types.Deposit{
			Pubkey:      depositMsg.Pubkey,
			Credentials: depositMsg.Credentials,
			Amount:      depositMsg.Amount,
			Signature:   signature,
		})
		if err != nil {
			return err
		}

		//nolint:forbidigo // simplifies output parsing
		fmt.Print(string(val))
		return nil
	}
}

func CreateDepositMessage(
	cs chain.Spec,
	blsSigner crypto.BLSSigner,
	genValRoot common.Root,
	creds types.WithdrawalCredentials,
	amount math.Gwei,
) (
	*types.DepositMessage,
	crypto.BLSSignature,
	error,
) {
	// All deposits are signed with the genesis version.
	genesisVersion := version.FromUint32[common.Version](constants.GenesisVersion)

	// Create and sign the deposit message.
	depositMsg, signature, err := types.CreateAndSignDepositMessage(
		types.NewForkData(genesisVersion, genValRoot),
		cs.DomainTypeDeposit(),
		blsSigner,
		creds,
		amount,
	)
	if err != nil {
		return nil, crypto.BLSSignature{}, fmt.Errorf("failed CreateAndSignDepositMessage: %w", err)
	}

	return depositMsg,
		signature,
		ValidateDeposit(
			cs,
			depositMsg.Pubkey,
			depositMsg.Credentials,
			depositMsg.Amount,
			genValRoot,
			signature,
		)
}

// getBLSSigner returns a BLS signer based on the override commands key flag.
func getBLSSigner(
	cmd *cobra.Command,
) (crypto.BLSSigner, error) {
	var legacyKey components.LegacyKey
	overrideFlag, err := cmd.Flags().GetBool(overrideNodeKey)
	if err != nil {
		return nil, err
	}

	// Build the BLS signer.
	if overrideFlag {
		var validatorPrivKey string
		validatorPrivKey, err = cmd.Flags().GetString(valPrivateKey)
		if err != nil {
			return nil, err
		}
		if validatorPrivKey == "" {
			return nil, ErrValidatorPrivateKeyRequired
		}
		legacyKey, err = signer.LegacyKeyFromString(validatorPrivKey)
		if err != nil {
			return nil, err
		}
	}

	return components.ProvideBlsSigner(
		components.BlsSignerInput{
			AppOpts: clicontext.GetViperFromCmd(cmd),
			PrivKey: legacyKey,
		},
	)
}
