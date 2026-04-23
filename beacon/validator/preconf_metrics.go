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

package validator

// Validator-side preconf fetch outcomes used as the `outcome` label on
// beacon_kit.preconf.client.fetch_total.
const (
	FetchOutcomeOK               = "ok"
	FetchOutcomeUnavailable      = "unavailable"
	FetchOutcomeNotFound         = "not_found"
	FetchOutcomeValidationFailed = "validation_failed"
	FetchOutcomeOther            = "other"
)

// Fallback-to-local reasons used as the `reason` label on
// beacon_kit.preconf.client.fallback_to_local_total.
const (
	FallbackReasonUnavailable      = "sequencer_unavailable"
	FallbackReasonNotFound         = "payload_not_found"
	FallbackReasonValidationFailed = "local_validation_failed"
	FallbackReasonOther            = "other"
)

// preconfMetrics wraps a TelemetrySink and emits validator-side preconf metrics.
type preconfMetrics struct {
	sink TelemetrySink
}

func newPreconfMetrics(sink TelemetrySink) *preconfMetrics {
	return &preconfMetrics{sink: sink}
}

// markFetch records the outcome of a sequencer payload fetch attempt.
func (m *preconfMetrics) markFetch(outcome string) {
	m.sink.IncrementCounter("beacon_kit.preconf.client.fetch_total", "outcome", outcome)
}

// markFallback records a fallback to local building with the reason that triggered it.
func (m *preconfMetrics) markFallback(reason string) {
	m.sink.IncrementCounter("beacon_kit.preconf.client.fallback_to_local_total", "reason", reason)
}

// setSequencerAvailable sets the sequencer availability gauge (1=available, 0=offline).
func (m *preconfMetrics) setSequencerAvailable(available bool) {
	v := int64(0)
	if available {
		v = 1
	}
	m.sink.SetGauge("beacon_kit.preconf.client.sequencer_available", v)
}
