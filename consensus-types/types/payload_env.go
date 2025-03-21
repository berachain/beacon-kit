// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package types

import (
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
)

// BuiltExecutionPayloadEnv is an interface for the execution payload envelope.
//
// TODO: move interface definition to packages where it is used.
type BuiltExecutionPayloadEnv interface {
	// GetExecutionPayload retrieves the associated execution payload.
	GetExecutionPayload() *ExecutionPayload
	// GetBlockValue returns the Wei value of the block in the execution payload.
	GetBlockValue() *math.U256
	// GetBlobsBundle fetches the associated BlobsBundleV1 if available.
	GetBlobsBundle() engineprimitives.BlobsBundle
	// GetEncodedExecutionRequests fetches the associated execution requests if available
	GetEncodedExecutionRequests() []EncodedExecutionRequest
	// ShouldOverrideBuilder indicates if the builder should be overridden.
	ShouldOverrideBuilder() bool
}

// executionPayloadEnvelope is a struct that holds the execution payload and
// its associated data.
type executionPayloadEnvelope[BlobsBundleT engineprimitives.BlobsBundle] struct {
	ExecutionPayload  *ExecutionPayload         `json:"executionPayload"`
	BlockValue        *math.U256                `json:"blockValue"`
	BlobsBundle       BlobsBundleT              `json:"blobsBundle"`
	ExecutionRequests []EncodedExecutionRequest `json:"executionRequests"`
	Override          bool                      `json:"shouldOverrideBuilder"`
}

// NewEmptyExecutionPayloadEnvelope returns an empty executionPayloadEnvelope
// for the given fork version.
func NewEmptyExecutionPayloadEnvelope[
	BlobsBundleT engineprimitives.BlobsBundle,
](forkVersion common.Version) BuiltExecutionPayloadEnv {
	return &executionPayloadEnvelope[BlobsBundleT]{
		ExecutionPayload: NewEmptyExecutionPayloadWithVersion(forkVersion),
	}
}

// GetExecutionPayload returns the execution payload of the
// executionPayloadEnvelope.
func (e *executionPayloadEnvelope[BlobsBundleT]) GetExecutionPayload() *ExecutionPayload {
	return e.ExecutionPayload
}

// GetBlockValue returns the block value of the executionPayloadEnvelope.
func (e *executionPayloadEnvelope[BlobsBundleT]) GetBlockValue() *math.U256 {
	return e.BlockValue
}

// GetBlobsBundle returns the blobs bundle of the executionPayloadEnvelope.
func (e *executionPayloadEnvelope[BlobsBundleT]) GetBlobsBundle() engineprimitives.BlobsBundle {
	return e.BlobsBundle
}

// GetEncodedExecutionRequests returns the encoded Execution Requests
func (e *executionPayloadEnvelope[BlobsBundleT]) GetEncodedExecutionRequests() []EncodedExecutionRequest {
	return e.ExecutionRequests
}

// ShouldOverrideBuilder returns whether the builder should be overridden.
func (e *executionPayloadEnvelope[BlobsBundleT]) ShouldOverrideBuilder() bool {
	return e.Override
}
