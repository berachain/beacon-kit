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
	"encoding/json"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// BuiltExecutionPayloadEnv is an interface for the execution payload envelope.
type BuiltExecutionPayloadEnv[ExecutionPayloadT any] interface {
	// GetExecutionPayload retrieves the associated execution payload.
	GetExecutionPayload() ExecutionPayloadT
	// GetValue returns the Wei value of the block in the execution payload.
	GetValue() math.Wei
	// GetBlobsBundle fetches the associated BlobsBundleV1 if available.
	GetBlobsBundle() BlobsBundle
	// ShouldOverrideBuilder indicates if the builder should be overridden.
	ShouldOverrideBuilder() bool
}

// BlobsBundle is an interface for the blobs bundle.
type BlobsBundle interface {
	// GetCommitments returns the commitments in the blobs bundle.
	GetCommitments() []eip4844.KZGCommitment
	// GetProofs returns the proofs in the blobs bundle.
	GetProofs() []eip4844.KZGProof
	// GetBlobs returns the blobs in the blobs bundle.
	GetBlobs() []*eip4844.Blob
}

// ExecutionPayloadEnvelope is a struct that holds the execution payload and
// its associated data.
// It utilizes a generic type ExecutionData to allow for different types of
// execution payloads depending on the active hard fork.
type ExecutionPayloadEnvelope[
	ExecutionPayloadT interface {
		json.Marshaler
		json.Unmarshaler
	},
	BlobsBundleT BlobsBundle,
] struct {
	ExecutionPayload ExecutionPayloadT `json:"executionPayload"`
	BlockValue       math.Wei          `json:"blockValue"`
	BlobsBundle      BlobsBundleT      `json:"blobsBundle"`
	Override         bool              `json:"shouldOverrideBuilder"`
}

// GetExecutionPayload returns the execution payload of the
// ExecutionPayloadEnvelope.
func (e *ExecutionPayloadEnvelope[
	ExecutionPayloadT, BlobsBundleT,
]) GetExecutionPayload() ExecutionPayloadT {
	return e.ExecutionPayload
}

// GetValue returns the value of the ExecutionPayloadEnvelope.
func (e *ExecutionPayloadEnvelope[
	ExecutionPayloadT, BlobsBundleT,
]) GetValue() math.Wei {
	return e.BlockValue
}

// GetBlobsBundle returns the blobs bundle of the ExecutionPayloadEnvelope.
func (e *ExecutionPayloadEnvelope[
	ExecutionPayloadT, BlobsBundleT,
]) GetBlobsBundle() BlobsBundle {
	return e.BlobsBundle
}

// ShouldOverrideBuilder returns whether the builder should be overridden.
func (e *ExecutionPayloadEnvelope[
	ExecutionPayloadT, BlobsBundleT,
]) ShouldOverrideBuilder() bool {
	return e.Override
}
