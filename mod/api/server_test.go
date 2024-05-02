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

package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	middleware "github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
)

type testcase struct {
	method   string
	endpoint string
	expected int
}

func TestEndpoints(t *testing.T) {
	e := NewServer(middleware.DefaultCORSConfig,
		middleware.DefaultLoggerConfig)

	for _, testcase := range getTestcases() {
		testcase.endpoint = remapParams(testcase.endpoint)
		req := httptest.NewRequest(testcase.method, testcase.endpoint, nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, testcase.expected, rec.Code,
			"Expected status code %d but got %d for path %s",
			testcase.expected, rec.Code, testcase.endpoint)
	}
}

//nolint:golint,maintidx // defining test cases for each route
func getTestcases() []testcase {
	return []testcase{
		{
			method:   "GET",
			endpoint: "/eth/v1/beacon/genesis",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/beacon/states/:state_id/root",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/beacon/states/:state_id/fork",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/beacon/states/:state_id/finality_checkpoints",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/beacon/states/:state_id/validators",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "POST",
			endpoint: "/eth/v1/beacon/states/:state_id/validators",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/beacon/states/:state_id/validators/:validator_id",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/beacon/states/:state_id/validators/validator_balances",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "POST",
			endpoint: "/eth/v1/beacon/states/:state_id/validators/validator_balances",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/beacon/states/:state_id/committees",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/beacon/states/:state_id/sync_committees",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/beacon/states/:state_id/randao",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/beacon/headers",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/beacon/headers/:block_id",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "POST",
			endpoint: "/eth/v1/beacon/blocks/blinded_blocks",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "POST",
			endpoint: "/eth/v2/beacon/blocks/blinded_blocks",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "POST",
			endpoint: "/eth/v1/beacon/blocks",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "POST",
			endpoint: "/eth/v2/beacon/blocks",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v2/beacon/blocks/:block_id",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/beacon/blocks/:block_id/root",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/beacon/blocks/:block_id/attestations",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/beacon/blob_sidecars/:block_id",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "POST",
			endpoint: "/eth/v1/beacon/rewards/sync_committee/:block_id",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/beacon/deposit_snapshot",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/beacon/rewards/block/:block_id",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "POST",
			endpoint: "/eth/v1/beacon/rewards/attestation/:epoch",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/beacon/blinded_blocks/:block_id",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/beacon/light_client/bootstrap/:block_root",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/beacon/light_client/updates",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/beacon/light_client/finality_update",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/beacon/light_client/optimistic_update",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/beacon/pool/attestations",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "POST",
			endpoint: "/eth/v1/beacon/pool/attestations",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/beacon/pool/attester_slashings",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "POST",
			endpoint: "/eth/v1/beacon/pool/attester_slashings",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/beacon/pool/proposer_slashings",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "POST",
			endpoint: "/eth/v1/beacon/pool/proposer_slashings",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "POST",
			endpoint: "/eth/v1/beacon/pool/sync_committees",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/beacon/pool/voluntary_exits",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "POST",
			endpoint: "/eth/v1/beacon/pool/voluntary_exits",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/beacon/pool/bls_to_execution_changes",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "POST",
			endpoint: "/eth/v1/beacon/pool/bls_to_execution_changes",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/builder/states/{state_id}/expected_withdrawals",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/config/fork_schedule",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/config/spec",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/config/deposit_contract",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v2/debug/beacon/states/:state_id",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v2/debug/beacon/states/heads",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/debug/fork_choice",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/events?topics=head&topics=proposer_slashing",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/node/identity",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/node/peers",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/node/peers/:peer_id",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/node/peers/peer_count",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/node/version",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/node/syncing",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/node/health",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "POST",
			endpoint: "/eth/v1/validator/duties/attester/:epoch",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/validator/duties/proposer/:epoch",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "POST",
			endpoint: "/eth/v1/validator/duties/sync/:epoch",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v3/validator/blocks/:slot",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/validator/attestation_data",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/validator/aggregate_attestation",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "POST",
			endpoint: "/eth/v1/validator/aggregate_and_proofs",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "POST",
			endpoint: "/eth/v1/validator/beacon_committee_subscriptions",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "POST",
			endpoint: "/eth/v1/validator/sync_committee_subscriptions",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "POST",
			endpoint: "/eth/v1/validator/beacon_committee_selections",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "GET",
			endpoint: "/eth/v1/validator/sync_committee_contribution",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "POST",
			endpoint: "/eth/v1/validator/sync_committee_selections",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "POST",
			endpoint: "/eth/v1/validator/contribution_and_proofs",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "POST",
			endpoint: "/eth/v1/validator/prepare_beacon_proposer",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "POST",
			endpoint: "/eth/v1/validator/register_validator",
			expected: http.StatusNotImplemented,
		},
		{
			method:   "POST",
			endpoint: "/eth/v1/validator/liveness/:epoch",
			expected: http.StatusNotImplemented,
		},
	}
}

func remapParams(url string) string {
	url = strings.ReplaceAll(url, ":epoch", "1")
	url = strings.ReplaceAll(url, ":state_id", "head")
	url = strings.ReplaceAll(url, ":block_id", "head")
	url = strings.ReplaceAll(url, ":slot", "1")
	url = strings.ReplaceAll(url, ":peer_id",
		"QmYyQSo1c1Ym7orWxLYvCrM2EmxFTANf8wXmmE7DWjhx5N")
	url = strings.ReplaceAll(url, ":block_root",
		"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2")
	url = strings.ReplaceAll(url, ":validator_id", "1")
	return url
}
