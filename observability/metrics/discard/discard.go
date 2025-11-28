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

package discard

import "github.com/berachain/beacon-kit/observability/metrics"

// Factory creates no-op metrics with zero runtime overhead.
type Factory struct{}

// NewFactory creates a new no-op metrics factory.
func NewFactory() metrics.Factory {
	return Factory{}
}

// NewCounter returns a no-op Counter.
func (Factory) NewCounter(metrics.CounterOpts, []string) metrics.Counter {
	return NewCounter()
}

// NewGauge returns a no-op Gauge.
func (Factory) NewGauge(metrics.GaugeOpts, []string) metrics.Gauge {
	return NewGauge()
}

// NewHistogram returns a no-op Histogram.
func (Factory) NewHistogram(metrics.HistogramOpts, []string) metrics.Histogram {
	return NewHistogram()
}

// NewSummary returns a no-op Summary.
func (Factory) NewSummary(metrics.SummaryOpts, []string) metrics.Summary {
	return NewSummary()
}

// noOpCounter is a no-op implementation of metrics.Counter.
type noOpCounter struct{}

// NewCounter returns a no-op Counter.
func NewCounter() metrics.Counter {
	return noOpCounter{}
}

// With returns the same no-op Counter.
func (noOpCounter) With(...string) metrics.Counter {
	return noOpCounter{}
}

// Add does nothing.
func (noOpCounter) Add(float64) {}

// noOpGauge is a no-op implementation of metrics.Gauge.
type noOpGauge struct{}

// NewGauge returns a no-op Gauge.
func NewGauge() metrics.Gauge {
	return noOpGauge{}
}

// With returns the same no-op Gauge.
func (noOpGauge) With(...string) metrics.Gauge {
	return noOpGauge{}
}

// Set does nothing.
func (noOpGauge) Set(float64) {}

// Add does nothing.
func (noOpGauge) Add(float64) {}

// noOpHistogram is a no-op implementation of metrics.Histogram.
type noOpHistogram struct{}

// NewHistogram returns a no-op Histogram.
func NewHistogram() metrics.Histogram {
	return noOpHistogram{}
}

// With returns the same no-op Histogram.
func (noOpHistogram) With(...string) metrics.Histogram {
	return noOpHistogram{}
}

// Observe does nothing.
func (noOpHistogram) Observe(float64) {}

// noOpSummary is a no-op implementation of metrics.Summary.
type noOpSummary struct{}

// NewSummary returns a no-op Summary.
func NewSummary() metrics.Summary {
	return noOpSummary{}
}

// With returns the same no-op Summary.
func (noOpSummary) With(...string) metrics.Summary {
	return noOpSummary{}
}

// Observe does nothing.
func (noOpSummary) Observe(float64) {}
