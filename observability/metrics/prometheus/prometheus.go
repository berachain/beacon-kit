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

package prometheus

import (
	"github.com/berachain/beacon-kit/observability/metrics"
	"github.com/berachain/beacon-kit/observability/metrics/lv"
	"github.com/prometheus/client_golang/prometheus"
)

// Factory creates Prometheus metrics and registers them with prometheus.DefaultRegisterer.
type Factory struct {
	namespace string
}

// NewFactory creates a new Prometheus metrics factory with the given namespace.
func NewFactory(namespace string) metrics.Factory {
	return &Factory{namespace: namespace}
}

// NewCounter creates a new Counter that registers with prometheus.DefaultRegisterer.
func (f *Factory) NewCounter(opts metrics.CounterOpts, labelNames []string) metrics.Counter {
	return NewCounter(prometheus.CounterOpts{
		Namespace: f.namespace,
		Name:      opts.Name,
		Help:      opts.Help,
	}, labelNames)
}

// NewGauge creates a new Gauge that registers with prometheus.DefaultRegisterer.
func (f *Factory) NewGauge(opts metrics.GaugeOpts, labelNames []string) metrics.Gauge {
	return NewGauge(prometheus.GaugeOpts{
		Namespace: f.namespace,
		Name:      opts.Name,
		Help:      opts.Help,
	}, labelNames)
}

// NewHistogram creates a new Histogram that registers with prometheus.DefaultRegisterer.
func (f *Factory) NewHistogram(opts metrics.HistogramOpts, labelNames []string) metrics.Histogram {
	return NewHistogram(prometheus.HistogramOpts{
		Namespace: f.namespace,
		Name:      opts.Name,
		Help:      opts.Help,
		Buckets:   opts.Buckets,
	}, labelNames)
}

// NewSummary creates a new Summary that registers with prometheus.DefaultRegisterer.
func (f *Factory) NewSummary(opts metrics.SummaryOpts, labelNames []string) metrics.Summary {
	return NewSummary(prometheus.SummaryOpts{
		Namespace:  f.namespace,
		Name:       opts.Name,
		Help:       opts.Help,
		Objectives: opts.Objectives,
	}, labelNames)
}

// counter wraps a prometheus.CounterVec and implements the metrics.Counter interface.
type counter struct {
	cv  *prometheus.CounterVec
	lvs lv.LabelValues
}

// NewCounter creates a new Counter that registers with prometheus.DefaultRegisterer.
func NewCounter(opts prometheus.CounterOpts, labelNames []string) metrics.Counter {
	cv := prometheus.NewCounterVec(opts, labelNames)
	prometheus.MustRegister(cv)
	return &counter{cv: cv}
}

// NewCounterFrom creates a new Counter from an existing prometheus.CounterVec.
// The CounterVec must already be registered.
func NewCounterFrom(cv *prometheus.CounterVec, labelValues ...string) metrics.Counter {
	return &counter{
		cv:  cv,
		lvs: lv.LabelValues(labelValues),
	}
}

// With returns a new Counter with the given label values applied.
func (c *counter) With(labelValues ...string) metrics.Counter {
	return &counter{
		cv:  c.cv,
		lvs: c.lvs.With(labelValues...),
	}
}

// Add increments the counter by the given delta.
func (c *counter) Add(delta float64) {
	c.cv.With(makeLabels(c.lvs...)).Add(delta)
}

// gauge wraps a prometheus.GaugeVec and implements the metrics.Gauge interface.
type gauge struct {
	gv  *prometheus.GaugeVec
	lvs lv.LabelValues
}

// NewGauge creates a new Gauge that registers with prometheus.DefaultRegisterer.
func NewGauge(opts prometheus.GaugeOpts, labelNames []string) metrics.Gauge {
	gv := prometheus.NewGaugeVec(opts, labelNames)
	prometheus.MustRegister(gv)
	return &gauge{gv: gv}
}

// NewGaugeFrom creates a new Gauge from an existing prometheus.GaugeVec.
// The GaugeVec must already be registered.
func NewGaugeFrom(gv *prometheus.GaugeVec, labelValues ...string) metrics.Gauge {
	return &gauge{
		gv:  gv,
		lvs: lv.LabelValues(labelValues),
	}
}

// With returns a new Gauge with the given label values applied.
func (g *gauge) With(labelValues ...string) metrics.Gauge {
	return &gauge{
		gv:  g.gv,
		lvs: g.lvs.With(labelValues...),
	}
}

// Set sets the gauge to the given value.
func (g *gauge) Set(value float64) {
	g.gv.With(makeLabels(g.lvs...)).Set(value)
}

// Add increments or decrements the gauge by the given delta.
func (g *gauge) Add(delta float64) {
	g.gv.With(makeLabels(g.lvs...)).Add(delta)
}

// histogram wraps a prometheus.HistogramVec and implements the metrics.Histogram interface.
type histogram struct {
	hv  *prometheus.HistogramVec
	lvs lv.LabelValues
}

// NewHistogram creates a new Histogram that registers with prometheus.DefaultRegisterer.
func NewHistogram(opts prometheus.HistogramOpts, labelNames []string) metrics.Histogram {
	hv := prometheus.NewHistogramVec(opts, labelNames)
	prometheus.MustRegister(hv)
	return &histogram{hv: hv}
}

// NewHistogramFrom creates a new Histogram from an existing prometheus.HistogramVec.
// The HistogramVec must already be registered.
func NewHistogramFrom(hv *prometheus.HistogramVec, labelValues ...string) metrics.Histogram {
	return &histogram{
		hv:  hv,
		lvs: lv.LabelValues(labelValues),
	}
}

// With returns a new Histogram with the given label values applied.
func (h *histogram) With(labelValues ...string) metrics.Histogram {
	return &histogram{
		hv:  h.hv,
		lvs: h.lvs.With(labelValues...),
	}
}

// Observe adds a single observation to the histogram.
func (h *histogram) Observe(value float64) {
	h.hv.With(makeLabels(h.lvs...)).Observe(value)
}

// summary wraps a prometheus.SummaryVec and implements the metrics.Summary interface.
type summary struct {
	sv  *prometheus.SummaryVec
	lvs lv.LabelValues
}

// NewSummary creates a new Summary that registers with prometheus.DefaultRegisterer.
func NewSummary(opts prometheus.SummaryOpts, labelNames []string) metrics.Summary {
	sv := prometheus.NewSummaryVec(opts, labelNames)
	prometheus.MustRegister(sv)
	return &summary{sv: sv}
}

// NewSummaryFrom creates a new Summary from an existing prometheus.SummaryVec.
// The SummaryVec must already be registered.
func NewSummaryFrom(sv *prometheus.SummaryVec, labelValues ...string) metrics.Summary {
	return &summary{
		sv:  sv,
		lvs: lv.LabelValues(labelValues),
	}
}

// With returns a new Summary with the given label values applied.
func (s *summary) With(labelValues ...string) metrics.Summary {
	return &summary{
		sv:  s.sv,
		lvs: s.lvs.With(labelValues...),
	}
}

// Observe adds a single observation to the summary.
func (s *summary) Observe(value float64) {
	s.sv.With(makeLabels(s.lvs...)).Observe(value)
}

// makeLabels converts a slice of label key-value pairs into a prometheus.Labels map.
// The input slice should contain alternating keys and values.
//
//nolint:mnd // 2 is for key-value pairs
func makeLabels(labelValues ...string) prometheus.Labels {
	labels := make(prometheus.Labels, len(labelValues)/2)
	for i := 0; i < len(labelValues); i += 2 {
		labels[labelValues[i]] = labelValues[i+1]
	}
	return labels
}
