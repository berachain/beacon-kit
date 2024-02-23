// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package ethclient

const (
	// MetricBaseKey is the base key for all metrics.
	MetricBaseKey = "beacon-kit.engine.ethclient."

	// MetricKeyParseErrorCount represents the metric key
	// for counting parse errors.
	MetricKeyParseErrorCount = MetricBaseKey + "parse_error_count"

	// MetricKeyInvalidRequestCount represents the metric key
	// for counting invalid requests.
	MetricKeyInvalidRequestCount = MetricBaseKey +
		"invalid_request_count"

	// MetricKeyMethodNotFoundCount represents the metric key
	// for counting instances
	// where a method is not found.
	MetricKeyMethodNotFoundCount = MetricBaseKey +
		"method_not_found_count"

	// MetricKeyInvalidParamsCount represents the metric key
	// for counting instances of
	// invalid parameters.
	MetricKeyInvalidParamsCount = MetricBaseKey +
		"invalid_params_count"

	// MetricKeyInternalErrorCount represents the metric key
	// for counting internal errors.
	MetricKeyInternalErrorCount = MetricBaseKey +
		"internal_error_count"

	// MetricKeyUnknownPayloadErrorCount represents the metric key
	// for counting unknown payload errors.
	MetricKeyUnknownPayloadErrorCount = MetricBaseKey +
		"unknown_payload_error_count"

	// MetricKeyInvalidForkchoiceStateCount represents the metric key
	// for counting invalid fork choice state errors.
	MetricKeyInvalidForkchoiceStateCount = MetricBaseKey +
		"invalid_forkchoice_state_count"

	// MetricKeyInvalidPayloadAttributesCount represents the metric key
	// for counting invalid payload attribute errors.
	MetricKeyInvalidPayloadAttributesCount = MetricBaseKey +
		"invalid_payload_attributes_count"

	// MetricKeyRequestTooLargeCount represents the metric key
	// for counting instances where a request is too large.
	MetricKeyRequestTooLargeCount = MetricBaseKey +
		"request_too_large_count"

	// MetricKeyInternalServerErrorCount represents the metric key
	// for counting internal server errors.
	MetricKeyInternalServerErrorCount = MetricBaseKey +
		"internal_server_error_count"
)
