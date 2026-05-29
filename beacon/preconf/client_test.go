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
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/berachain/beacon-kit/beacon/preconf"
	"github.com/berachain/beacon-kit/log/noop"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/net/jwt"
	"github.com/stretchr/testify/require"
)

// TestClient_RejectsOversizedResponse verifies the client caps the sequencer
// response body so a malicious/compromised sequencer cannot OOM the validator.
func TestClient_RejectsOversizedResponse(t *testing.T) {
	t.Parallel()

	// Cap the client at 1MiB and return a 2MiB response, over the cap.
	const maxSize = 1024 * 1024
	huge := make([]byte, 2*maxSize)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		//nolint:errcheck // test handler
		w.Write(huge)
	}))
	defer srv.Close()

	secret, err := jwt.NewFromHex(secretAHex)
	require.NoError(t, err)

	client := preconf.NewClient(noop.NewLogger[any](), srv.URL, secret, 5*time.Second, nil, maxSize)

	_, err = client.GetPayloadBySlot(t.Context(), 100, common.Root{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "exceeds")
}
