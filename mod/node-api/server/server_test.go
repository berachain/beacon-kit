// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
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

package server

// type testcase struct {
// 	method         string
// 	endpoint       string
// 	body           string
// 	expectedStatus int
// 	expectedBody   string
// }

// func TestEndpoints(t *testing.T) {
// 	e := New(
// 		DefaultConfig(),
// 		backend.NewMockBackend(),
// 		middleware.DefaultCORSConfig,
// 		middleware.DefaultLoggerConfig,
// 	)

// 	for _, testcase := range getTestcases() {
// 		testcase.endpoint = remapParams(testcase.endpoint)
// 		req := buildRequest(testcase.method, testcase.endpoint, &testcase.body)
// 		rec := httptest.NewRecorder()
// 		e.ServeHTTP(rec, req)
// 		assert.Equal(t, testcase.expectedStatus, rec.Code,
// 			"Expected status code %d but got %d for path %s",
// 			testcase.expectedStatus, rec.Code, testcase.endpoint)
// 		if testcase.expectedBody != "" {
// 			assert.Equal(t, testcase.expectedBody, rec.Body.String(),
// 				"Expected response body %d but got %d for path %s",
// 				testcase.expectedBody, rec.Body.String(), testcase.endpoint)
// 		}
// 	}
// }

// func buildRequest(method, endpoint string, body *string) *http.Request {
// 	req := httptest.NewRequest(method, endpoint, nil)
// 	if method != "GET" && body != nil {
// 		req = httptest.NewRequest(method, endpoint,
// 			io.NopCloser(strings.NewReader(*body)))
// 		req.Header.Set("Content-Type", "application/json")
// 	}
// 	return req
// }

// //nolint:golint,maintidx,lll // defining test cases for each route
// func getTestcases() []testcase {
// 	return []testcase{
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/beacon/genesis",
// 			expectedStatus: http.StatusOK,
// 			expectedBody:   "{\"data\":{\"genesis_time\":\"1590832934\",\"genesis_validators_root\":\"0x0100000000000000000000000000000000000000000000000000000000000000\",\"genesis_fork_version\":\"0x00000000\"}}\n",
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/beacon/states/:state_id/root",
// 			expectedStatus: http.StatusOK,
// 			expectedBody:   "{\"execution_optimistic\":false,\"finalized\":false,\"data\":{\"data\":{\"root\":\"0x0100000000000000000000000000000000000000000000000000000000000000\"}}}\n",
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/beacon/states/:state_id/fork",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/beacon/states/:state_id/finality_checkpoints",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/beacon/states/:state_id/validators?id=1",
// 			expectedStatus: http.StatusOK,
// 			expectedBody:   "{\"execution_optimistic\":false,\"finalized\":false,\"data\":[{\"index\":\"1\",\"balance\":\"1\",\"status\":\"active\",\"validator\":{\"pubkey\":\"0x010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000\",\"withdrawalCredentials\":\"0x0100000000000000000000000000000000000000000000000000000000000000\",\"effectiveBalance\":\"0x0\",\"slashed\":false,\"activationEligibilityEpoch\":\"0x0\",\"activationEpoch\":\"0x0\",\"exitEpoch\":\"0x0\",\"withdrawableEpoch\":\"0x0\"}}]}\n",
// 		},
// 		{
// 			method:         "POST",
// 			endpoint:       "/eth/v1/beacon/states/:state_id/validators",
// 			body:           `{"ids":["1"]}`,
// 			expectedStatus: http.StatusOK,
// 			expectedBody:   "{\"execution_optimistic\":false,\"finalized\":false,\"data\":[{\"index\":\"1\",\"balance\":\"1\",\"status\":\"active\",\"validator\":{\"pubkey\":\"0x010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000\",\"withdrawalCredentials\":\"0x0100000000000000000000000000000000000000000000000000000000000000\",\"effectiveBalance\":\"0x0\",\"slashed\":false,\"activationEligibilityEpoch\":\"0x0\",\"activationEpoch\":\"0x0\",\"exitEpoch\":\"0x0\",\"withdrawableEpoch\":\"0x0\"}}]}\n",
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/beacon/states/:state_id/validators/:validator_id",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/beacon/states/:state_id/validator_balances?id=1",
// 			expectedStatus: http.StatusOK,
// 			expectedBody:   "{\"execution_optimistic\":false,\"finalized\":false,\"data\":[{\"index\":\"1\",\"balance\":\"1\"}]}\n",
// 		},
// 		{
// 			method:         "POST",
// 			endpoint:       "/eth/v1/beacon/states/:state_id/validator_balances",
// 			body:           `["1"]`,
// 			expectedStatus: http.StatusOK,
// 			expectedBody:   "{\"execution_optimistic\":false,\"finalized\":false,\"data\":[{\"index\":\"1\",\"balance\":\"1\"}]}\n",
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/beacon/states/:state_id/committees",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/beacon/states/:state_id/sync_committees",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/beacon/states/:state_id/randao",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/beacon/headers",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/beacon/headers/:block_id",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "POST",
// 			endpoint:       "/eth/v1/beacon/blocks/blinded_blocks",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "POST",
// 			endpoint:       "/eth/v2/beacon/blocks/blinded_blocks",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "POST",
// 			endpoint:       "/eth/v1/beacon/blocks",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "POST",
// 			endpoint:       "/eth/v2/beacon/blocks",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v2/beacon/blocks/:block_id",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/beacon/blocks/:block_id/root",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/beacon/blocks/:block_id/attestations",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/beacon/blob_sidecars/:block_id",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "POST",
// 			endpoint:       "/eth/v1/beacon/rewards/sync_committee/:block_id",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/beacon/deposit_snapshot",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/beacon/rewards/blocks/:block_id",
// 			expectedStatus: http.StatusOK,
// 			expectedBody:   "{\"execution_optimistic\":false,\"finalized\":false,\"data\":{\"proposer_index\":\"1\",\"total\":\"1\",\"attestations\":\"1\",\"sync_aggregate\":\"1\",\"proposer_slashings\":\"1\",\"attester_slashings\":\"1\"}}\n",
// 		},
// 		{
// 			method:         "POST",
// 			endpoint:       "/eth/v1/beacon/rewards/attestation/:epoch",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/beacon/blinded_blocks/:block_id",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/beacon/light_client/bootstrap/:block_root",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/beacon/light_client/updates",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/beacon/light_client/finality_update",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/beacon/light_client/optimistic_update",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/beacon/pool/attestations",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "POST",
// 			endpoint:       "/eth/v1/beacon/pool/attestations",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/beacon/pool/attester_slashings",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "POST",
// 			endpoint:       "/eth/v1/beacon/pool/attester_slashings",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/beacon/pool/proposer_slashings",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "POST",
// 			endpoint:       "/eth/v1/beacon/pool/proposer_slashings",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "POST",
// 			endpoint:       "/eth/v1/beacon/pool/sync_committees",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/beacon/pool/voluntary_exits",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "POST",
// 			endpoint:       "/eth/v1/beacon/pool/voluntary_exits",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/beacon/pool/bls_to_execution_changes",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "POST",
// 			endpoint:       "/eth/v1/beacon/pool/bls_to_execution_changes",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/builder/states/{state_id}/expected_withdrawals",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/config/fork_schedule",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/config/spec",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/config/deposit_contract",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v2/debug/beacon/states/:state_id",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v2/debug/beacon/states/heads",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/debug/fork_choice",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/events?topics=head&topics=proposer_slashing",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/node/identity",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/node/peers",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/node/peers/:peer_id",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/node/peers/peer_count",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/node/version",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/node/syncing",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/node/health",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "POST",
// 			endpoint:       "/eth/v1/validator/duties/attester/:epoch",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/validator/duties/proposer/:epoch",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "POST",
// 			endpoint:       "/eth/v1/validator/duties/sync/:epoch",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v3/validator/blocks/:slot",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/validator/attestation_data",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/validator/aggregate_attestation",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "POST",
// 			endpoint:       "/eth/v1/validator/aggregate_and_proofs",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "POST",
// 			endpoint:       "/eth/v1/validator/beacon_committee_subscriptions",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "POST",
// 			endpoint:       "/eth/v1/validator/sync_committee_subscriptions",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "POST",
// 			endpoint:       "/eth/v1/validator/beacon_committee_selections",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "GET",
// 			endpoint:       "/eth/v1/validator/sync_committee_contribution",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "POST",
// 			endpoint:       "/eth/v1/validator/sync_committee_selections",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "POST",
// 			endpoint:       "/eth/v1/validator/contribution_and_proofs",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "POST",
// 			endpoint:       "/eth/v1/validator/prepare_beacon_proposer",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "POST",
// 			endpoint:       "/eth/v1/validator/register_validator",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 		{
// 			method:         "POST",
// 			endpoint:       "/eth/v1/validator/liveness/:epoch",
// 			expectedStatus: http.StatusNotImplemented,
// 		},
// 	}
// }

// func remapParams(url string) string {
// 	url = strings.ReplaceAll(url, ":epoch", "1")
// 	url = strings.ReplaceAll(url, ":state_id", "head")
// 	url = strings.ReplaceAll(url, ":block_id", "head")
// 	url = strings.ReplaceAll(url, ":slot", "1")
// 	url = strings.ReplaceAll(url, ":peer_id",
// 		"QmYyQSo1c1Ym7orWxLYvCrM2EmxFTANf8wXmmE7DWjhx5N")
// 	url = strings.ReplaceAll(url, ":block_root",
// 		"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2")
// 	url = strings.ReplaceAll(url, ":validator_id", "1")
// 	return url
// }
