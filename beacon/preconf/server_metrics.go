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

package preconf

import "strconv"

// ServerResult is the outcome of a preconf GetPayload request, used as the
// `result` label on beacon_kit.preconf.server.payload_request_total.
type ServerResult string

const (
	ServerResultOK                 ServerResult = "ok"
	ServerResultUnauthorized       ServerResult = "unauthorized"
	ServerResultNotWhitelisted     ServerResult = "not_whitelisted"
	ServerResultWrongProposer      ServerResult = "wrong_proposer"
	ServerResultPayloadNotFound    ServerResult = "payload_not_found"
	ServerResultInternalError      ServerResult = "internal_error"
	ServerResultResponseWriteError ServerResult = "response_write_error"
	ServerResultMethodNotAllowed   ServerResult = "method_not_allowed"
	ServerResultBadRequest         ServerResult = "bad_request"
)

// TelemetrySink is a minimal sink interface used by the preconf server for
// emitting counters.
type TelemetrySink interface {
	IncrementCounter(key string, args ...string)
}

// serverMetrics wraps a TelemetrySink and emits sequencer-server metrics.
type serverMetrics struct {
	sink TelemetrySink
}

func newServerMetrics(sink TelemetrySink) *serverMetrics {
	return &serverMetrics{sink: sink}
}

// markPayloadRequest records the outcome of a GetPayload request.
func (m *serverMetrics) markPayloadRequest(result ServerResult) {
	m.sink.IncrementCounter("beacon_kit.preconf.server.payload_request_total", "result", string(result))
}

// markProposerCheck records the outcome of an expected-proposer check.
func (m *serverMetrics) markProposerCheck(matched bool) {
	m.sink.IncrementCounter("beacon_kit.preconf.proposer_tracker.check_total", "matched", strconv.FormatBool(matched))
}
