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
	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/encoding/json"
	"github.com/berachain/beacon-kit/primitives/math"
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
			depositsSlice := make(types.DepositsSlice, depositCount)
			elems := make([]*types.Deposit, depositCount)

			for i, deposit := range genesis.AppState.Beacon.Deposits {
				sig := bytes.B96{}
				err := sig.UnmarshalText([]byte(deposit.Signature))
				if err != nil {
					return errors.Wrap(err, "failed to convert signature to bytes96")
				}

				elems[i] = &types.Deposit{
					Pubkey:      deposit.Pubkey,
					Credentials: types.WithdrawalCredentials(deposit.Credentials),
					Amount:      math.U64(deposit.Amount),
					Signature:   sig,
					Index:       uint64(deposit.Index),
				}
				depositsSlice[i] = &types.Deposit{
					Pubkey:      deposit.Pubkey,
					Credentials: types.WithdrawalCredentials(deposit.Credentials),
					Amount:      math.U64(deposit.Amount),
					Signature:   sig,
					Index:       uint64(deposit.Index),
				}
			}

			deposits := &types.Deposits{Elems: elems}
			depositsRoot, err := deposits.HashTreeRoot()
			if err != nil {
				return errors.Wrap(err, "failed to get hash tree root")
			}
			cmd.Printf("deposits fastssz root: %s\n", bytes.B32(depositsRoot))

			cmd.Printf("deposits karalabe root: %s\n", depositsSlice.HashTreeRoot())

			return nil
		},
	}

	return cmd
}
