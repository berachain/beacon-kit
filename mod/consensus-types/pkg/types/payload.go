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

package types

import (
	"context"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/errors"
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/merkle"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	"golang.org/x/sync/errgroup"
)

// ExecutionPayload is the execution payload for Deneb.
//
//nolint:lll
type ExecutionPayload struct {
	// ParentHash is the hash of the parent block.
	ParentHash gethprimitives.ExecutionHash `json:"parentHash"       gencodec:"required"`
	// FeeRecipient is the address of the fee recipient.
	FeeRecipient gethprimitives.ExecutionAddress `json:"feeRecipient"     gencodec:"required"`
	// StateRoot is the root of the state trie.
	StateRoot common.Bytes32 `json:"stateRoot"        gencodec:"required"`
	// ReceiptsRoot is the root of the receipts trie.
	ReceiptsRoot common.Bytes32 `json:"receiptsRoot"     gencodec:"required"`
	// LogsBloom is the bloom filter for the logs.
	LogsBloom bytes.B256 `json:"logsBloom"        gencodec:"required"`
	// Random is the prevRandao value.
	Random common.Bytes32 `json:"prevRandao"       gencodec:"required"`
	// Number is the block number.
	Number math.U64 `json:"blockNumber"      gencodec:"required"`
	// GasLimit is the gas limit for the block.
	GasLimit math.U64 `json:"gasLimit"         gencodec:"required"`
	// GasUsed is the amount of gas used in the block.
	GasUsed math.U64 `json:"gasUsed"          gencodec:"required"`
	// Timestamp is the timestamp of the block.
	Timestamp math.U64 `json:"timestamp"        gencodec:"required"`
	// ExtraData is the extra data of the block.
	ExtraData bytes.Bytes `json:"extraData"        gencodec:"required"`
	// BaseFeePerGas is the base fee per gas.
	BaseFeePerGas math.Wei `json:"baseFeePerGas"    gencodec:"required"`
	// BlockHash is the hash of the block.
	BlockHash gethprimitives.ExecutionHash `json:"blockHash"        gencodec:"required"`
	// Transactions is the list of transactions in the block.
	Transactions [][]byte `json:"transactions"  ssz-size:"?,?" gencodec:"required" ssz-max:"1048576,1073741824"`
	// Withdrawals is the list of withdrawals in the block.
	Withdrawals []*engineprimitives.Withdrawal `json:"withdrawals"                                      ssz-max:"16"`
	// BlobGasUsed is the amount of blob gas used in the block.
	BlobGasUsed math.U64 `json:"blobGasUsed"`
	// ExcessBlobGas is the amount of excess blob gas in the block.
	ExcessBlobGas math.U64 `json:"excessBlobGas"`
}

// Empty returns an empty ExecutionPayload for the given fork version.
func (e *ExecutionPayload) Empty(forkVersion uint32) *ExecutionPayload {
	e = new(ExecutionPayload)
	switch forkVersion {
	case version.Deneb, version.DenebPlus:
		e = &ExecutionPayload{}
	default:
		panic("unknown fork version")
	}
	return e
}

// Version returns the version of the ExecutionPayload.
func (d *ExecutionPayload) Version() uint32 {
	return version.Deneb
}

// IsNil checks if the ExecutionPayload is nil.
func (d *ExecutionPayload) IsNil() bool {
	return d == nil
}

// IsBlinded checks if the ExecutionPayload is blinded.
func (d *ExecutionPayload) IsBlinded() bool {
	return false
}

// GetParentHash returns the parent hash of the ExecutionPayload.
func (d *ExecutionPayload) GetParentHash() gethprimitives.ExecutionHash {
	return d.ParentHash
}

// GetFeeRecipient returns the fee recipient address of the ExecutionPayload.
func (
	d *ExecutionPayload,
) GetFeeRecipient() gethprimitives.ExecutionAddress {
	return d.FeeRecipient
}

// GetStateRoot returns the state root of the ExecutionPayload.
func (d *ExecutionPayload) GetStateRoot() common.Bytes32 {
	return d.StateRoot
}

// GetReceiptsRoot returns the receipts root of the ExecutionPayload.
func (d *ExecutionPayload) GetReceiptsRoot() common.Bytes32 {
	return d.ReceiptsRoot
}

// GetLogsBloom returns the logs bloom of the ExecutionPayload.
func (d *ExecutionPayload) GetLogsBloom() []byte {
	return d.LogsBloom[:]
}

// GetPrevRandao returns the previous Randao value of the ExecutionPayload.
func (d *ExecutionPayload) GetPrevRandao() common.Bytes32 {
	return d.Random
}

// GetNumber returns the block number of the ExecutionPayload.
func (d *ExecutionPayload) GetNumber() math.U64 {
	return d.Number
}

// GetGasLimit returns the gas limit of the ExecutionPayload.
func (d *ExecutionPayload) GetGasLimit() math.U64 {
	return d.GasLimit
}

// GetGasUsed returns the gas used of the ExecutionPayload.
func (d *ExecutionPayload) GetGasUsed() math.U64 {
	return d.GasUsed
}

// GetTimestamp returns the timestamp of the ExecutionPayload.
func (d *ExecutionPayload) GetTimestamp() math.U64 {
	return d.Timestamp
}

// GetExtraData returns the extra data of the ExecutionPayload.
func (d *ExecutionPayload) GetExtraData() []byte {
	return d.ExtraData
}

// GetBaseFeePerGas returns the base fee per gas of the ExecutionPayload.
func (d *ExecutionPayload) GetBaseFeePerGas() math.Wei {
	return d.BaseFeePerGas
}

// GetBlockHash returns the block hash of the ExecutionPayload.
func (d *ExecutionPayload) GetBlockHash() gethprimitives.ExecutionHash {
	return d.BlockHash
}

// GetTransactions returns the transactions of the ExecutionPayload.
func (d *ExecutionPayload) GetTransactions() [][]byte {
	return d.Transactions
}

// GetWithdrawals returns the withdrawals of the ExecutionPayload.
func (d *ExecutionPayload) GetWithdrawals() []*engineprimitives.Withdrawal {
	return d.Withdrawals
}

// GetBlobGasUsed returns the blob gas used of the ExecutionPayload.
func (d *ExecutionPayload) GetBlobGasUsed() math.U64 {
	return d.BlobGasUsed
}

// GetExcessBlobGas returns the excess blob gas of the ExecutionPayload.
func (d *ExecutionPayload) GetExcessBlobGas() math.U64 {
	return d.ExcessBlobGas
}

// ToHeader converts the ExecutionPayload to an ExecutionPayloadHeader.
func (e *ExecutionPayload) ToHeader(
	txsMerkleizer *merkle.Merkleizer[[32]byte, common.Root],
	maxWithdrawalsPerPayload uint64,
) (*ExecutionPayloadHeader, error) {
	// Get the merkle roots of transactions and withdrawals in parallel.
	var (
		g, _            = errgroup.WithContext(context.Background())
		txsRoot         common.Root
		withdrawalsRoot common.Root
	)

	g.Go(func() error {
		var txsRootErr error
		// TODO: This is live on bArtio with a bug and needs to be hardforked
		// off of.
		txsRoot, txsRootErr = engineprimitives.Transactions(
			e.GetTransactions(),
		).HashTreeRootWith(txsMerkleizer)
		return txsRootErr
	})

	g.Go(func() error {
		var withdrawalsRootErr error
		wds := ssz.ListFromElements(
			maxWithdrawalsPerPayload,
			e.GetWithdrawals()...)
		withdrawalsRoot, withdrawalsRootErr = wds.HashTreeRoot()
		return withdrawalsRootErr
	})

	// Wait for the goroutines to finish.
	if err := g.Wait(); err != nil {
		return nil, err
	}

	switch e.Version() {
	case version.Deneb, version.DenebPlus:
		return &ExecutionPayloadHeader{
			// version:          e.Version(),
			ParentHash:       e.GetParentHash(),
			FeeRecipient:     e.GetFeeRecipient(),
			StateRoot:        e.GetStateRoot(),
			ReceiptsRoot:     e.GetReceiptsRoot(),
			LogsBloom:        [256]byte(e.GetLogsBloom()),
			Random:           e.GetPrevRandao(),
			Number:           e.GetNumber(),
			GasLimit:         e.GetGasLimit(),
			GasUsed:          e.GetGasUsed(),
			Timestamp:        e.GetTimestamp(),
			ExtraData:        e.GetExtraData(),
			BaseFeePerGas:    e.GetBaseFeePerGas(),
			BlockHash:        e.GetBlockHash(),
			TransactionsRoot: txsRoot,
			WithdrawalsRoot:  withdrawalsRoot,
			BlobGasUsed:      e.GetBlobGasUsed(),
			ExcessBlobGas:    e.GetExcessBlobGas(),
		}, nil
	default:
		return nil, errors.New("unknown fork version")
	}
}
