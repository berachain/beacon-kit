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

package preconf_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/berachain/beacon-kit/beacon/preconf"
	"github.com/berachain/beacon-kit/cli/utils/parser"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/log/noop"
	"github.com/berachain/beacon-kit/node-core/components/metrics"
	payloadbuilder "github.com/berachain/beacon-kit/payload/builder"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/net/jwt"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/stretchr/testify/require"
)

// Test fixtures - valid BLS pubkeys (48 bytes = 96 hex chars)
const (
	pubkeyAHex = "0x93247f2209abcacf57b75a51dafae777f9dd38bc7053d1af526f220a7489a6d3a2753e5f3e8b1cfe39b56f43611df74a"
	pubkeyBHex = "0xa572cbea904d67468808c8eb50a9450c9721db309128012543902d0ac358a62ae28f75bb8f1c7c42c39a8c5529bf0f4e"
	secretAHex = "0x0102030405060708091011121314151617181920212223242526272829303132"
	secretBHex = "0x3132333435363738394041424344454647484950515253545556575859606162"
)

func TestServer_HandleGetPayload(t *testing.T) {
	t.Parallel()

	validatorA, _ := parser.ConvertPubkey(pubkeyAHex)
	validatorB, _ := parser.ConvertPubkey(pubkeyBHex)
	secretA, _ := jwt.NewFromHex(secretAHex)
	secretB, _ := jwt.NewFromHex(secretBHex)

	tests := []struct {
		name          string
		requestSlot   math.Slot
		requestJWT    *jwt.Secret
		payloadExists bool
		wantStatus    int
		wantContains  string
	}{
		{
			name:          "success - authenticated whitelisted validator",
			requestSlot:   100,
			requestJWT:    secretA,
			payloadExists: true,
			wantStatus:    http.StatusOK,
		},
		{
			name:          "success - different authenticated validator",
			requestSlot:   100,
			requestJWT:    secretB,
			payloadExists: true,
			wantStatus:    http.StatusOK,
		},
		{
			name:          "not found - no payload for slot",
			requestSlot:   999,
			requestJWT:    secretA,
			payloadExists: false,
			wantStatus:    http.StatusNotFound,
			wantContains:  "payload not available",
		},
		{
			name:         "unauthorized - missing auth",
			requestSlot:  100,
			requestJWT:   nil,
			wantStatus:   http.StatusUnauthorized,
			wantContains: "missing Authorization",
		},
	}

	jwtToValidator := map[*jwt.Secret]crypto.BLSPubkey{
		secretA: validatorA,
		secretB: validatorB,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			provider := &mockPayloadProvider{
				hasPayload: tt.payloadExists,
			}
			tracker := preconf.NewProposerTracker()
			if pubkey, ok := jwtToValidator[tt.requestJWT]; ok {
				tracker.SetExpectedProposer(tt.requestSlot, pubkey)
			}
			server := preconf.NewServer(
				noop.NewLogger[any](),
				preconf.ValidatorJWTs{validatorA: secretA, validatorB: secretB},
				newTestWhitelist(t, pubkeyAHex, pubkeyBHex),
				tracker,
				provider,
				0,
				metrics.NewNoOpTelemetrySink(),
			)

			body, _ := json.Marshal(preconf.GetPayloadRequest{Slot: tt.requestSlot})
			req := httptest.NewRequest(http.MethodPost, preconf.PayloadEndpoint, bytes.NewReader(body))
			if tt.requestJWT != nil {
				token, _ := tt.requestJWT.BuildSignedToken()
				req.Header.Set("Authorization", "Bearer "+token)
			}

			rec := httptest.NewRecorder()
			server.Handler().ServeHTTP(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantContains != "" {
				require.Contains(t, rec.Body.String(), tt.wantContains)
			}
		})
	}
}

func TestServer_ProposerCheck(t *testing.T) {
	t.Parallel()

	validatorA, _ := parser.ConvertPubkey(pubkeyAHex)
	validatorB, _ := parser.ConvertPubkey(pubkeyBHex)
	secretA, _ := jwt.NewFromHex(secretAHex)
	secretB, _ := jwt.NewFromHex(secretBHex)

	const targetSlot math.Slot = 100

	tests := []struct {
		name         string
		setupTracker func() preconf.ProposerTracker
		requestJWT   *jwt.Secret
		wantStatus   int
		wantContains string
	}{
		{
			name: "expected proposer gets payload",
			setupTracker: func() preconf.ProposerTracker {
				tr := preconf.NewProposerTracker()
				tr.SetExpectedProposer(targetSlot, validatorA)
				return tr
			},
			requestJWT: secretA,
			wantStatus: http.StatusOK,
		},
		{
			name: "non-expected proposer is rejected",
			setupTracker: func() preconf.ProposerTracker {
				tr := preconf.NewProposerTracker()
				tr.SetExpectedProposer(targetSlot, validatorA)
				return tr
			},
			requestJWT:   secretB,
			wantStatus:   http.StatusForbidden,
			wantContains: "not the expected proposer",
		},
		{
			name:         "no tracked proposer for slot is rejected",
			setupTracker: preconf.NewProposerTracker,
			requestJWT:   secretA,
			wantStatus:   http.StatusForbidden,
			wantContains: "not the expected proposer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := preconf.NewServer(
				noop.NewLogger[any](),
				preconf.ValidatorJWTs{validatorA: secretA, validatorB: secretB},
				newTestWhitelist(t, pubkeyAHex, pubkeyBHex),
				tt.setupTracker(),
				&mockPayloadProvider{hasPayload: true},
				0,
				metrics.NewNoOpTelemetrySink(),
			)

			body, _ := json.Marshal(preconf.GetPayloadRequest{Slot: targetSlot})
			req := httptest.NewRequest(http.MethodPost, preconf.PayloadEndpoint, bytes.NewReader(body))
			token, _ := tt.requestJWT.BuildSignedToken()
			req.Header.Set("Authorization", "Bearer "+token)

			rec := httptest.NewRecorder()
			server.Handler().ServeHTTP(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantContains != "" {
				require.Contains(t, rec.Body.String(), tt.wantContains)
			}
		})
	}
}

func TestServer_RejectsNonPostMethods(t *testing.T) {
	t.Parallel()

	server := preconf.NewServer(noop.NewLogger[any](), nil, nil, nil, nil, 0, metrics.NewNoOpTelemetrySink())

	for _, method := range []string{http.MethodGet, http.MethodPut, http.MethodDelete} {
		req := httptest.NewRequest(method, preconf.PayloadEndpoint, nil)
		rec := httptest.NewRecorder()
		server.Handler().ServeHTTP(rec, req)
		require.Equal(t, http.StatusMethodNotAllowed, rec.Code, "method: %s", method)
	}
}

// newTestWhitelist writes a temp JSON whitelist file from the given hex pubkeys
// and returns a Whitelist loaded from it.
func newTestWhitelist(t *testing.T, pubkeyHexes ...string) preconf.Whitelist {
	t.Helper()
	content, err := json.Marshal(pubkeyHexes)
	require.NoError(t, err)
	tmpFile := filepath.Join(t.TempDir(), "whitelist.json")
	err = os.WriteFile(tmpFile, content, 0o644)
	require.NoError(t, err)
	wl, err := preconf.NewWhitelist(tmpFile)
	require.NoError(t, err)
	return wl
}

func TestServer_OnSIGHUP(t *testing.T) {
	t.Parallel()

	pkA, err := parser.ConvertPubkey(pubkeyAHex)
	require.NoError(t, err)
	pkB, err := parser.ConvertPubkey(pubkeyBHex)
	require.NoError(t, err)

	tmpFile := filepath.Join(t.TempDir(), "whitelist.json")

	// Write initial whitelist with only key A.
	content, err := json.Marshal([]string{pubkeyAHex})
	require.NoError(t, err)
	err = os.WriteFile(tmpFile, content, 0o644)
	require.NoError(t, err)

	wl, err := preconf.NewWhitelist(tmpFile)
	require.NoError(t, err)

	server := preconf.NewServer(noop.NewLogger[any](), nil, wl, nil, nil, 0, metrics.NewNoOpTelemetrySink())

	require.True(t, wl.IsWhitelisted(pkA))
	require.False(t, wl.IsWhitelisted(pkB))

	// Add key B to file and trigger hot-reload via OnSIGHUP.
	content, err = json.Marshal([]string{pubkeyAHex, pubkeyBHex})
	require.NoError(t, err)
	err = os.WriteFile(tmpFile, content, 0o644)
	require.NoError(t, err)

	server.OnSIGHUP()

	require.True(t, wl.IsWhitelisted(pkA))
	require.True(t, wl.IsWhitelisted(pkB))
}

// mockPayloadProvider implements PayloadProvider for tests.
type mockPayloadProvider struct {
	hasPayload bool
	returnErr  error
}

func (m *mockPayloadProvider) GetPayloadBySlot(
	_ context.Context,
	_ math.Slot,
	_ common.Root,
) (ctypes.BuiltExecutionPayloadEnv, error) {
	if m.returnErr != nil {
		return nil, m.returnErr
	}
	if !m.hasPayload {
		return nil, payloadbuilder.ErrPayloadIDNotFound
	}
	return &mockPayloadEnvelope{forkVersion: version.Deneb1()}, nil
}

// mockPayloadEnvelope implements BuiltExecutionPayloadEnv.
type mockPayloadEnvelope struct {
	forkVersion common.Version
}

func (m *mockPayloadEnvelope) GetExecutionPayload() *ctypes.ExecutionPayload {
	return ctypes.NewEmptyExecutionPayloadWithVersion(m.forkVersion)
}

func (m *mockPayloadEnvelope) GetBlockValue() *math.U256 { return nil }

func (m *mockPayloadEnvelope) GetBlobsBundle() engineprimitives.BlobsBundle { return nil }

func (m *mockPayloadEnvelope) GetEncodedExecutionRequests() []ctypes.EncodedExecutionRequest {
	return nil
}

func (m *mockPayloadEnvelope) ShouldOverrideBuilder() bool { return false }

// recordingSink captures counter increments for label assertions.
type recordingSink struct {
	counters map[string][]string
}

func newRecordingSink() *recordingSink {
	return &recordingSink{counters: make(map[string][]string)}
}

func (r *recordingSink) IncrementCounter(key string, args ...string) {
	// Expect labels as k,v pairs; record the first value label.
	if len(args) >= 2 {
		r.counters[key] = append(r.counters[key], args[1])
	}
}

func TestServer_MetricsLabels(t *testing.T) {
	t.Parallel()

	validatorA, _ := parser.ConvertPubkey(pubkeyAHex)
	validatorB, _ := parser.ConvertPubkey(pubkeyBHex)
	secretA, _ := jwt.NewFromHex(secretAHex)
	secretB, _ := jwt.NewFromHex(secretBHex)

	const (
		payloadKey  = "beacon_kit.preconf.server.payload_request_total"
		proposerKey = "beacon_kit.preconf.proposer_tracker.check_total"
		targetSlot  = math.Slot(100)
	)

	tests := []struct {
		name         string
		requestJWT   *jwt.Secret
		hasPayload   bool
		providerErr  error
		wantResult   preconf.ServerResult
		wantProposer string // empty when the proposer check isn't reached
		wantStatus   int
	}{
		{
			name:         "happy path emits ok",
			requestJWT:   secretA,
			hasPayload:   true,
			wantResult:   preconf.ServerResultOK,
			wantProposer: "true",
			wantStatus:   http.StatusOK,
		},
		{
			name:         "wrong proposer emits wrong_proposer",
			requestJWT:   secretB,
			hasPayload:   true,
			wantResult:   preconf.ServerResultWrongProposer,
			wantProposer: "false",
			wantStatus:   http.StatusForbidden,
		},
		{
			name:         "missing payload emits payload_not_found",
			requestJWT:   secretA,
			hasPayload:   false,
			wantResult:   preconf.ServerResultPayloadNotFound,
			wantProposer: "true",
			wantStatus:   http.StatusNotFound,
		},
		{
			name:         "provider internal error emits internal_error with 500",
			requestJWT:   secretA,
			providerErr:  errors.New("EL exploded"),
			wantResult:   preconf.ServerResultInternalError,
			wantProposer: "true",
			wantStatus:   http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sink := newRecordingSink()
			tracker := preconf.NewProposerTracker()
			tracker.SetExpectedProposer(targetSlot, validatorA)

			server := preconf.NewServer(
				noop.NewLogger[any](),
				preconf.ValidatorJWTs{validatorA: secretA, validatorB: secretB},
				newTestWhitelist(t, pubkeyAHex, pubkeyBHex),
				tracker,
				&mockPayloadProvider{hasPayload: tt.hasPayload, returnErr: tt.providerErr},
				0,
				sink,
			)

			body, _ := json.Marshal(preconf.GetPayloadRequest{Slot: targetSlot})
			req := httptest.NewRequest(http.MethodPost, preconf.PayloadEndpoint, bytes.NewReader(body))
			token, _ := tt.requestJWT.BuildSignedToken()
			req.Header.Set("Authorization", "Bearer "+token)

			rec := httptest.NewRecorder()
			server.Handler().ServeHTTP(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)
			require.Equal(t, []string{string(tt.wantResult)}, sink.counters[payloadKey])
			if tt.wantProposer == "" {
				require.Empty(t, sink.counters[proposerKey])
			} else {
				require.Equal(t, []string{tt.wantProposer}, sink.counters[proposerKey])
			}
		})
	}
}
