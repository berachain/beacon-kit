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
	"context"
	"errors"
	"fmt"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/engine/client"
	"github.com/berachain/beacon-kit/engine/types"
	"github.com/berachain/beacon-kit/primitives"
	coretypes "github.com/ethereum/go-ethereum/core/types"
)

// ExecutionEngine is Beacon-Kit's implementation of the `ExecutionEngine`
// from the Ethereum 2.0 Specification.
type ExecutionEngine struct {
	ec     *client.EngineClient
	logger log.Logger
}

// NewExecuitionEngine creates a new ExecutionEngine.
func NewExecutionEngine(
	ec *client.EngineClient,
	logger log.Logger,
) *ExecutionEngine {
	return &ExecutionEngine{
		ec:     ec,
		logger: logger,
	}
}

// Start spawns any goroutines required by the service.
func (ee *ExecutionEngine) Start(ctx context.Context) {
	go ee.ec.Start(ctx)
}

// Status returns error if the service is not considered healthy.
func (ee *ExecutionEngine) Status() error {
	return ee.ec.Status()
}

// TODO move.
func (ee *ExecutionEngine) GetLogs(
	ctx context.Context,
	blockHash primitives.ExecutionHash,
	addrs []primitives.ExecutionAddress,
) ([]coretypes.Log, error) {
	return ee.ec.GetLogs(ctx, blockHash, addrs)
}

// GetPayload returns the payload and blobs bundle for the given slot.
func (ee *ExecutionEngine) GetPayload(
	ctx context.Context,
	req *NewGetPayloadRequest,
) (types.ExecutionPayload, *types.BlobsBundleV1, bool, error) {
	return ee.ec.GetPayload(
		ctx, req.PayloadID,
		req.ForkVersion,
	)
}

// NotifyForkchoiceUpdate notifies the execution client of a forkchoice update.
func (ee *ExecutionEngine) NotifyForkchoiceUpdate(
	ctx context.Context,
	req *NewForkchoiceUpdateRequest,
) (*types.PayloadID, *primitives.ExecutionHash, error) {
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
	case errors.Is(err, client.ErrAcceptedSyncingPayloadStatus):
		ee.logger.Info("forkchoice updated with optimistic block",
			"head_eth1_hash", req.State.HeadBlockHash,
		)
		// telemetry.IncrCounter(1, MetricsKeyAcceptedSyncingPayloadStatus)
		return payloadID, nil, nil
	case errors.Is(err, client.ErrInvalidPayloadStatus):
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

// VerifyAndNotifyNewPayload verifies the new payload and notifies the execution
// client.
// It implictly handles:
// - IsValidBlockHash
// - IsValidVersionedHashes
// from the Ethereum 2.0 Specification from within the NewPayload call.
func (ee *ExecutionEngine) VerifyAndNotifyNewPayload(
	ctx context.Context,
	req *NewPayloadRequest,
) (bool, error) {
	payload := req.ExecutionPayload
	lastValidHash, err := ee.ec.NewPayload(
		ctx,
		payload,
		req.VersionedHashes,
		(*[32]byte)(req.ParentBeaconBlockRoot),
	)
	switch {
	case errors.Is(err, client.ErrAcceptedSyncingPayloadStatus):
		ee.logger.Info("new payload called with optimistic block",
			"payload_block_hash", (payload.GetBlockHash()),
			"parent_hash", (payload.GetParentHash()),
		)
		return false, nil
	case errors.Is(err, client.ErrInvalidPayloadStatus):
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
