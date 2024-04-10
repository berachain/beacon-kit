// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package genesis

import (
	"encoding/json"
	"unsafe"

	"github.com/berachain/beacon-kit/mod/core/state/deneb"
	enginetypes "github.com/berachain/beacon-kit/mod/execution/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/constants"
	"github.com/cockroachdb/errors"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"
	ethenginetypes "github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/core"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func AddExecutionPayloadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execution-payload [eth/genesis/file.json]",
		Short: "adds the eth1 genesis execution payload to the genesis file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Read the genesis file.
			genesisBz, err := afero.ReadFile(afero.NewOsFs(), args[0])
			if err != nil {
				return errors.Wrap(err, "failed to read eth1 genesis file")
			}

			// Unmarshal the genesis file.
			ethGenesis := &core.Genesis{}
			if err = ethGenesis.UnmarshalJSON(genesisBz); err != nil {
				return errors.Wrap(err, "failed to unmarshal eth1 genesis")
			}
			genesisBlock := ethGenesis.ToBlock()

			// Create the execution payload.
			payload := ethenginetypes.BlockToExecutableData(
				genesisBlock,
				nil,
				nil,
			).ExecutionPayload

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			genesis, err := types.AppGenesisFromFile(config.GenesisFile())
			if err != nil {
				return errors.Wrap(err, "failed to read genesis doc from file")
			}

			// create the app state
			appGenesisState, err := types.GenesisStateFromAppGenesis(genesis)
			if err != nil {
				return err
			}

			beaconState := &deneb.BeaconState{}
			if err = json.Unmarshal(
				appGenesisState["beacon"], beaconState,
			); err != nil {
				return errors.Wrap(err, "failed to unmarshal beacon state")
			}

			// Inject the execution payload.
			beaconState.LatestExecutionPayload = executableDataToExecutionPayload(
				payload,
			)

			appGenesisState["beacon"], err = json.Marshal(beaconState)
			if err != nil {
				return errors.Wrap(err, "failed to marshal beacon state")
			}

			if genesis.AppState, err = json.MarshalIndent(
				appGenesisState, "", "  ",
			); err != nil {
				return err
			}

			return genutil.ExportGenesisFile(genesis, config.GenesisFile())
		},
	}

	return cmd
}

// Converts the eth executable data type to the beacon execution payload
// interface.
func executableDataToExecutionPayload(
	data *ethenginetypes.ExecutableData,
) *enginetypes.ExecutableDataDeneb {
	withdrawals := make([]*primitives.Withdrawal, len(data.Withdrawals))
	for i, withdrawal := range data.Withdrawals {
		// #nosec:G103 // primitives.Withdrawals is data.Withdrawals with hard
		// types.
		withdrawals[i] = (*primitives.Withdrawal)(unsafe.Pointer(withdrawal))
	}

	if len(data.ExtraData) > constants.ExtraDataLength {
		data.ExtraData = data.ExtraData[:constants.ExtraDataLength]
	}

	var baseFeePerGas []byte
	if data.BaseFeePerGas != nil {
		baseFeePerGas = data.BaseFeePerGas.Bytes()
	}

	var blobGasUsed uint64
	if data.BlobGasUsed != nil {
		blobGasUsed = *data.BlobGasUsed
	}

	var excessBlobGas uint64
	if data.ExcessBlobGas != nil {
		excessBlobGas = *data.ExcessBlobGas
	}

	return &enginetypes.ExecutableDataDeneb{
		ParentHash:    data.ParentHash,
		FeeRecipient:  data.FeeRecipient,
		StateRoot:     data.StateRoot,
		ReceiptsRoot:  data.ReceiptsRoot,
		LogsBloom:     data.LogsBloom,
		Random:        data.Random,
		Number:        data.Number,
		GasLimit:      data.GasLimit,
		GasUsed:       data.GasUsed,
		Timestamp:     data.Timestamp,
		ExtraData:     data.ExtraData,
		BaseFeePerGas: baseFeePerGas,
		BlockHash:     data.BlockHash,
		Transactions:  data.Transactions,
		Withdrawals:   withdrawals,
		BlobGasUsed:   blobGasUsed,
		ExcessBlobGas: excessBlobGas,
	}
}
