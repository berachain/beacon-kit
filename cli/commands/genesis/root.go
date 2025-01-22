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

package genesis

import (
	"github.com/berachain/beacon-kit/chain"
	"github.com/berachain/beacon-kit/cli/utils/genesis"
	"github.com/spf13/cobra"
)

// GetGenesisValidatorRootCmd returns a cobra command that computes the genesis validator root
// from the given genesis file.
func GetGenesisValidatorRootCmd(cs chain.Spec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-root [beacond/genesis.json]",
		Short: "gets and returns the genesis validator root",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			genesisValidatorsRoot, err := genesis.ComputeValidatorsRootFromFile(args[0], cs)
			if err != nil {
				return err
			}

			cmd.Printf("%s\n", genesisValidatorsRoot)
			return nil
		},
	}

	return cmd
}
