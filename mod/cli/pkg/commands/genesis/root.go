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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type Genesis struct {
	AppState struct {
		Beacon struct {
			ForkVersion string `json:"fork_version"`
			Deposits    []struct {
				Pubkey      string `json:"pubkey"`
				Credentials string `json:"credentials"`
				Amount      string `json:"amount"`
				Signature   string `json:"signature"`
				Index       int    `json:"index"`
			} `json:"deposits"`
			ExecutionPayloadHeader struct {
				BlockHash string `json:"blockHash"`
			} `json:"execution_payload_header"`
		} `json:"beacon"`
	} `json:"app_state"`
}

func GetGenesisValidatorRootCmd() *cobra.Command {
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
				fmt.Println("Error unmarshaling JSON REEEEEE:", err)
				// return
			}

			depositCount := uint64(len(genesis.AppState.Beacon.Deposits))
			blockHashBytes := common.HexToHash(genesis.AppState.Beacon.ExecutionPayloadHeader.BlockHash)
			var eth1Data *types.Eth1Data
			eth1Data = eth1Data.New(
				bytes.B32{},
				math.U64(depositCount),
				blockHashBytes,
			)

			state := deneb.BeaconState{
				Eth1Data: eth1Data,
			}

			// fmt.Println("Genesis struct:", genesis)
			for _, deposit := range genesis.AppState.Beacon.Deposits {
				fmt.Println("Deposit:", deposit)
			}

			return nil

			// // Unmarshal the genesis file.
			// ethGenesis := &core.Genesis{}
			// if err = ethGenesis.UnmarshalJSON(genesisBz); err != nil {
			// 	return errors.Wrap(err, "failed to unmarshal eth1 genesis")
			// }
			// genesisBlock := ethGenesis.ToBlock()

			// // Create the execution payload.
			// payload := ethengineprimitives.BlockToExecutableData(
			// 	genesisBlock,
			// 	nil,
			// 	nil,
			// ).ExecutionPayload

			// serverCtx := server.GetServerContextFromCmd(cmd)
			// config := serverCtx.Config

			// appGenesis, err := genutiltypes.AppGenesisFromFile(
			// 	config.GenesisFile(),
			// )
			// if err != nil {
			// 	return errors.Wrap(err, "failed to read genesis doc from file")
			// }

			// // create the app state
			// appGenesisState, err := genutiltypes.GenesisStateFromAppGenesis(
			// 	appGenesis,
			// )
			// if err != nil {
			// 	return err
			// }

			// genesisInfo := &genesis.Genesis[
			// 	*types.Deposit, *types.ExecutionPayloadHeader,
			// ]{}

			// if err = json.Unmarshal(
			// 	appGenesisState["beacon"], genesisInfo,
			// ); err != nil {
			// 	return errors.Wrap(err, "failed to unmarshal beacon state")
			// }

			// // Inject the execution payload.
			// header, err := executableDataToExecutionPayloadHeader(
			// 	version.ToUint32(genesisInfo.ForkVersion),
			// 	payload,
			// )
			// if err != nil {
			// 	return errors.Wrap(
			// 		err,
			// 		"failed to convert executable data to execution payload header",
			// 	)
			// }
			// genesisInfo.ExecutionPayloadHeader = header

			// appGenesisState["beacon"], err = json.Marshal(genesisInfo)
			// if err != nil {
			// 	return errors.Wrap(err, "failed to marshal beacon state")
			// }

			// if appGenesis.AppState, err = json.MarshalIndent(
			// 	appGenesisState, "", "  ",
			// ); err != nil {
			// 	return err
			// }

			// return genutil.ExportGenesisFile(appGenesis, config.GenesisFile())
		},
	}

	return cmd

}
