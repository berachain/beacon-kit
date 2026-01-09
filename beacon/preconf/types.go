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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package preconf

import (
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
)

const (
	// PayloadEndpoint is the API endpoint for fetching/serving preconf payloads.
	PayloadEndpoint = "/eth/v1/preconf/payload"
)

// GetPayloadRequest is the request body for the GetPayload endpoint.
type GetPayloadRequest struct {
	// Slot is the slot number for which to retrieve the payload.
	Slot math.Slot `json:"slot"`
}

// GetPayloadResponse is the response body for the GetPayload endpoint.
type GetPayloadResponse struct {
	// ForkVersion is the fork version of the payload.
	// This is needed because the Versionable field is not JSON-serialized.
	ForkVersion common.Version `json:"fork_version"`

	// ExecutionPayload is the execution payload for the requested slot.
	ExecutionPayload *ctypes.ExecutionPayload `json:"execution_payload"`

	// BlobsBundle contains the blobs, commitments, and proofs.
	BlobsBundle *engineprimitives.BlobsBundleV1 `json:"blobs_bundle"`

	// BlockValue is the Wei value of the block.
	BlockValue *math.U256 `json:"block_value"`

	// ExecutionRequests contains the encoded execution requests (Electra+).
	ExecutionRequests []ctypes.EncodedExecutionRequest `json:"execution_requests,omitempty"`
}

// ToExecutionPayloadEnvelope converts the response to a BuiltExecutionPayloadEnv.
// This method sets the Versionable field which is not JSON-serialized.
func (r *GetPayloadResponse) ToExecutionPayloadEnvelope() ctypes.BuiltExecutionPayloadEnv {
	// Set the Versionable field from the serialized fork version
	if r.ExecutionPayload != nil {
		r.ExecutionPayload.Versionable = ctypes.NewVersionable(r.ForkVersion)
	}

	return ctypes.NewExecutionPayloadEnvelope[*engineprimitives.BlobsBundleV1](
		r.ExecutionPayload,
		r.BlobsBundle,
		r.ExecutionRequests,
	)
}

// NewGetPayloadResponseFromEnvelope creates a GetPayloadResponse from a BuiltExecutionPayloadEnv.
func NewGetPayloadResponseFromEnvelope(env ctypes.BuiltExecutionPayloadEnv) *GetPayloadResponse {
	// Type assert blobs bundle to concrete type (nil-safe: assertion returns nil if input is nil)
	blobsBundle, _ := env.GetBlobsBundle().(*engineprimitives.BlobsBundleV1)

	var forkVersion common.Version
	if payload := env.GetExecutionPayload(); payload != nil {
		forkVersion = payload.GetForkVersion()
	}

	return &GetPayloadResponse{
		ForkVersion:       forkVersion,
		ExecutionPayload:  env.GetExecutionPayload(),
		BlobsBundle:       blobsBundle,
		BlockValue:        env.GetBlockValue(),
		ExecutionRequests: env.GetEncodedExecutionRequests(),
	}
}

// ErrorResponse is the error response body.
type ErrorResponse struct {
	// Code is the error code.
	Code int `json:"code"`

	// Message is the error message.
	Message string `json:"message"`
}
