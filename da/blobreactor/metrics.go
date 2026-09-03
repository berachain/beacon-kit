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
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package blobreactor

import (
	"errors"
	"strconv"
	"time"
)

// Fetch attempt/completion statuses used in metrics labels and error classification.
const (
	statusSuccess         = "success"
	statusTimeout         = "timeout"
	statusSendFailed      = "send_failed"
	statusPeerNotFound    = "peer_not_found"
	statusInvalidResponse = "invalid_response"
	statusVerifyFailed    = "verify_failed"
	statusMiss            = "miss"
	statusAllPeersFailed  = "all_peers_failed"
)

// fetchError carries a metrics status alongside the underlying error.
type fetchError struct {
	err    error
	status string
}

func (e *fetchError) Error() string { return e.err.Error() }
func (e *fetchError) Unwrap() error { return e.err }

func newFetchError(err error, status string) error {
	return &fetchError{err: err, status: status}
}

func fetchErrStatus(err error) string {
	var fe *fetchError
	if errors.As(err, &fe) {
		return fe.status
	}
	return statusInvalidResponse
}

type blobReactorMetrics struct {
	sink TelemetrySink
}

func newBlobReactorMetrics(sink TelemetrySink) *blobReactorMetrics {
	return &blobReactorMetrics{sink: sink}
}

func (m *blobReactorMetrics) setPeerCount(n int) {
	m.sink.SetGauge("beacon_kit.da.blobreactor.peers", int64(n))
}

func (m *blobReactorMetrics) setActiveRequests(n int) {
	m.sink.SetGauge("beacon_kit.da.blobreactor.active_requests", int64(n))
}

func (m *blobReactorMetrics) observeMessageReceived(msgType string) {
	m.sink.IncrementCounter("beacon_kit.da.blobreactor.messages_received", "type", msgType)
}

func (m *blobReactorMetrics) observeRateLimited(msgType string) {
	m.sink.IncrementCounter("beacon_kit.da.blobreactor.rate_limited", "type", msgType)
}

// observeQueueFull counts an inbound message dropped because its lane queue was saturated.
func (m *blobReactorMetrics) observeQueueFull(taskType string) {
	m.sink.IncrementCounter("beacon_kit.da.blobreactor.queue_full", "task_type", taskType)
}

// observePush counts an inbound push by outcome ("accepted", "invalid", "slot_out_of_window", "duplicate").
func (m *blobReactorMetrics) observePush(result string) {
	m.sink.IncrementCounter("beacon_kit.da.blobreactor.push_received", "result", result)
}

func (m *blobReactorMetrics) observePushSent() {
	m.sink.IncrementCounter("beacon_kit.da.blobreactor.pushes_sent")
}

func (m *blobReactorMetrics) observeUnsolicitedResponse() {
	m.sink.IncrementCounter("beacon_kit.da.blobreactor.unsolicited_responses")
}

func (m *blobReactorMetrics) observeServed(kind string, found bool) {
	m.sink.IncrementCounter("beacon_kit.da.blobreactor.requests_served",
		"kind", kind, "found", strconv.FormatBool(found))
}

// recordFetchAttempt counts one per-peer fetch attempt by outcome.
func (m *blobReactorMetrics) recordFetchAttempt(kind, status string) {
	m.sink.IncrementCounter("beacon_kit.da.blobreactor.fetch_attempts", "kind", kind, "status", status)
}

// recordFetchDone measures the overall duration of one fetch operation by outcome; the histogram's count doubles
// as the per-operation counter.
func (m *blobReactorMetrics) recordFetchDone(kind, status string, start time.Time) {
	m.sink.MeasureSince("beacon_kit.da.blobreactor.fetch_duration", start, "kind", kind, "status", status)
}
