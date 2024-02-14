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

package capella

import (
	"github.com/itsdevbear/bolaris/math"
	"github.com/itsdevbear/bolaris/types/consensus/version"
	"github.com/itsdevbear/bolaris/types/engine/interfaces"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
	"google.golang.org/protobuf/proto"
)

// WrappedExecutionPayloadCapella ensures compatibility with the
// engine.ExecutionPayload interface.
var _ interfaces.ExecutionPayload = (*WrappedExecutionPayloadCapella)(nil)

// WrappedExecutionPayloadCapella wraps the ExecutionPayloadCapella
// from Prysmatic Labs' EngineAPI v1 protobuf definitions.
type WrappedExecutionPayloadCapella struct {
	*enginev1.ExecutionPayloadCapella
	value math.Wei
}

// NewWrappedExecutionPayloadCapella creates a new WrappedExecutionPayloadCapella.
func NewWrappedExecutionPayloadCapella(
	payload *enginev1.ExecutionPayloadCapella,
	value math.Wei,
) *WrappedExecutionPayloadCapella {
	return &WrappedExecutionPayloadCapella{
		ExecutionPayloadCapella: payload,
		value:                   value,
	}
}

// Version returns the version identifier for the ExecutionPayloadCapella.
func (p *WrappedExecutionPayloadCapella) Version() int {
	return version.Capella
}

// IsBlinded indicates whether the payload is blinded. For ExecutionPayloadCapella,
// this is always false.
func (p *WrappedExecutionPayloadCapella) IsBlinded() bool {
	return false
}

// ToProto returns the ExecutionPayloadCapella as a proto.Message.
func (p *WrappedExecutionPayloadCapella) ToProto() proto.Message {
	return p.ExecutionPayloadCapella
}

// ToPayload returns itself as it implements the engine.ExecutionPayload interface.
func (p *WrappedExecutionPayloadCapella) ToPayload() interfaces.ExecutionPayload {
	return p
}

// ToHeader produces an ExecutionPayloadHeader.
func (p *WrappedExecutionPayloadCapella) ToHeader() interfaces.ExecutionPayloadHeader {
	// TODO: @ocnc
	panic("TODO: Implement slice merkalization for ExecutionPayloadCapella")
}

// GetValue returns the value of the payload.
func (p *WrappedExecutionPayloadCapella) GetValue() math.Wei {
	if p.value == nil {
		return math.ZeroWei()
	}
	return p.value
}
