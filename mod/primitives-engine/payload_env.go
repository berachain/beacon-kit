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
	"github.com/ethereum/go-ethereum/beacon/engine"
)

// ExecutionPayloadEnvelope is an interface for the execution payload envelope.
type ExecutionPayloadEnvelope interface {
	// GetExecutionPayload retrieves the execution payload associated with the
	// envelope.
	GetExecutionPayload() ExecutionPayload
	// GetValue returns the Wei value of the block within the execution payload
	// envelope.
	GetValue() primitives.Wei
	// GetBlobsBundle fetches the BlobsBundleV1 associated with the execution
	// payload, if available.
	GetBlobsBundle() *engine.BlobsBundleV1
	// ShouldOverrideBuilder indicates whether the builder should be overridden
	// in the execution environment.
	ShouldOverrideBuilder() bool
}

// TODO: this can be updated with generics to allow for different types of
// execution payloads based on the current hardfork. This should reduce
// code-duplication.
type ExecutionPayloadEnvelopeDeneb struct {
	ExecutionPayload *ExecutableDataDeneb  `json:"executionPayload"`
	BlockValue       primitives.Wei        `json:"blockValue"`
	BlobsBundle      *engine.BlobsBundleV1 `json:"blobsBundle"`
	Override         bool                  `json:"shouldOverrideBuilder"`
}

// GetExecutionPayload returns the execution payload of the
// ExecutionPayloadEnvelope.
func (e *ExecutionPayloadEnvelopeDeneb) GetExecutionPayload() ExecutionPayload {
	return e.ExecutionPayload
}

// GetValue returns the value of the ExecutionPayloadEnvelope.
func (e *ExecutionPayloadEnvelopeDeneb) GetValue() primitives.Wei {
	return e.BlockValue
}

// GetBlobsBundle returns the blobs bundle of the ExecutionPayloadEnvelope.
func (e *ExecutionPayloadEnvelopeDeneb) GetBlobsBundle() *engine.BlobsBundleV1 {
	return e.BlobsBundle
}

// ShouldOverrideBuilder returns whether the builder should be overridden.
func (e *ExecutionPayloadEnvelopeDeneb) ShouldOverrideBuilder() bool {
	return e.Override
}
