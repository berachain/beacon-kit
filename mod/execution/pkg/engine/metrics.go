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
	"strconv"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
)

// engineMetrics is a struct that contains metrics for the engine.
type engineMetrics struct {
	// TelemetrySink is the sink for the metrics.
	sink TelemetrySink
	// logger is the logger for the engineMetrics.
	logger log.Logger[any]
}

// newEngineMetrics creates a new engineMetrics.
func newEngineMetrics(
	sink TelemetrySink,
	logger log.Logger[any],
) *engineMetrics {
	return &engineMetrics{
		sink:   sink,
		logger: logger,
	}
}

// markNewPayloadCalled increments the counter for new payload calls.
func (em *engineMetrics) markNewPayloadCalled(
	payload ExecutionPayload,
	isOptimistic bool,
) {
	em.logger.Info(
		"calling new payload",
		"payload_block_hash", payload.GetBlockHash(),
		"payload_parent_block_hash", payload.GetParentHash(),
		"is_optimistic", isOptimistic,
	)

	em.sink.IncrementCounter(
		"beacon_kit.execution.engine.new_payload",
		"payload_block_hash", payload.GetBlockHash().Hex(),
		"payload_parent_block_hash", payload.GetParentHash().Hex(),
		"is_optimistic", strconv.FormatBool(isOptimistic),
	)
}

// markNewPayloadAcceptedSyncingPayloadStatus increments
// the counter for accepted syncing payload status.
func (em *engineMetrics) markNewPayloadAcceptedSyncingPayloadStatus(
	payloadHash common.ExecutionHash,
	parentHash common.ExecutionHash,
	isOptimistic bool,
) {
	em.errorLoggerFn(isOptimistic)(
		"received accepted syncing payload status",
		"payload_block_hash", payloadHash,
		"parent_hash", payloadHash,
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
func (em *engineMetrics) markNewPayloadInvalidPayloadStatus(
	payloadHash common.ExecutionHash,
	isOptimistic bool,
) {
	em.errorLoggerFn(isOptimistic)(
		"received invalid payload status during new payload call",
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
func (em *engineMetrics) markNewPayloadJSONRPCError(
	payloadHash common.ExecutionHash,
	lastValidHash common.ExecutionHash,
	isOptimistic bool,
	err error,
) {
	em.errorLoggerFn(isOptimistic)(
		"received JSON-RPC error during new payload call",
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
func (em *engineMetrics) markNewPayloadUndefinedError(
	payloadHash common.ExecutionHash,
	isOptimistic bool,
	err error,
) {
	em.errorLoggerFn(isOptimistic)(
		"received undefined error during new payload call",
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
func (em *engineMetrics) markNotifyForkchoiceUpdateCalled(
	state *engineprimitives.ForkchoiceStateV1,
	hasPayloadAttributes bool,
) {
	em.logger.Info("notifying forkchoice update",
		"head_eth1_hash", state.HeadBlockHash,
		"safe_eth1_hash", state.SafeBlockHash,
		"finalized_eth1_hash", state.FinalizedBlockHash,
		"has_attributes", hasPayloadAttributes,
	)

	em.sink.IncrementCounter(
		"beacon_kit.execution.engine.forkchoice_update",
		"has_payload_attributes", strconv.FormatBool(hasPayloadAttributes),
	)
}

// markForkchoiceUpdateAcceptedSyncing increments
// the counter for accepted syncing forkchoice updates.
func (em *engineMetrics) markForkchoiceUpdateAcceptedSyncing(
	state *engineprimitives.ForkchoiceStateV1,
) {
	em.errorLoggerFn(true)(
		"received accepted syncing payload status during forkchoice update call",
		"head_block_hash",
		state.HeadBlockHash,
		"safe_block_hash",
		state.SafeBlockHash,
		"finalized_block_hash",
		state.FinalizedBlockHash,
	)

	em.sink.IncrementCounter(
		"beacon_kit.execution.engine.forkchoice_update_accepted_syncing",
	)
}

// markForkchoiceUpdateInvalid increments the counter
// for invalid forkchoice updates.
func (em *engineMetrics) markForkchoiceUpdateInvalid(
	state *engineprimitives.ForkchoiceStateV1,
) {
	em.logger.Error(
		"received invalid payload status during forkchoice update call",
		"head_block_hash", state.HeadBlockHash,
		"safe_block_hash", state.SafeBlockHash,
		"finalized_block_hash", state.FinalizedBlockHash,
	)

	em.sink.IncrementCounter(
		"beacon_kit.execution.engine.forkchoice_update_invalid",
	)
}

// markForkchoiceUpdateJSONRPCError increments the counter for JSON-RPC errors
// during forkchoice updates.
func (em *engineMetrics) markForkchoiceUpdateJSONRPCError(err error) {
	em.logger.Error(
		"received json-rpc error during forkchoice update call",
		"error", err,
	)

	em.sink.IncrementCounter(
		"beacon_kit.execution.engine.forkchoice_update_json_rpc_error",
		"error", err.Error(),
	)
}

// markForkchoiceUpdateUndefinedError increments the counter for undefined
// errors during forkchoice updates.
func (em *engineMetrics) markForkchoiceUpdateUndefinedError(err error) {
	em.logger.Error(
		"received undefined execution engine error during forkchoice update call",
		"error",
		err,
	)

	em.sink.IncrementCounter(
		"beacon_kit.execution.engine.forkchoice_update_undefined_error",
		"error", err.Error(),
	)
}

// errorLoggerFn returns a logger fn based on the optimistic flag.
func (em *engineMetrics) errorLoggerFn(
	isOptimistic bool,
) func(msg string, keyVals ...any) {
	if isOptimistic {
		return em.logger.Warn
	}
	return em.logger.Error
}
