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
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package engine

import (
	"fmt"
	"strconv"

	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	engineerrors "github.com/berachain/beacon-kit/engine-primitives/errors"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/primitives/common"
)

// engineMetrics is a struct that contains metrics for the engine.
type engineMetrics struct {
	// TelemetrySink is the sink for the metrics.
	sink TelemetrySink
	// logger is the logger for the engineMetrics.
	logger log.Logger
}

// newEngineMetrics creates a new engineMetrics.
func newEngineMetrics(
	sink TelemetrySink,
	logger log.Logger,
) *engineMetrics {
	return &engineMetrics{
		sink:   sink,
		logger: logger,
	}
}

// markNewPayloadCalled increments the counter for new payload calls.
func (em *engineMetrics) markNewPayloadCalled(
	payloadHash common.ExecutionHash,
	parentHash common.ExecutionHash,
) {
	em.sink.IncrementCounter(
		"beacon_kit.execution.engine.new_payload",
		"payload_block_hash", payloadHash.Hex(),
		"payload_parent_block_hash", parentHash.Hex(),
	)
}

// markNewPayloadValid increments the counter for valid payloads.
func (em *engineMetrics) markNewPayloadValid(
	payloadHash common.ExecutionHash,
	parentHash common.ExecutionHash,
) {
	em.logger.Debug(
		"Inserted new payload into execution chain",
		"payload_block_hash", payloadHash,
		"payload_parent_block_hash", parentHash,
	)

	em.sink.IncrementCounter(
		"beacon_kit.execution.engine.new_payload_valid",
	)
}

// markNewPayloadAcceptedSyncingPayloadStatus increments
// the counter for accepted syncing payload status.
func (em *engineMetrics) markNewPayloadAcceptedSyncingPayloadStatus(
	errStatus error,
	payloadHash common.ExecutionHash,
	parentHash common.ExecutionHash,
) {
	status := "accepted"
	if errors.Is(errStatus, engineerrors.ErrSyncingPayloadStatus) {
		status = "syncing"
	}
	em.logger.Warn(
		fmt.Sprintf("Received %s payload status during new payload. Awaiting execution client to finish sync.", status),
		"payload_block_hash", payloadHash,
		"parent_hash", parentHash,
	)

	em.sink.IncrementCounter(
		fmt.Sprintf("beacon_kit.execution.engine.new_payload_%s_payload_status", status),
	)
}

// markNewPayloadInvalidPayloadStatus increments the counter
// for invalid payload status.
func (em *engineMetrics) markNewPayloadInvalidPayloadStatus(
	payloadHash common.ExecutionHash,
) {
	em.logger.Error(
		"Received invalid payload status during new payload call",
		"payload_block_hash", payloadHash,
		"parent_hash", payloadHash,
	)

	em.sink.IncrementCounter(
		"beacon_kit.execution.engine.new_payload_invalid_payload_status",
	)
}

// markNewPayloadFatalError increments the counter for JSON-RPC errors.
func (em *engineMetrics) markNewPayloadNonFatalError(
	payloadHash common.ExecutionHash,
	lastValidHash common.ExecutionHash,
	err error,
) {
	em.logger.Error(
		"Received non-fatal error during new payload call",
		"payload_block_hash", payloadHash,
		"parent_hash", payloadHash,
		"last_valid_hash", lastValidHash,
		"error", err,
	)

	em.sink.IncrementCounter(
		"beacon_kit.execution.engine.new_payload_non_fatal_error",
		"error", err.Error(),
	)
}

// markNewPayloadFatalError increments the counter for JSON-RPC errors.
func (em *engineMetrics) markNewPayloadFatalError(
	payloadHash common.ExecutionHash,
	lastValidHash common.ExecutionHash,
	err error,
) {
	em.logger.Error(
		"Received fatal error during new payload call",
		"payload_block_hash", payloadHash,
		"parent_hash", payloadHash,
		"last_valid_hash", lastValidHash,
		"error", err,
	)

	em.sink.IncrementCounter(
		"beacon_kit.execution.engine.new_payload_fatal_error",
		"error", err.Error(),
	)
}

// markNewPayloadUndefinedError increments the counter for undefined errors.
func (em *engineMetrics) markNewPayloadUndefinedError(
	payloadHash common.ExecutionHash,
	err error,
) {
	em.logger.Error(
		"Received undefined error during new payload call",
		"payload_block_hash", payloadHash,
		"parent_hash", payloadHash,
		"error", err,
	)

	em.sink.IncrementCounter(
		"beacon_kit.execution.engine.new_payload_undefined_error",
		"error", err.Error(),
	)
}

// markNotifyForkchoiceUpdateCalled increments the counter for
// notify forkchoice update calls.
func (em *engineMetrics) markNotifyForkchoiceUpdateCalled(
	hasPayloadAttributes bool,
) {
	em.sink.IncrementCounter(
		"beacon_kit.execution.engine.forkchoice_update",
		"has_payload_attributes", strconv.FormatBool(hasPayloadAttributes),
	)
}

// markForkchoiceUpdateValid increments the counter for valid forkchoice
// updates.
func (em *engineMetrics) markForkchoiceUpdateValid(
	state *engineprimitives.ForkchoiceStateV1,
	hasPayloadAttributes bool,
	payloadID *engineprimitives.PayloadID,
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
	em.logger.Debug("Forkchoice updated", args...)

	em.sink.IncrementCounter(
		"beacon_kit.execution.engine.forkchoice_update_valid",
	)
}

// markForkchoiceUpdateSyncing increments
// the counter for accepted syncing forkchoice updates.
func (em *engineMetrics) markForkchoiceUpdateSyncing(
	state *engineprimitives.ForkchoiceStateV1,
	err error,
) {
	em.logger.Warn(
		"Received syncing payload status during forkchoice update. Awaiting execution client to finish sync.",
		"head_block_hash",
		state.HeadBlockHash,
		"safe_block_hash",
		state.SafeBlockHash,
		"finalized_block_hash",
		state.FinalizedBlockHash,
	)

	em.sink.IncrementCounter(
		"beacon_kit.execution.engine.forkchoice_update_syncing",
		"error",
		err.Error(),
	)
}

// markForkchoiceUpdateInvalid increments the counter
// for invalid forkchoice updates.
func (em *engineMetrics) markForkchoiceUpdateInvalid(
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

// markForkchoiceUpdateFatalError increments the counter for JSON-RPC errors
// during forkchoice updates.
func (em *engineMetrics) markForkchoiceUpdateFatalError(err error) {
	em.logger.Error(
		"Received fatal error during forkchoice update call",
		"error", err,
	)

	em.sink.IncrementCounter(
		"beacon_kit.execution.engine.forkchoice_update_fatal_error",
		"error", err.Error(),
	)
}

// markForkchoiceUpdateNonFatalError increments the counter for JSON-RPC errors
// during forkchoice updates.
func (em *engineMetrics) markForkchoiceUpdateNonFatalError(err error) {
	em.logger.Error(
		"Received non-fatal error during forkchoice update call",
		"error", err,
	)

	em.sink.IncrementCounter(
		"beacon_kit.execution.engine.forkchoice_update_non_fatal_error",
		"error", err.Error(),
	)
}

// markForkchoiceUpdateUndefinedError increments the counter for undefined
// errors during forkchoice updates.
func (em *engineMetrics) markForkchoiceUpdateUndefinedError(err error) {
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
