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

package genesis

import (
	"encoding/json"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/merkle"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type Genesis struct {
	AppState struct {
		Beacon struct {
			Deposits []struct {
				Pubkey      bytes.B48 `json:"pubkey"`
				Credentials bytes.B32 `json:"credentials"`
				Amount      math.U64  `json:"amount"`
				Signature   string    `json:"signature"`
				Index       int       `json:"index"`
			} `json:"deposits"`
		} `json:"beacon"`
	} `json:"app_state"`
}

func GetGenesisValidatorRootCmd(cs common.ChainSpec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-root [beacond/genesis.json]",
		Short: "gets and returns the genesis validator root",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Read the genesis file.
			genesisBz, err := afero.ReadFile(afero.NewOsFs(), args[0])
			if err != nil {
				return errors.Wrap(err, "failed to genesis json file")
			}

			var genesis Genesis
			// Unmarshal JSON data into the Genesis struct
			err = json.Unmarshal(genesisBz, &genesis)
			if err != nil {
				return errors.Wrap(err, "failed to unmarshal JSON")
			}

			depositCount := uint64(len(genesis.AppState.Beacon.Deposits))
			validators := make(
				[]*types.Validator,
				depositCount,
			)
			for i, deposit := range genesis.AppState.Beacon.Deposits {
				var val *types.Validator
				validators[i] = val.New(
					deposit.Pubkey,
					types.WithdrawalCredentials(deposit.Credentials),
					deposit.Amount,
					math.Gwei(cs.EffectiveBalanceIncrement()),
					math.Gwei(cs.MaxEffectiveBalance()),
				)
			}

			var validatorsRoot common.Root
			validatorsRoot, err = merkle.NewMerkleizer[
				[32]byte, *types.Validator,
			]().MerkleizeListComposite(
				validators, uint64(len(validators)),
			)
			if err != nil {
				return errors.Wrap(err, "failed to get validators root")
			}

			cmd.Printf("%s\n", validatorsRoot)
			return nil
		},
	}

	return cmd
}
