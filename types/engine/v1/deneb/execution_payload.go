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

package deneb

import (
	"errors"

	"github.com/itsdevbear/bolaris/crypto/sha256"
	"github.com/itsdevbear/bolaris/math"
	"github.com/itsdevbear/bolaris/types/consensus/primitives"
	"github.com/itsdevbear/bolaris/types/consensus/version"
	"github.com/itsdevbear/bolaris/types/engine/interfaces"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
	"google.golang.org/protobuf/proto"
)

// WrappedExecutionPayloadDeneb ensures compatibility with the
// engine.ExecutionPayload interface.
var _ interfaces.ExecutionPayload = (*WrappedExecutionPayloadDeneb)(nil)

// WrappedExecutionPayloadDeneb wraps the ExecutionPayloadDeneb
// from Prysmatic Labs' EngineAPI v1 protobuf definitions.
type WrappedExecutionPayloadDeneb struct {
	*enginev1.ExecutionPayloadDeneb
	value math.Wei
}

// NewWrappedExecutionPayloadDeneb creates a new WrappedExecutionPayloadDeneb.
func NewWrappedExecutionPayloadDeneb(
	payload *enginev1.ExecutionPayloadDeneb,
	value math.Wei,
) *WrappedExecutionPayloadDeneb {
	return &WrappedExecutionPayloadDeneb{
		ExecutionPayloadDeneb: payload,
		value:                 value,
	}
}

// Version returns the version identifier for the ExecutionPayloadDeneb.
func (p *WrappedExecutionPayloadDeneb) Version() int {
	return version.Deneb
}

// IsBlinded indicates whether the payload is blinded. For ExecutionPayloadDeneb,
// this is always false.
func (p *WrappedExecutionPayloadDeneb) IsBlinded() bool {
	return false
}

// ToProto returns the ExecutionPayloadDeneb as a proto.Message.
func (p *WrappedExecutionPayloadDeneb) ToProto() proto.Message {
	return p.ExecutionPayloadDeneb
}

// ToPayload returns itself as it implements the engine.ExecutionPayload interface.
func (p *WrappedExecutionPayloadDeneb) ToPayload() interfaces.ExecutionPayload {
	return p
}

// ToHeader produces an ExecutionPayloadHeader.
func (p *WrappedExecutionPayloadDeneb) ToHeader() (interfaces.ExecutionPayloadHeader, error) {
	if len(p.Transactions) > primitives.MaxTxsPerPayloadLength {
		return nil, errors.New("too many transactions")
	}

	if len(p.Withdrawals) > primitives.MaxWithdrawalsPerPayload {
		return nil, errors.New("too many withdrawals")
	}

	return &WrappedExecutionPayloadHeaderDeneb{
		ExecutionPayloadHeaderDeneb: enginev1.ExecutionPayloadHeaderDeneb{
			ParentHash:       p.ParentHash,
			FeeRecipient:     p.FeeRecipient,
			StateRoot:        p.StateRoot,
			ReceiptsRoot:     p.ReceiptsRoot,
			LogsBloom:        p.LogsBloom,
			PrevRandao:       p.PrevRandao,
			BlockNumber:      p.BlockNumber,
			GasLimit:         p.GasLimit,
			GasUsed:          p.GasUsed,
			Timestamp:        p.Timestamp,
			ExtraData:        p.ExtraData,
			BaseFeePerGas:    p.BaseFeePerGas,
			BlockHash:        p.BlockHash,
			TransactionsRoot: sha256.HashRootAndMixinLengthAsBzSlice(p.Transactions),
			WithdrawalsRoot:  sha256.HashRootAndMixinLengthAsSlice(p.Withdrawals),
			BlobGasUsed:      p.BlobGasUsed,
			ExcessBlobGas:    p.ExcessBlobGas,
		},
		value: p.GetValue(),
	}, nil
}

// GetValue returns the value of the payload.
func (p *WrappedExecutionPayloadDeneb) GetValue() math.Wei {
	if p.value == nil {
		return math.ZeroWei()
	}
	return p.value
}
