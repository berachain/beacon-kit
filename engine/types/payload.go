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
	"github.com/itsdevbear/bolaris/primitives"
)

// ExecutableData is the data necessary to execute an EL payload.
//
// go:generate go run github.com/fjl/gencodec -type ExecutableData -field-override executableDataMarshaling -out payload.json.go
//
//go:generate go run github.com/itsdevbear/fastssz/sszgen -path . -objs ExecutableData -include ../../primitives,$HOME/go/pkg/mod/github.com/ethereum/go-ethereum@v1.13.14/common
//nolint:lll // struct tags
type ExecutableData struct {
	ParentHash   common.Hash          `json:"parentHash"    ssz-size:"32"  gencodec:"required"`
	FeeRecipient common.Address       `json:"feeRecipient"  ssz-size:"20"  gencodec:"required"`
	StateRoot    common.Hash          `json:"stateRoot"     ssz-size:"32"  gencodec:"required"`
	ReceiptsRoot common.Hash          `json:"receiptsRoot"  ssz-size:"32"  gencodec:"required"`
	LogsBloom    []byte               `json:"logsBloom"     ssz-size:"256" gencodec:"required"`
	Random       [32]byte             `json:"prevRandao"    ssz-size:"32"  gencodec:"required"`
	Number       primitives.SSZUint64 `json:"blockNumber"   ssz-size:"8"               gencodec:"required"`
	GasLimit     primitives.SSZUint64 `json:"gasLimit"      ssz-size:"8"               gencodec:"required"`
	GasUsed      primitives.SSZUint64 `json:"gasUsed"       ssz-size:"8"               gencodec:"required"`
	Timestamp    primitives.SSZUint64 `json:"timestamp"     ssz-size:"8"               gencodec:"required"`
	ExtraData    []byte               `json:"extraData"     ssz-size:"32"  gencodec:"required"`
	// // BaseFeePerGas *big.Int              `json:"baseFeePerGas" ssz-size:"32"  gencodec:"required"`
	BlockHash    common.Hash   `json:"blockHash"     ssz-size:"32"  gencodec:"required"`
	Transactions [][]byte      `json:"transactions"  ssz-size:"?,?" gencodec:"required" ssz-max:"1048576,1073741824"`
	Withdrawals  []*Withdrawal `json:"withdrawals"                                      ssz-max:"16"`
	// BlobGasUsed   *primitives.SSZUint64 `json:"blobGasUsed" ssz-size:"32"`
	// ExcessBlobGas *primitives.SSZUint64 `json:"excessBlobGas" ssz-size:"32"`
}
