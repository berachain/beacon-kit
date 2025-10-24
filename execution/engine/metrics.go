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
	"fmt"
	"strconv"

	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	engineerrors "github.com/berachain/beacon-kit/engine-primitives/errors"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/observability/metrics"
	"github.com/berachain/beacon-kit/primitives/common"
)

// Metrics is a struct that contains metrics for the execution engine.
type Metrics struct {
	// New payload metrics
	NewPayload                      metrics.Counter
	NewPayloadValid                 metrics.Counter
	NewPayloadAcceptedPayloadStatus metrics.Counter
	NewPayloadSyncingPayloadStatus  metrics.Counter
	NewPayloadInvalidPayloadStatus  metrics.Counter
	NewPayloadNonFatalError         metrics.Counter
	NewPayloadFatalError            metrics.Counter
	NewPayloadUndefinedError        metrics.Counter

	// Forkchoice update metrics
	ForkchoiceUpdate               metrics.Counter
	ForkchoiceUpdateValid          metrics.Counter
	ForkchoiceUpdateSyncing        metrics.Counter
	ForkchoiceUpdateInvalid        metrics.Counter
	ForkchoiceUpdateFatalError     metrics.Counter
	ForkchoiceUpdateNonFatalError  metrics.Counter
	ForkchoiceUpdateUndefinedError metrics.Counter

	logger log.Logger
}

// NewMetrics returns a new Metrics instance.
// Metric names are kept identical to cosmos-sdk/telemetry output for Grafana compatibility.
//
//nolint:funlen
func NewMetrics(factory metrics.Factory, logger log.Logger) *Metrics {
	return &Metrics{
		// New payload metrics
		NewPayload: factory.NewCounter(
			metrics.CounterOpts{
				Name: "beacon_kit_execution_engine_new_payload",
				Help: "Number of new payload calls",
			},
			[]string{"payload_block_hash", "payload_parent_block_hash"},
		),
		NewPayloadValid: factory.NewCounter(
			metrics.CounterOpts{
				Name: "beacon_kit_execution_engine_new_payload_valid",
				Help: "Number of valid new payload responses",
			},
			nil,
		),
		NewPayloadAcceptedPayloadStatus: factory.NewCounter(
			metrics.CounterOpts{
				Name: "beacon_kit_execution_engine_new_payload_accepted_payload_status",
				Help: "Number of accepted payload status responses",
			},
			nil,
		),
		NewPayloadSyncingPayloadStatus: factory.NewCounter(
			metrics.CounterOpts{
				Name: "beacon_kit_execution_engine_new_payload_syncing_payload_status",
				Help: "Number of syncing payload status responses",
			},
			nil,
		),
		NewPayloadInvalidPayloadStatus: factory.NewCounter(
			metrics.CounterOpts{
				Name: "beacon_kit_execution_engine_new_payload_invalid_payload_status",
				Help: "Number of invalid payload status responses",
			},
			nil,
		),
		NewPayloadNonFatalError: factory.NewCounter(
			metrics.CounterOpts{
				Name: "beacon_kit_execution_engine_new_payload_non_fatal_error",
				Help: "Number of non-fatal errors during new payload",
			},
			[]string{"error"},
		),
		NewPayloadFatalError: factory.NewCounter(
			metrics.CounterOpts{
				Name: "beacon_kit_execution_engine_new_payload_fatal_error",
				Help: "Number of fatal errors during new payload",
			},
			[]string{"error"},
		),
		NewPayloadUndefinedError: factory.NewCounter(
			metrics.CounterOpts{
				Name: "beacon_kit_execution_engine_new_payload_undefined_error",
				Help: "Number of undefined errors during new payload",
			},
			[]string{"error"},
		),

		// Forkchoice update metrics
		ForkchoiceUpdate: factory.NewCounter(
			metrics.CounterOpts{
				Name: "beacon_kit_execution_engine_forkchoice_update",
				Help: "Number of forkchoice update calls",
			},
			[]string{"has_payload_attributes"},
		),
		ForkchoiceUpdateValid: factory.NewCounter(
			metrics.CounterOpts{
				Name: "beacon_kit_execution_engine_forkchoice_update_valid",
				Help: "Number of valid forkchoice update responses",
			},
			nil,
		),
		ForkchoiceUpdateSyncing: factory.NewCounter(
			metrics.CounterOpts{
				Name: "beacon_kit_execution_engine_forkchoice_update_syncing",
				Help: "Number of syncing forkchoice update responses",
			},
			[]string{"error"},
		),
		ForkchoiceUpdateInvalid: factory.NewCounter(
			metrics.CounterOpts{
				Name: "beacon_kit_execution_engine_forkchoice_update_invalid",
				Help: "Number of invalid forkchoice update responses",
			},
			[]string{"error"},
		),
		ForkchoiceUpdateFatalError: factory.NewCounter(
			metrics.CounterOpts{
				Name: "beacon_kit_execution_engine_forkchoice_update_fatal_error",
				Help: "Number of fatal errors during forkchoice update",
			},
			[]string{"error"},
		),
		ForkchoiceUpdateNonFatalError: factory.NewCounter(
			metrics.CounterOpts{
				Name: "beacon_kit_execution_engine_forkchoice_update_non_fatal_error",
				Help: "Number of non-fatal errors during forkchoice update",
			},
			[]string{"error"},
		),
		ForkchoiceUpdateUndefinedError: factory.NewCounter(
			metrics.CounterOpts{
				Name: "beacon_kit_execution_engine_forkchoice_update_undefined_error",
				Help: "Number of undefined errors during forkchoice update",
			},
			[]string{"error"},
		),

		logger: logger,
	}
}

// markNewPayloadCalled increments the counter for new payload calls.
func (m *Metrics) markNewPayloadCalled(
	payloadHash common.ExecutionHash,
	parentHash common.ExecutionHash,
) {
	m.NewPayload.With(
		"payload_block_hash", payloadHash.Hex(),
		"payload_parent_block_hash", parentHash.Hex(),
	).Add(1)
}

// markNewPayloadValid increments the counter for valid payloads.
func (m *Metrics) markNewPayloadValid(
	payloadHash common.ExecutionHash,
	parentHash common.ExecutionHash,
) {
	m.logger.Debug(
		"Inserted new payload into execution chain",
		"payload_block_hash", payloadHash,
		"payload_parent_block_hash", parentHash,
	)
	m.NewPayloadValid.Add(1)
}

// markNewPayloadAcceptedSyncingPayloadStatus increments
// the counter for accepted syncing payload status.
func (m *Metrics) markNewPayloadAcceptedSyncingPayloadStatus(
	errStatus error,
	payloadHash common.ExecutionHash,
	parentHash common.ExecutionHash,
) {
	status := "accepted"
	if errors.Is(errStatus, engineerrors.ErrSyncingPayloadStatus) {
		status = "syncing"
	}
	m.logger.Warn(
		fmt.Sprintf("Received %s payload status during new payload. Awaiting execution client to finish sync.", status),
		"payload_block_hash", payloadHash,
		"parent_hash", parentHash,
	)

	if status == "accepted" {
		m.NewPayloadAcceptedPayloadStatus.Add(1)
	} else {
		m.NewPayloadSyncingPayloadStatus.Add(1)
	}
}

// markNewPayloadInvalidPayloadStatus increments the counter
// for invalid payload status.
func (m *Metrics) markNewPayloadInvalidPayloadStatus(
	payloadHash common.ExecutionHash,
) {
	m.logger.Error(
		"Received invalid payload status during new payload call",
		"payload_block_hash", payloadHash,
		"parent_hash", payloadHash,
	)
	m.NewPayloadInvalidPayloadStatus.Add(1)
}

// markNewPayloadNonFatalError increments the counter for non-fatal errors.
func (m *Metrics) markNewPayloadNonFatalError(
	payloadHash common.ExecutionHash,
	lastValidHash common.ExecutionHash,
	err error,
) {
	m.logger.Error(
		"Received non-fatal error during new payload call",
		"payload_block_hash", payloadHash,
		"parent_hash", payloadHash,
		"last_valid_hash", lastValidHash,
		"error", err,
	)
	m.NewPayloadNonFatalError.With("error", err.Error()).Add(1)
}

// markNewPayloadFatalError increments the counter for fatal errors.
func (m *Metrics) markNewPayloadFatalError(
	payloadHash common.ExecutionHash,
	lastValidHash common.ExecutionHash,
	err error,
) {
	m.logger.Error(
		"Received fatal error during new payload call",
		"payload_block_hash", payloadHash,
		"parent_hash", payloadHash,
		"last_valid_hash", lastValidHash,
		"error", err,
	)
	m.NewPayloadFatalError.With("error", err.Error()).Add(1)
}

// markNewPayloadUndefinedError increments the counter for undefined errors.
func (m *Metrics) markNewPayloadUndefinedError(
	payloadHash common.ExecutionHash,
	err error,
) {
	m.logger.Error(
		"Received undefined error during new payload call",
		"payload_block_hash", payloadHash,
		"parent_hash", payloadHash,
		"error", err,
	)
	m.NewPayloadUndefinedError.With("error", err.Error()).Add(1)
}

// markNotifyForkchoiceUpdateCalled increments the counter for
// notify forkchoice update calls.
func (m *Metrics) markNotifyForkchoiceUpdateCalled(
	hasPayloadAttributes bool,
) {
	m.ForkchoiceUpdate.With("has_payload_attributes", strconv.FormatBool(hasPayloadAttributes)).Add(1)
}

// markForkchoiceUpdateValid increments the counter for valid forkchoice updates.
func (m *Metrics) markForkchoiceUpdateValid(
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
	m.logger.Debug("Forkchoice updated", args...)
	m.ForkchoiceUpdateValid.Add(1)
}

// markForkchoiceUpdateSyncing increments the counter for syncing forkchoice updates.
func (m *Metrics) markForkchoiceUpdateSyncing(
	state *engineprimitives.ForkchoiceStateV1,
	err error,
) {
	m.logger.Warn(
		"Received syncing payload status during forkchoice update. Awaiting execution client to finish sync.",
		"head_block_hash",
		state.HeadBlockHash,
		"safe_block_hash",
		state.SafeBlockHash,
		"finalized_block_hash",
		state.FinalizedBlockHash,
	)
	m.ForkchoiceUpdateSyncing.With("error", err.Error()).Add(1)
}

// markForkchoiceUpdateInvalid increments the counter for invalid forkchoice updates.
func (m *Metrics) markForkchoiceUpdateInvalid(
	state *engineprimitives.ForkchoiceStateV1,
	err error,
) {
	m.logger.Error(
		"Received invalid payload status during forkchoice update call",
		"head_block_hash", state.HeadBlockHash,
		"safe_block_hash", state.SafeBlockHash,
		"finalized_block_hash", state.FinalizedBlockHash,
		"error", err,
	)
	m.ForkchoiceUpdateInvalid.With("error", err.Error()).Add(1)
}

// markForkchoiceUpdateFatalError increments the counter for fatal errors
// during forkchoice updates.
func (m *Metrics) markForkchoiceUpdateFatalError(err error) {
	m.logger.Error(
		"Received fatal error during forkchoice update call",
		"error", err,
	)
	m.ForkchoiceUpdateFatalError.With("error", err.Error()).Add(1)
}

// markForkchoiceUpdateNonFatalError increments the counter for non-fatal errors
// during forkchoice updates.
func (m *Metrics) markForkchoiceUpdateNonFatalError(err error) {
	m.logger.Error(
		"Received non-fatal error during forkchoice update call",
		"error", err,
	)
	m.ForkchoiceUpdateNonFatalError.With("error", err.Error()).Add(1)
}

// markForkchoiceUpdateUndefinedError increments the counter for undefined
// errors during forkchoice updates.
func (m *Metrics) markForkchoiceUpdateUndefinedError(err error) {
	m.logger.Error(
		"Received undefined execution engine error during forkchoice update call",
		"error",
		err,
	)
	m.ForkchoiceUpdateUndefinedError.With("error", err.Error()).Add(1)
}
