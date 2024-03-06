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
	"github.com/ethereum/go-ethereum/common"
	"github.com/itsdevbear/bolaris/config/version"
	"github.com/itsdevbear/bolaris/math"
)

//go:generate go run github.com/fjl/gencodec -type ExecutableDataDeneb -field-override executableDataDenebMarshaling -out payload.json.go
//nolint:lll
type ExecutableDataDeneb struct {
	ParentHash    common.Hash    `json:"parentHash"    ssz-size:"32"  gencodec:"required"`
	FeeRecipient  common.Address `json:"feeRecipient"  ssz-size:"20"  gencodec:"required"`
	StateRoot     common.Hash    `json:"stateRoot"     ssz-size:"32"  gencodec:"required"`
	ReceiptsRoot  common.Hash    `json:"receiptsRoot"  ssz-size:"32"  gencodec:"required"`
	LogsBloom     []byte         `json:"logsBloom"     ssz-size:"256" gencodec:"required"`
	Random        [32]byte       `json:"prevRandao"    ssz-size:"32"  gencodec:"required"`
	Number        uint64         `json:"blockNumber"   ssz-size:"8"   gencodec:"required"`
	GasLimit      uint64         `json:"gasLimit"      ssz-size:"8"   gencodec:"required"`
	GasUsed       uint64         `json:"gasUsed"       ssz-size:"8"   gencodec:"required"`
	Timestamp     uint64         `json:"timestamp"     ssz-size:"8"   gencodec:"required"`
	ExtraData     []byte         `json:"extraData"     ssz-size:"32"  gencodec:"required"`
	BaseFeePerGas []byte         `json:"baseFeePerGas" ssz-size:"32"  gencodec:"required"`
	BlockHash     common.Hash    `json:"blockHash"     ssz-size:"32"  gencodec:"required"`
	Transactions  [][]byte       `json:"transactions"  ssz-size:"?,?" gencodec:"required" ssz-max:"1048576,1073741824"`
	Withdrawals   []*Withdrawal  `json:"withdrawals"                                      ssz-max:"16"`
	BlobGasUsed   uint64         `json:"blobGasUsed"   ssz-size:"32"`
	ExcessBlobGas uint64         `json:"excessBlobGas" ssz-size:"32"`
}

// Version returns the version of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) Version() int {
	return version.Deneb
}

// IsBlinded checks if the ExecutableDataDeneb is blinded.
func (d *ExecutableDataDeneb) IsBlinded() bool {
	return false
}

// GetParentHash returns the parent hash of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetParentHash() common.Hash {
	return d.ParentHash
}

// GetFeeRecipient returns the fee recipient address of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetFeeRecipient() common.Address {
	return d.FeeRecipient
}

// GetBlockHash returns the block hash of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetBlockHash() common.Hash {
	return d.StateRoot
}

// GetReceiptsRoot returns the receipts root of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetReceiptsRoot() common.Hash {
	return d.ReceiptsRoot
}

// GetValue returns the value of the ExecutableDataDeneb in Wei.
// TODO: Needs to be on the envelope.
func (d *ExecutableDataDeneb) GetValue() math.Wei {
	return math.Wei{}
}

// GetTransactions returns the transactions of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetTransactions() [][]byte {
	return d.Transactions
}

// GetWithdrawals returns the withdrawals of the ExecutableDataDeneb.
func (d *ExecutableDataDeneb) GetWithdrawals() []*Withdrawal {
	return d.Withdrawals
}
