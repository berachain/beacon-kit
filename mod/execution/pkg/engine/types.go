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

package engine

import (
	"encoding/json"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
)

// ExecutionPayload represents the payload of an execution block.
type ExecutionPayload[ExecutionPayloadT any] interface {
	json.Marshaler
	json.Unmarshaler
	// Empty creates an empty execution payload.
	Empty(uint32) ExecutionPayloadT
	// GetTransactions returns the transactions included in the payload.
	GetTransactions() [][]byte
	// GetBlockHash returns the hash of the block.
	GetBlockHash() common.ExecutionHash
	// GetParentHash returns the hash of the parent block.
	GetParentHash() common.ExecutionHash
	// Version returns the version of the payload.
	Version() uint32
}

// TelemetrySink is an interface for sending metrics to a telemetry backend.
type TelemetrySink interface {
	// IncrementCounter increments a counter metric identified by the provided
	// keys.
	IncrementCounter(key string, args ...string)
	// SetGauge sets a gauge metric to the specified value, identified by the
	// provided keys.
	SetGauge(key string, value int64, args ...string)
}
