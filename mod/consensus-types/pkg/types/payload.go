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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
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
	ParentHash gethprimitives.ExecutionHash `json:"parentHash"    gencodec:"required"`
	// FeeRecipient is the address of the fee recipient.
	FeeRecipient gethprimitives.ExecutionAddress `json:"feeRecipient"  gencodec:"required"`
	// StateRoot is the root of the state trie.
	StateRoot common.Bytes32 `json:"stateRoot"     gencodec:"required"`
	// ReceiptsRoot is the root of the receipts trie.
	ReceiptsRoot common.Bytes32 `json:"receiptsRoot"  gencodec:"required"`
	// LogsBloom is the bloom filter for the logs.
	LogsBloom bytes.B256 `json:"logsBloom"     gencodec:"required"`
	// Random is the prevRandao value.
	Random common.Bytes32 `json:"prevRandao"    gencodec:"required"`
	// Number is the block number.
	Number math.U64 `json:"blockNumber"   gencodec:"required"`
	// GasLimit is the gas limit for the block.
	GasLimit math.U64 `json:"gasLimit"      gencodec:"required"`
	// GasUsed is the amount of gas used in the block.
	GasUsed math.U64 `json:"gasUsed"       gencodec:"required"`
	// Timestamp is the timestamp of the block.
	Timestamp math.U64 `json:"timestamp"     gencodec:"required"`
	// ExtraData is the extra data of the block.
	ExtraData bytes.Bytes `json:"extraData"     gencodec:"required"`
	// BaseFeePerGas is the base fee per gas.
	BaseFeePerGas math.Wei `json:"baseFeePerGas" gencodec:"required"`
	// BlockHash is the hash of the block.
	BlockHash gethprimitives.ExecutionHash `json:"blockHash"     gencodec:"required"`
	// Transactions is the list of transactions in the block.
	Transactions [][]byte `json:"transactions"  gencodec:"required"`
	// Withdrawals is the list of withdrawals in the block.
	Withdrawals []*engineprimitives.Withdrawal `json:"withdrawals"`
	// BlobGasUsed is the amount of blob gas used in the block.
	BlobGasUsed math.U64 `json:"blobGasUsed"`
	// ExcessBlobGas is the amount of excess blob gas in the block.
	ExcessBlobGas math.U64 `json:"excessBlobGas"`
}

// Empty returns an empty ExecutionPayload for the given fork version.
func (p *ExecutionPayload) Empty(_ uint32) *ExecutionPayload {
	return &ExecutionPayload{}
}

// Version returns the version of the ExecutionPayload.
func (p *ExecutionPayload) Version() uint32 {
	return version.Deneb
}

// IsNil checks if the ExecutionPayload is nil.
func (p *ExecutionPayload) IsNil() bool {
	return p == nil
}

// IsBlinded checks if the ExecutionPayload is blinded.
func (p *ExecutionPayload) IsBlinded() bool {
	return false
}

// GetParentHash returns the parent hash of the ExecutionPayload.
func (p *ExecutionPayload) GetParentHash() gethprimitives.ExecutionHash {
	return p.ParentHash
}

// GetFeeRecipient returns the fee recipient address of the ExecutionPayload.
func (
	p *ExecutionPayload,
) GetFeeRecipient() gethprimitives.ExecutionAddress {
	return p.FeeRecipient
}

// GetStateRoot returns the state root of the ExecutionPayload.
func (p *ExecutionPayload) GetStateRoot() common.Bytes32 {
	return p.StateRoot
}

// GetReceiptsRoot returns the receipts root of the ExecutionPayload.
func (p *ExecutionPayload) GetReceiptsRoot() common.Bytes32 {
	return p.ReceiptsRoot
}

// GetLogsBloom returns the logs bloom of the ExecutionPayload.
func (p *ExecutionPayload) GetLogsBloom() []byte {
	return p.LogsBloom[:]
}

// GetPrevRandao returns the previous Randao value of the ExecutionPayload.
func (p *ExecutionPayload) GetPrevRandao() common.Bytes32 {
	return p.Random
}

// GetNumber returns the block number of the ExecutionPayload.
func (p *ExecutionPayload) GetNumber() math.U64 {
	return p.Number
}

// GetGasLimit returns the gas limit of the ExecutionPayload.
func (p *ExecutionPayload) GetGasLimit() math.U64 {
	return p.GasLimit
}

// GetGasUsed returns the gas used of the ExecutionPayload.
func (p *ExecutionPayload) GetGasUsed() math.U64 {
	return p.GasUsed
}

// GetTimestamp returns the timestamp of the ExecutionPayload.
func (p *ExecutionPayload) GetTimestamp() math.U64 {
	return p.Timestamp
}

// GetExtraData returns the extra data of the ExecutionPayload.
func (p *ExecutionPayload) GetExtraData() []byte {
	return p.ExtraData
}

// GetBaseFeePerGas returns the base fee per gas of the ExecutionPayload.
func (p *ExecutionPayload) GetBaseFeePerGas() math.Wei {
	return p.BaseFeePerGas
}

// GetBlockHash returns the block hash of the ExecutionPayload.
func (p *ExecutionPayload) GetBlockHash() gethprimitives.ExecutionHash {
	return p.BlockHash
}

// GetTransactions returns the transactions of the ExecutionPayload.
func (p *ExecutionPayload) GetTransactions() [][]byte {
	return p.Transactions
}

// GetWithdrawals returns the withdrawals of the ExecutionPayload.
func (p *ExecutionPayload) GetWithdrawals() []*engineprimitives.Withdrawal {
	return p.Withdrawals
}

// GetBlobGasUsed returns the blob gas used of the ExecutionPayload.
func (p *ExecutionPayload) GetBlobGasUsed() math.U64 {
	return p.BlobGasUsed
}

// GetExcessBlobGas returns the excess blob gas of the ExecutionPayload.
func (p *ExecutionPayload) GetExcessBlobGas() math.U64 {
	return p.ExcessBlobGas
}

// ToHeader converts the ExecutionPayload to an ExecutionPayloadHeader.
func (p *ExecutionPayload) ToHeader(
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
			p.GetTransactions(),
		).HashTreeRootWith(txsMerkleizer)
		return txsRootErr
	})

	g.Go(func() error {
		var withdrawalsRootErr error
		wds := ssz.ListFromElements(
			maxWithdrawalsPerPayload,
			p.GetWithdrawals()...)
		withdrawalsRoot, withdrawalsRootErr = wds.HashTreeRoot()
		return withdrawalsRootErr
	})

	// Wait for the goroutines to finish.
	if err := g.Wait(); err != nil {
		return nil, err
	}

	switch p.Version() {
	case version.Deneb, version.DenebPlus:
		return &ExecutionPayloadHeader{
			ParentHash:       p.GetParentHash(),
			FeeRecipient:     p.GetFeeRecipient(),
			StateRoot:        p.GetStateRoot(),
			ReceiptsRoot:     p.GetReceiptsRoot(),
			LogsBloom:        [256]byte(p.GetLogsBloom()),
			Random:           p.GetPrevRandao(),
			Number:           p.GetNumber(),
			GasLimit:         p.GetGasLimit(),
			GasUsed:          p.GetGasUsed(),
			Timestamp:        p.GetTimestamp(),
			ExtraData:        p.GetExtraData(),
			BaseFeePerGas:    p.GetBaseFeePerGas(),
			BlockHash:        p.GetBlockHash(),
			TransactionsRoot: txsRoot,
			WithdrawalsRoot:  withdrawalsRoot,
			BlobGasUsed:      p.GetBlobGasUsed(),
			ExcessBlobGas:    p.GetExcessBlobGas(),
		}, nil
	default:
		return nil, errors.New("unknown fork version")
	}
}
