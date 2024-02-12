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

package v3

import (
	"github.com/itsdevbear/bolaris/types/consensus/version"
	"github.com/itsdevbear/bolaris/types/engine"
	v1 "github.com/itsdevbear/bolaris/types/engine/v1"
)

var (
	// ExecutionPayloadDeneb ensures compatibility with the engine.ExecutionPayload interface.
	_ engine.ExecutionPayload = (*ExecutionPayloadDeneb)(nil)
	_ engine.ExecutionPayload = (*ExecutionPayloadDenebWithValueAndBlobsBundle)(nil)
)

// Version returns the version identifier for the ExecutionPayloadDeneb.
func (p *ExecutionPayloadDeneb) Version() int {
	return version.Capella
}

// IsBlinded indicates whether the payload is blinded. For ExecutionPayloadDeneb,
// this is always false.
func (p *ExecutionPayloadDeneb) IsBlinded() bool {
	return false
}

// ToPayload returns itself as it implements the engine.ExecutionPayload interface.
func (p *ExecutionPayloadDeneb) ToPayload() engine.ExecutionPayload {
	return p
}

// ToHeader is intended to convert the ExecutionPayloadDeneb to an ExecutionPayloadHeader.
// Currently, it panics as the slice merkalization is yet to be implemented.
func (p *ExecutionPayloadDeneb) ToHeader() engine.ExecutionPayloadHeader {
	panic("TODO: Implement slice merkalization for ExecutionPayloadDeneb")
}

func (p *ExecutionPayloadDenebWithValueAndBlobsBundle) Version() int {
	return version.Capella
}

func (p *ExecutionPayloadDenebWithValueAndBlobsBundle) IsBlinded() bool {
	return false
}

func (p *ExecutionPayloadDenebWithValueAndBlobsBundle) ToPayload() engine.ExecutionPayload {
	return p
}

func (p *ExecutionPayloadDenebWithValueAndBlobsBundle) ToHeader() engine.ExecutionPayloadHeader {
	panic("TODO: Implement slice merkalization for ExecutionPayloadDenebWithValue")
}

// GetTransactions returns the transactions of the payload.
func (p *ExecutionPayloadDenebWithValueAndBlobsBundle) GetTransactions() [][]byte {
	return p.GetPayload().GetTransactions()
}

// GetWithdrawals returns the withdrawals of the payload.
func (p *ExecutionPayloadDenebWithValueAndBlobsBundle) GetWithdrawals() []*v1.Withdrawal {
	return p.GetPayload().GetWithdrawals()
}
