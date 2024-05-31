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

package blockchain

import (
	"strconv"
	"time"
)

// chainMetrics is a struct that contains metrics for the chain.
type chainMetrics struct {
	// sink is the sink for the metrics.
	sink TelemetrySink
}

// newChainMetrics creates a new chainMetrics.
func newChainMetrics(
	sink TelemetrySink,
) *chainMetrics {
	return &chainMetrics{
		sink: sink,
	}
}

// measureStateTransitionDuration measures the time to process
// the state transition for a block.
func (cm *chainMetrics) measureStateTransitionDuration(
	start time.Time, skipPayloadVerification bool,
) {
	cm.sink.MeasureSince(
		"beacon_kit.beacon.blockchain.state_transition_duration",
		start,
		"skip_payload_verification",
		strconv.FormatBool(skipPayloadVerification),
	)
}

// measureBlobProcessingDuration measures the time to process
// the blobs for a block.
func (cm *chainMetrics) measureBlobProcessingDuration(start time.Time) {
	cm.sink.MeasureSince(
		"beacon_kit.beacon.blockchain.blob_processing_duration", start,
	)
}
