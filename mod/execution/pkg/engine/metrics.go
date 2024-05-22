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

import "github.com/berachain/beacon-kit/mod/primitives/pkg/common"

// engineMetrics is a struct that contains metrics for the engine.
type engineMetrics struct {
	TelemetrySink
}

// newEngineMetrics creates a new engineMetrics.
func newEngineMetrics(
	ts TelemetrySink,
) *engineMetrics {
	return &engineMetrics{
		TelemetrySink: ts,
	}
}

// MarkAcceptedSyncingPayloadStatus increments
// the counter for accepted syncing payload status.
func (em *engineMetrics) MarkNewPayloadAcceptedSyncingPayloadStatus(
	requestedHead common.ExecutionHash,
) {
	em.TelemetrySink.IncrementCounter(
		"beacon-kit.execution.engine.new_payload_accepted_syncing_payload_status",
	)
}

// MarkInvalidPayloadStatus increments the counter
// for invalid payload status.
func (em *engineMetrics) MarkNewPayloadInvalidPayloadStatus() {
	em.TelemetrySink.IncrementCounter(
		"beacon-kit.execution.engine.new_payload_invalid_payload_status",
	)
}

// MarkJSONRPCError increments the counter for JSON-RPC errors.
func (em *engineMetrics) MarkNewPayloadJSONRPCError() {
	em.TelemetrySink.IncrementCounter(
		"beacon-kit.execution.engine.new_payload_json_rpc_error",
	)
}

// MarkNewPayloadUndefinedError increments the counter for undefined errors.
func (em *engineMetrics) MarkNewPayloadUndefinedError() {
	em.TelemetrySink.IncrementCounter(
		"beacon-kit.execution.engine.new_payload_undefined_error",
	)
}

// MarkForkchoiceUpdateAcceptedSyncing increments
// the counter for accepted syncing forkchoice updates.
func (em *engineMetrics) MarkForkchoiceUpdateAcceptedSyncing() {
	em.TelemetrySink.IncrementCounter(
		"beacon-kit.execution.engine.forkchoice_update_accepted_syncing",
	)
}

// MarkForkchoiceUpdateInvalid increments the counter
// for invalid forkchoice updates.
func (em *engineMetrics) MarkForkchoiceUpdateInvalid() {
	em.TelemetrySink.IncrementCounter(
		"beacon-kit.execution.engine.forkchoice_update_invalid",
	)
}

// MarkForkchoiceUpdateJSONRPCError increments the counter for JSON-RPC errors
// during forkchoice updates.
func (em *engineMetrics) MarkForkchoiceUpdateJSONRPCError() {
	em.TelemetrySink.IncrementCounter(
		"beacon-kit.execution.engine.forkchoice_update_json_rpc_error",
	)
}

// MarkForkchoiceUpdateUndefinedError increments the counter for undefined
// errors during forkchoice updates.
func (em *engineMetrics) MarkForkchoiceUpdateUndefinedError() {
	em.TelemetrySink.IncrementCounter(
		"beacon-kit.execution.engine.forkchoice_update_undefined_error",
	)
}
