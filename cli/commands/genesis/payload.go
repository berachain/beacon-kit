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

	servertypes "github.com/berachain/beacon-kit/cli/commands/server/types"
	"github.com/berachain/beacon-kit/cli/context"
	"github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/errors"
	gethprimitives "github.com/berachain/beacon-kit/geth-primitives"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/encoding/json"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func AddExecutionPayloadCmd(chainSpecCreator servertypes.ChainSpecCreator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execution-payload [eth/genesis/file.json]",
		Short: "adds the eth1 genesis execution payload to the genesis file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Read the genesis file.
			elGenesisPath := args[0]
			config := context.GetConfigFromCmd(cmd)
			v := context.GetViperFromCmd(cmd)
			chainSpec, err := chainSpecCreator(v)
			if err != nil {
				return err
			}
			return AddExecutionPayload(chainSpec, elGenesisPath, config)
		},
	}

	return cmd
}

func AddExecutionPayload(chainSpec ChainSpec, elGenesisPath string, config *cmtcfg.Config) error {
	genesisBz, err := afero.ReadFile(afero.NewOsFs(), elGenesisPath)
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
		nil,
	).ExecutionPayload

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

	genesisInfo := &types.Genesis{}

	if err = json.Unmarshal(
		appGenesisState["beacon"], genesisInfo,
	); err != nil {
		return errors.Wrap(err, "failed to unmarshal beacon state")
	}
	// Inject the execution payload.
	eph, err := executableDataToExecutionPayloadHeader(
		chainSpec.GenesisForkVersion(),
		payload,
		chainSpec.MaxWithdrawalsPerPayload(),
	)
	if err != nil {
		return errors.Wrap(err, "failed to convert executable data to execution payload header")
	}
	if eph == nil {
		return errors.New("failed to get execution payload header")
	}
	genesisInfo.ExecutionPayloadHeader = eph

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
}

// Converts the eth executable data type to the beacon execution payload header
// interface.
func executableDataToExecutionPayloadHeader(
	forkVersion common.Version,
	data *gethprimitives.ExecutableData,
	// todo: re-enable when codec supports.
	_ uint64,
) (*types.ExecutionPayloadHeader, error) {
	eph := &types.ExecutionPayloadHeader{}

	// We do not support fork versions before Deneb and after Electra.
	if version.IsAfter(forkVersion, version.Electra()) ||
		version.IsBefore(forkVersion, version.Deneb()) {
		return nil, types.ErrForkVersionNotSupported
	}

	withdrawals := make(
		engineprimitives.Withdrawals,
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

	baseFeePerGas, err := math.NewU256FromBigInt(data.BaseFeePerGas)
	if err != nil {
		return nil, fmt.Errorf("failed baseFeePerGas conversion: %w", err)
	}

	eph.Versionable = types.NewVersionable(forkVersion)
	eph.ParentHash = common.ExecutionHash(data.ParentHash)
	eph.FeeRecipient = common.ExecutionAddress(data.FeeRecipient)
	eph.StateRoot = common.Bytes32(data.StateRoot)
	eph.ReceiptsRoot = common.Bytes32(data.ReceiptsRoot)
	eph.LogsBloom = [256]byte(data.LogsBloom)
	eph.Random = common.Bytes32(data.Random)
	eph.Number = math.U64(data.Number)
	eph.GasLimit = math.U64(data.GasLimit)
	eph.GasUsed = math.U64(data.GasUsed)
	eph.Timestamp = math.U64(data.Timestamp)
	eph.ExtraData = data.ExtraData
	eph.BaseFeePerGas = baseFeePerGas
	eph.BlockHash = common.ExecutionHash(data.BlockHash)
	eph.TransactionsRoot = engineprimitives.Transactions(data.Transactions).HashTreeRoot()
	eph.WithdrawalsRoot = withdrawals.HashTreeRoot()
	eph.BlobGasUsed = math.U64(blobGasUsed)
	eph.ExcessBlobGas = math.U64(excessBlobGas)

	return eph, nil
}
