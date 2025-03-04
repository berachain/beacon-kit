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
	mathpkg "github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/berachain/beacon-kit/testing/simulated/execution"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
)

// transformWithdrawalsToGethWithdrawals converts a SimulatedBlock to a geth Block
func transformSimulatedBlockToGethBlock(simBlock *execution.SimulatedBlock, txs []*gethtypes.Transaction, parentBeaconRoot libcommon.Root) *gethtypes.Block {
	// Construct a new execution block header with the provided transactions.
	excessBlobGas := simBlock.ExcessBlobGas.ToInt().Uint64()
	blobGasUsed := simBlock.BlobGasUsed.ToInt().Uint64()
	baseFeePerGas := simBlock.BaseFeePerGas.ToInt()
	withdrawalsHash := gethprimitives.DeriveSha(
		simBlock.Withdrawals,
		gethprimitives.NewStackTrie(nil),
	)
	executionBlock := gethprimitives.NewBlockWithHeader(
		&gethprimitives.Header{
			ParentHash: simBlock.ParentHash,
			UncleHash:  gethprimitives.EmptyUncleHash,
			Coinbase:   simBlock.Miner,
			Root:       simBlock.StateRoot,
			// We cannot use the receipts from the simulation as the simulation does not have access to signature, resulting
			// in incorrect transaction hash calculation.
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
		},
	).WithBody(gethprimitives.Body{
		Transactions: txs,
		Uncles:       nil,
		Withdrawals:  simBlock.Withdrawals,
	})
	return executionBlock
}

// transformWithdrawalsToGethWithdrawals converts beaconkit withdrawals to geth withdrawals
func transformWithdrawalsToGethWithdrawals(withdrawals engineprimitives.Withdrawals) gethtypes.Withdrawals {
	// This one-line conversion is valid only if the underlying types have the same memory layout.
	return *(*gethtypes.Withdrawals)(unsafe.Pointer(&withdrawals))
}

// transformExecutableDataToExecutionPayload converts Ethereum executable data to a beacon execution payload.
// It supports fork versions before Deneb1 and returns an error if the fork version is not supported.
func transformExecutableDataToExecutionPayload(
	forkVersion libcommon.Version,
	data *gethprimitives.ExecutableData,
) (*ctypes.ExecutionPayload, error) {
	// Only support fork versions before Deneb1.
	if version.IsBefore(forkVersion, version.Deneb1()) {
		// Convert withdrawals from gethprimitives to engineprimitives.
		withdrawals := make(engineprimitives.Withdrawals, len(data.Withdrawals))
		for i, withdrawal := range data.Withdrawals {
			// #nosec:G103 -- safe conversion assuming the underlying types are compatible.
			withdrawals[i] = (*engineprimitives.Withdrawal)(unsafe.Pointer(withdrawal))
		}

		// Truncate ExtraData if it exceeds the allowed length.
		if len(data.ExtraData) > constants.ExtraDataLength {
			data.ExtraData = data.ExtraData[:constants.ExtraDataLength]
		}

		// Dereference optional fields safely.
		var blobGasUsed, excessBlobGas uint64
		if data.BlobGasUsed != nil {
			blobGasUsed = *data.BlobGasUsed
		}
		if data.ExcessBlobGas != nil {
			excessBlobGas = *data.ExcessBlobGas
		}

		// Convert BaseFeePerGas into a U256 value.
		baseFeePerGas, err := mathpkg.NewU256FromBigInt(data.BaseFeePerGas)
		if err != nil {
			return nil, fmt.Errorf("failed baseFeePerGas conversion: %w", err)
		}

		executionPayload := &ctypes.ExecutionPayload{
			ParentHash:    libcommon.ExecutionHash(data.ParentHash),
			FeeRecipient:  libcommon.ExecutionAddress(data.FeeRecipient),
			StateRoot:     libcommon.Bytes32(data.StateRoot),
			ReceiptsRoot:  libcommon.Bytes32(data.ReceiptsRoot),
			LogsBloom:     [256]byte(data.LogsBloom),
			Random:        libcommon.Bytes32(data.Random),
			Number:        mathpkg.U64(data.Number),
			GasLimit:      mathpkg.U64(data.GasLimit),
			GasUsed:       mathpkg.U64(data.GasUsed),
			Timestamp:     mathpkg.U64(data.Timestamp),
			Withdrawals:   withdrawals,
			ExtraData:     data.ExtraData,
			BaseFeePerGas: baseFeePerGas,
			BlockHash:     libcommon.ExecutionHash(data.BlockHash),
			Transactions:  data.Transactions,
			BlobGasUsed:   mathpkg.U64(blobGasUsed),
			ExcessBlobGas: mathpkg.U64(excessBlobGas),
			EpVersion:     forkVersion,
		}
		return executionPayload, nil
	}
	return nil, ctypes.ErrForkVersionNotSupported
}

// splitTxs iterates over txs and returns two slices:
// 1. The transactions with their blob sidecars removed.
// 2. Any blob sidecars that were present.
func splitTxs(txs []*gethtypes.Transaction) (txsWithoutSidecars []*gethtypes.Transaction, txSidecars []*gethtypes.BlobTxSidecar) {
	// Preallocate slices (optional but can improve performance)
	txsWithoutSidecars = make([]*gethtypes.Transaction, 0, len(txs))
	txSidecars = make([]*gethtypes.BlobTxSidecar, 0, len(txs))
	for _, tx := range txs {
		// Append the transformed transaction.
		txsWithoutSidecars = append(txsWithoutSidecars, tx.WithoutBlobTxSidecar())
		// If there's a blob sidecar, append it to the sidecars slice.
		if sidecar := tx.BlobTxSidecar(); sidecar != nil {
			txSidecars = append(txSidecars, sidecar)
		}
	}
	return
}
