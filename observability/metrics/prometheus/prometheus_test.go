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

package prometheus_test

import (
	"testing"

	bkprometheus "github.com/berachain/beacon-kit/observability/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/require"
)

func TestCounter(t *testing.T) {
	t.Parallel()
	// Create a new registry for isolated testing
	reg := prometheus.NewRegistry()

	cv := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "test",
		Subsystem: "counter",
		Name:      "total",
		Help:      "Test counter",
	}, []string{"label1"})
	reg.MustRegister(cv)

	counter := bkprometheus.NewCounterFrom(cv)

	// Test With and Add
	counter.With("label1", "value1").Add(5)
	counter.With("label1", "value1").Add(2)
	counter.With("label1", "value2").Add(10)

	// Gather metrics
	metrics, err := reg.Gather()
	require.NoError(t, err)
	require.Len(t, metrics, 1)

	// Verify metric family
	mf := metrics[0]
	require.Equal(t, "test_counter_total", mf.GetName())
	require.Equal(t, dto.MetricType_COUNTER, mf.GetType())

	// Verify we have metrics with different labels
	require.Len(t, mf.GetMetric(), 2)
}

func TestGauge(t *testing.T) {
	t.Parallel()
	// Create a new registry for isolated testing
	reg := prometheus.NewRegistry()

	gv := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "test",
		Subsystem: "gauge",
		Name:      "current",
		Help:      "Test gauge",
	}, []string{"label1"})
	reg.MustRegister(gv)

	gauge := bkprometheus.NewGaugeFrom(gv)

	// Test Set
	gauge.With("label1", "value1").Set(42)

	// Test Add (increment)
	gauge.With("label1", "value1").Add(8)

	// Test Add (decrement)
	gauge.With("label1", "value1").Add(-10)

	// Gather metrics
	metrics, err := reg.Gather()
	require.NoError(t, err)
	require.Len(t, metrics, 1)

	// Verify metric family
	mf := metrics[0]
	require.Equal(t, "test_gauge_current", mf.GetName())
	require.Equal(t, dto.MetricType_GAUGE, mf.GetType())

	// Verify value (42 + 8 - 10 = 40)
	require.Len(t, mf.GetMetric(), 1)
	require.InDelta(t, float64(40), mf.GetMetric()[0].GetGauge().GetValue(), 0.001)
}

func TestHistogram(t *testing.T) {
	t.Parallel()
	// Create a new registry for isolated testing
	reg := prometheus.NewRegistry()

	hv := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "test",
		Subsystem: "histogram",
		Name:      "duration_seconds",
		Help:      "Test histogram",
		Buckets:   []float64{0.1, 0.5, 1.0, 5.0},
	}, []string{"label1"})
	reg.MustRegister(hv)

	histogram := bkprometheus.NewHistogramFrom(hv)

	// Test Observe
	histogram.With("label1", "value1").Observe(0.3)
	histogram.With("label1", "value1").Observe(0.7)
	histogram.With("label1", "value1").Observe(2.0)

	// Gather metrics
	metrics, err := reg.Gather()
	require.NoError(t, err)
	require.Len(t, metrics, 1)

	// Verify metric family
	mf := metrics[0]
	require.Equal(t, "test_histogram_duration_seconds", mf.GetName())
	require.Equal(t, dto.MetricType_HISTOGRAM, mf.GetType())

	// Verify histogram properties
	require.Len(t, mf.GetMetric(), 1)
	h := mf.GetMetric()[0].GetHistogram()
	require.Equal(t, uint64(3), h.GetSampleCount())
	require.InDelta(t, 3.0, h.GetSampleSum(), 0.001)

	// Verify buckets
	// Observations: 0.3, 0.7, 2.0
	// Buckets: 0.1, 0.5, 1.0, 5.0
	require.Len(t, h.GetBucket(), 4)
	require.Equal(t, uint64(0), h.GetBucket()[0].GetCumulativeCount()) // <= 0.1: none
	require.Equal(t, uint64(1), h.GetBucket()[1].GetCumulativeCount()) // <= 0.5: 0.3
	require.Equal(t, uint64(2), h.GetBucket()[2].GetCumulativeCount()) // <= 1.0: 0.3, 0.7
	require.Equal(t, uint64(3), h.GetBucket()[3].GetCumulativeCount()) // <= 5.0: 0.3, 0.7, 2.0
}

func TestCounterLabelChaining(t *testing.T) {
	t.Parallel()
	// Create a new registry for isolated testing
	reg := prometheus.NewRegistry()

	cv := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "test",
		Subsystem: "counter",
		Name:      "chaining",
		Help:      "Test counter chaining",
	}, []string{"label1", "label2"})
	reg.MustRegister(cv)

	counter := bkprometheus.NewCounterFrom(cv)

	// Test chaining With calls
	counter.With("label1", "value1").With("label2", "value2").Add(5)

	// Gather metrics
	metrics, err := reg.Gather()
	require.NoError(t, err)
	require.Len(t, metrics, 1)

	mf := metrics[0]
	require.Len(t, mf.GetMetric(), 1)
	require.InDelta(t, float64(5), mf.GetMetric()[0].GetCounter().GetValue(), 0.001)

	// Verify labels
	labels := mf.GetMetric()[0].GetLabel()
	require.Len(t, labels, 2)
}
