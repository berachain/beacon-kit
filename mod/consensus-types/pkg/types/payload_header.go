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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
)

// ExecutionPayloadHeader represents an execution payload across
// all fork versions.
type ExecutionPayloadHeader struct {
	InnerExecutionPayloadHeader
}

// InnerExecutionPayloadHeader is the interface for the execution
// data of a block.
// It contains all the fields that are part of both an execution payload header
// and a full execution payload.
type InnerExecutionPayloadHeader interface {
	executionPayloadBody
	ssz.SSZTreeable
	GetTransactionsRoot() common.Root
	GetWithdrawalsRoot() common.Root
}

// Empty returns an empty ExecutionPayload for the given fork version.
func (e *ExecutionPayloadHeader) Empty(
	forkVersion uint32,
) *ExecutionPayloadHeader {
	e = new(ExecutionPayloadHeader)
	switch forkVersion {
	case version.Deneb:
		e.InnerExecutionPayloadHeader = &ExecutionPayloadHeaderDeneb{}
	default:
		panic(
			"unknown fork version, cannot create empty ExecutionPayloadHeader",
		)
	}
	return e
}

// NewFromSSZ returns a new ExecutionPayloadHeader from the given SSZ bytes.
func (e *ExecutionPayloadHeader) NewFromSSZ(
	bz []byte, forkVersion uint32,
) (*ExecutionPayloadHeader, error) {
	e = e.Empty(forkVersion)
	if err := e.UnmarshalSSZ(bz); err != nil {
		return nil, err
	}
	return e, nil
}

// NewFromJSON returns a new ExecutionPayloadHeader from the given JSON bytes.
func (e *ExecutionPayloadHeader) NewFromJSON(
	bz []byte, forkVersion uint32,
) (*ExecutionPayloadHeader, error) {
	e = e.Empty(forkVersion)
	if err := e.UnmarshalJSON(bz); err != nil {
		return nil, err
	}
	return e, nil
}

// ExecutionPayloadHeaderDeneb is the execution header payload of Deneb.
//
//go:generate go run github.com/fjl/gencodec -type ExecutionPayloadHeaderDeneb -out payload_header.json.go -field-override executionPayloadHeaderDenebMarshaling
//go:generate go run github.com/ferranbt/fastssz/sszgen -path payload_header.go -objs ExecutionPayloadHeaderDeneb -include ../../../primitives/pkg/bytes,../../../primitives/pkg/common,../../../primitives/pkg/math,$GETH_PKG_INCLUDE/common,$GETH_PKG_INCLUDE/common/hexutil -output payload_header.ssz.go
//nolint:lll
type ExecutionPayloadHeaderDeneb struct {
	ParentHash       common.ExecutionHash    `json:"parentHash"       ssz-size:"32"  gencodec:"required"`
	FeeRecipient     common.ExecutionAddress `json:"feeRecipient"     ssz-size:"20"  gencodec:"required"`
	StateRoot        common.Bytes32          `json:"stateRoot"        ssz-size:"32"  gencodec:"required"`
	ReceiptsRoot     common.Bytes32          `json:"receiptsRoot"     ssz-size:"32"  gencodec:"required"`
	LogsBloom        []byte                  `json:"logsBloom"        ssz-size:"256" gencodec:"required"`
	Random           common.Bytes32          `json:"prevRandao"       ssz-size:"32"  gencodec:"required"`
	Number           math.U64                `json:"blockNumber"                     gencodec:"required"`
	GasLimit         math.U64                `json:"gasLimit"                        gencodec:"required"`
	GasUsed          math.U64                `json:"gasUsed"                         gencodec:"required"`
	Timestamp        math.U64                `json:"timestamp"                       gencodec:"required"`
	ExtraData        []byte                  `json:"extraData"                       gencodec:"required" ssz-max:"32"`
	BaseFeePerGas    math.Wei                `json:"baseFeePerGas"    ssz-size:"32"  gencodec:"required"`
	BlockHash        common.ExecutionHash    `json:"blockHash"        ssz-size:"32"  gencodec:"required"`
	TransactionsRoot common.Root             `json:"transactionsRoot" ssz-size:"32"  gencodec:"required"`
	WithdrawalsRoot  common.Root             `json:"withdrawalsRoot"  ssz-size:"32"`
	BlobGasUsed      math.U64                `json:"blobGasUsed"`
	ExcessBlobGas    math.U64                `json:"excessBlobGas"`
}

type executionPayloadHeaderDenebMarshaling struct {
	ExtraData bytes.Bytes
	LogsBloom bytes.Bytes
}

// Version returns the version of the ExecutionPayloadHeaderDeneb.
func (d *ExecutionPayloadHeaderDeneb) Version() uint32 {
	return version.Deneb
}

// IsNil checks if the ExecutionPayloadHeaderDeneb is nil.
func (d *ExecutionPayloadHeaderDeneb) IsNil() bool {
	return d == nil
}

// IsBlinded checks if the ExecutionPayloadHeaderDeneb is blinded.
func (d *ExecutionPayloadHeaderDeneb) IsBlinded() bool {
	return false
}

// GetParentHash returns the parent hash of the ExecutionPayloadHeaderDeneb.
func (d *ExecutionPayloadHeaderDeneb) GetParentHash() common.ExecutionHash {
	return d.ParentHash
}

// GetFeeRecipient returns the fee recipient address of the
// ExecutionPayloadHeaderDeneb.
//
//nolint:lll // long variable names.
func (d *ExecutionPayloadHeaderDeneb) GetFeeRecipient() common.ExecutionAddress {
	return d.FeeRecipient
}

// GetStateRoot returns the state root of the ExecutionPayloadHeaderDeneb.
func (d *ExecutionPayloadHeaderDeneb) GetStateRoot() common.Bytes32 {
	return d.StateRoot
}

// GetReceiptsRoot returns the receipts root of the ExecutionPayloadHeaderDeneb.
func (d *ExecutionPayloadHeaderDeneb) GetReceiptsRoot() common.Bytes32 {
	return d.ReceiptsRoot
}

// GetLogsBloom returns the logs bloom of the ExecutionPayloadHeaderDeneb.
func (d *ExecutionPayloadHeaderDeneb) GetLogsBloom() []byte {
	return d.LogsBloom
}

// GetPrevRandao returns the previous Randao value of the
// ExecutionPayloadHeaderDeneb.
func (d *ExecutionPayloadHeaderDeneb) GetPrevRandao() common.Bytes32 {
	return d.Random
}

// GetNumber returns the block number of the ExecutionPayloadHeaderDeneb.
func (d *ExecutionPayloadHeaderDeneb) GetNumber() math.U64 {
	return d.Number
}

// GetGasLimit returns the gas limit of the ExecutionPayloadHeaderDeneb.
func (d *ExecutionPayloadHeaderDeneb) GetGasLimit() math.U64 {
	return d.GasLimit
}

// GetGasUsed returns the gas used of the ExecutionPayloadHeaderDeneb.
func (d *ExecutionPayloadHeaderDeneb) GetGasUsed() math.U64 {
	return d.GasUsed
}

// GetTimestamp returns the timestamp of the ExecutionPayloadHeaderDeneb.
func (d *ExecutionPayloadHeaderDeneb) GetTimestamp() math.U64 {
	return d.Timestamp
}

// GetExtraData returns the extra data of the ExecutionPayloadHeaderDeneb.
func (d *ExecutionPayloadHeaderDeneb) GetExtraData() []byte {
	return d.ExtraData
}

// GetBaseFeePerGas returns the base fee per gas of the
// ExecutionPayloadHeaderDeneb.
func (d *ExecutionPayloadHeaderDeneb) GetBaseFeePerGas() math.Wei {
	return d.BaseFeePerGas
}

// GetBlockHash returns the block hash of the ExecutionPayloadHeaderDeneb.
func (d *ExecutionPayloadHeaderDeneb) GetBlockHash() common.ExecutionHash {
	return d.BlockHash
}

// GetTransactionsRoot returns the transactions root of the
// ExecutionPayloadHeaderDeneb.
func (d *ExecutionPayloadHeaderDeneb) GetTransactionsRoot() common.Root {
	return d.TransactionsRoot
}

// GetWithdrawalsRoot returns the withdrawals root of the
// ExecutionPayloadHeaderDeneb.
func (d *ExecutionPayloadHeaderDeneb) GetWithdrawalsRoot() common.Root {
	return d.WithdrawalsRoot
}

// GetBlobGasUsed returns the blob gas used of the ExecutionPayloadHeaderDeneb.
func (d *ExecutionPayloadHeaderDeneb) GetBlobGasUsed() math.U64 {
	return d.BlobGasUsed
}

// GetExcessBlobGas returns the excess blob gas of the
// ExecutionPayloadHeaderDeneb.
func (d *ExecutionPayloadHeaderDeneb) GetExcessBlobGas() math.U64 {
	return d.ExcessBlobGas
}

func (d *ExecutionPayloadHeaderDeneb) GetRootNode() (*ssz.Node, error) {
	return ssz.NewTreeFromFastSSZ(d)
}
