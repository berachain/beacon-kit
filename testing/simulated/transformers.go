//go:build simulated

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

package simulated

import (
	"fmt"
	"math/big"
	"unsafe"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	gethprimitives "github.com/berachain/beacon-kit/geth-primitives"
	libcommon "github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/berachain/beacon-kit/testing/simulated/execution"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
)

// transformSimulatedBlockToGethBlock converts a simulated execution block into a Geth-style block.
// It uses the provided transactions and parent beacon root to construct a new execution block header.
func transformSimulatedBlockToGethBlock(
	simBlock *execution.SimulatedBlock,
	txs []*gethtypes.Transaction,
	parentBeaconRoot libcommon.Root,
) *gethtypes.Block {
	// Convert numeric fields.
	excessBlobGas := simBlock.ExcessBlobGas.ToInt().Uint64()
	blobGasUsed := simBlock.BlobGasUsed.ToInt().Uint64()
	baseFeePerGas := simBlock.BaseFeePerGas.ToInt()

	// Compute the withdrawals hash from the simulated block's withdrawals.
	withdrawalsHash := gethprimitives.DeriveSha(simBlock.Withdrawals, gethprimitives.NewStackTrie(nil))

	// Create a new header using values from the simulated block.
	header := &gethprimitives.Header{
		ParentHash: simBlock.ParentHash,
		UncleHash:  gethprimitives.EmptyUncleHash,
		Coinbase:   simBlock.Miner,
		Root:       simBlock.StateRoot,
		// TxHash is computed from the provided transactions since simulation does not have signatures
		// which is required for correct hash calculation.
		TxHash:           gethprimitives.DeriveSha(gethprimitives.Transactions(txs), gethprimitives.NewStackTrie(nil)),
		ReceiptHash:      simBlock.ReceiptsRoot,
		Bloom:            gethtypes.Bloom(simBlock.LogsBloom),
		Difficulty:       big.NewInt(0),
		Number:           (*big.Int)(simBlock.Number),
		GasLimit:         (uint64)(*simBlock.GasLimit),
		GasUsed:          (uint64)(*simBlock.GasUsed),
		Time:             (uint64)(*simBlock.Timestamp),
		BaseFee:          baseFeePerGas,
		Extra:            simBlock.ExtraData,
		MixDigest:        simBlock.MixHash,
		WithdrawalsHash:  &withdrawalsHash,
		ExcessBlobGas:    &excessBlobGas,
		BlobGasUsed:      &blobGasUsed,
		ParentBeaconRoot: (*common.Hash)(&parentBeaconRoot),
	}

	// Create the block body using the transactions and withdrawals from the simulation.
	body := gethprimitives.Body{
		Transactions: txs,
		Uncles:       nil,
		Withdrawals:  simBlock.Withdrawals,
	}

	return gethprimitives.NewBlockWithHeader(header).WithBody(body)
}

// transformExecutableDataToExecutionPayload converts Ethereum executable data into a beacon execution payload.
// This function supports fork versions prior to Deneb1. For unsupported fork versions, it returns an error.
func transformExecutableDataToExecutionPayload(
	forkVersion libcommon.Version,
	data *gethprimitives.ExecutableData,
) (*ctypes.ExecutionPayload, error) {
	// Check that the fork version is supported (pre-Deneb1).
	// TODO(pectra): Extended simulated test support for pectra
	if !version.IsBefore(forkVersion, version.Deneb1()) {
		return nil, ctypes.ErrForkVersionNotSupported
	}

	// Convert withdrawals
	withdrawals := *(*engineprimitives.Withdrawals)(unsafe.Pointer(&data.Withdrawals))

	// Truncate ExtraData if it exceeds the allowed length.
	if len(data.ExtraData) > constants.ExtraDataLength {
		data.ExtraData = data.ExtraData[:constants.ExtraDataLength]
	}

	// Safely dereference optional fields.
	var blobGasUsed, excessBlobGas uint64
	if data.BlobGasUsed != nil {
		blobGasUsed = *data.BlobGasUsed
	}
	if data.ExcessBlobGas != nil {
		excessBlobGas = *data.ExcessBlobGas
	}

	// Convert BaseFeePerGas into a U256 value.
	baseFeePerGas, err := math.NewU256FromBigInt(data.BaseFeePerGas)
	if err != nil {
		return nil, fmt.Errorf("failed baseFeePerGas conversion: %w", err)
	}

	// Construct the execution payload.
	executionPayload := &ctypes.ExecutionPayload{
		Versionable:   ctypes.NewVersionable(forkVersion),
		ParentHash:    libcommon.ExecutionHash(data.ParentHash),
		FeeRecipient:  libcommon.ExecutionAddress(data.FeeRecipient),
		StateRoot:     libcommon.Bytes32(data.StateRoot),
		ReceiptsRoot:  libcommon.Bytes32(data.ReceiptsRoot),
		LogsBloom:     [256]byte(data.LogsBloom),
		Random:        libcommon.Bytes32(data.Random),
		Number:        math.U64(data.Number),
		GasLimit:      math.U64(data.GasLimit),
		GasUsed:       math.U64(data.GasUsed),
		Timestamp:     math.U64(data.Timestamp),
		Withdrawals:   withdrawals,
		ExtraData:     data.ExtraData,
		BaseFeePerGas: baseFeePerGas,
		BlockHash:     libcommon.ExecutionHash(data.BlockHash),
		Transactions:  data.Transactions,
		BlobGasUsed:   math.U64(blobGasUsed),
		ExcessBlobGas: math.U64(excessBlobGas),
	}
	return executionPayload, nil
}

// splitTxs separates transactions into two slices:
// 1. Transactions with blob sidecars removed.
// 2. The extracted blob sidecars, if any.
func splitTxs(txs []*gethtypes.Transaction) (txsWithoutSidecars []*gethtypes.Transaction, txSidecars []*gethtypes.BlobTxSidecar) {
	txsWithoutSidecars = make([]*gethtypes.Transaction, 0, len(txs))
	txSidecars = make([]*gethtypes.BlobTxSidecar, 0, len(txs))

	for _, tx := range txs {
		// Append the transaction with its blob sidecar removed.
		txsWithoutSidecars = append(txsWithoutSidecars, tx.WithoutBlobTxSidecar())
		// If a blob sidecar exists, collect it.
		if sidecar := tx.BlobTxSidecar(); sidecar != nil {
			txSidecars = append(txSidecars, sidecar)
		}
	}
	return
}
