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
	"testing"

	"github.com/berachain/beacon-kit/beacon/preconf"
	"github.com/berachain/beacon-kit/cli/utils/parser"
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			provider := &mockPayloadProvider{
				hasPayload: tt.payloadExists,
			}
			server := preconf.NewServer(
				noop.NewLogger[any](),
				preconf.ValidatorJWTs{validatorA: secretA, validatorB: secretB},
				preconf.NewWhitelist([]crypto.BLSPubkey{validatorA, validatorB}),
				provider,
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

func TestServer_RejectsNonPostMethods(t *testing.T) {
	t.Parallel()

	server := preconf.NewServer(noop.NewLogger[any](), nil, nil, nil, 0)

	for _, method := range []string{http.MethodGet, http.MethodPut, http.MethodDelete} {
		req := httptest.NewRequest(method, preconf.PayloadEndpoint, nil)
		rec := httptest.NewRecorder()
		server.Handler().ServeHTTP(rec, req)
		require.Equal(t, http.StatusMethodNotAllowed, rec.Code, "method: %s", method)
	}
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
