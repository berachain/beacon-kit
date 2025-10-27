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

package client

import (
	"time"

	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/observability/metrics"
)

// Metrics is a struct that contains metrics for the execution client.
type Metrics struct {
	// Duration histograms
	ForkchoiceUpdateDuration metrics.Summary
	NewPayloadDuration       metrics.Summary
	GetPayloadDuration       metrics.Summary

	// Timeout counters
	EngineAPITimeout                metrics.Counter
	ForkchoiceUpdateDurationTimeout metrics.Counter
	NewPayloadDurationTimeout       metrics.Counter
	GetPayloadDurationTimeout       metrics.Counter
	HTTPTimeout                     metrics.Counter

	// Error counters
	ParseError               metrics.Counter
	InvalidRequest           metrics.Counter
	MethodNotFound           metrics.Counter
	InvalidParams            metrics.Counter
	InternalError            metrics.Counter
	UnknownPayloadError      metrics.Counter
	InvalidForkchoiceState   metrics.Counter
	InvalidPayloadAttributes metrics.Counter
	RequestTooLarge          metrics.Counter
	InternalServerError      metrics.Counter

	logger log.Logger
}

// NewPrometheusMetrics returns a new Metrics instance with Prometheus metrics.
// Metric names are kept identical to cosmos-sdk/telemetry output for Grafana compatibility.
//
//nolint:funlen
func NewMetrics(factory metrics.Factory, logger log.Logger) *Metrics {
	return &Metrics{
		// Duration histograms
		ForkchoiceUpdateDuration: factory.NewSummary(
			metrics.SummaryOpts{
				Name:       "beacon_kit_execution_client_forkchoice_update_duration",
				Help:       "Time taken for forkchoice update in milliseconds",
				Objectives: metrics.QuantilesP50P90P99,
			},
			nil,
		),
		NewPayloadDuration: factory.NewSummary(
			metrics.SummaryOpts{
				Name:       "beacon_kit_execution_client_new_payload_duration",
				Help:       "Time taken for new payload in milliseconds",
				Objectives: metrics.QuantilesP50P90P99,
			},
			nil,
		),
		GetPayloadDuration: factory.NewSummary(
			metrics.SummaryOpts{
				Name:       "beacon_kit_execution_client_get_payload_duration",
				Help:       "Time taken for get payload in milliseconds",
				Objectives: metrics.QuantilesP50P90P99,
			},
			nil,
		),

		// Timeout counters
		EngineAPITimeout: factory.NewCounter(
			metrics.CounterOpts{
				Name: "beacon_kit_execution_client_engine_api_timeout",
				Help: "Number of engine API timeouts",
			},
			nil,
		),
		ForkchoiceUpdateDurationTimeout: factory.NewCounter(
			metrics.CounterOpts{
				Name: "beacon_kit_execution_client_forkchoice_update_duration_timeout",
				Help: "Number of forkchoice update timeouts",
			},
			nil,
		),
		NewPayloadDurationTimeout: factory.NewCounter(
			metrics.CounterOpts{
				Name: "beacon_kit_execution_client_new_payload_duration_timeout",
				Help: "Number of new payload timeouts",
			},
			nil,
		),
		GetPayloadDurationTimeout: factory.NewCounter(
			metrics.CounterOpts{
				Name: "beacon_kit_execution_client_get_payload_duration_timeout",
				Help: "Number of get payload timeouts",
			},
			nil,
		),
		HTTPTimeout: factory.NewCounter(
			metrics.CounterOpts{
				Name: "beacon_kit_execution_client_http_timeout",
				Help: "Number of HTTP timeouts",
			},
			nil,
		),

		// Error counters
		ParseError: factory.NewCounter(
			metrics.CounterOpts{
				Name: "beacon_kit_execution_client_parse_error",
				Help: "Number of parse errors",
			},
			nil,
		),
		InvalidRequest: factory.NewCounter(
			metrics.CounterOpts{
				Name: "beacon_kit_execution_client_invalid_request",
				Help: "Number of invalid requests",
			},
			nil,
		),
		MethodNotFound: factory.NewCounter(
			metrics.CounterOpts{
				Name: "beacon_kit_execution_client_method_not_found",
				Help: "Number of method not found errors",
			},
			nil,
		),
		InvalidParams: factory.NewCounter(
			metrics.CounterOpts{
				Name: "beacon_kit_execution_client_invalid_params",
				Help: "Number of invalid params errors",
			},
			nil,
		),
		InternalError: factory.NewCounter(
			metrics.CounterOpts{
				Name: "beacon_kit_execution_client_internal_error",
				Help: "Number of internal errors",
			},
			nil,
		),
		UnknownPayloadError: factory.NewCounter(
			metrics.CounterOpts{
				Name: "beacon_kit_execution_client_unknown_payload_error",
				Help: "Number of unknown payload errors",
			},
			nil,
		),
		InvalidForkchoiceState: factory.NewCounter(
			metrics.CounterOpts{
				Name: "beacon_kit_execution_client_invalid_forkchoice_state",
				Help: "Number of invalid forkchoice state errors",
			},
			nil,
		),
		InvalidPayloadAttributes: factory.NewCounter(
			metrics.CounterOpts{
				Name: "beacon_kit_execution_client_invalid_payload_attributes",
				Help: "Number of invalid payload attributes errors",
			},
			nil,
		),
		RequestTooLarge: factory.NewCounter(
			metrics.CounterOpts{
				Name: "beacon_kit_execution_client_request_too_large",
				Help: "Number of request too large errors",
			},
			nil,
		),
		InternalServerError: factory.NewCounter(
			metrics.CounterOpts{
				Name: "beacon_kit_execution_client_internal_server_error",
				Help: "Number of internal server errors",
			},
			nil,
		),

		logger: logger,
	}
}

// measureForkchoiceUpdateDuration measures the duration of the forkchoice update.
func (m *Metrics) measureForkchoiceUpdateDuration(startTime time.Time) {
	m.ForkchoiceUpdateDuration.Observe(float64(time.Since(startTime).Milliseconds()))
}

// measureNewPayloadDuration measures the duration of the new payload.
func (m *Metrics) measureNewPayloadDuration(startTime time.Time) {
	m.NewPayloadDuration.Observe(float64(time.Since(startTime).Milliseconds()))
}

// measureGetPayloadDuration measures the duration of the get payload.
func (m *Metrics) measureGetPayloadDuration(startTime time.Time) {
	m.GetPayloadDuration.Observe(float64(time.Since(startTime).Milliseconds()))
}

// incrementEngineAPITimeout increments the timeout counter for general engine api timeouts.
func (m *Metrics) incrementEngineAPITimeout() {
	m.EngineAPITimeout.Add(1)
}

// incrementForkchoiceUpdateTimeout increments the timeout counter for forkchoice update.
func (m *Metrics) incrementForkchoiceUpdateTimeout() {
	m.ForkchoiceUpdateDurationTimeout.Add(1)
}

// incrementNewPayloadTimeout increments the timeout counter for new payload.
func (m *Metrics) incrementNewPayloadTimeout() {
	m.NewPayloadDurationTimeout.Add(1)
}

// incrementGetPayloadTimeout increments the timeout counter for get payload.
func (m *Metrics) incrementGetPayloadTimeout() {
	m.GetPayloadDurationTimeout.Add(1)
}

// incrementHTTPTimeoutCounter increments the timeout counter for HTTP.
func (m *Metrics) incrementHTTPTimeoutCounter() {
	m.HTTPTimeout.Add(1)
}

// incrementParseErrorCounter increments the parse error counter.
func (m *Metrics) incrementParseErrorCounter() {
	m.ParseError.Add(1)
}

// incrementInvalidRequestCounter increments the invalid request counter.
func (m *Metrics) incrementInvalidRequestCounter() {
	m.InvalidRequest.Add(1)
}

// incrementMethodNotFoundCounter increments the method not found counter.
func (m *Metrics) incrementMethodNotFoundCounter() {
	m.MethodNotFound.Add(1)
}

// incrementInvalidParamsCounter increments the invalid params counter.
func (m *Metrics) incrementInvalidParamsCounter() {
	m.InvalidParams.Add(1)
}

// incrementInternalErrorCounter increments the internal error counter.
func (m *Metrics) incrementInternalErrorCounter() {
	m.InternalError.Add(1)
}

// incrementUnknownPayloadErrorCounter increments the unknown payload error counter.
func (m *Metrics) incrementUnknownPayloadErrorCounter() {
	m.UnknownPayloadError.Add(1)
}

// incrementInvalidForkchoiceStateCounter increments the invalid forkchoice state counter.
func (m *Metrics) incrementInvalidForkchoiceStateCounter() {
	m.InvalidForkchoiceState.Add(1)
}

// incrementInvalidPayloadAttributesCounter increments the invalid payload attributes counter.
func (m *Metrics) incrementInvalidPayloadAttributesCounter() {
	m.InvalidPayloadAttributes.Add(1)
}

// incrementRequestTooLargeCounter increments the request too large counter.
func (m *Metrics) incrementRequestTooLargeCounter() {
	m.RequestTooLarge.Add(1)
}

// incrementInternalServerErrorCounter increments the internal server error counter.
func (m *Metrics) incrementInternalServerErrorCounter() {
	m.InternalServerError.Add(1)
}
