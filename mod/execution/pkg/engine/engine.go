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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package engine

import (
	"bytes"
	"context"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	engineerrors "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/errors"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/execution/pkg/client"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	jsonrpc "github.com/berachain/beacon-kit/mod/primitives/pkg/net/json-rpc"
)

// Engine is Beacon-Kit's implementation of the `ExecutionEngine`
// from the Ethereum 2.0 Specification.
type Engine[
	ExecutionPayloadT ExecutionPayload[ExecutionPayloadT, WithdrawalsT],
	PayloadAttributesT engineprimitives.PayloadAttributer,
	PayloadIDT ~[8]byte,
	WithdrawalsT interface {
		Len() int
		EncodeIndex(int, *bytes.Buffer)
	},
] struct {
	// ec is the engine client that the engine will use to
	// interact with the execution layer.
	ec *client.EngineClient[ExecutionPayloadT, PayloadAttributesT]
	// logger is the logger for the engine.
	logger log.Logger
	// metrics is the metrics for the engine.
	metrics *engineMetrics
}

// New creates a new Engine.
func New[
	ExecutionPayloadT ExecutionPayload[ExecutionPayloadT, WithdrawalsT],
	PayloadAttributesT engineprimitives.PayloadAttributer,
	PayloadIDT ~[8]byte,
	WithdrawalsT interface {
		Len() int
		EncodeIndex(int, *bytes.Buffer)
	},
](
	engineClient *client.EngineClient[ExecutionPayloadT, PayloadAttributesT],
	logger log.Logger,
	telemtrySink TelemetrySink,
) *Engine[
	ExecutionPayloadT, PayloadAttributesT,
	PayloadIDT, WithdrawalsT,
] {
	return &Engine[
		ExecutionPayloadT, PayloadAttributesT, PayloadIDT,
		WithdrawalsT,
	]{
		ec:      engineClient,
		logger:  logger,
		metrics: newEngineMetrics(telemtrySink, logger),
	}
}

// Start spawns any goroutines required by the service.
func (ee *Engine[_, _, _, _]) Start(
	ctx context.Context,
) error {
	go func() {
		// TODO: handle better
		if err := ee.ec.Start(ctx); err != nil {
			panic(err)
		}
	}()
	return nil
}

// GetPayload returns the payload and blobs bundle for the given slot.
func (ee *Engine[
	ExecutionPayloadT, _, _, _,
]) GetPayload(
	ctx context.Context,
	req *engineprimitives.GetPayloadRequest[engineprimitives.PayloadID],
) (engineprimitives.BuiltExecutionPayloadEnv[ExecutionPayloadT], error) {
	return ee.ec.GetPayload(
		ctx, req.PayloadID,
		req.ForkVersion,
	)
}

// NotifyForkchoiceUpdate notifies the execution client of a forkchoice update.
func (ee *Engine[
	_, PayloadAttributesT, _, _,
]) NotifyForkchoiceUpdate(
	ctx context.Context,
	req *engineprimitives.ForkchoiceUpdateRequest[PayloadAttributesT],
) (*engineprimitives.PayloadID, *common.ExecutionHash, error) {
	// Log the forkchoice update attempt.
	hasPayloadAttributes := !req.PayloadAttributes.IsNil()
	ee.metrics.markNotifyForkchoiceUpdateCalled(hasPayloadAttributes)

	// Notify the execution engine of the forkchoice update.
	payloadID, latestValidHash, err := ee.ec.ForkchoiceUpdated(
		ctx,
		req.State,
		req.PayloadAttributes,
		req.ForkVersion,
	)

	switch {
	// We do not bubble the error up, since we want to handle it
	// in the same way as the other cases.
	case errors.IsAny(
		err,
		engineerrors.ErrAcceptedPayloadStatus,
		engineerrors.ErrSyncingPayloadStatus,
	):
		ee.metrics.markForkchoiceUpdateAcceptedSyncing(req.State, err)
		return payloadID, nil, nil

	// If we get invalid payload status, we will need to find a valid
	// ancestor block and force a recovery.
	//
	// These two cases are semantically the same:
	// https://github.com/ethereum/execution-apis/issues/270
	case errors.IsAny(
		err,
		engineerrors.ErrInvalidPayloadStatus,
		engineerrors.ErrInvalidBlockHashPayloadStatus,
	):
		ee.metrics.markForkchoiceUpdateInvalid(req.State, err)
		return payloadID, latestValidHash, ErrBadBlockProduced

	// JSON-RPC errors are predefined and should be handled as such.
	case jsonrpc.IsPreDefinedError(err):
		ee.metrics.markForkchoiceUpdateJSONRPCError(err)
		return nil, nil, errors.Join(err, engineerrors.ErrPreDefinedJSONRPC)

	// All other errors are handled as undefined errors.
	case err != nil:
		ee.metrics.markForkchoiceUpdateUndefinedError(err)
		return nil, nil, err
	default:
		ee.metrics.markForkchoiceUpdateValid(
			req.State, hasPayloadAttributes, payloadID,
		)
	}

	// If we reached here, and we have a nil payload ID, we should log a
	// warning.
	if payloadID == nil && hasPayloadAttributes {
		ee.logger.Warn(
			"Received nil payload ID on VALID engine response",
			"head_eth1_hash", req.State.HeadBlockHash,
			"safe_eth1_hash", req.State.SafeBlockHash,
			"finalized_eth1_hash", req.State.FinalizedBlockHash,
		)
		return payloadID, latestValidHash, ErrNilPayloadOnValidResponse
	}

	return payloadID, latestValidHash, nil
}

// VerifyAndNotifyNewPayload verifies the new payload and notifies the
// execution client.
func (ee *Engine[
	ExecutionPayloadT, _, _, WithdrawalsT,
]) VerifyAndNotifyNewPayload(
	ctx context.Context,
	req *engineprimitives.NewPayloadRequest[
		ExecutionPayloadT, WithdrawalsT,
	],
) error {
	// Log the new payload attempt.
	ee.metrics.markNewPayloadCalled(
		req.ExecutionPayload.GetBlockHash(),
		req.ExecutionPayload.GetParentHash(),
		req.Optimistic,
	)

	// First we verify the block hash and versioned hashes are valid.
	//
	// TODO: is this required? Or will the EL handle this for us during
	// new payload?
	if err := req.HasValidVersionedAndBlockHashes(); err != nil {
		return err
	}

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
	// If we get accepted or syncing, we are going to optimistically
	// say that the block is valid, this is utilized during syncing
	// to allow the beacon-chain to continue processing blocks, while
	// its execution client is fetching things over it's p2p layer.
	case errors.IsAny(
		err,
		engineerrors.ErrAcceptedPayloadStatus,
		engineerrors.ErrSyncingPayloadStatus,
	):
		ee.metrics.markNewPayloadAcceptedSyncingPayloadStatus(
			req.ExecutionPayload.GetBlockHash(),
			req.ExecutionPayload.GetParentHash(),
			req.Optimistic,
		)

	// These two cases are semantically the same:
	// https://github.com/ethereum/execution-apis/issues/270
	case errors.IsAny(
		err,
		engineerrors.ErrInvalidPayloadStatus,
		engineerrors.ErrInvalidBlockHashPayloadStatus,
	):
		ee.metrics.markNewPayloadInvalidPayloadStatus(
			req.ExecutionPayload.GetBlockHash(),
			req.Optimistic,
		)

		// We want to return bad block irrespective of
		// if we are running in optimistic mode or not.
		//
		// TODO: should we still nillify the error in optimistic mode?
		return ErrBadBlockProduced

	case jsonrpc.IsPreDefinedError(err):
		// Protect against possible nil value.
		if lastValidHash == nil {
			lastValidHash = &common.ExecutionHash{}
		}

		ee.metrics.markNewPayloadJSONRPCError(
			req.ExecutionPayload.GetBlockHash(),
			*lastValidHash,
			req.Optimistic,
			err,
		)

		err = errors.Join(err, engineerrors.ErrPreDefinedJSONRPC)
	case err != nil:
		ee.metrics.markNewPayloadUndefinedError(
			req.ExecutionPayload.GetBlockHash(),
			req.Optimistic,
			err,
		)
	default:
		ee.metrics.markNewPayloadValid(
			req.ExecutionPayload.GetBlockHash(),
			req.ExecutionPayload.GetParentHash(),
			req.Optimistic,
		)
	}

	// Under the optimistic condition, we are fine ignoring the error. This
	// is mainly to allow us to safely call the execution client
	// during abci.FinalizeBlock. If we are in abci.FinalizeBlock and
	// we get an error here, we make the assumption that
	// abci.ProcessProposal
	// has deemed that the BeaconBlock containing the given ExecutionPayload
	// was marked as valid by an honest majority of validators, and we
	// don't want to halt the chain because of an error here.
	//
	// The practical reason we want to handle this edge case
	// is to protect against an awkward shutdown condition in which an
	// execution client dies between the end of abci.ProcessProposal
	// and the beginning of abci.FinalizeBlock. Without handling this case
	// it would cause a failure of abci.FinalizeBlock and a
	// "CONSENSUS FAILURE!!!!" at the CometBFT layer.
	if req.Optimistic {
		return nil
	}
	return err
}
