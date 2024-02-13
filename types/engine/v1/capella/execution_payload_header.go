// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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
	"errors"

	"github.com/itsdevbear/bolaris/math"
	"github.com/itsdevbear/bolaris/types/consensus/version"
	"github.com/itsdevbear/bolaris/types/engine/interfaces"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
	"google.golang.org/protobuf/proto"
)

// WrappedExecutionPayloadHeaderCapella ensures compatibility with the
// engine.ExecutionPayload interface.
var _ interfaces.ExecutionPayloadHeader = (*WrappedExecutionPayloadHeaderCapella)(nil)

// WrappedExecutionPayloadHeaderCapella wraps the ExecutionPayloadHeaderCapella
// from Prysmatic Labs' EngineAPI v1 protobuf definitions.
type WrappedExecutionPayloadHeaderCapella struct {
	enginev1.ExecutionPayloadHeaderCapella
}

// Version returns the version identifier for the WrappedExecutionPayloadHeaderCapella.
func (p *WrappedExecutionPayloadHeaderCapella) Version() int {
	return version.Capella
}

// IsBlinded indicates whether the payload is blinded. For WrappedExecutionPayloadHeaderCapella,
// this is always false.
func (p *WrappedExecutionPayloadHeaderCapella) IsBlinded() bool {
	return false
}

// ToProto returns the WrappedExecutionPayloadHeaderCapella as a proto.Message.
func (p *WrappedExecutionPayloadHeaderCapella) ToProto() proto.Message {
	return &p.ExecutionPayloadHeaderCapella
}

// ToPayload returns itself as it implements the engine.ExecutionPayload interface.
func (p *WrappedExecutionPayloadHeaderCapella) ToPayload() (interfaces.ExecutionPayload, error) {
	return nil, errors.New("cannot go from header to payload without consulting execution client")
}

// ToHeader returns itself as it implements the engine.ExecutionPayloadHeader interface.
func (p *WrappedExecutionPayloadHeaderCapella) ToHeader() (
	interfaces.ExecutionPayloadHeader, error,
) {
	return p, nil
}

// GetValue returns the value of the payload.
func (p *WrappedExecutionPayloadHeaderCapella) GetValue() math.Wei {
	panic("TODO")
}
