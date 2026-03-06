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
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/berachain/beacon-kit/beacon/preconf"
	"github.com/berachain/beacon-kit/log/noop"
	"github.com/berachain/beacon-kit/primitives/net/jwt"
	"github.com/stretchr/testify/require"
)

func newTestClient(t *testing.T, serverURL string, probeInterval time.Duration) *preconf.Client {
	t.Helper()
	secret, err := jwt.NewFromHex(secretAHex)
	require.NoError(t, err)
	return preconf.NewClient(
		noop.NewLogger[any](),
		serverURL,
		secret,
		500*time.Millisecond,
		probeInterval,
	)
}

func TestClient_IsAvailableByDefault(t *testing.T) {
	t.Parallel()
	client := newTestClient(t, "http://localhost:9999", time.Second)
	require.True(t, client.IsAvailable(), "client should be available before any call")
}

func TestClient_MarksUnavailableOnConnectionFailure(t *testing.T) {
	t.Parallel()

	// Point at a port with nothing listening.
	client := newTestClient(t, "http://"+closedAddr(t), time.Hour)

	_, err := client.GetPayloadBySlot(t.Context(), 1, [32]byte{})
	require.ErrorIs(t, err, preconf.ErrSequencerUnavailable)
	require.False(t, client.IsAvailable(), "client should be unavailable after connection failure")
}

func TestClient_RestoresAvailabilityOnSequencerRecover(t *testing.T) {
	t.Parallel()

	// Start a server that's initially down (nil handler = not started yet).
	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Use a short probe interval so the test doesn't wait long.
	client := newTestClient(t, "http://"+srv.Listener.Addr().String(), 20*time.Millisecond)

	// Trigger a failure to start the health monitor.
	_, err := client.GetPayloadBySlot(t.Context(), 1, [32]byte{})
	require.ErrorIs(t, err, preconf.ErrSequencerUnavailable)
	require.False(t, client.IsAvailable())

	// Bring the server up — monitor should detect it within probeInterval.
	srv.Start()
	defer srv.Close()

	require.Eventually(t,
		func() bool { return client.IsAvailable() },
		2*time.Second,
		10*time.Millisecond,
		"client should restore availability once sequencer is back",
	)
}

// closedAddr returns an address that is not currently listening, for testing connection failure.
func closedAddr(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	addr := l.Addr().String()
	l.Close()
	return addr
}
