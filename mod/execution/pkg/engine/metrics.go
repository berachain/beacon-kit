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

	"github.com/berachain/beacon-kit/mod/log"
	engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine"
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

// MarkNewPayloadCalled increments the counter for new payload calls.
func (em *engineMetrics) MarkNewPayloadCalled(
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

// MarkNewPayloadAcceptedSyncingPayloadStatus increments
// the counter for accepted syncing payload status.
func (em *engineMetrics) MarkNewPayloadAcceptedSyncingPayloadStatus(
	payloadHash common.ExecutionHash,
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
		"payload_block_hash",
		payloadHash.Hex(),
		"is_optimistic",
		strconv.FormatBool(isOptimistic),
	)
}

// MarkNewPayloadInvalidPayloadStatus increments the counter
// for invalid payload status.
func (em *engineMetrics) MarkNewPayloadInvalidPayloadStatus(
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
		"payload_block_hash", payloadHash.Hex(),
		"is_optimistic", strconv.FormatBool(isOptimistic),
	)
}

// MarkNewPayloadJSONRPCError increments the counter for JSON-RPC errors.
func (em *engineMetrics) MarkNewPayloadJSONRPCError(
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
		"payload_block_hash", payloadHash.Hex(),
		"last_valid_hash", lastValidHash.Hex(),
		"is_optimistic", strconv.FormatBool(isOptimistic),
		"error", err.Error(),
	)
}

// MarkNewPayloadUndefinedError increments the counter for undefined errors.
func (em *engineMetrics) MarkNewPayloadUndefinedError(
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
		"payload_block_hash", payloadHash.Hex(),
		"is_optimistic", strconv.FormatBool(isOptimistic),
		"error", err.Error(),
	)
}

// MarkNotifyForkchoiceUpdateCalled increments the counter for
// notify forkchoice update calls.
func (em *engineMetrics) MarkNotifyForkchoiceUpdateCalled(
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
		"head_block_hash", state.HeadBlockHash.Hex(),
		"safe_block_hash", state.SafeBlockHash.Hex(),
		"finalized_block_hash", state.FinalizedBlockHash.Hex(),
		"has_payload_attributes", strconv.FormatBool(hasPayloadAttributes),
	)
}

// MarkForkchoiceUpdateAcceptedSyncing increments
// the counter for accepted syncing forkchoice updates.
func (em *engineMetrics) MarkForkchoiceUpdateAcceptedSyncing(
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
		"head_block_hash", state.HeadBlockHash.Hex(),
		"safe_block_hash", state.SafeBlockHash.Hex(),
		"finalized_block_hash", state.FinalizedBlockHash.Hex(),
	)
}

// MarkForkchoiceUpdateInvalid increments the counter
// for invalid forkchoice updates.
func (em *engineMetrics) MarkForkchoiceUpdateInvalid(
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
		"head_block_hash", state.HeadBlockHash.Hex(),
		"safe_block_hash", state.SafeBlockHash.Hex(),
		"finalized_block_hash", state.FinalizedBlockHash.Hex(),
	)
}

// MarkForkchoiceUpdateJSONRPCError increments the counter for JSON-RPC errors
// during forkchoice updates.
func (em *engineMetrics) MarkForkchoiceUpdateJSONRPCError(err error) {
	em.logger.Error(
		"received json-rpc error during forkchoice update call",
		"error", err,
	)

	em.sink.IncrementCounter(
		"beacon_kit.execution.engine.forkchoice_update_json_rpc_error",
		"error", err.Error(),
	)
}

// MarkForkchoiceUpdateUndefinedError increments the counter for undefined
// errors during forkchoice updates.
func (em *engineMetrics) MarkForkchoiceUpdateUndefinedError(err error) {
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
