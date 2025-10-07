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

package blobreactor

import (
	"time"
)

// Metric status constants for blob reactor requests.
const (
	statusSuccess         = "success"
	statusTimeout         = "timeout"
	statusPeerNotFound    = "peer_not_found"
	statusSendFailed      = "send_failed"
	statusAllPeersFailed  = "all_peers_failed"
	statusMarshalFailed   = "marshal_failed"
	statusInvalidResponse = "invalid_response"
	statusVerifyFailed    = "verification_failed"
	messageTypeRequest    = "request"
	messageTypeResponse   = "response"
)

// blobReactorMetrics contains metrics for the blob reactor P2P operations.
type blobReactorMetrics struct {
	sink TelemetrySink
}

// newBlobReactorMetrics creates a new blobReactorMetrics instance.
func newBlobReactorMetrics(sink TelemetrySink) *blobReactorMetrics {
	return &blobReactorMetrics{sink: sink}
}

// recordOverallRequestComplete records completion of entire blob request (may try multiple peers).
func (m *blobReactorMetrics) recordOverallRequestComplete(status string, start time.Time) {
	m.sink.IncrementCounter("beacon_kit.blobreactor.request_total", "status", status)
	m.sink.MeasureSince("beacon_kit.blobreactor.request_duration", start, "status", status)
}

// recordPeerAttempt records a single peer attempt with status (no duration to avoid high cardinality).
func (m *blobReactorMetrics) recordPeerAttempt(status string) {
	m.sink.IncrementCounter("beacon_kit.blobreactor.peer_attempts_total", "status", status)
}

// observeWorkerPoolFull increments counter when worker pool is full and messages are dropped.
func (m *blobReactorMetrics) observeWorkerPoolFull(messageType string) {
	m.sink.IncrementCounter("beacon_kit.blobreactor.worker_pool_full_total", "message_type", messageType)
}

// setActiveRequests sets gauge for currently active blob requests.
func (m *blobReactorMetrics) setActiveRequests(count int) {
	m.sink.SetGauge("beacon_kit.blobreactor.active_requests", int64(count))
}

// setPeerPoolSize sets gauges for peer pool statistics.
func (m *blobReactorMetrics) setPeerPoolSize(available, total int) {
	m.sink.SetGauge("beacon_kit.blobreactor.peers_available", int64(available))
	m.sink.SetGauge("beacon_kit.blobreactor.peers_total", int64(total))
}
