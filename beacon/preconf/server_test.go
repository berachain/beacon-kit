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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

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
				&mockSyncChecker{ready: true},
				&mockELChecker{connected: true},
				0,
				preconf.TLSPaths{},
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
				&mockSyncChecker{ready: true},
				&mockELChecker{connected: true},
				0,
				preconf.TLSPaths{},
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

func TestServer_RejectsOversizedRequestBody(t *testing.T) {
	t.Parallel()

	validatorA, _ := parser.ConvertPubkey(pubkeyAHex)
	secretA, _ := jwt.NewFromHex(secretAHex)

	server := preconf.NewServer(
		noop.NewLogger[any](),
		preconf.ValidatorJWTs{validatorA: secretA},
		newTestWhitelist(t, pubkeyAHex),
		preconf.NewProposerTracker(),
		&mockPayloadProvider{hasPayload: true},
		&mockSyncChecker{ready: true},
		&mockELChecker{connected: true},
		0,
		preconf.TLSPaths{},
		metrics.NewNoOpTelemetrySink(),
	)

	// A valid-JSON body well over the 2KB cap: a long string value forces the decoder to read past the limit, where MaxBytesReader trips.
	// Authenticated and whitelisted so the request reaches body parsing.
	body := append([]byte(`{"slot":1,"pad":"`), bytes.Repeat([]byte("a"), 4096)...)
	body = append(body, '"', '}')
	req := httptest.NewRequest(http.MethodPost, preconf.PayloadEndpoint, bytes.NewReader(body))
	token, _ := secretA.BuildSignedToken()
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)

	require.Equal(t, http.StatusRequestEntityTooLarge, rec.Code)
	require.Contains(t, rec.Body.String(), "request body too large")
}

func TestServer_RejectsNonPostMethods(t *testing.T) {
	t.Parallel()

	server := preconf.NewServer(noop.NewLogger[any](), nil, nil, nil, nil, nil, nil, 0, preconf.TLSPaths{}, metrics.NewNoOpTelemetrySink())

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

	server := preconf.NewServer(noop.NewLogger[any](), nil, wl, nil, nil, nil, nil, 0, preconf.TLSPaths{}, metrics.NewNoOpTelemetrySink())

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
				&mockSyncChecker{ready: true},
				&mockELChecker{connected: true},
				0,
				preconf.TLSPaths{},
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

// TestConfig_Validate checks the TLS config rules enforced at startup: cert and
// key must be set together, and a pinned CA requires an https sequencer URL.
func TestConfig_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cfg     preconf.Config
		wantErr string
	}{
		{
			name: "both TLS paths set is valid",
			cfg:  preconf.Config{TLSCertPath: "/cert.pem", TLSKeyPath: "/key.pem"},
		},
		{
			name: "neither TLS path set is valid",
			cfg:  preconf.Config{},
		},
		{
			name:    "only cert path set",
			cfg:     preconf.Config{TLSCertPath: "/cert.pem"},
			wantErr: "tls-cert-path and tls-key-path must both be set or both be empty",
		},
		{
			name:    "only key path set",
			cfg:     preconf.Config{TLSKeyPath: "/key.pem"},
			wantErr: "tls-cert-path and tls-key-path must both be set or both be empty",
		},
		{
			name:    "CA cert without sequencer URL",
			cfg:     preconf.Config{SequencerCACertPath: "/ca.pem"},
			wantErr: "sequencer-ca-cert-path requires sequencer-url to be set",
		},
		{
			name:    "CA cert with http URL",
			cfg:     preconf.Config{SequencerCACertPath: "/ca.pem", SequencerURL: "http://seq:9090"},
			wantErr: "sequencer-url must use https:// scheme",
		},
		{
			name: "CA cert with https URL is valid",
			cfg:  preconf.Config{SequencerCACertPath: "/ca.pem", SequencerURL: "https://seq:9090"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.cfg.Validate()
			if tt.wantErr != "" {
				require.ErrorContains(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestConfig_TLSEnabled checks that TLS counts as enabled only when both the
// cert and key paths are present.
func TestConfig_TLSEnabled(t *testing.T) {
	t.Parallel()

	require.False(t, (&preconf.Config{}).TLSEnabled())
	require.False(t, (&preconf.Config{TLSCertPath: "/cert.pem"}).TLSEnabled())
	require.False(t, (&preconf.Config{TLSKeyPath: "/key.pem"}).TLSEnabled())
	require.True(t, (&preconf.Config{TLSCertPath: "/cert.pem", TLSKeyPath: "/key.pem"}).TLSEnabled())
}

// TestServer_TLS_RejectsPlaintext asserts that once TLS is configured the server
// only speaks HTTPS: a plaintext HTTP request to the same port fails.
func TestServer_TLS_RejectsPlaintext(t *testing.T) {
	t.Parallel()

	certPath, keyPath := generateSelfSignedCert(t)
	port := freePort(t)

	server := preconf.NewServer(
		noop.NewLogger[any](),
		nil,
		newTestWhitelist(t, pubkeyAHex),
		preconf.NewProposerTracker(),
		&mockPayloadProvider{},
		&mockSyncChecker{ready: true, latestHeight: 0, syncToHeight: 0},
		&mockELChecker{connected: true},
		port,
		preconf.TLSPaths{Cert: certPath, Key: keyPath},
		metrics.NewNoOpTelemetrySink(),
	)
	require.NoError(t, server.Start(t.Context()))
	t.Cleanup(func() { _ = server.Stop() })

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	// Confirm the TLS listener is up before probing with plaintext, so the
	// failure below is a plaintext rejection and not a connection-refused race.
	require.NotEmpty(t, fetchServedCert(t, addr))

	// Send a plaintext HTTP request to the TLS port. The server speaks TLS and
	// never falls back to cleartext (no dual-mode listener).
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://"+addr+preconf.HealthEndpoint, nil)
	require.NoError(t, err)
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	// Go's TLS server answers a cleartext request with a 400 ("client sent an
	// HTTP request to an HTTPS server") rather than a transport error. Either a
	// transport error or a non-200 response proves plaintext was not served.
	if err == nil {
		require.NotEqual(t, http.StatusOK, resp.StatusCode,
			"plaintext request to a TLS-only server must not be served")
	}
}

// generateSelfSignedCert creates a self-signed cert/key pair in a temp directory.
func generateSelfSignedCert(t *testing.T) (string, string) {
	t.Helper()
	dir := t.TempDir()
	certPath := filepath.Join(dir, "cert.pem")
	keyPath := filepath.Join(dir, "key.pem")
	writeSelfSignedCert(t, certPath, keyPath)
	return certPath, keyPath
}

// writeSelfSignedCert writes a fresh self-signed cert/key pair to the given
// paths. Each call generates a new key, so the resulting cert differs from any
// prior one written to the same paths (useful for exercising rotation).
func writeSelfSignedCert(t *testing.T, certPath, keyPath string) {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	template := &x509.Certificate{
		SerialNumber:          big.NewInt(time.Now().UnixNano()),
		Subject:               pkix.Name{Organization: []string{"test"}},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		// SAN so a pinned client verifying against 127.0.0.1 / localhost passes.
		IPAddresses: []net.IP{net.ParseIP("127.0.0.1")},
		DNSNames:    []string{"localhost"},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	require.NoError(t, err)

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	require.NoError(t, os.WriteFile(certPath, certPEM, 0o600))

	keyDER, err := x509.MarshalECPrivateKey(key)
	require.NoError(t, err)
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})
	require.NoError(t, os.WriteFile(keyPath, keyPEM, 0o600))
}

// freePort returns a TCP port that is free at the time of the call.
func freePort(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	tcpAddr, ok := ln.Addr().(*net.TCPAddr)
	require.True(t, ok)
	require.NoError(t, ln.Close())
	return tcpAddr.Port
}

// fetchServedCert dials the TLS server and returns the DER bytes of the leaf
// certificate it presents. It retries until the server is accepting.
func fetchServedCert(t *testing.T, addr string) []byte {
	t.Helper()
	var raw []byte
	require.Eventually(t, func() bool {
		// InsecureSkipVerify: the test inspects the served cert rather than verifying it.
		conn, err := tls.Dial("tcp", addr, &tls.Config{InsecureSkipVerify: true})
		if err != nil {
			return false
		}
		defer conn.Close()
		certs := conn.ConnectionState().PeerCertificates
		if len(certs) == 0 {
			return false
		}
		raw = certs[0].Raw
		return true
	}, 2*time.Second, 20*time.Millisecond)
	return raw
}

// TestServer_TLS_CertReload verifies SIGHUP swaps the served certificate: after
// rotating the cert files on disk, the server presents a different cert.
func TestServer_TLS_CertReload(t *testing.T) {
	t.Parallel()

	certPath, keyPath := generateSelfSignedCert(t)
	port := freePort(t)

	server := preconf.NewServer(
		noop.NewLogger[any](),
		nil,
		newTestWhitelist(t, pubkeyAHex),
		preconf.NewProposerTracker(),
		&mockPayloadProvider{},
		&mockSyncChecker{ready: true, latestHeight: 0, syncToHeight: 0},
		&mockELChecker{connected: true},
		port,
		preconf.TLSPaths{Cert: certPath, Key: keyPath},
		metrics.NewNoOpTelemetrySink(),
	)
	require.NoError(t, server.Start(t.Context()))
	t.Cleanup(func() { _ = server.Stop() })

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	before := fetchServedCert(t, addr)
	require.NotEmpty(t, before)

	// Rotate the cert on the same paths and signal a reload.
	writeSelfSignedCert(t, certPath, keyPath)
	server.OnSIGHUP()

	after := fetchServedCert(t, addr)
	require.NotEqual(t, before, after, "served cert should change after SIGHUP reload")
}

// TestServer_TLS_CertReload_KeepsOldOnError verifies a SIGHUP reload of a corrupt
// cert file is ignored, leaving the previously served cert in place.
func TestServer_TLS_CertReload_KeepsOldOnError(t *testing.T) {
	t.Parallel()

	certPath, keyPath := generateSelfSignedCert(t)
	port := freePort(t)

	server := preconf.NewServer(
		noop.NewLogger[any](),
		nil,
		newTestWhitelist(t, pubkeyAHex),
		preconf.NewProposerTracker(),
		&mockPayloadProvider{},
		&mockSyncChecker{ready: true, latestHeight: 0, syncToHeight: 0},
		&mockELChecker{connected: true},
		port,
		preconf.TLSPaths{Cert: certPath, Key: keyPath},
		metrics.NewNoOpTelemetrySink(),
	)
	require.NoError(t, server.Start(t.Context()))
	t.Cleanup(func() { _ = server.Stop() })

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	before := fetchServedCert(t, addr)
	require.NotEmpty(t, before)

	// A corrupt cert file must not break the listener: keep the old cert.
	require.NoError(t, os.WriteFile(certPath, []byte("not a certificate"), 0o600))
	server.OnSIGHUP()

	after := fetchServedCert(t, addr)
	require.Equal(t, before, after, "served cert should be unchanged after a failed reload")
}

// TestServer_TLS_ClientPinning verifies CA pinning end to end: a client pinned to
// the server's cert connects, while one pinned to a different CA is rejected.
func TestServer_TLS_ClientPinning(t *testing.T) {
	t.Parallel()

	certPath, keyPath := generateSelfSignedCert(t)
	port := freePort(t)

	server := preconf.NewServer(
		noop.NewLogger[any](),
		nil,
		newTestWhitelist(t, pubkeyAHex),
		preconf.NewProposerTracker(),
		&mockPayloadProvider{},
		&mockSyncChecker{ready: true, latestHeight: 0, syncToHeight: 0},
		&mockELChecker{connected: true},
		port,
		preconf.TLSPaths{Cert: certPath, Key: keyPath},
		metrics.NewNoOpTelemetrySink(),
	)
	require.NoError(t, server.Start(t.Context()))
	t.Cleanup(func() { _ = server.Stop() })

	url := fmt.Sprintf("https://127.0.0.1:%d", port)
	secret, err := jwt.NewFromHex(secretAHex)
	require.NoError(t, err)

	// A client pinned to the server's own cert trusts it. The health check
	// succeeds once the listener is accepting (retried to absorb startup).
	pinnedPool, err := preconf.LoadCACert(certPath)
	require.NoError(t, err)
	pinnedClient := preconf.NewClient(noop.NewLogger[any](), url, secret, 2*time.Second, pinnedPool, 0)
	require.Eventually(t, func() bool {
		return pinnedClient.CheckHealth(t.Context()) == nil
	}, 2*time.Second, 20*time.Millisecond, "pinned client should trust the server's cert")

	// A client pinned to a different CA rejects the server's cert. The server is
	// already up (the pinned check passed), so this fails at cert verification
	// rather than connection setup.
	otherCertPath, _ := generateSelfSignedCert(t)
	otherPool, err := preconf.LoadCACert(otherCertPath)
	require.NoError(t, err)
	otherClient := preconf.NewClient(noop.NewLogger[any](), url, secret, 2*time.Second, otherPool, 0)
	err = otherClient.CheckHealth(t.Context())
	require.ErrorContains(t, err, "certificate signed by unknown authority")
}

// TestServer_Start_BindFailure verifies a bind failure (port already in use)
// fails Start synchronously instead of leaving a dead endpoint behind.
func TestServer_Start_BindFailure(t *testing.T) {
	t.Parallel()

	// Occupy a port on all interfaces, matching what Start binds (":<port>").
	occupied, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = occupied.Close() })
	tcpAddr, ok := occupied.Addr().(*net.TCPAddr)
	require.True(t, ok)

	server := preconf.NewServer(
		noop.NewLogger[any](),
		nil,
		newTestWhitelist(t, pubkeyAHex),
		preconf.NewProposerTracker(),
		&mockPayloadProvider{},
		&mockSyncChecker{ready: true, latestHeight: 0, syncToHeight: 0},
		&mockELChecker{connected: true},
		tcpAddr.Port,
		preconf.TLSPaths{},
		metrics.NewNoOpTelemetrySink(),
	)
	t.Cleanup(func() { _ = server.Stop() })

	require.Error(t, server.Start(t.Context()), "Start should fail when the port is already in use")
}

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
		name        string
		syncChecker *mockSyncChecker
		elChecker   *mockELChecker
		wantStatus  int
		wantReady   bool
		wantSync    bool
		wantELConn  bool
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
			name:       "unhealthy - nil checkers",
			wantStatus: http.StatusServiceUnavailable,
			wantReady:  false,
			wantSync:   false,
			wantELConn: false,
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
				syncChecker, elChecker, 0, preconf.TLSPaths{}, metrics.NewNoOpTelemetrySink(),
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
