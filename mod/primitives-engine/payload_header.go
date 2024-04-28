// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package engineprimitives

import (
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/math"
	"github.com/berachain/beacon-kit/mod/primitives/version"
)

var _ ExecutionPayloadHeader = (*ExecutionHeaderDeneb)(nil)

// ExecutionHeaderDeneb is the execution header payload of Deneb.
//
//go:generate go run github.com/ferranbt/fastssz/sszgen -path payload_header.go -objs ExecutionHeaderDeneb -include ../primitives,../primitives/math,./withdrawal.go,$GETH_PKG_INCLUDE/common,$GETH_PKG_INCLUDE/common/hexutil,$GOPATH/pkg/mod/github.com/holiman/uint256@v1.2.4 -output payload_header.ssz.go
//nolint:lll
type ExecutionHeaderDeneb struct {
	ParentHash       primitives.ExecutionHash    `json:"parentHash"        ssz-size:"32"  gencodec:"required"`
	FeeRecipient     primitives.ExecutionAddress `json:"feeRecipient"      ssz-size:"20"  gencodec:"required"`
	StateRoot        primitives.Bytes32          `json:"stateRoot"         ssz-size:"32"  gencodec:"required"`
	ReceiptsRoot     primitives.Bytes32          `json:"receiptsRoot"      ssz-size:"32"  gencodec:"required"`
	LogsBloom        []byte                      `json:"logsBloom"         ssz-size:"256" gencodec:"required"`
	Random           primitives.Bytes32          `json:"prevRandao"        ssz-size:"32"  gencodec:"required"`
	Number           math.U64                    `json:"blockNumber"                      gencodec:"required"`
	GasLimit         math.U64                    `json:"gasLimit"                         gencodec:"required"`
	GasUsed          math.U64                    `json:"gasUsed"                          gencodec:"required"`
	Timestamp        math.U64                    `json:"timestamp"                        gencodec:"required"`
	ExtraData        []byte                      `json:"extraData"                        gencodec:"required" ssz-max:"32"`
	BaseFeePerGas    math.Wei                    `json:"baseFeePerGas"     ssz-size:"32"  gencodec:"required"`
	BlockHash        primitives.ExecutionHash    `json:"blockHash"         ssz-size:"32"  gencodec:"required"`
	TransactionsRoot primitives.Root             `json:"transactionsRoot"  ssz-size:"32"  gencodec:"required"`
	WithdrawalsRoot  primitives.Root             `json:"withdrawalsRoot"   ssz-size:"32"`
	BlobGasUsed      math.U64                    `json:"blobGasUsed"`
	ExcessBlobGas    math.U64                    `json:"excessBlobGas"`
}

// Version returns the version of the ExecutionHeaderDeneb.
func (d *ExecutionHeaderDeneb) Version() uint32 {
	return version.Deneb
}

// IsNil checks if the ExecutionHeaderDeneb is nil.
func (d *ExecutionHeaderDeneb) IsNil() bool {
	return d == nil
}

// IsBlinded checks if the ExecutionHeaderDeneb is blinded.
func (d *ExecutionHeaderDeneb) IsBlinded() bool {
	return false
}

// GetParentHash returns the parent hash of the ExecutionHeaderDeneb.
func (d *ExecutionHeaderDeneb) GetParentHash() primitives.ExecutionHash {
	return d.ParentHash
}

// GetFeeRecipient returns the fee recipient address of the
// ExecutionHeaderDeneb.
func (d *ExecutionHeaderDeneb) GetFeeRecipient() primitives.ExecutionAddress {
	return d.FeeRecipient
}

// GetStateRoot returns the state root of the ExecutionHeaderDeneb.
func (d *ExecutionHeaderDeneb) GetStateRoot() primitives.Bytes32 {
	return d.StateRoot
}

// GetReceiptsRoot returns the receipts root of the ExecutionHeaderDeneb.
func (d *ExecutionHeaderDeneb) GetReceiptsRoot() primitives.Bytes32 {
	return d.ReceiptsRoot
}

// GetLogsBloom returns the logs bloom of the ExecutionHeaderDeneb.
func (d *ExecutionHeaderDeneb) GetLogsBloom() []byte {
	return d.LogsBloom
}

// GetPrevRandao returns the previous Randao value of the ExecutionHeaderDeneb.
func (d *ExecutionHeaderDeneb) GetPrevRandao() primitives.Bytes32 {
	return d.Random
}

// GetNumber returns the block number of the ExecutionHeaderDeneb.
func (d *ExecutionHeaderDeneb) GetNumber() math.U64 {
	return d.Number
}

// GetGasLimit returns the gas limit of the ExecutionHeaderDeneb.
func (d *ExecutionHeaderDeneb) GetGasLimit() math.U64 {
	return d.GasLimit
}

// GetGasUsed returns the gas used of the ExecutionHeaderDeneb.
func (d *ExecutionHeaderDeneb) GetGasUsed() math.U64 {
	return d.GasUsed
}

// GetTimestamp returns the timestamp of the ExecutionHeaderDeneb.
func (d *ExecutionHeaderDeneb) GetTimestamp() math.U64 {
	return d.Timestamp
}

// GetExtraData returns the extra data of the ExecutionHeaderDeneb.
func (d *ExecutionHeaderDeneb) GetExtraData() []byte {
	return d.ExtraData
}

// GetBaseFeePerGas returns the base fee per gas of the ExecutionHeaderDeneb.
func (d *ExecutionHeaderDeneb) GetBaseFeePerGas() math.Wei {
	return d.BaseFeePerGas
}

// GetBlockHash returns the block hash of the ExecutionHeaderDeneb.
func (d *ExecutionHeaderDeneb) GetBlockHash() primitives.ExecutionHash {
	return d.BlockHash
}

// GetTransactionsRoot returns the transactions root of the
// ExecutionHeaderDeneb.
func (d *ExecutionHeaderDeneb) GetTransactionsRoot() primitives.Root {
	return d.TransactionsRoot
}

// GetWithdrawalsRoot returns the withdrawals root of the ExecutionHeaderDeneb.
func (d *ExecutionHeaderDeneb) GetWithdrawalsRoot() primitives.Root {
	return d.WithdrawalsRoot
}

// GetBlobGasUsed returns the blob gas used of the ExecutionHeaderDeneb.
func (d *ExecutionHeaderDeneb) GetBlobGasUsed() math.U64 {
	return d.BlobGasUsed
}

// GetExcessBlobGas returns the excess blob gas of the ExecutionHeaderDeneb.
func (d *ExecutionHeaderDeneb) GetExcessBlobGas() math.U64 {
	return d.ExcessBlobGas
}
