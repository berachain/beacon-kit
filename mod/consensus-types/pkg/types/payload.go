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

// ExecutionPayload represents an execution payload across
// all fork versions.
type ExecutionPayload struct {
	InnerExecutionPayload
}

// InnerExecutionPayload represents the inner execution payload.
type InnerExecutionPayload interface {
	executionPayloadBody
	GetTransactions() [][]byte
	GetWithdrawals() []*engineprimitives.Withdrawal
}

// Empty returns an empty ExecutionPayload for the given fork version.
func (e *ExecutionPayload) Empty(forkVersion uint32) *ExecutionPayload {
	e = new(ExecutionPayload)
	switch forkVersion {
	case version.Deneb:
		e.InnerExecutionPayload = &ExecutableDataDeneb{}
	default:
		panic("unknown fork version")
	}
	return e
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
	case version.Deneb:
		return &ExecutionPayloadHeader{
			InnerExecutionPayloadHeader: &ExecutionPayloadHeaderDeneb{
				ParentHash:       e.GetParentHash(),
				FeeRecipient:     e.GetFeeRecipient(),
				StateRoot:        e.GetStateRoot(),
				ReceiptsRoot:     e.GetReceiptsRoot(),
				LogsBloom:        e.GetLogsBloom(),
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
			},
		}, nil
	default:
		return nil, errors.New("unknown fork version")
	}
}

// ExecutableDataDeneb is the execution payload for Deneb.
//
//go:generate go run github.com/ferranbt/fastssz/sszgen -path payload.go -objs ExecutableDataDeneb -include ../../../primitives/pkg/common,../../../primitives/pkg/bytes,../../../engine-primitives/pkg/engine-primitives/withdrawal.go,../../../primitives/pkg/common,../../../primitives/pkg/math,../../../primitives/pkg/bytes,$GETH_PKG_INCLUDE/common,$GETH_PKG_INCLUDE/common/hexutil -output payload.ssz.go
//go:generate go run github.com/fjl/gencodec -type ExecutableDataDeneb -field-override executableDataDenebMarshaling -out payload.json.go
//nolint:lll
type ExecutableDataDeneb struct {
	ParentHash    gethprimitives.ExecutionHash    `json:"parentHash"    ssz-size:"32"  gencodec:"required"`
	FeeRecipient  gethprimitives.ExecutionAddress `json:"feeRecipient"  ssz-size:"20"  gencodec:"required"`
	StateRoot     common.Bytes32                  `json:"stateRoot"     ssz-size:"32"  gencodec:"required"`
	ReceiptsRoot  common.Bytes32                  `json:"receiptsRoot"  ssz-size:"32"  gencodec:"required"`
	LogsBloom     []byte                          `json:"logsBloom"     ssz-size:"256" gencodec:"required"`
	Random        common.Bytes32                  `json:"prevRandao"    ssz-size:"32"  gencodec:"required"`
	Number        math.U64                        `json:"blockNumber"                  gencodec:"required"`
	GasLimit      math.U64                        `json:"gasLimit"                     gencodec:"required"`
	GasUsed       math.U64                        `json:"gasUsed"                      gencodec:"required"`
	Timestamp     math.U64                        `json:"timestamp"                    gencodec:"required"`
	ExtraData     []byte                          `json:"extraData"                    gencodec:"required" ssz-max:"32"`
	BaseFeePerGas math.Wei                        `json:"baseFeePerGas" ssz-size:"32"  gencodec:"required"`
	BlockHash     gethprimitives.ExecutionHash    `json:"blockHash"     ssz-size:"32"  gencodec:"required"`
	Transactions  [][]byte                        `json:"transactions"  ssz-size:"?,?" gencodec:"required" ssz-max:"1048576,1073741824"`
	Withdrawals   []*engineprimitives.Withdrawal  `json:"withdrawals"                                      ssz-max:"16"`
	BlobGasUsed   math.U64                        `json:"blobGasUsed"`
	ExcessBlobGas math.U64                        `json:"excessBlobGas"`
}

// JSON type overrides for ExecutableDataDeneb.
type executableDataDenebMarshaling struct {
	ExtraData    bytes.Bytes
	LogsBloom    bytes.Bytes
	Transactions []bytes.Bytes
}

// Version returns the version of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) Version() uint32 {
	return version.Deneb
}

// IsNil checks if the ExecutableDataDeneb is nil.
func (d *ExecutableDataDeneb) IsNil() bool {
	return d == nil
}

// IsBlinded checks if the ExecutableDataDeneb is blinded.
func (d *ExecutableDataDeneb) IsBlinded() bool {
	return false
}

// GetParentHash returns the parent hash of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetParentHash() gethprimitives.ExecutionHash {
	return d.ParentHash
}

// GetFeeRecipient returns the fee recipient address of the ExecutableDataDeneb.
func (
	d *ExecutableDataDeneb,
) GetFeeRecipient() gethprimitives.ExecutionAddress {
	return d.FeeRecipient
}

// GetStateRoot returns the state root of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetStateRoot() common.Bytes32 {
	return d.StateRoot
}

// GetReceiptsRoot returns the receipts root of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetReceiptsRoot() common.Bytes32 {
	return d.ReceiptsRoot
}

// GetLogsBloom returns the logs bloom of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetLogsBloom() []byte {
	return d.LogsBloom
}

// GetPrevRandao returns the previous Randao value of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetPrevRandao() common.Bytes32 {
	return d.Random
}

// GetNumber returns the block number of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetNumber() math.U64 {
	return d.Number
}

// GetGasLimit returns the gas limit of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetGasLimit() math.U64 {
	return d.GasLimit
}

// GetGasUsed returns the gas used of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetGasUsed() math.U64 {
	return d.GasUsed
}

// GetTimestamp returns the timestamp of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetTimestamp() math.U64 {
	return d.Timestamp
}

// GetExtraData returns the extra data of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetExtraData() []byte {
	return d.ExtraData
}

// GetBaseFeePerGas returns the base fee per gas of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetBaseFeePerGas() math.Wei {
	return d.BaseFeePerGas
}

// GetBlockHash returns the block hash of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetBlockHash() gethprimitives.ExecutionHash {
	return d.BlockHash
}

// GetTransactions returns the transactions of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetTransactions() [][]byte {
	return d.Transactions
}

// GetWithdrawals returns the withdrawals of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetWithdrawals() []*engineprimitives.Withdrawal {
	return d.Withdrawals
}

// GetBlobGasUsed returns the blob gas used of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetBlobGasUsed() math.U64 {
	return d.BlobGasUsed
}

// GetExcessBlobGas returns the excess blob gas of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetExcessBlobGas() math.U64 {
	return d.ExcessBlobGas
}
