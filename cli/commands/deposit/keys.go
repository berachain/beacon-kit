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
	"encoding/base64"

	"github.com/berachain/beacon-kit/cli/context"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/node-core/components"
	"github.com/berachain/beacon-kit/node-core/components/signer"
	"github.com/cometbft/cometbft/crypto/bls12381"
	"github.com/spf13/cobra"
)

// GetValidatorKeysCmd returns a command that returns the validator public key in different formats
// for the given private key files.
//
//nolint:lll // reads better if long description is one line.
func GetValidatorKeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-keys",
		Short: "Outputs the validator public key in different formats.",
		Long:  `Outputs the validator public key in formats of Comet address, Comet pubkey, and Eth/Beacon pubkey. Uses the private key file specified as the value of "priv_validator_key_file" in the config.toml file in the beacond HOMEDIR.`,
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get the BLS signer.
			blsSignerI, err := components.ProvideBlsSigner(
				components.BlsSignerInput{
					AppOpts: context.GetViperFromCmd(cmd),
				},
			)
			if err != nil {
				return errors.Wrap(err, "failed to initialize BLS signer from validator files")
			}
			blsSigner, ok := blsSignerI.(*signer.BLSSigner)
			if !ok {
				return errors.New("failed to assert BLS signer type")
			}

			// Get the comet public key.
			cometKey, err := blsSigner.PrivValidator.GetPubKey()
			if err != nil {
				return errors.Wrap(err, "failed to get comet public key from bls signer")
			}

			// Get the comet BLS public key.
			blsKey, err := bls12381.NewPublicKeyFromBytes(cometKey.Bytes())
			if err != nil {
				return errors.Wrap(err, "failed to create BLS key from bytes")
			}

			// Output the validator public key in different formats.
			cmd.Printf(
				"Comet Address: \"%s\"\n",
				cometKey.Address(),
			)
			cmd.Printf(
				"Comet Pubkey (Base64): \"%s\"\n",
				base64.StdEncoding.EncodeToString(blsKey.Bytes()),
			)
			cmd.Printf(
				"Eth/Beacon Pubkey (Compressed 48-byte Hex): \"%s\"\n",
				blsSigner.PublicKey().String(),
			)
			return nil
		},
	}

	return cmd
}
