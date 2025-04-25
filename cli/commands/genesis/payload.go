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
	"fmt"
	"unsafe"

	"github.com/berachain/beacon-kit/cli/context"
	"github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/errors"
	gethprimitives "github.com/berachain/beacon-kit/geth-primitives"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/encoding/json"
	"github.com/berachain/beacon-kit/primitives/math"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func AddExecutionPayloadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execution-payload [eth/genesis/file.json]",
		Short: "adds the eth1 genesis execution payload to the genesis file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return AddExecutionPayload(args[0], context.GetConfigFromCmd(cmd))
		},
	}

	return cmd
}

func AddExecutionPayload(elGenesisPath string, config *cmtcfg.Config) error {
	genesisBz, err := afero.ReadFile(afero.NewOsFs(), elGenesisPath)
	if err != nil {
		return errors.Wrap(err, "failed to read eth1 genesis file")
	}

	// Unmarshal the EL genesis file.
	ethGenesis := &gethprimitives.Genesis{}
	if err = ethGenesis.UnmarshalJSON(genesisBz); err != nil {
		return errors.Wrap(err, "failed to unmarshal eth1 genesis")
	}
	genesisBlock := ethGenesis.ToBlock()

	// Create the executable data from the EL genesis file.
	payload := gethprimitives.BlockToExecutableData(
		genesisBlock,
		nil,
		nil,
		nil,
	).ExecutionPayload

	appGenesis, err := genutiltypes.AppGenesisFromFile(config.GenesisFile())
	if err != nil {
		return errors.Wrap(err, "failed to read genesis doc from file")
	}

	// Create the app state.
	appGenesisState, err := genutiltypes.GenesisStateFromAppGenesis(appGenesis)
	if err != nil {
		return err
	}

	genesisInfo := &types.Genesis{}
	if err = json.Unmarshal(appGenesisState["beacon"], genesisInfo); err != nil {
		return errors.Wrap(err, "failed to unmarshal beacon state")
	}

	// Inject the execution payload headerfrom the executable data.
	eph, err := executableDataToExecutionPayloadHeader(payload)
	if err != nil {
		return errors.Wrap(err, "failed to convert executable data to execution payload header")
	}
	genesisInfo.ExecutionPayloadHeader = eph

	if appGenesisState["beacon"], err = json.Marshal(genesisInfo); err != nil {
		return errors.Wrap(err, "failed to marshal beacon state")
	}

	if appGenesis.AppState, err = json.MarshalIndent(appGenesisState, "", "  "); err != nil {
		return err
	}

	return genutil.ExportGenesisFile(appGenesis, config.GenesisFile())
}

// Converts the eth executable data type to the beacon execution payload header.
func executableDataToExecutionPayloadHeader(
	data *gethprimitives.ExecutableData,
) (*types.ExecutionPayloadHeader, error) {
	eph := &types.ExecutionPayloadHeader{
		ParentHash:       common.ExecutionHash(data.ParentHash),
		FeeRecipient:     common.ExecutionAddress(data.FeeRecipient),
		StateRoot:        common.Bytes32(data.StateRoot),
		ReceiptsRoot:     common.Bytes32(data.ReceiptsRoot),
		LogsBloom:        [256]byte(data.LogsBloom),
		Random:           common.Bytes32(data.Random),
		Number:           math.U64(data.Number),
		GasLimit:         math.U64(data.GasLimit),
		GasUsed:          math.U64(data.GasUsed),
		Timestamp:        math.U64(data.Timestamp),
		ExtraData:        data.ExtraData,
		BlockHash:        common.ExecutionHash(data.BlockHash),
		TransactionsRoot: engineprimitives.Transactions(data.Transactions).HashTreeRoot(),
	}

	// #nosec:G103 // engineprimitives.Withdrawals is data.Withdrawals with hard types.
	withdrawals := *(*engineprimitives.Withdrawals)(unsafe.Pointer(&data.Withdrawals))
	eph.WithdrawalsRoot = withdrawals.HashTreeRoot()

	if len(data.ExtraData) > constants.ExtraDataLength {
		data.ExtraData = data.ExtraData[:constants.ExtraDataLength]
	}
	eph.ExtraData = data.ExtraData

	if data.BlobGasUsed != nil {
		eph.BlobGasUsed = math.U64(*data.BlobGasUsed)
	}

	if data.ExcessBlobGas != nil {
		eph.ExcessBlobGas = math.U64(*data.ExcessBlobGas)
	}

	baseFeePerGas, err := math.NewU256FromBigInt(data.BaseFeePerGas)
	if err != nil {
		return nil, fmt.Errorf("failed baseFeePerGas conversion: %w", err)
	}
	eph.BaseFeePerGas = baseFeePerGas

	return eph, nil
}
