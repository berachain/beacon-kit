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
	"encoding/json"
	"math/big"

	"github.com/berachain/beacon-kit/mod/config/version"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/uint256"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common/hexutil"
	huint256 "github.com/holiman/uint256"
)

// ExecutionPayloadEnvelope is an interface for the execution payload envelope.
type ExecutionPayloadEnvelope interface {
	// GetExecutionPayload retrieves the execution payload associated with the
	// envelope.
	GetExecutionPayload() ExecutionPayload
	// GetValue returns the Wei value of the block within the execution payload
	// envelope.
	GetValue() primitives.Wei
	// GetBlobsBundle fetches the BlobsBundleV1 associated with the execution
	// payload, if available.
	GetBlobsBundle() *engine.BlobsBundleV1
	// ShouldOverrideBuilder indicates whether the builder should be overridden
	// in the execution environment.
	ShouldOverrideBuilder() bool
}

//go:generate go run github.com/fjl/gencodec -type ExecutionPayloadEnvelopeDeneb -field-override executionPayloadEnvelopeMarshaling -out payload_env.json.go
//nolint:lll
type ExecutionPayloadEnvelopeDeneb struct {
	ExecutionPayload *ExecutableDataDeneb  `json:"executionPayload"      gencodec:"required"`
	BlockValue       *big.Int              `json:"blockValue"            gencodec:"required"`
	BlobsBundle      *engine.BlobsBundleV1 `json:"blobsBundle"`
	Override         bool                  `json:"shouldOverrideBuilder"`
}

// GetExecutionPayload returns the execution payload of the
// ExecutionPayloadEnvelope.
func (e *ExecutionPayloadEnvelopeDeneb) GetExecutionPayload() ExecutionPayload {
	return e.ExecutionPayload
}

// GetValue returns the value of the ExecutionPayloadEnvelope.
func (e *ExecutionPayloadEnvelopeDeneb) GetValue() primitives.Wei {
	val, ok := huint256.FromBig(e.BlockValue)
	if !ok {
		return primitives.Wei{}
	}
	return primitives.Wei{Int: val}
}

// GetBlobsBundle returns the blobs bundle of the ExecutionPayloadEnvelope.
func (e *ExecutionPayloadEnvelopeDeneb) GetBlobsBundle() *engine.BlobsBundleV1 {
	return e.BlobsBundle
}

// ShouldOverrideBuilder returns whether the builder should be overridden.
func (e *ExecutionPayloadEnvelopeDeneb) ShouldOverrideBuilder() bool {
	return e.Override
}

// String returns the string representation of the ExecutionPayloadEnvelope.
func (e *ExecutionPayloadEnvelopeDeneb) String() string {
	return spew.Sdump(e)
}

// JSON type overrides for ExecutionPayloadEnvelope.
type executionPayloadEnvelopeMarshaling struct {
	BlockValue *hexutil.Big
}

//
//go:generate go run github.com/fjl/gencodec -type ExecutableDataDeneb -field-override executableDataDenebMarshaling -out payload.json.go
//go:generate go run github.com/ferranbt/fastssz/sszgen -path payload.go -objs ExecutableDataDeneb -include ../../primitives,../../primitives/uint256,$GETH_PKG_INCLUDE/common,$GOPATH/pkg/mod/github.com/holiman/uint256@v1.2.4 -output payload.ssz.go
//nolint:lll
type ExecutableDataDeneb struct {
	ParentHash    primitives.ExecutionHash    `json:"parentHash"    ssz-size:"32"  gencodec:"required"`
	FeeRecipient  primitives.ExecutionAddress `json:"feeRecipient"  ssz-size:"20"  gencodec:"required"`
	StateRoot     primitives.ExecutionHash    `json:"stateRoot"     ssz-size:"32"  gencodec:"required"`
	ReceiptsRoot  primitives.ExecutionHash    `json:"receiptsRoot"  ssz-size:"32"  gencodec:"required"`
	LogsBloom     []byte                      `json:"logsBloom"     ssz-size:"256" gencodec:"required"`
	Random        primitives.ExecutionHash    `json:"prevRandao"    ssz-size:"32"  gencodec:"required"`
	Number        uint64                      `json:"blockNumber"                  gencodec:"required"`
	GasLimit      uint64                      `json:"gasLimit"                     gencodec:"required"`
	GasUsed       uint64                      `json:"gasUsed"                      gencodec:"required"`
	Timestamp     uint64                      `json:"timestamp"                    gencodec:"required"`
	ExtraData     []byte                      `json:"extraData"                    gencodec:"required" ssz-max:"32"`
	BaseFeePerGas uint256.LittleEndian        `json:"baseFeePerGas" ssz-size:"32"  gencodec:"required"`
	BlockHash     primitives.ExecutionHash    `json:"blockHash"     ssz-size:"32"  gencodec:"required"`
	Transactions  [][]byte                    `json:"transactions"  ssz-size:"?,?" gencodec:"required" ssz-max:"1048576,1073741824"`
	Withdrawals   []*primitives.Withdrawal    `json:"withdrawals"                                      ssz-max:"16"`
	BlobGasUsed   uint64                      `json:"blobGasUsed"`
	ExcessBlobGas uint64                      `json:"excessBlobGas"`
}

// JSON type overrides for ExecutableDataDeneb.
type executableDataDenebMarshaling struct {
	Number        hexutil.Uint64
	GasLimit      hexutil.Uint64
	GasUsed       hexutil.Uint64
	Timestamp     hexutil.Uint64
	Random        primitives.ExecutionHash
	ExtraData     hexutil.Bytes
	LogsBloom     hexutil.Bytes
	Transactions  []hexutil.Bytes
	BlobGasUsed   hexutil.Uint64
	ExcessBlobGas hexutil.Uint64
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
func (d *ExecutableDataDeneb) GetParentHash() primitives.ExecutionHash {
	return d.ParentHash
}

// GetFeeRecipient returns the fee recipient address of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetFeeRecipient() primitives.ExecutionAddress {
	return d.FeeRecipient
}

// GetStateRoot returns the state root of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetStateRoot() primitives.ExecutionHash {
	return d.StateRoot
}

// GetReceiptsRoot returns the receipts root of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetReceiptsRoot() primitives.ExecutionHash {
	return d.ReceiptsRoot
}

// GetLogsBloom returns the logs bloom of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetLogsBloom() []byte {
	return d.LogsBloom
}

// GetPrevRandao returns the previous Randao value of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetPrevRandao() [32]byte {
	return d.Random
}

// GetNumber returns the block number of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetNumber() uint64 {
	return d.Number
}

// GetGasLimit returns the gas limit of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetGasLimit() uint64 {
	return d.GasLimit
}

// GetGasUsed returns the gas used of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetGasUsed() uint64 {
	return d.GasUsed
}

// GetTimestamp returns the timestamp of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetTimestamp() uint64 {
	return d.Timestamp
}

// GetExtraData returns the extra data of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetExtraData() []byte {
	return d.ExtraData
}

// GetBaseFeePerGas returns the base fee per gas of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetBaseFeePerGas() uint256.LittleEndian {
	return d.BaseFeePerGas
}

// GetBlockHash returns the block hash of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetBlockHash() primitives.ExecutionHash {
	return d.BlockHash
}

// GetTransactions returns the transactions of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetTransactions() [][]byte {
	return d.Transactions
}

// GetWithdrawals returns the withdrawals of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetWithdrawals() []*primitives.Withdrawal {
	return d.Withdrawals
}

// GetBlobGasUsed returns the blob gas used of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetBlobGasUsed() *uint64 {
	return &d.BlobGasUsed
}

// GetExcessBlobGas returns the excess blob gas of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetExcessBlobGas() *uint64 {
	return &d.ExcessBlobGas
}

// String returns the string representation of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) String() string {
	//#nosec:G703 // ignore potential marshalling failure.
	output, _ := json.Marshal(d)
	return string(output)
}
