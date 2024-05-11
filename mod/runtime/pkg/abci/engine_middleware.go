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

package abci

import (
	"context"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
)

// ExecutionEngineMiddleware is a middleware for the execution engine.
// It filters error messages to account for the current ABCI lifecycle context.
type ExecutionEngineMiddleware[
	ExecutionPayloadT types.ExecutionPayload,
] struct {
	ExecutionEngine[ExecutionPayloadT]
}

// NewExecutionEngineMiddleware wraps the provided ExecutionEngine with
// additional middleware functionalities.
func NewExecutionEngineMiddleware[ExecutionPayloadT types.ExecutionPayload](
	engine ExecutionEngine[ExecutionPayloadT],
) *ExecutionEngineMiddleware[ExecutionPayloadT] {
	return &ExecutionEngineMiddleware[ExecutionPayloadT]{
		ExecutionEngine: engine,
	}
}

// GetPayload returns the payload and blobs bundle for the given slot.
func (eem *ExecutionEngineMiddleware[ExecutionPayloadT]) GetPayload(
	ctx context.Context,
	req *engineprimitives.GetPayloadRequest,
) (engineprimitives.BuiltExecutionPayloadEnv, error) {
	// TODO: Filter Errors
	return eem.ExecutionEngine.GetPayload(ctx, req)
}

// NotifyForkchoiceUpdate notifies the execution client of a forkchoice
// update.
func (eem *ExecutionEngineMiddleware[ExecutionPayloadT]) NotifyForkchoiceUpdate(
	ctx context.Context,
	req *engineprimitives.ForkchoiceUpdateRequest,
) (*engineprimitives.PayloadID, *common.ExecutionHash, error) {
	// TODO: Filter Errors
	return eem.ExecutionEngine.NotifyForkchoiceUpdate(ctx, req)
}

// VerifyAndNotifyNewPayload verifies the new payload and notifies the
// execution client.
//
//nolint:lll // formatter getting mad.
func (eem *ExecutionEngineMiddleware[ExecutionPayloadT]) VerifyAndNotifyNewPayload(
	ctx context.Context,
	req *engineprimitives.NewPayloadRequest[ExecutionPayloadT],
) error {
	// TODO: Filter Errors
	return eem.ExecutionEngine.VerifyAndNotifyNewPayload(ctx, req)
}
