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

package engine

import (
	"context"
	"time"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	engineerrors "github.com/berachain/beacon-kit/engine-primitives/errors"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/execution/client"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/primitives/common"
	jsonrpc "github.com/berachain/beacon-kit/primitives/net/json-rpc"
)

const engineAPITimeout = time.Minute * 5

// Engine is Beacon-Kit's implementation of the `ExecutionEngine`
// from the Ethereum 2.0 Specification.
type Engine struct {
	// ec is the engine client that the engine will use to
	// interact with the execution layer.
	ec *client.EngineClient
	// logger is the logger for the engine.
	logger log.Logger
	// metrics is the metrics for the engine.
	metrics *engineMetrics
}

// New creates a new Engine.
func New(
	engineClient *client.EngineClient,
	logger log.Logger,
	telemtrySink TelemetrySink,
) *Engine {
	return &Engine{
		ec:      engineClient,
		logger:  logger,
		metrics: newEngineMetrics(telemtrySink, logger),
	}
}

// GetPayload returns the payload and blobs bundle for the given slot.
func (ee *Engine) GetPayload(
	ctx context.Context,
	req *ctypes.GetPayloadRequest,
) (ctypes.BuiltExecutionPayloadEnv, error) {
	return ee.ec.GetPayload(
		ctx, req.PayloadID,
		req.ForkVersion,
	)
}

// NotifyForkchoiceUpdate notifies the execution client of a forkchoice update.
func (ee *Engine) NotifyForkchoiceUpdate(
	ctx context.Context,
	req *ctypes.ForkchoiceUpdateRequest,
) (*engineprimitives.PayloadID, *common.ExecutionHash, error) {
	var (
		payloadID       *engineprimitives.PayloadID
		latestValidHash *common.ExecutionHash
	)
	hasPayloadAttributes := !req.PayloadAttributes.IsNil()
	err := retryWithTimeout(ctx, engineAPITimeout,
		func(ctx context.Context) (bool, error) {
			// Log the forkchoice update attempt.
			ee.metrics.markNotifyForkchoiceUpdateCalled(hasPayloadAttributes)

			// Notify the execution engine of the forkchoice update.
			var innerErr error
			payloadID, latestValidHash, innerErr = ee.ec.ForkchoiceUpdated(
				ctx,
				req.State,
				req.PayloadAttributes,
				req.ForkVersion,
			)

			switch {
			case innerErr == nil:
				ee.metrics.markForkchoiceUpdateValid(
					req.State, hasPayloadAttributes, payloadID,
				)

				// If we reached here, we have a VALID status and a nil payload ID,
				// we should log a warning.
				if payloadID == nil && hasPayloadAttributes {
					ee.logger.Warn(
						"Received nil payload ID on VALID engine response",
						"head_eth1_hash", req.State.HeadBlockHash,
						"safe_eth1_hash", req.State.SafeBlockHash,
						"finalized_eth1_hash", req.State.FinalizedBlockHash,
					)
					// Do not retry, return the error.
					return false, ErrNilPayloadOnValidResponse
				}

				// We've received a valid response, no more retries.
				return true, nil

			case errors.IsAny(innerErr, engineerrors.ErrSyncingPayloadStatus):
				ee.metrics.markForkchoiceUpdateSyncing(req.State, innerErr)
				// Retry on SYNCING to give EL opportunity to catch up.
				return false, nil

			case errors.Is(innerErr, engineerrors.ErrInvalidPayloadStatus):
				// If we get invalid payload status, we will need to find a valid
				// ancestor block and force a recovery.
				ee.metrics.markForkchoiceUpdateInvalid(req.State, innerErr)
				// Do not retry on INVALID, return the error.
				return false, innerErr

			case jsonrpc.IsPreDefinedError(innerErr):
				// JSON-RPC errors are predefined and should be handled as such.
				ee.metrics.markForkchoiceUpdateJSONRPCError(innerErr)
				// Retry on JSON-RPC errors.
				return false, nil

			default:
				// All other errors are handled as undefined errors.
				ee.metrics.markForkchoiceUpdateUndefinedError(innerErr)
				// Retry on unknown errors, we'll log the error and retry.
				return false, nil
			}
		},
	)
	if err != nil {
		return nil, nil, err
	}

	return payloadID, latestValidHash, nil
}

// NotifyNewPayload notifies the execution client of the new payload.
func (ee *Engine) NotifyNewPayload(
	ctx context.Context,
	req *ctypes.NewPayloadRequest,
) error {
	// Otherwise we will send the payload to the execution client.
	err := retryWithTimeout(ctx, engineAPITimeout,
		func(ctx context.Context) (bool, error) {
			// Log the new payload attempt.
			ee.metrics.markNewPayloadCalled(
				req.ExecutionPayload.GetBlockHash(),
				req.ExecutionPayload.GetParentHash(),
			)

			lastValidHash, innerErr := ee.ec.NewPayload(
				ctx,
				req.ExecutionPayload,
				req.VersionedHashes,
				req.ParentBeaconBlockRoot,
			)

			// We abstract away some of the complexity and categorize status codes
			// to make it easier to reason about.
			switch {
			case errors.Is(innerErr, engineerrors.ErrSyncingPayloadStatus):
				ee.metrics.markNewPayloadSyncingPayloadStatus(
					req.ExecutionPayload.GetBlockHash(),
					req.ExecutionPayload.GetParentHash(),
				)
				// Retry on SYNCING to give EL opportunity to catch up.
				return false, nil

			case errors.IsAny(innerErr, engineerrors.ErrAcceptedPayloadStatus):
				ee.metrics.markNewPayloadAcceptedPayloadStatus(
					req.ExecutionPayload.GetBlockHash(),
					req.ExecutionPayload.GetParentHash(),
				)
				// Retry on ACCEPTED to give EL opportunity to catch up.
				return false, nil

			case errors.Is(innerErr, engineerrors.ErrInvalidPayloadStatus):
				ee.metrics.markNewPayloadInvalidPayloadStatus(
					req.ExecutionPayload.GetBlockHash(),
				)
				// Do not retry on INVALID, return the error.
				return false, innerErr

			case jsonrpc.IsPreDefinedError(innerErr):
				// Protect against possible nil value.
				if lastValidHash == nil {
					lastValidHash = &common.ExecutionHash{}
				}

				ee.metrics.markNewPayloadJSONRPCError(
					req.ExecutionPayload.GetBlockHash(),
					*lastValidHash,
					innerErr,
				)

				// Retry on JSON-RPC errors.
				return false, nil
			case innerErr != nil:
				ee.metrics.markNewPayloadUndefinedError(
					req.ExecutionPayload.GetBlockHash(),
					innerErr,
				)
				// Retry on unknown errors, we'll log the error and retry.
				return false, nil
			default:
				ee.metrics.markNewPayloadValid(
					req.ExecutionPayload.GetBlockHash(),
					req.ExecutionPayload.GetParentHash(),
				)
				return true, nil
			}
		},
	)
	return err
}
