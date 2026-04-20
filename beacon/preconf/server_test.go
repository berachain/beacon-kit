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
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/berachain/beacon-kit/beacon/preconf"
	"github.com/berachain/beacon-kit/cli/utils/parser"
	"github.com/berachain/beacon-kit/errors"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/log/noop"
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
				&mockSyncChecker{ready: true},
				&mockELChecker{connected: true},
				0,
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
				&mockSyncChecker{ready: true},
				&mockELChecker{connected: true},
				0,
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

	server := preconf.NewServer(noop.NewLogger[any](), nil, nil, nil, nil, nil, nil, 0)

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

	server := preconf.NewServer(noop.NewLogger[any](), nil, wl, nil, nil, nil, nil, 0)

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
}

func (m *mockPayloadProvider) GetPayloadBySlot(
	_ context.Context,
	_ math.Slot,
	_ common.Root,
) (ctypes.BuiltExecutionPayloadEnv, error) {
	if !m.hasPayload {
		return nil, preconf.ErrPayloadNotFound
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

// mockSyncChecker implements preconf.SyncChecker for tests.
type mockSyncChecker struct {
	ready        bool
	latestHeight int64
	syncToHeight int64
}

func (m *mockSyncChecker) IsAppReady() error {
	if !m.ready {
		return errors.New("app not ready")
	}
	return nil
}

func (m *mockSyncChecker) GetSyncData() (int64, int64) {
	return m.latestHeight, m.syncToHeight
}

// mockELChecker implements preconf.ELChecker for tests.
type mockELChecker struct {
	connected bool
}

func (m *mockELChecker) IsConnected() bool {
	return m.connected
}

func TestServer_HealthEndpoint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		syncChecker   *mockSyncChecker
		elChecker     *mockELChecker
		wantStatus    int
		wantReady     bool
		wantSync      bool
		wantELConn    bool
	}{
		{
			name:        "healthy - synced, ready, EL connected",
			syncChecker: &mockSyncChecker{ready: true, latestHeight: 100, syncToHeight: 100},
			elChecker:   &mockELChecker{connected: true},
			wantStatus:  http.StatusOK,
			wantReady:   true,
			wantSync:    false,
			wantELConn:  true,
		},
		{
			name:        "unhealthy - still syncing",
			syncChecker: &mockSyncChecker{ready: true, latestHeight: 50, syncToHeight: 100},
			elChecker:   &mockELChecker{connected: true},
			wantStatus:  http.StatusServiceUnavailable,
			wantReady:   true,
			wantSync:    true,
			wantELConn:  true,
		},
		{
			name:        "unhealthy - app not ready",
			syncChecker: &mockSyncChecker{ready: false, latestHeight: 0, syncToHeight: 0},
			elChecker:   &mockELChecker{connected: true},
			wantStatus:  http.StatusServiceUnavailable,
			wantReady:   false,
			wantSync:    false,
			wantELConn:  true,
		},
		{
			name:        "unhealthy - EL disconnected",
			syncChecker: &mockSyncChecker{ready: true, latestHeight: 100, syncToHeight: 100},
			elChecker:   &mockELChecker{connected: false},
			wantStatus:  http.StatusServiceUnavailable,
			wantReady:   true,
			wantSync:    false,
			wantELConn:  false,
		},
		{
			name:       "healthy - nil checkers",
			wantStatus: http.StatusOK,
			wantReady:  true,
			wantSync:   false,
			wantELConn: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var syncChecker preconf.SyncChecker
			if tt.syncChecker != nil {
				syncChecker = tt.syncChecker
			}
			var elChecker preconf.ELChecker
			if tt.elChecker != nil {
				elChecker = tt.elChecker
			}

			server := preconf.NewServer(
				noop.NewLogger[any](), nil, nil, nil, nil,
				syncChecker, elChecker, 0,
			)

			req := httptest.NewRequest(http.MethodGet, preconf.HealthEndpoint, nil)
			rec := httptest.NewRecorder()
			server.Handler().ServeHTTP(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)

			var resp preconf.HealthResponse
			err := json.NewDecoder(rec.Body).Decode(&resp)
			require.NoError(t, err)
			require.Equal(t, tt.wantReady, resp.IsReady)
			require.Equal(t, tt.wantSync, resp.IsSyncing)
			require.Equal(t, tt.wantELConn, resp.ELConnected)
		})
	}
}
