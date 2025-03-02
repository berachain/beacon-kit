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
	"math/big"
	"unsafe"

	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	gethprimitives "github.com/berachain/beacon-kit/geth-primitives"
	libcommon "github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/testing/simulated/execution"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
)

func TransformSimulatedBlockToGethBlock(simBlock *execution.SimulatedBlock, txs []*gethtypes.Transaction, parentBeaconRoot libcommon.Root) *gethtypes.Block {
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
			BaseFee:          baseFeePerGas, // Simulation sets BaseFee to zero which is incorrect
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

func TransformWithdrawalsToGethWithdrawals(withdrawals engineprimitives.Withdrawals) gethtypes.Withdrawals {
	w := make([]*gethtypes.Withdrawal, len(withdrawals))
	for i, withdrawal := range withdrawals {
		gethWithdrawal := (*gethtypes.Withdrawal)(unsafe.Pointer(withdrawal))
		w[i] = gethWithdrawal
	}
	return w
}

// SplitTxs iterates over txs and returns two slices:
// 1. The transactions with their blob sidecars removed.
// 2. Any blob sidecars that were present.
func SplitTxs(txs []*gethtypes.Transaction) (txsWithoutSidecars []*gethtypes.Transaction, txSidecars []*gethtypes.BlobTxSidecar) {
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
