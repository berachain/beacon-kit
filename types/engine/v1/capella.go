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

package enginev1

import (
	"github.com/itsdevbear/bolaris/types/consensus/version"
	"github.com/itsdevbear/bolaris/types/engine"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
)

type ExecutionPayloadCapella struct {
	enginev1.ExecutionPayloadCapella
}

var (
	// ExecutionPayloadCapella ensures compatibility with the engine.ExecutionPayload interface.
	_ engine.ExecutionPayload = (*ExecutionPayloadCapella)(nil)
)

// Version returns the version identifier for the ExecutionPayloadCapella.
func (p *ExecutionPayloadCapella) Version() int {
	return version.Capella
}

// IsBlinded indicates whether the payload is blinded. For ExecutionPayloadCapella,
// this is always false.
func (p *ExecutionPayloadCapella) IsBlinded() bool {
	return false
}

// ToPayload returns itself as it implements the engine.ExecutionPayload interface.
func (p *ExecutionPayloadCapella) ToPayload() engine.ExecutionPayload {
	return p
}

// ToHeader is intended to convert the ExecutionPayloadCapella to an ExecutionPayloadHeader.
// Currently, it panics as the slice merkalization is yet to be implemented.
func (p *ExecutionPayloadCapella) ToHeader() engine.ExecutionPayloadHeader {
	panic("TODO: Implement slice merkalization for ExecutionPayloadCapella")
}
