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

package metrics

// Counter represents a monotonically increasing metric.
type Counter interface {
	// Add increments the counter by the given delta.
	// Delta must be non-negative.
	Add(delta float64)

	// With returns a new Counter with the given label values applied.
	// Label values are provided as key-value pairs.
	// If the number of label values is odd, "unknown" is appended.
	With(labelValues ...string) Counter
}

// Gauge represents a metric that can increase or decrease.
type Gauge interface {
	// Set sets the gauge to the given value.
	Set(value float64)

	// Add increments or decrements the gauge by the given delta.
	// Delta can be positive or negative.
	Add(delta float64)

	// With returns a new Gauge with the given label values applied.
	// Label values are provided as key-value pairs.
	// If the number of label values is odd, "unknown" is appended.
	With(labelValues ...string) Gauge
}

// Histogram represents a metric that samples observations and counts them
// in configurable buckets. Histograms are used to track distributions of
// values such as request durations or response sizes.
//
// A histogram automatically provides:
//   - Sum of all observed values
//   - Count of observations
//   - Distribution across configured buckets
type Histogram interface {
	// Observe adds a single observation to the histogram.
	Observe(value float64)

	// With returns a new Histogram with the given label values applied.
	// Label values are provided as key-value pairs.
	// If the number of label values is odd, "unknown" is appended.
	With(labelValues ...string) Histogram
}

// CounterOpts defines options for creating a counter metric.
type CounterOpts struct {
	Subsystem string
	Name      string
	Help      string
}

// GaugeOpts defines options for creating a gauge metric.
type GaugeOpts struct {
	Subsystem string
	Name      string
	Help      string
}

// HistogramOpts defines options for creating a histogram metric.
type HistogramOpts struct {
	Subsystem string
	Name      string
	Help      string
	Buckets   []float64
}

// Factory creates metrics instances (Counter, Gauge, Histogram).
// Implementations include PrometheusFactory and NoOpFactory.
type Factory interface {
	// NewCounter creates a new Counter with the given options and label names.
	NewCounter(opts CounterOpts, labelNames []string) Counter

	// NewGauge creates a new Gauge with the given options and label names.
	NewGauge(opts GaugeOpts, labelNames []string) Gauge

	// NewHistogram creates a new Histogram with the given options and label names.
	NewHistogram(opts HistogramOpts, labelNames []string) Histogram
}
