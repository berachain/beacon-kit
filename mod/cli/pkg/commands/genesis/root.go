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
	"fmt"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/state/deneb"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type Genesis struct {
	AppState struct {
		Beacon struct {
			ForkVersion bytes.B4 `json:"fork_version"`
			Deposits    []struct {
				Pubkey      bytes.B48 `json:"pubkey"`
				Credentials bytes.B32 `json:"credentials"`
				Amount      math.U64  `json:"amount"`
				Signature   string    `json:"signature"`
				Index       int       `json:"index"`
			} `json:"deposits"`
			ExecutionPayloadHeader types.ExecutionPayloadHeaderDeneb `json:"execution_payload_header"`
		} `json:"beacon"`
	} `json:"app_state"`
}

func GetGenesisValidatorRootCmd(cs primitives.ChainSpec) *cobra.Command {
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

			var fork *types.Fork
			fork = fork.New(
				genesis.AppState.Beacon.ForkVersion,
				genesis.AppState.Beacon.ForkVersion,
				math.U64(0),
			)

			var blkBody *types.BeaconBlockBody
			bodyRoot, err := blkBody.Empty(
				version.ToUint32(genesis.AppState.Beacon.ForkVersion)).HashTreeRoot()
			if err != nil {
				return errors.Wrap(err, "failed to get body root")
			}

			var blkHeader *types.BeaconBlockHeader
			blkHeader = blkHeader.New(
				0, 0, common.Root{}, common.Root{}, bodyRoot,
			)

			var eth1Data *types.Eth1Data
			depositCount := uint64(len(genesis.AppState.Beacon.Deposits))
			eth1BlockHash := genesis.AppState.Beacon.ExecutionPayloadHeader.BlockHash
			eth1Data = eth1Data.New(
				bytes.B32{},
				math.U64(depositCount),
				genesis.AppState.Beacon.ExecutionPayloadHeader.BlockHash,
			)

			var randaoMixes []primitives.Bytes32
			epochsPerHistoricalVector := cs.EpochsPerHistoricalVector()
			randaoMixes = make([]primitives.Bytes32, epochsPerHistoricalVector)
			for i := range randaoMixes {
				randaoMixes[i] = bytes.B32(eth1BlockHash)
			}

			validators := make([]types.Validator, depositCount)
			for i, deposit := range genesis.AppState.Beacon.Deposits {
				validators[i] = types.Validator{
					Pubkey:                     deposit.Pubkey,
					WithdrawalCredentials:      types.WithdrawalCredentials(deposit.Credentials),
					EffectiveBalance:           deposit.Amount,
					Slashed:                    false,
					ActivationEligibilityEpoch: math.U64(0),
					ActivationEpoch:            math.U64(0),
					ExitEpoch:                  math.U64(0),
					WithdrawableEpoch:          math.U64(0),
				}
			}

			st := deneb.BeaconState{
				Fork:                         fork,
				LatestBlockHeader:            blkHeader,
				Eth1Data:                     eth1Data,
				RandaoMixes:                  randaoMixes,
				LatestExecutionPayloadHeader: &genesis.AppState.Beacon.ExecutionPayloadHeader,
			}

			root, err := st.HashTreeRoot()
			if err != nil {
				return errors.Wrap(err, "failed to get hash tree root")
			}

			rootHex := fmt.Sprintf("%x", root)
			cmd.Printf("%s\n", rootHex)
			return nil
		},
	}

	return cmd

}
