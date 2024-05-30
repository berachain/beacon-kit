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

package types

import (
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
)

//

var _ engineprimitives.ExecutionPayloadHeader = (*ExecutionPayloadHeaderDeneb)(
	nil,
)

// ExecutionPayloadHeaderDeneb is the execution header payload of Deneb.
//
//go:generate go run github.com/fjl/gencodec -type ExecutionPayloadHeaderDeneb -out payload_header.json.go -field-override executionPayloadHeaderDenebMarshaling
//go:generate go run github.com/ferranbt/fastssz/sszgen -path payload_header.go -objs ExecutionPayloadHeaderDeneb -include ../../../primitives/pkg/bytes,../../../primitives/mod.go,../../../primitives/pkg/common,../../../primitives/pkg/math,$GETH_PKG_INCLUDE/common,$GETH_PKG_INCLUDE/common/hexutil,$GOPATH/pkg/mod/github.com/holiman/uint256@v1.2.4 -output payload_header.ssz.go
//nolint:lll
type ExecutionPayloadHeaderDeneb struct {
	ParentHash       common.ExecutionHash    `json:"parentHash"       ssz-size:"32"  gencodec:"required"`
	FeeRecipient     common.ExecutionAddress `json:"feeRecipient"     ssz-size:"20"  gencodec:"required"`
	StateRoot        primitives.Bytes32      `json:"stateRoot"        ssz-size:"32"  gencodec:"required"`
	ReceiptsRoot     primitives.Bytes32      `json:"receiptsRoot"     ssz-size:"32"  gencodec:"required"`
	LogsBloom        []byte                  `json:"logsBloom"        ssz-size:"256" gencodec:"required"`
	Random           primitives.Bytes32      `json:"prevRandao"       ssz-size:"32"  gencodec:"required"`
	Number           math.U64                `json:"blockNumber"                     gencodec:"required"`
	GasLimit         math.U64                `json:"gasLimit"                        gencodec:"required"`
	GasUsed          math.U64                `json:"gasUsed"                         gencodec:"required"`
	Timestamp        math.U64                `json:"timestamp"                       gencodec:"required"`
	ExtraData        []byte                  `json:"extraData"                       gencodec:"required" ssz-max:"32"`
	BaseFeePerGas    math.Wei                `json:"baseFeePerGas"    ssz-size:"32"  gencodec:"required"`
	BlockHash        common.ExecutionHash    `json:"blockHash"        ssz-size:"32"  gencodec:"required"`
	TransactionsRoot primitives.Root         `json:"transactionsRoot" ssz-size:"32"  gencodec:"required"`
	WithdrawalsRoot  primitives.Root         `json:"withdrawalsRoot"  ssz-size:"32"`
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
func (d *ExecutionPayloadHeaderDeneb) GetStateRoot() primitives.Bytes32 {
	return d.StateRoot
}

// GetReceiptsRoot returns the receipts root of the ExecutionPayloadHeaderDeneb.
func (d *ExecutionPayloadHeaderDeneb) GetReceiptsRoot() primitives.Bytes32 {
	return d.ReceiptsRoot
}

// GetLogsBloom returns the logs bloom of the ExecutionPayloadHeaderDeneb.
func (d *ExecutionPayloadHeaderDeneb) GetLogsBloom() []byte {
	return d.LogsBloom
}

// GetPrevRandao returns the previous Randao value of the
// ExecutionPayloadHeaderDeneb.
func (d *ExecutionPayloadHeaderDeneb) GetPrevRandao() primitives.Bytes32 {
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
func (d *ExecutionPayloadHeaderDeneb) GetTransactionsRoot() primitives.Root {
	return d.TransactionsRoot
}

// GetWithdrawalsRoot returns the withdrawals root of the
// ExecutionPayloadHeaderDeneb.
func (d *ExecutionPayloadHeaderDeneb) GetWithdrawalsRoot() primitives.Root {
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
