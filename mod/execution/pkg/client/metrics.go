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

// MeasureForkchoiceUpdateDuration measures the duration of the forkchoice
// update.
func (cm *clientMetrics) MeasureForkchoiceUpdateDuration(startTime time.Time) {
	// TODO: Add Labels.
	cm.sink.MeasureSince(
		"beacon-kit.execution.client.forkchoice_update_duration",
		startTime,
	)
}

// MeasureNewPayloadDuration measures the duration of the new payload.
func (cm *clientMetrics) MeasureNewPayloadDuration(startTime time.Time) {
	// TODO: Add Labels.
	cm.sink.MeasureSince(
		"beacon-kit.execution.client.new_payload_duration",
		startTime,
	)
}

// MeasureGetPayloadDuration measures the duration of the get payload.
func (cm *clientMetrics) MeasureGetPayloadDuration(startTime time.Time) {
	// TODO: Add Labels.
	cm.sink.MeasureSince(
		"beacon-kit.execution.client.get_payload_duration",
		startTime,
	)
}
