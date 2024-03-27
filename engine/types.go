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
	"github.com/berachain/beacon-kit/engine/types"
	"github.com/berachain/beacon-kit/primitives"
)

// NewPayloadRequest as per the Ethereum 2.0 specification:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/deneb/beacon-chain.md#modified-newpayloadrequest
//
//nolint:lll
type NewPayloadRequest struct {
	// ExecutionPayload is the payload to the execution client.
	ExecutionPayload types.ExecutionPayload
	// VersionedHashes is the versioned hashes of the execution payload.
	VersionedHashes []primitives.ExecutionHash
	// ParentBeaconBlockRoot is the root of the parent beacon block.
	ParentBeaconBlockRoot *primitives.Root
}

// NewForkchoiceUpdateRequest.
type NewForkchoiceUpdateRequest struct {
	// State is the forkchoice state.
	State *types.ForkchoiceState
	// PayloadAttributes is the payload attributer.
	PayloadAttributes types.PayloadAttributer
	// ForkVersion is the fork version that we
	// are going to be submitting for.
	ForkVersion uint32
}

// NewGetPayloadRequest represents a request to get a payload.
type NewGetPayloadRequest struct {
	// PayloadID is the payload ID.
	PayloadID types.PayloadID
	// ForkVersion is the fork version that we are
	// currently on.
	ForkVersion uint32
}
