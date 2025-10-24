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
// AN "AS IS" BASIS. LICSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package discard_test

import (
	"testing"

	"github.com/berachain/beacon-kit/observability/metrics/discard"
	"github.com/stretchr/testify/require"
)

func TestNoOpCounter(t *testing.T) {
	t.Parallel()
	counter := discard.NewCounter()

	// All operations should complete without panic
	counter.Add(1)
	counter.Add(100)

	// With should return a counter (itself)
	counter2 := counter.With("label", "value")
	require.NotNil(t, counter2)

	counter2.Add(50)

	// Chaining should work
	counter.With("label1", "value1").With("label2", "value2").Add(10)
}

func TestNoOpGauge(t *testing.T) {
	t.Parallel()
	gauge := discard.NewGauge()

	// All operations should complete without panic
	gauge.Set(42)
	gauge.Add(10)
	gauge.Add(-5)

	// With should return a gauge (itself)
	gauge2 := gauge.With("label", "value")
	require.NotNil(t, gauge2)

	gauge2.Set(100)

	// Chaining should work
	gauge.With("label1", "value1").With("label2", "value2").Set(50)
}

func TestNoOpHistogram(t *testing.T) {
	t.Parallel()
	histogram := discard.NewHistogram()

	// All operations should complete without panic
	histogram.Observe(0.5)
	histogram.Observe(1.0)
	histogram.Observe(10.0)

	// With should return a histogram (itself)
	histogram2 := histogram.With("label", "value")
	require.NotNil(t, histogram2)

	histogram2.Observe(2.5)

	// Chaining should work
	histogram.With("label1", "value1").With("label2", "value2").Observe(3.3)
}

// Benchmark to verify no-op operations have zero cost
func BenchmarkNoOpCounter(b *testing.B) {
	counter := discard.NewCounter()
	b.ResetTimer()
	for range b.N {
		counter.Add(1)
	}
}

func BenchmarkNoOpGauge(b *testing.B) {
	gauge := discard.NewGauge()
	b.ResetTimer()
	for i := range b.N {
		gauge.Set(float64(i))
	}
}

func BenchmarkNoOpHistogram(b *testing.B) {
	histogram := discard.NewHistogram()
	b.ResetTimer()
	for i := range b.N {
		histogram.Observe(float64(i))
	}
}
