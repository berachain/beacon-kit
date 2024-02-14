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

package interfaces

import (
	"github.com/itsdevbear/bolaris/math"
	ssz "github.com/prysmaticlabs/fastssz"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
	"google.golang.org/protobuf/proto"
)

// ExecutionPayloadBody is the interface for the execution data of a block.
// It contains all the fields that are part of both an execution payload header
// and a full execution payload.
type ExecutionPayloadBody interface {
	ssz.Marshaler
	ssz.Unmarshaler
	ssz.HashRoot
	Version() int
	IsBlinded() bool
	ToProto() proto.Message
	GetBlockHash() []byte
	GetParentHash() []byte
	GetValue() math.Wei
}

// ExecutionPayload is the interface for the execution data of a block.
type ExecutionPayload interface {
	ExecutionPayloadBody
	GetTransactions() [][]byte
	GetWithdrawals() []*enginev1.Withdrawal
}

// ExecutionPayloadHeader is the interface representing an execution payload header.
type ExecutionPayloadHeader interface {
	ExecutionPayloadBody
	GetTransactionsRoot() []byte
	GetWithdrawalsRoot() []byte
}

// PayloadAttributer is the interface for the payload attributes of a block.
type PayloadAttributer interface {
	Version() int
	IsEmpty() bool
	ToProto() proto.Message
	GetPrevRandao() []byte
	GetTimestamp() uint64
	GetSuggestedFeeRecipient() []byte
	GetWithdrawals() []*enginev1.Withdrawal
}
