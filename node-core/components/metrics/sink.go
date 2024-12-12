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

package metrics

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/hashicorp/go-metrics"
)

type TelemetrySink struct{}

// NewTelemetrySink creates a new TelemetrySink.
func NewTelemetrySink() TelemetrySink {
	return TelemetrySink{}
}

// IncrementCounter increments a counter metric identified by the provided
// keys.
func (TelemetrySink) IncrementCounter(key string, args ...string) {
	telemetry.IncrCounterWithLabels([]string{key}, 1, argsToLabels(args...))
}

// SetGauge sets a gauge metric to the specified value, identified by the
// provided keys.
func (TelemetrySink) SetGauge(key string, value int64, args ...string) {
	telemetry.SetGaugeWithLabels(
		[]string{key},
		float32(value),
		argsToLabels(args...),
	)
}

// MeasureSince measures the time since the provided start time and records
// the duration in a metric identified by the provided key.
func (TelemetrySink) MeasureSince(key string, start time.Time, args ...string) {
	if !telemetry.IsTelemetryEnabled() {
		return
	}

	// TODO: Make PR to SDK, currently this will not have any globalLabels.
	metrics.MeasureSinceWithLabels(
		[]string{key},
		start.UTC(),
		argsToLabels(args...),
	)
}

// argsToLabels converts a list of key-value pairs to a list of metrics labels.
//
//nolint:mnd // its okay.
func argsToLabels(args ...string) []metrics.Label {
	labels := make([]metrics.Label, len(args)/2)
	for i := 0; i < len(args); i += 2 {
		labels[i/2] = metrics.Label{
			Name:  args[i],
			Value: args[i+1],
		}
	}
	return labels
}

// NoOpTelemetrySink is a no-op implementation of the TelemetrySink interface.
type NoOpTelemetrySink struct{}

// NewNoOpTelemetrySink creates a new NoOpTelemetrySink.
func NewNoOpTelemetrySink() NoOpTelemetrySink {
	return NoOpTelemetrySink{}
}

// IncrementCounter is a no-op implementation of the TelemetrySink interface.
func (NoOpTelemetrySink) IncrementCounter(string, ...string) {}

// SetGauge is a no-op implementation of the TelemetrySink interface.
func (NoOpTelemetrySink) SetGauge(string, int64, ...string) {}

// MeasureSince is a no-op implementation of the TelemetrySink interface.
func (NoOpTelemetrySink) MeasureSince(string, time.Time, ...string) {}
