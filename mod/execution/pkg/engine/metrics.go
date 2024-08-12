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

package engine

import (
	"strconv"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
)

// engineMetrics is a struct that contains metrics for the engine.
type engineMetrics[PayloadIDT any] struct {
	// TelemetrySink is the sink for the metrics.
	sink TelemetrySink
	// logger is the logger for the engineMetrics.
	logger log.Logger[any]
}

// newEngineMetrics creates a new engineMetrics.
func newEngineMetrics[PayloadIDT any](
	sink TelemetrySink,
	logger log.Logger[any],
) *engineMetrics[PayloadIDT] {
	return &engineMetrics[PayloadIDT]{
		sink:   sink,
		logger: logger,
	}
}

// markNewPayloadCalled increments the counter for new payload calls.
func (em *engineMetrics[_]) markNewPayloadCalled(
	payloadHash common.ExecutionHash,
	parentHash common.ExecutionHash,
	isOptimistic bool,
) {
	em.sink.IncrementCounter(
		"beacon_kit.execution.engine.new_payload",
		"payload_block_hash", payloadHash.Hex(),
		"payload_parent_block_hash", parentHash.Hex(),
		"is_optimistic", strconv.FormatBool(isOptimistic),
	)
}

// markNewPayloadValid increments the counter for valid payloads.
func (em *engineMetrics[_]) markNewPayloadValid(
	payloadHash common.ExecutionHash,
	parentHash common.ExecutionHash,
	isOptimistic bool,
) {
	em.logger.Info(
		"Inserted new payload into execution chain",
		"payload_block_hash", payloadHash,
		"payload_parent_block_hash", parentHash,
		"is_optimistic", isOptimistic,
	)

	em.sink.IncrementCounter(
		"beacon_kit.execution.engine.new_payload_valid",
		"is_optimistic", strconv.FormatBool(isOptimistic),
	)
}

// markNewPayloadAcceptedSyncingPayloadStatus increments
// the counter for accepted syncing payload status.
func (em *engineMetrics[_]) markNewPayloadAcceptedSyncingPayloadStatus(
	payloadHash common.ExecutionHash,
	parentHash common.ExecutionHash,
	isOptimistic bool,
) {
	em.errorLoggerFn(isOptimistic)(
		"Received accepted syncing payload status",
		"payload_block_hash", payloadHash,
		"parent_hash", parentHash,
		"is_optimistic", isOptimistic,
	)

	em.sink.IncrementCounter(
		"beacon_kit.execution.engine.new_payload_accepted_syncing_payload_status",
		"is_optimistic",
		strconv.FormatBool(isOptimistic),
	)
}

// markNewPayloadInvalidPayloadStatus increments the counter
// for invalid payload status.
func (em *engineMetrics[_]) markNewPayloadInvalidPayloadStatus(
	payloadHash common.ExecutionHash,
	isOptimistic bool,
) {
	em.errorLoggerFn(isOptimistic)(
		"Received invalid payload status during new payload call",
		"payload_block_hash", payloadHash,
		"parent_hash", payloadHash,
		"is_optimistic", isOptimistic,
	)

	em.sink.IncrementCounter(
		"beacon_kit.execution.engine.new_payload_invalid_payload_status",
		"is_optimistic", strconv.FormatBool(isOptimistic),
	)
}

// markNewPayloadJSONRPCError increments the counter for JSON-RPC errors.
func (em *engineMetrics[_]) markNewPayloadJSONRPCError(
	payloadHash common.ExecutionHash,
	lastValidHash common.ExecutionHash,
	isOptimistic bool,
	err error,
) {
	em.errorLoggerFn(isOptimistic)(
		"Received JSON-RPC error during new payload call",
		"payload_block_hash", payloadHash,
		"parent_hash", payloadHash,
		"last_valid_hash", lastValidHash,
		"is_optimistic", isOptimistic,
		"error", err,
	)

	em.sink.IncrementCounter(
		"beacon_kit.execution.engine.new_payload_json_rpc_error",
		"is_optimistic", strconv.FormatBool(isOptimistic),
		"error", err.Error(),
	)
}

// markNewPayloadUndefinedError increments the counter for undefined errors.
func (em *engineMetrics[_]) markNewPayloadUndefinedError(
	payloadHash common.ExecutionHash,
	isOptimistic bool,
	err error,
) {
	em.errorLoggerFn(isOptimistic)(
		"Received undefined error during new payload call",
		"payload_block_hash", payloadHash,
		"parent_hash", payloadHash,
		"is_optimistic", isOptimistic,
		"error", err,
	)

	em.sink.IncrementCounter(
		"beacon_kit.execution.engine.new_payload_undefined_error",
		"is_optimistic", strconv.FormatBool(isOptimistic),
		"error", err.Error(),
	)
}

// markNotifyForkchoiceUpdateCalled increments the counter for
// notify forkchoice update calls.
func (em *engineMetrics[_]) markNotifyForkchoiceUpdateCalled(
	hasPayloadAttributes bool,
) {
	em.sink.IncrementCounter(
		"beacon_kit.execution.engine.forkchoice_update",
		"has_payload_attributes", strconv.FormatBool(hasPayloadAttributes),
	)
}

// markForkchoiceUpdateValid increments the counter for valid forkchoice
// updates.
func (em *engineMetrics[PayloadIDT]) markForkchoiceUpdateValid(
	state *engineprimitives.ForkchoiceStateV1,
	hasPayloadAttributes bool,
	payloadID PayloadIDT,
) {
	args := []any{
		"head_block_hash", state.HeadBlockHash,
		"safe_block_hash", state.SafeBlockHash,
		"finalized_block_hash", state.FinalizedBlockHash,
		"with_attributes", hasPayloadAttributes,
	}
	if hasPayloadAttributes {
		args = append(args, "payload_id", payloadID)
	}
	em.logger.Info("Forkchoice updated", args...)

	em.sink.IncrementCounter(
		"beacon_kit.execution.engine.forkchoice_update_valid",
	)
}

// markForkchoiceUpdateAcceptedSyncing increments
// the counter for accepted syncing forkchoice updates.
func (em *engineMetrics[_]) markForkchoiceUpdateAcceptedSyncing(
	state *engineprimitives.ForkchoiceStateV1,
	err error,
) {
	em.errorLoggerFn(true)(
		"Received accepted syncing payload status during forkchoice update call",
		"head_block_hash",
		state.HeadBlockHash,
		"safe_block_hash",
		state.SafeBlockHash,
		"finalized_block_hash",
		state.FinalizedBlockHash,
		"error",
		err,
	)

	em.sink.IncrementCounter(
		"beacon_kit.execution.engine.forkchoice_update_accepted_syncing",
		"error",
		err.Error(),
	)
}

// markForkchoiceUpdateInvalid increments the counter
// for invalid forkchoice updates.
func (em *engineMetrics[_]) markForkchoiceUpdateInvalid(
	state *engineprimitives.ForkchoiceStateV1,
	err error,
) {
	em.logger.Error(
		"Received invalid payload status during forkchoice update call",
		"head_block_hash", state.HeadBlockHash,
		"safe_block_hash", state.SafeBlockHash,
		"finalized_block_hash", state.FinalizedBlockHash,
		"error", err,
	)

	em.sink.IncrementCounter(
		"beacon_kit.execution.engine.forkchoice_update_invalid",
		"error",
		err.Error(),
	)
}

// markForkchoiceUpdateJSONRPCError increments the counter for JSON-RPC errors
// during forkchoice updates.
func (em *engineMetrics[_]) markForkchoiceUpdateJSONRPCError(err error) {
	em.logger.Error(
		"Received json-rpc error during forkchoice update call",
		"error", err,
	)

	em.sink.IncrementCounter(
		"beacon_kit.execution.engine.forkchoice_update_json_rpc_error",
		"error", err.Error(),
	)
}

// markForkchoiceUpdateUndefinedError increments the counter for undefined
// errors during forkchoice updates.
func (em *engineMetrics[_]) markForkchoiceUpdateUndefinedError(err error) {
	em.logger.Error(
		"Received undefined execution engine error during forkchoice update call",
		"error",
		err,
	)

	em.sink.IncrementCounter(
		"beacon_kit.execution.engine.forkchoice_update_undefined_error",
		"error", err.Error(),
	)
}

// errorLoggerFn returns a logger fn based on the optimistic flag.
func (em *engineMetrics[_]) errorLoggerFn(
	isOptimistic bool,
) func(msg string, keyVals ...any) {
	if isOptimistic {
		return em.logger.Warn
	}
	return em.logger.Error
}
