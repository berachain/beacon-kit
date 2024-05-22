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
//
//nolint:mnd // trivial.
func (TelemetrySink) IncrementCounter(key string, args ...string) {
	labels := make([]metrics.Label, len(args)/2)
	for i := 0; i < len(args); i += 2 {
		labels[i/2] = metrics.Label{
			Name:  args[i],
			Value: args[i+1],
		}
	}
	telemetry.IncrCounterWithLabels([]string{key}, 1, labels)
}

// SetGauge sets a gauge metric to the specified value, identified by the
// provided keys.
//
//nolint:mnd // trivial.
func (TelemetrySink) SetGauge(key string, value int64, args ...string) {
	labels := make([]metrics.Label, len(args)/2)
	for i := 0; i < len(args); i += 2 {
		labels[i/2] = metrics.Label{
			Name:  args[i],
			Value: args[i+1],
		}
	}
	telemetry.SetGaugeWithLabels([]string{key}, float32(value), labels)
}

// MeasureSince measures the time since the provided start time and records
// the duration in a metric identified by the provided key.
func (TelemetrySink) MeasureSince(key string, start time.Time) {
	telemetry.MeasureSince(start, key)
}
