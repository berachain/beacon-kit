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
	"unsafe"

	serverContext "github.com/berachain/beacon-kit/mod/cli/pkg/utils/context"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/errors"
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/json"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func AddExecutionPayloadCmd[
	LoggerT log.AdvancedLogger[any, LoggerT],
](chainSpec common.ChainSpec) *cobra.Command {
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
			ethGenesis := &gethprimitives.Genesis{}
			if err = ethGenesis.UnmarshalJSON(genesisBz); err != nil {
				return errors.Wrap(err, "failed to unmarshal eth1 genesis")
			}
			genesisBlock := ethGenesis.ToBlock()

			// Create the execution payload.
			payload := gethprimitives.BlockToExecutableData(
				genesisBlock,
				nil,
				nil,
			).ExecutionPayload

			serverCtx := serverContext.GetServerContextFromCmd[LoggerT](cmd)
			config := serverCtx.Config

			appGenesis, err := genutiltypes.AppGenesisFromFile(
				config.GenesisFile(),
			)
			if err != nil {
				return errors.Wrap(err, "failed to read genesis doc from file")
			}

			// create the app state
			appGenesisState, err := genutiltypes.GenesisStateFromAppGenesis(
				appGenesis,
			)
			if err != nil {
				return err
			}

			genesisInfo := &types.Genesis[
				*types.Deposit, *types.ExecutionPayloadHeader,
			]{}

			if err = json.Unmarshal(
				appGenesisState["beacon"], genesisInfo,
			); err != nil {
				return errors.Wrap(err, "failed to unmarshal beacon state")
			}

			// Inject the execution payload.
			header, err := executableDataToExecutionPayloadHeader(
				version.ToUint32(genesisInfo.ForkVersion),
				payload,
				chainSpec.MaxWithdrawalsPerPayload(),
			)
			if err != nil {
				return errors.Wrap(
					err,
					"failed to convert executable data to execution payload header",
				)
			}
			genesisInfo.ExecutionPayloadHeader = header

			appGenesisState["beacon"], err = json.Marshal(genesisInfo)
			if err != nil {
				return errors.Wrap(err, "failed to marshal beacon state")
			}

			if appGenesis.AppState, err = json.MarshalIndent(
				appGenesisState, "", "  ",
			); err != nil {
				return err
			}

			return genutil.ExportGenesisFile(appGenesis, config.GenesisFile())
		},
	}

	return cmd
}

// Converts the eth executable data type to the beacon execution payload header
// interface.
func executableDataToExecutionPayloadHeader(
	forkVersion uint32,
	data *gethprimitives.ExecutableData,
	// todo: re-enable when codec supports.
	_ uint64,
) (*types.ExecutionPayloadHeader, error) {
	var executionPayloadHeader *types.ExecutionPayloadHeader
	switch forkVersion {
	case version.Deneb, version.DenebPlus:
		withdrawals := make(
			[]*engineprimitives.Withdrawal,
			len(data.Withdrawals),
		)
		for i, withdrawal := range data.Withdrawals {
			// #nosec:G103 // primitives.Withdrawals is data.Withdrawals with
			// hard
			// types.
			withdrawals[i] = (*engineprimitives.Withdrawal)(
				unsafe.Pointer(withdrawal),
			)
		}

		if len(data.ExtraData) > constants.ExtraDataLength {
			data.ExtraData = data.ExtraData[:constants.ExtraDataLength]
		}

		var blobGasUsed uint64
		if data.BlobGasUsed != nil {
			blobGasUsed = *data.BlobGasUsed
		}

		var excessBlobGas uint64
		if data.ExcessBlobGas != nil {
			excessBlobGas = *data.ExcessBlobGas
		}

		executionPayloadHeader = &types.ExecutionPayloadHeader{
			ParentHash:    common.ExecutionHash(data.ParentHash),
			FeeRecipient:  common.ExecutionAddress(data.FeeRecipient),
			StateRoot:     common.Bytes32(data.StateRoot),
			ReceiptsRoot:  common.Bytes32(data.ReceiptsRoot),
			LogsBloom:     [256]byte(data.LogsBloom),
			Random:        common.Bytes32(data.Random),
			Number:        math.U64(data.Number),
			GasLimit:      math.U64(data.GasLimit),
			GasUsed:       math.U64(data.GasUsed),
			Timestamp:     math.U64(data.Timestamp),
			ExtraData:     data.ExtraData,
			BaseFeePerGas: math.NewU256FromBigInt(data.BaseFeePerGas),
			BlockHash:     common.ExecutionHash(data.BlockHash),
			// TODO: Decouple from broken bArtio.
			TransactionsRoot: engineprimitives.
				BartioTransactions(
					data.Transactions,
				).HashTreeRoot(),
			WithdrawalsRoot: engineprimitives.Withdrawals(withdrawals).
				HashTreeRoot(),
			BlobGasUsed:   math.U64(blobGasUsed),
			ExcessBlobGas: math.U64(excessBlobGas),
		}
	default:
		return nil, errors.Newf("unsupported fork version %d", forkVersion)
	}

	return executionPayloadHeader, nil
}
