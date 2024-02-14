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

	"github.com/itsdevbear/bolaris/math"
	"github.com/itsdevbear/bolaris/types/consensus/version"
	"github.com/itsdevbear/bolaris/types/engine/interfaces"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
	"google.golang.org/protobuf/proto"
)

// WrappedExecutionPayloadHeaderDeneb ensures compatibility with the
// engine.ExecutionPayload interface.
var _ interfaces.ExecutionPayloadHeader = (*WrappedExecutionPayloadHeaderDeneb)(nil)

// WrappedExecutionPayloadHeaderDeneb wraps the ExecutionPayloadHeaderDeneb
// from Prysmatic Labs' EngineAPI v1 protobuf definitions.
type WrappedExecutionPayloadHeaderDeneb struct {
	enginev1.ExecutionPayloadHeaderDeneb
}

// Version returns the version identifier for the WrappedExecutionPayloadHeaderDeneb.
func (p *WrappedExecutionPayloadHeaderDeneb) Version() int {
	return version.Deneb
}

// IsBlinded indicates whether the payload is blinded. For WrappedExecutionPayloadHeaderDeneb,
// this is always false.
func (p *WrappedExecutionPayloadHeaderDeneb) IsBlinded() bool {
	return false
}

// ToProto returns the WrappedExecutionPayloadHeaderDeneb as a proto.Message.
func (p *WrappedExecutionPayloadHeaderDeneb) ToProto() proto.Message {
	return &p.ExecutionPayloadHeaderDeneb
}

// ToPayload returns itself as it implements the engine.ExecutionPayload interface.
func (p *WrappedExecutionPayloadHeaderDeneb) ToPayload() (interfaces.ExecutionPayload, error) {
	return nil, errors.New("cannot go from header to payload without consulting execution client")
}

// ToHeader returns itself as it implements the engine.ExecutionPayloadHeader interface.
func (p *WrappedExecutionPayloadHeaderDeneb) ToHeader() (
	interfaces.ExecutionPayloadHeader, error,
) {
	return p, nil
}

// GetValue returns the value of the payload.
func (p *WrappedExecutionPayloadHeaderDeneb) GetValue() math.Wei {
	panic("TODO")
}
