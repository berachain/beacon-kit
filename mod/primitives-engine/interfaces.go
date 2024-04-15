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
	ssz "github.com/ferranbt/fastssz"
)

// ExecutionPayloadBody is the interface for the execution data of a block.
// It contains all the fields that are part of both an execution payload header
// and a full execution payload.
type ExecutionPayloadBody interface {
	ssz.Marshaler
	ssz.Unmarshaler
	ssz.HashRoot
	IsNil() bool
	Version() uint32
	IsBlinded() bool
	GetPrevRandao() [32]byte
	GetBlockHash() primitives.ExecutionHash
	GetParentHash() primitives.ExecutionHash
	GetNumber() uint64
	GetGasLimit() uint64
	GetGasUsed() uint64
	GetTimestamp() uint64
	GetExtraData() []byte
	GetBaseFeePerGas() primitives.Wei
	GetFeeRecipient() primitives.ExecutionAddress
	GetStateRoot() primitives.ExecutionHash
	GetReceiptsRoot() primitives.ExecutionHash
	GetLogsBloom() []byte
	GetBlobGasUsed() *uint64
	GetExcessBlobGas() *uint64
}

// ExecutionPayload represents the execution data of a block.
type ExecutionPayload interface {
	ExecutionPayloadBody
	GetTransactions() [][]byte
	GetWithdrawals() []*Withdrawal
}

// PayloadAttributer represents payload attributes of a block.
type PayloadAttributer interface {
	Version() uint32
	Validate() error
}
