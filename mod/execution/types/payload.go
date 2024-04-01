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
	"fmt"
	"math/big"

	"github.com/berachain/beacon-kit/mod/forks/version"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/holiman/uint256"
)

type ExecutionPayloadEnvelope interface {
	GetExecutionPayload() ExecutionPayload
	GetValue() primitives.Wei
	GetBlobsBundle() *engine.BlobsBundleV1
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

func (e *ExecutionPayloadEnvelopeDeneb) GetExecutionPayload() ExecutionPayload {
	return e.ExecutionPayload
}

func (e *ExecutionPayloadEnvelopeDeneb) GetValue() primitives.Wei {
	val, ok := uint256.FromBig(e.BlockValue)
	if !ok {
		return primitives.Wei{}
	}
	return primitives.Wei{Int: val}
}

func (e *ExecutionPayloadEnvelopeDeneb) GetBlobsBundle() *engine.BlobsBundleV1 {
	return e.BlobsBundle
}

func (e *ExecutionPayloadEnvelopeDeneb) ShouldOverrideBuilder() bool {
	return e.Override
}

func (e *ExecutionPayloadEnvelopeDeneb) String() string {
	return fmt.Sprintf(`
ExecutionPayloadEnvelopeDeneb{
	ExecutionPayload: %s,
	BlockValue: %s,
	BlobsBundle: %s,
	Override: %v,
}`, e.ExecutionPayload.String(),
		e.BlockValue.String(),
		e.GetBlobsBundle().Blobs,
		e.Override,
	)
}

// JSON type overrides for ExecutionPayloadEnvelope.
type executionPayloadEnvelopeMarshaling struct {
	BlockValue *hexutil.Big
}

//go:generate go run github.com/fjl/gencodec -type ExecutableDataDeneb -field-override executableDataDenebMarshaling -out payload.json.go
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
	BaseFeePerGas []byte                      `json:"baseFeePerGas" ssz-size:"32"  gencodec:"required"`
	BlockHash     primitives.ExecutionHash    `json:"blockHash"     ssz-size:"32"  gencodec:"required"`
	Transactions  [][]byte                    `json:"transactions"  ssz-size:"?,?" gencodec:"required" ssz-max:"1048576,1073741824"`
	Withdrawals   []*primitives.Withdrawal    `json:"withdrawals"                                      ssz-max:"16"`
	BlobGasUsed   uint64                      `json:"blobGasUsed"`
	ExcessBlobGas uint64                      `json:"excessBlobGas"`
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

// GetBlockHash returns the block hash of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetBlockHash() primitives.ExecutionHash {
	return d.BlockHash
}

// GetTransactions returns the transactions of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetTransactions() [][]byte {
	return d.Transactions
}

// GetGasUsed returns the gas used of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetPrevRandao() [32]byte {
	return d.Random
}

// GetWithdrawals returns the withdrawals of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetWithdrawals() []*primitives.Withdrawal {
	return d.Withdrawals
}

// String returns the string representation of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) String() string {
	//#nosec:G703 // ignore potential marshalling failure.
	output, _ := json.Marshal(d)
	return string(output)
}
