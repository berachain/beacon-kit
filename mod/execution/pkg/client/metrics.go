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

package client

import (
	"time"

	"github.com/berachain/beacon-kit/mod/log"
)

// clientMetrics is a struct that contains metrics for the engine.
type clientMetrics struct {
	// TelemetrySink is the sink for the metrics.
	sink TelemetrySink
	// logger is the logger for the engineMetrics.
	logger log.Logger[any]
}

// newClientMetrics creates a new engineMetrics.
func newClientMetrics(
	sink TelemetrySink,
	logger log.Logger[any],
) *clientMetrics {
	return &clientMetrics{
		sink:   sink,
		logger: logger,
	}
}

// measureForkchoiceUpdateDuration measures the duration of the forkchoice
// update.
func (cm *clientMetrics) measureForkchoiceUpdateDuration(startTime time.Time) {
	// TODO: Add Labels.
	cm.sink.MeasureSince(
		"beacon_kit.execution.client.forkchoice_update_duration",
		startTime,
	)
}

// measureNewPayloadDuration measures the duration of the new payload.
func (cm *clientMetrics) measureNewPayloadDuration(startTime time.Time) {
	// TODO: Add Labels.
	cm.sink.MeasureSince(
		"beacon_kit.execution.client.new_payload_duration",
		startTime,
	)
}

// measureGetPayloadDuration measures the duration of the get payload.
func (cm *clientMetrics) measureGetPayloadDuration(startTime time.Time) {
	// TODO: Add Labels.
	cm.sink.MeasureSince(
		"beacon_kit.execution.client.get_payload_duration",
		startTime,
	)
}

// incrementForkchoiceUpdateTimeout increments the timeout counter
// for forkchoice update.
func (cm *clientMetrics) incrementForkchoiceUpdateTimeout() {
	cm.incrementTimeoutCounter(
		"beacon_kit.execution.client.forkchoice_update_duration")
}

// incrementNewPayloadTimeout increments the timeout counter for
// new payload.
func (cm *clientMetrics) incrementNewPayloadTimeout() {
	cm.incrementTimeoutCounter(
		"beacon_kit.execution.client.new_payload_duration")
}

// incrementGetPayloadTimeout increments the timeout counter for
// get payload.
func (cm *clientMetrics) incrementGetPayloadTimeout() {
	cm.incrementTimeoutCounter(
		"beacon_kit.execution.client.get_payload_duration")
}

// incrementHTTPTimeout increments the timeout counter for HTTP.
func (cm *clientMetrics) incrementHTTPTimeoutCounter() {
	cm.incrementTimeoutCounter("beacon_kit.execution.client.http")
}

// incrementTimeoutCounter increments the timeout counter for
// the given metric.
func (cm *clientMetrics) incrementTimeoutCounter(metricName string) {
	cm.sink.IncrementCounter(metricName + "_timeout")
}

// incrementParseErrorCounter increments the parse error counter
// for the given metric.
func (cm *clientMetrics) incrementParseErrorCounter() {
	cm.sink.IncrementCounter("beacon_kit.execution.client.parse_error")
}

// incrementInvalidRequestCounter increments the invalid request counter
// for the given metric.
func (cm *clientMetrics) incrementInvalidRequestCounter() {
	cm.incrementErrorCounter("beacon_kit.execution.client.invalid_request")
}

// incrementMethodNotFoundCounter increments the method not found counter
// for the given metric.
func (cm *clientMetrics) incrementMethodNotFoundCounter() {
	cm.incrementErrorCounter("beacon_kit.execution.client.method_not_found")
}

// incrementInvalidParamsCounter increments the invalid params counter
// for the given metric.
func (cm *clientMetrics) incrementInvalidParamsCounter() {
	cm.incrementErrorCounter("beacon_kit.execution.client.invalid_params")
}

// incrementInternalErrorCounter increments the internal error counter
// for the given metric.
func (cm *clientMetrics) incrementInternalErrorCounter() {
	cm.incrementErrorCounter("beacon_kit.execution.client.internal_error")
}

// incrementUnknownPayloadErrorCounter increments the unknown payload error
// counter
// for the given metric.
func (cm *clientMetrics) incrementUnknownPayloadErrorCounter() {
	cm.incrementErrorCounter(
		"beacon_kit.execution.client.unknown_payload_error",
	)
}

// incrementInvalidForkchoiceStateCounter increments the invalid forkchoice
// state counter
// for the given metric.
func (cm *clientMetrics) incrementInvalidForkchoiceStateCounter() {
	cm.incrementErrorCounter(
		"beacon_kit.execution.client.invalid_forkchoice_state",
	)
}

// incrementInvalidPayloadAttributesCounter increments the invalid payload
// attributes counter
// for the given metric.
func (cm *clientMetrics) incrementInvalidPayloadAttributesCounter() {
	cm.incrementErrorCounter(
		"beacon_kit.execution.client.invalid_payload_attributes",
	)
}

// incrementRequestTooLargeCounter increments the request too large counter
// for the given metric.
func (cm *clientMetrics) incrementRequestTooLargeCounter() {
	cm.incrementErrorCounter("beacon_kit.execution.client.request_too_large")
}

// incrementInternalServerErrorCounter increments the internal server error
// counter
// for the given metric.
func (cm *clientMetrics) incrementInternalServerErrorCounter() {
	cm.incrementErrorCounter(
		"beacon_kit.execution.client.internal_server_error",
	)
}

// incrementErrorCounter increments the error counter for
// the given metric.
func (cm *clientMetrics) incrementErrorCounter(metricName string) {
	cm.sink.IncrementCounter(metricName)
}
