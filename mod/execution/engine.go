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

package execution

import (
	"context"
	"fmt"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/execution/client"
	"github.com/berachain/beacon-kit/mod/primitives"
	engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine"
	"github.com/cockroachdb/errors"
)

// Engine is Beacon-Kit's implementation of the `ExecutionEngine`
// from the Ethereum 2.0 Specification.
type Engine struct {
	ec     *client.EngineClient
	logger log.Logger
}

// NewExecuitionEngine creates a new Engine.
func NewEngine(
	ec *client.EngineClient,
	logger log.Logger,
) *Engine {
	return &Engine{
		ec:     ec,
		logger: logger,
	}
}

// Start spawns any goroutines required by the service.
func (ee *Engine) Start(ctx context.Context) {
	go ee.ec.Start(ctx)
}

// Status returns error if the service is not considered healthy.
func (ee *Engine) Status() error {
	return ee.ec.Status()
}

// TODO move.
func (ee *Engine) GetLogs(
	ctx context.Context,
	blockHash primitives.ExecutionHash,
	addrs []primitives.ExecutionAddress,
) ([]engineprimitives.Log, error) {
	return ee.ec.GetLogs(ctx, blockHash, addrs)
}

// GetPayload returns the payload and blobs bundle for the given slot.
func (ee *Engine) GetPayload(
	ctx context.Context,
	req *engineprimitives.GetPayloadRequest,
) (engineprimitives.BuiltExecutionPayload, error) {
	return ee.ec.GetPayload(
		ctx, req.PayloadID,
		req.ForkVersion,
	)
}

// NotifyForkchoiceUpdate notifies the execution client of a forkchoice update.
func (ee *Engine) NotifyForkchoiceUpdate(
	ctx context.Context,
	req *engineprimitives.ForkchoiceUpdateRequest,
) (*engineprimitives.PayloadID, *primitives.ExecutionHash, error) {
	ee.logger.Info("notifying forkchoice update",
		"head_eth1_hash", req.State.HeadBlockHash,
		"safe_eth1_hash", req.State.SafeBlockHash,
		"finalized_eth1_hash", req.State.FinalizedBlockHash,
		"has_attributes", req.PayloadAttributes != nil,
	)

	// Notify the execution engine of the forkchoice update.
	payloadID, latestValidHash, err := ee.ec.ForkchoiceUpdated(
		ctx,
		req.State,
		req.PayloadAttributes,
		req.ForkVersion,
	)
	switch {
	case errors.Is(err, client.ErrAcceptedPayloadStatus) ||
		errors.Is(err, client.ErrSyncingPayloadStatus):
		ee.logger.Info("forkchoice updated with optimistic block",
			"head_eth1_hash", req.State.HeadBlockHash,
		)
		// telemetry.IncrCounter(1, MetricsKeyAcceptedSyncingPayloadStatus)
		return payloadID, nil, nil
	case errors.Is(err, client.ErrInvalidPayloadStatus) ||
		errors.Is(err, client.ErrInvalidBlockHashPayloadStatus):
		// Attempt to get the chain back into a valid state, by
		// getting finding an ancestor block with a valid payload and
		// forcing a recovery.
		req.State.HeadBlockHash = req.State.SafeBlockHash
		payloadID, latestValidHash, err = ee.NotifyForkchoiceUpdate(ctx, req)
		if err != nil {
			// We have to return the error here since this function
			// is recursive.
			return nil, nil, err
		}
		return payloadID, latestValidHash, ErrBadBlockProduced
	case err != nil:
		ee.logger.Error("undefined execution engine error", "error", err)
		return nil, nil, err
	}

	return payloadID, latestValidHash, nil
}

// VerifyAndNotifyNewPayload verifies the new payload and notifies the
// execution client.
func (ee *Engine) VerifyAndNotifyNewPayload(
	ctx context.Context,
	req *engineprimitives.NewPayloadRequest,
) (bool, error) {
	// First we verify the block hash and versioned hashes are valid.
	if err := req.HasValidVersionedAndBlockHashes(); err != nil {
		return false, err
	}

	// If the block already exists, we can skip sending the payload to the
	// execution client.
	if req.SkipIfExists {
		header, err := ee.ec.HeaderByHash(
			ctx,
			req.ExecutionPayload.GetBlockHash(),
		)
		if header != nil && err != nil {
			ee.logger.Info("skipping new payload, block already available",
				"block_hash", req.ExecutionPayload.GetBlockHash(),
			)
		}
		return true, nil
	}

	// Then we can ask the EL to process the new payload.
	lastValidHash, err := ee.ec.NewPayload(
		ctx,
		req.ExecutionPayload,
		req.VersionedHashes,
		req.ParentBeaconBlockRoot,
	)
	switch {
	case errors.Is(err, client.ErrAcceptedPayloadStatus) ||
		errors.Is(err, client.ErrSyncingPayloadStatus):
		ee.logger.Info("new payload called with optimistic block",
			"payload_block_hash", (req.ExecutionPayload.GetBlockHash()),
			"parent_hash", (req.ExecutionPayload.GetParentHash()),
		)
		return false, nil
	case errors.Is(err, client.ErrInvalidPayloadStatus) ||
		errors.Is(err, client.ErrInvalidBlockHashPayloadStatus):
		ee.logger.Error(
			"invalid payload status",
			"last_valid_hash", fmt.Sprintf("%#x", lastValidHash),
		)
		return false, ErrBadBlockProduced
	case err != nil:
		return false, err
	}
	return true, nil
}
