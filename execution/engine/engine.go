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
	jsonrpc "github.com/berachain/beacon-kit/primitives/net/json-rpc"
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
	// Log the forkchoice update attempt.
	hasPayloadAttributes := !req.PayloadAttributes.IsNil()
	ee.metrics.markNotifyForkchoiceUpdateCalled(hasPayloadAttributes)

	// Notify the execution engine of the forkchoice update.
	payloadID, err := ee.ec.ForkchoiceUpdated(
		ctx,
		req.State,
		req.PayloadAttributes,
		req.ForkVersion,
	)

	switch {
	case err == nil:
		ee.metrics.markForkchoiceUpdateValid(
			req.State, hasPayloadAttributes, payloadID,
		)

		// If we reached here, and we have a nil payload ID, we should log a
		// warning.
		if payloadID == nil && hasPayloadAttributes {
			ee.logger.Warn(
				"Received nil payload ID on VALID engine response",
				"head_eth1_hash", req.State.HeadBlockHash,
				"safe_eth1_hash", req.State.SafeBlockHash,
				"finalized_eth1_hash", req.State.FinalizedBlockHash,
			)
			return nil, ErrNilPayloadOnValidResponse
		}

		return payloadID, nil

	case errors.Is(err, engineerrors.ErrSyncingPayloadStatus):
		// We bubble up syncing as an error, to be able to stop
		// bootstrapping from progressing in CL while EL is syncing.
		ee.metrics.markForkchoiceUpdateSyncing(req.State, err)
		return nil, err

	case errors.Is(err, engineerrors.ErrInvalidPayloadStatus):
		// If we get invalid payload status, we will need to find a valid
		// ancestor block and force a recovery.
		ee.metrics.markForkchoiceUpdateInvalid(req.State, err)
		return nil, ErrBadBlockProduced

	case jsonrpc.IsPreDefinedError(err):
		// JSON-RPC errors are predefined and should be handled as such.
		ee.metrics.markForkchoiceUpdateJSONRPCError(err)
		return nil, errors.Join(err, engineerrors.ErrPreDefinedJSONRPC)

	default:
		// All other errors are handled as undefined errors.
		ee.metrics.markForkchoiceUpdateUndefinedError(err)
		return nil, err
	}
}

// NotifyNewPayload notifies the execution client of the new payload.
func (ee *Engine) NotifyNewPayload(
	ctx context.Context,
	req *ctypes.NewPayloadRequest,
) error {
	// Log the new payload attempt.
	ee.metrics.markNewPayloadCalled(
		req.ExecutionPayload.GetBlockHash(),
		req.ExecutionPayload.GetParentHash(),
	)

	// Otherwise we will send the payload to the execution client.
	lastValidHash, err := ee.ec.NewPayload(
		ctx,
		req.ExecutionPayload,
		req.VersionedHashes,
		req.ParentBeaconBlockRoot,
	)

	// We abstract away some of the complexity and categorize status codes
	// to make it easier to reason about.
	switch {
	case err == nil:
		ee.metrics.markNewPayloadValid(
			req.ExecutionPayload.GetBlockHash(),
			req.ExecutionPayload.GetParentHash(),
		)

	case errors.Is(err, engineerrors.ErrSyncingPayloadStatus):
		// If we get accepted or syncing, we are going to optimistically
		// say that the block is valid, this is utilized during syncing
		// to allow the beacon-chain to continue processing blocks, while
		// its execution client is fetching things over it's p2p layer.
		ee.metrics.markNewPayloadSyncingPayloadStatus(
			req.ExecutionPayload.GetBlockHash(),
			req.ExecutionPayload.GetParentHash(),
		)

	case errors.IsAny(err, engineerrors.ErrAcceptedPayloadStatus):
		ee.metrics.markNewPayloadAcceptedPayloadStatus(
			req.ExecutionPayload.GetBlockHash(),
			req.ExecutionPayload.GetParentHash(),
		)

	case errors.Is(err, engineerrors.ErrInvalidPayloadStatus):
		ee.metrics.markNewPayloadInvalidPayloadStatus(
			req.ExecutionPayload.GetBlockHash(),
		)

	case jsonrpc.IsPreDefinedError(err):
		// Protect against possible nil value.
		if lastValidHash == nil {
			lastValidHash = &common.ExecutionHash{}
		}

		ee.metrics.markNewPayloadJSONRPCError(
			req.ExecutionPayload.GetBlockHash(),
			*lastValidHash,
			err,
		)
		err = errors.Join(err, engineerrors.ErrPreDefinedJSONRPC)

	default:
		ee.metrics.markNewPayloadUndefinedError(
			req.ExecutionPayload.GetBlockHash(),
			err,
		)
	}
	return err
}
