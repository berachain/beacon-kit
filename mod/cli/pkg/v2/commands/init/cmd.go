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

package initcli

import (
	"encoding/json"
	"os"

	"github.com/berachain/beacon-kit/mod/cli/pkg/v2/flags"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/spf13/cobra"
)

func Command[
	ConsensusParamsT interface {
		json.Marshaler
		Default() ConsensusParamsT
	},
	GenesisStateT interface {
		json.Marshaler
		Default() GenesisStateT
	},
]() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize the beacon node",
		Long:  "Initialize the beacon node",
		RunE: func(cmd *cobra.Command, args []string) error {
			var cfg cmtcfg.Config

			// TODO: recovery mnemonic

			// TODO: intial height

			// TODO: consensus key

			// TODO: initialize node validator

			// Build genesis state
			genesisFilePath := cfg.GenesisFile()
			overwrite, err := cmd.Flags().GetBool(flags.FlagOverwrite)
			if err != nil {
				overwrite = false
			}
			// Check if the genesis file exists and we're not overwriting it
			if _, err := os.Stat(
				genesisFilePath,
			); !overwrite && !os.IsNotExist(err) {
				return ErrGenesisFileExists
			}

			// TODO: do we need more genesis data than state????
			var cp ConsensusParamsT
			var gs GenesisStateT
			genesis := &Genesis[ConsensusParamsT, GenesisStateT]{
				State:           gs.Default(),
				ConsensusParams: cp.Default(),
			}
			return genesis.Save(genesisFilePath)
		},
	}

	//nolint:lll // it's honestly more clear this way
	cmd.Flags().
		BoolP(flags.FlagOverwrite, "o", flags.DefaultOverwrite, flags.OverwriteDescription)
	cmd.Flags().
		BoolP(flags.FlagRecover, "r", flags.DefaultRecover, flags.RecoverDescription)

	return cmd
}
