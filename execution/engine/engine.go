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

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	engineerrors "github.com/berachain/beacon-kit/engine-primitives/errors"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/execution/client"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/cenkalti/backoff/v5"
)

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
) (*engineprimitives.PayloadID, error) {
	hasPayloadAttributes := !req.PayloadAttributes.IsNil()

	// Configure backoff. This will retry maxRetries number of times.
	// Specifying 0 maxRetries will retry infinitely. Between each retry, it
	// will wait RPCRetryInterval amount of time. This backoff will increase
	// exponentially until it reaches RPCMaxRetryInterval.
	engineAPIBackoff := backoff.NewExponentialBackOff()
	engineAPIBackoff.InitialInterval = ee.ec.GetRPCRetryInterval()
	engineAPIBackoff.MaxInterval = ee.ec.GetRPCMaxRetryInterval()
	maxRetries := uint(ee.ec.GetRPCRetries())

	return backoff.Retry(ctx, func() (*engineprimitives.PayloadID, error) {
		// Log and call the forkchoice update.
		ee.metrics.markNotifyForkchoiceUpdateCalled(hasPayloadAttributes)
		payloadID, err := ee.ec.ForkchoiceUpdated(
			ctx, req.State, req.PayloadAttributes, req.ForkVersion,
		)

		// NotifyForkchoiceUpdate gets called under two circumstances:
		// 1. Payload Building (During PrepareProposal or
		//    optimistically in ProcessProposal)
		// 2. FinalizeBlock
		// We'll discriminate error handling based on these.
		switch {
		case err == nil:
			ee.metrics.markForkchoiceUpdateValid(
				req.State, hasPayloadAttributes, payloadID,
			)

			// If we reached here, we have a VALID status and a nil payload ID,
			// we should log a warning and error.
			if payloadID == nil && hasPayloadAttributes {
				ee.logger.Warn(
					"Received nil payload ID on VALID engine response",
					"head_eth1_hash", req.State.HeadBlockHash,
					"safe_eth1_hash", req.State.SafeBlockHash,
					"finalized_eth1_hash", req.State.FinalizedBlockHash,
				)
				// Do not retry, return the error.
				return nil, ErrNilPayloadOnValidResponse
			}

			// We've received a valid response, no more retries.
			return payloadID, nil

		case errors.IsAny(err, engineerrors.ErrSyncingPayloadStatus):
			ee.metrics.markForkchoiceUpdateSyncing(req.State, err)
			// In all circumstances, keep retrying until the EVM is synced.
			return nil, err

		case errors.Is(err, engineerrors.ErrInvalidPayloadStatus):
			ee.metrics.markForkchoiceUpdateInvalid(req.State, err)
			// During payload building, then there is an invalid
			// payload and should error.
			// During FinalizeBlock, something is broken because
			// this should never happen.
			return nil, backoff.Permanent(err)

		case client.IsNonFatalError(err):
			ee.metrics.markForkchoiceUpdateNonFatalError(err)
			return nil, err

		case client.IsFatalError(err):
			ee.metrics.markForkchoiceUpdateFatalError(err)
			return nil, backoff.Permanent(err)

		default:
			ee.metrics.markForkchoiceUpdateUndefinedError(err)
			// Retry on unknown errors, we'll log the error and retry.
			return nil, backoff.Permanent(err)
		}
	},
		backoff.WithBackOff(engineAPIBackoff),
		backoff.WithMaxTries(maxRetries),
		// Set 0 max elapsed time so we don't check it.
		backoff.WithMaxElapsedTime(0),
	)
}

// NotifyNewPayload notifies the execution client of the new payload.
func (ee *Engine) NotifyNewPayload(
	ctx context.Context,
	req *ctypes.NewPayloadRequest,
	retryOnSyncingStatus bool,
) error {
	// Configure backoff. This will retry maxRetries number of times.
	// Specifying 0 maxRetries will retry infinitely. Between each retry, it
	// will wait RPCRetryInterval amount of time. This backoff will increase
	// exponentially until it reaches RPCMaxRetryInterval.
	engineAPIBackoff := backoff.NewExponentialBackOff()
	engineAPIBackoff.InitialInterval = ee.ec.GetRPCRetryInterval()
	engineAPIBackoff.MaxInterval = ee.ec.GetRPCMaxRetryInterval()
	maxRetries := uint(ee.ec.GetRPCRetries())

	// Otherwise we will send the payload to the execution client.
	_, err := backoff.Retry(ctx, func() (*common.ExecutionHash, error) {
		// Log the new payload attempt.
		ee.metrics.markNewPayloadCalled(
			req.ExecutionPayload.GetBlockHash(), req.ExecutionPayload.GetParentHash(),
		)
		lastValidHash, err := ee.ec.NewPayload(
			ctx, req.ExecutionPayload, req.VersionedHashes, req.ParentBeaconBlockRoot,
		)

		// NotifyNewPayload gets called under three circumstances:
		// 1. ProcessProposal state transition
		// 2. FinalizeBlock state transition
		// We'll discriminate error handling based on these.
		switch {
		case err == nil:
			ee.metrics.markNewPayloadValid(
				req.ExecutionPayload.GetBlockHash(), req.ExecutionPayload.GetParentHash(),
			)
			// We've received a valid response, no more retries.
			return lastValidHash, nil

		case errors.IsAny(err, engineerrors.ErrSyncingPayloadStatus, engineerrors.ErrAcceptedPayloadStatus):
			ee.metrics.markNewPayloadAcceptedSyncingPayloadStatus(
				err, req.ExecutionPayload.GetBlockHash(), req.ExecutionPayload.GetParentHash(),
			)
			// During ProcessProposal, we must be able to verify the
			// block. Since we do not send a NotifyForkchoiceUpdate
			// during ProcessProposal, we must retry here until EL is
			// synced.
			if retryOnSyncingStatus {
				return nil, err
			}
			// During FinalizeBlock, we do not need to verify the block.
			// We do not need to retry here, as the following call to
			// NotifyForkchoiceUpdate will inform the EL of the new head
			// and then wait for it to sync.
			// Don't return error here, because we want to send the forkchoice update regardless.
			ee.logger.Warn("Pushed new payload to SYNCING node.", "error", err,
				"blockNum", req.ExecutionPayload.GetNumber(), "blockHash", req.ExecutionPayload.GetBlockHash(),
			)
			return &common.ExecutionHash{}, nil

		case errors.Is(err, engineerrors.ErrInvalidPayloadStatus):
			ee.metrics.markNewPayloadInvalidPayloadStatus(
				req.ExecutionPayload.GetBlockHash(),
			)
			// During payload building, then there is an invalid
			// payload and should error.
			// During FinalizeBlock, something is broken because
			// this should never happen.
			return nil, backoff.Permanent(err)

		case client.IsNonFatalError(err):
			// Protect against possible nil value.
			if lastValidHash == nil {
				lastValidHash = &common.ExecutionHash{}
			}
			ee.metrics.markNewPayloadNonFatalError(
				req.ExecutionPayload.GetBlockHash(),
				*lastValidHash, err,
			)
			return nil, err

		case client.IsFatalError(err):
			// Protect against possible nil value.
			if lastValidHash == nil {
				lastValidHash = &common.ExecutionHash{}
			}
			ee.metrics.markNewPayloadFatalError(
				req.ExecutionPayload.GetBlockHash(),
				*lastValidHash, err,
			)
			return nil, backoff.Permanent(err)
		default:
			ee.metrics.markNewPayloadUndefinedError(
				req.ExecutionPayload.GetBlockHash(), err,
			)
			// Do not retry on unknown errors.
			return nil, backoff.Permanent(err)
		}
	}, backoff.WithBackOff(engineAPIBackoff), backoff.WithMaxTries(maxRetries),
		// Set 0 max elapsed time so we don't check it.
		backoff.WithMaxElapsedTime(0),
	)
	return err
}
