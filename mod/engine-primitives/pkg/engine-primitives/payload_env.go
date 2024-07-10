// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
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

package engineprimitives

import (
	engineprimitives "github.com/berachain/beacon-kit/mod/interfaces/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// BuiltExecutionPayloadEnv is an interface for the execution payload envelope.
// type BuiltExecutionPayloadEnv[ExecutionPayloadT any] interface {
// 	// GetExecutionPayload retrieves the associated execution payload.
// 	GetExecutionPayload() ExecutionPayloadT
// 	// GetValue returns the Wei value of the block in the execution payload.
// 	GetValue() math.Wei
// 	// GetBlobsBundle fetches the associated BlobsBundleV1 if available.
// 	GetBlobsBundle() BlobsBundle
// 	// ShouldOverrideBuilder indicates if the builder should be overridden.
// 	ShouldOverrideBuilder() bool
// }

// ExecutionPayloadEnvelope is a struct that holds the execution payload and
// its associated data.
// It utilizes a generic type ExecutionData to allow for different types of
// execution payloads depending on the active hard fork.
type ExecutionPayloadEnvelope[
	BlobsBundleT engineprimitives.BlobsBundle[
		eip4844.KZGCommitment, eip4844.KZGProof, eip4844.Blob,
	],
	ExecutionPayloadT constraints.JSONMarshallable,
] struct {
	ExecutionPayload ExecutionPayloadT `json:"executionPayload"`
	BlockValue       math.Wei          `json:"blockValue"`
	BlobsBundle      BlobsBundleT      `json:"blobsBundle"`
	Override         bool              `json:"shouldOverrideBuilder"`
}

func (e *ExecutionPayloadEnvelope[
	BlobsBundleT, ExecutionPayloadT,
]) New(
	blobsBundle BlobsBundleT, payload ExecutionPayloadT,
) *ExecutionPayloadEnvelope[
	BlobsBundleT, ExecutionPayloadT,
] {
	return &ExecutionPayloadEnvelope[BlobsBundleT, ExecutionPayloadT]{
		BlobsBundle:      blobsBundle,
		ExecutionPayload: payload,
	}
}

func (e *ExecutionPayloadEnvelope[
	BlobsBundleT, ExecutionPayloadT,
]) IsNil() bool {
	return e == nil
}

// GetExecutionPayload returns the execution payload of the
// ExecutionPayloadEnvelope.
func (e *ExecutionPayloadEnvelope[
	BlobsBundleT, ExecutionPayloadT,
]) GetExecutionPayload() ExecutionPayloadT {
	return e.ExecutionPayload
}

// GetValue returns the value of the ExecutionPayloadEnvelope.
func (e *ExecutionPayloadEnvelope[
	BlobsBundleT, ExecutionPayloadT,
]) GetValue() math.Wei {
	return e.BlockValue
}

// GetBlobsBundle returns the blobs bundle of the ExecutionPayloadEnvelope.
func (e *ExecutionPayloadEnvelope[
	BlobsBundleT, ExecutionPayloadT,
]) GetBlobsBundle() BlobsBundleT {
	return e.BlobsBundle
}

// ShouldOverrideBuilder returns whether the builder should be overridden.
func (e *ExecutionPayloadEnvelope[
	BlobsBundleT, ExecutionPayloadT,
]) ShouldOverrideBuilder() bool {
	return e.Override
}
