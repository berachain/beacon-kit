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

package enginetypes

import (
	"github.com/berachain/beacon-kit/config/version"
	"github.com/berachain/beacon-kit/primitives"
	"github.com/cockroachdb/errors"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ssz "github.com/prysmaticlabs/fastssz"
)

var (
	_ ssz.Marshaler   = (*ExecutableData)(nil)
	_ ssz.Unmarshaler = (*ExecutableData)(nil)
	_ ssz.HashRoot    = (*ExecutableData)(nil)
)

//go:generate go run github.com/fjl/gencodec -type ExecutableData -field-override executableDataMarshaling -out payload.json.go
//nolint:lll
type ExecutableData struct {
	version       int
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
	Withdrawals   []*Withdrawal               `json:"withdrawals"                                      ssz-max:"16"`
	BlobGasUsed   uint64                      `json:"blobGasUsed"`
	ExcessBlobGas uint64                      `json:"excessBlobGas"`
}

// JSON type overrides for ExecutableDataDeneb.
type executableDataMarshaling struct {
	Number        hexutil.Uint64
	GasLimit      hexutil.Uint64
	GasUsed       hexutil.Uint64
	Timestamp     hexutil.Uint64
	BaseFeePerGas primitives.SSZUInt256
	Random        primitives.ExecutionHash
	ExtraData     hexutil.Bytes
	LogsBloom     hexutil.Bytes
	Transactions  []hexutil.Bytes
	BlobGasUsed   hexutil.Uint64
	ExcessBlobGas hexutil.Uint64
}

// ExecutableDataDeneb is the ExecutableDataDeneb.
func (d *ExecutableData) SetVersion(v int) {
	if v == 0 {
		d.version = v
	}
}

// Version returns the version of the ExecutableDataDeneb.
func (d *ExecutableData) Version() int {
	return version.Deneb
}

// SizeSSZ returns the SSZ size of the ExecutableData. It varies based on the
// version of the data.
func (d *ExecutableData) SizeSSZ() int {
	//nolint:gocritic // future versions needed.
	switch d.version {
	case version.Deneb:
		return d.toDenebExecutionData().SizeSSZ()
	}
	return 0
}

// MarshalSSZ marshals the ExecutableData into a byte slice using SSZ encoding.
func (d *ExecutableData) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(d)
}

// MarshalSSZTo marshals the ExecutableData into the provided buffer and returns
// the result.
// It returns an error for unsupported versions.
func (d *ExecutableData) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	//nolint:gocritic // future versions needed.
	switch d.version {
	case version.Deneb:
		return d.toDenebExecutionData().MarshalSSZTo(buf)
	default:
		return nil, errors.New("unsupported version")
	}
}

// UnmarshalSSZ unmarshals the provided byte slice into the ExecutableData using
// SSZ encoding.
// It returns an error for unsupported versions.
func (d *ExecutableData) UnmarshalSSZ(buf []byte) error {
	//nolint:gocritic // future versions needed.
	switch d.version {
	case version.Deneb:
		data := new(ExecutableDataDeneb)
		if err := data.UnmarshalSSZ(buf); err != nil {
			return err
		}
		d.fromDenebExecutionData(data)
	default:
		return errors.New("unsupported version")
	}
	return nil
}

// HashTreeRoot computes the hash tree root of the ExecutableData. It returns an
// error for unsupported versions.
func (d *ExecutableData) HashTreeRoot() ([32]byte, error) {
	//nolint:gocritic // future versions needed.
	switch d.version {
	case version.Deneb:
		return d.toDenebExecutionData().HashTreeRoot()
	default:
		return [32]byte{}, errors.New("unsupported version")
	}
}

// HashTreeRootWith computes the hash tree root of the ExecutableData using the
// provided hasher.
// It returns an error for unsupported versions.
func (d *ExecutableData) HashTreeRootWith(hh *ssz.Hasher) error {
	//nolint:gocritic // future versions needed.
	switch d.version {
	case version.Deneb:
		return d.toDenebExecutionData().HashTreeRootWith(hh)
	default:
		return errors.New("unsupported version")
	}
}

// ExecutableDataDeneb is the ExecutableDataDeneb.
func (d *ExecutableData) toDenebExecutionData() *ExecutableDataDeneb {
	return &ExecutableDataDeneb{
		ParentHash:    d.ParentHash,
		FeeRecipient:  d.FeeRecipient,
		StateRoot:     d.StateRoot,
		ReceiptsRoot:  d.ReceiptsRoot,
		LogsBloom:     d.LogsBloom,
		Random:        d.Random,
		Number:        d.Number,
		GasLimit:      d.GasLimit,
		GasUsed:       d.GasUsed,
		Timestamp:     d.Timestamp,
		ExtraData:     d.ExtraData,
		BaseFeePerGas: d.BaseFeePerGas,
		BlockHash:     d.BlockHash,
		Transactions:  d.Transactions,
		Withdrawals:   d.Withdrawals,
		BlobGasUsed:   d.BlobGasUsed,
		ExcessBlobGas: d.ExcessBlobGas,
	}
}

// ExecutableDataDeneb is the ExecutableDataDeneb.
func (d *ExecutableData) fromDenebExecutionData(data *ExecutableDataDeneb) {
	d.ParentHash = data.ParentHash
	d.FeeRecipient = data.FeeRecipient
	d.StateRoot = data.StateRoot
	d.ReceiptsRoot = data.ReceiptsRoot
	d.LogsBloom = data.LogsBloom
	d.Random = data.Random
	d.Number = data.Number
	d.GasLimit = data.GasLimit
	d.GasUsed = data.GasUsed
	d.Timestamp = data.Timestamp
	d.ExtraData = data.ExtraData
	d.BaseFeePerGas = data.BaseFeePerGas
	d.BlockHash = data.BlockHash
	d.Transactions = data.Transactions
	d.Withdrawals = data.Withdrawals
	d.BlobGasUsed = data.BlobGasUsed
	d.ExcessBlobGas = data.ExcessBlobGas
}
