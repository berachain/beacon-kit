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

//nolint:testpackage // we test the unexported handleRPCError method.
package client

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	engineerrors "github.com/berachain/beacon-kit/engine-primitives/errors"
	"github.com/berachain/beacon-kit/execution/client/ethclient/rpc"
	"github.com/berachain/beacon-kit/log/phuslu"
	beaconhttp "github.com/berachain/beacon-kit/primitives/net/http"
	"github.com/stretchr/testify/require"
)

// noopSink implements TelemetrySink with no side effects.
type noopSink struct{}

func (noopSink) IncrementCounter(string, ...string)        {}
func (noopSink) MeasureSince(string, time.Time, ...string) {}

func newTestEngineClient() *EngineClient {
	return &EngineClient{
		logger:  phuslu.NewLogger(io.Discard, nil),
		metrics: newClientMetrics(noopSink{}, phuslu.NewLogger(io.Discard, nil)),
	}
}

func TestHandleRPCError_Classification(t *testing.T) {
	t.Parallel()
	s := newTestEngineClient()

	tests := []struct {
		name      string
		in        error
		wantIs    error // expected target for errors.Is (nil to skip)
		wantFatal bool  // expected IsFatalError verdict
		wantRetry bool  // expected IsNonFatalError verdict
	}{
		{
			name:      "nil passes through",
			in:        nil,
			wantIs:    nil,
			wantFatal: false,
			wantRetry: false,
		},
		{
			name:      "HTTP 413 (PoC) → ErrHTTPClientError (fatal)",
			in:        &rpc.HTTPStatusError{StatusCode: 413, Body: `{"code":-32007,"message":"Request is too big"}`},
			wantIs:    ErrHTTPClientError,
			wantFatal: true,
		},
		{
			name:      "HTTP 499 → ErrHTTPClientError (fatal)",
			in:        &rpc.HTTPStatusError{StatusCode: 499, Body: ""},
			wantIs:    ErrHTTPClientError,
			wantFatal: true,
		},
		{
			name:      "HTTP 401 → ErrHTTPClientError (fatal)",
			in:        &rpc.HTTPStatusError{StatusCode: 401, Body: ""},
			wantIs:    ErrHTTPClientError,
			wantFatal: true,
		},
		{
			name:      "HTTP 500 → ErrBadConnection (retryable, jsonrpsee body-stream blip)",
			in:        &rpc.HTTPStatusError{StatusCode: 500, Body: `{"error":{"code":-32603}}`},
			wantIs:    ErrBadConnection,
			wantFatal: false,
			wantRetry: true,
		},
		{
			name:      "engine API timeout passes through (retryable)",
			in:        engineerrors.ErrEngineAPITimeout,
			wantIs:    engineerrors.ErrEngineAPITimeout,
			wantFatal: false,
			wantRetry: true,
		},
		{
			name:      "plain transport error → ErrBadConnection (retryable, EL down)",
			in:        errors.New("dial tcp: connection refused"),
			wantIs:    ErrBadConnection,
			wantFatal: false,
			wantRetry: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := s.handleRPCError(tc.in)

			if tc.wantIs == nil {
				require.NoError(t, got)
			} else {
				require.Error(t, got)
				require.ErrorIs(t, got, tc.wantIs)
			}

			require.Equal(t, tc.wantFatal, IsFatalError(got), "IsFatalError mismatch")
			require.Equal(t, tc.wantRetry, IsNonFatalError(got), "IsNonFatalError mismatch")
		})
	}
}

// TestHandleRPCError_HTTPTimeout exercises the http.IsTimeoutError branch with
// a real timeout produced by an unresponsive HTTP server.
func TestHandleRPCError_HTTPTimeout(t *testing.T) {
	t.Parallel()
	s := newTestEngineClient()

	// Server that never responds within the client's timeout.
	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer srv.Close()

	httpClient := &http.Client{Timeout: 50 * time.Millisecond}
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, srv.URL, nil)
	require.NoError(t, err)
	_, doErr := httpClient.Do(req)
	require.Error(t, doErr)

	got := s.handleRPCError(doErr)
	require.ErrorIs(t, got, beaconhttp.ErrTimeout)
	require.True(t, IsNonFatalError(got), "timeouts must be retryable")
	require.False(t, IsFatalError(got))
}
