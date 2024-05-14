package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	middleware "github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
)

func TestEndpoints(t *testing.T) {
	e := NewServer(middleware.DefaultCORSConfig, middleware.DefaultLoggerConfig)

	tests := []struct {
		method         string
		target         string
		expectedStatus int
	}{
		{method: "GET", target: "/eth/v1/beacon/genesis", expectedStatus: http.StatusOK},
		{method: "GET", target: "/eth/v1/beacon/states/:state_id/root", expectedStatus: http.StatusOK},
		{method: "GET", target: "/eth/v1/beacon/states/:state_id/fork", expectedStatus: http.StatusOK},
		{method: "GET", target: "/eth/v1/beacon/states/:state_id/finality_checkpoints", expectedStatus: http.StatusOK},
		{method: "GET", target: "/eth/v1/beacon/states/:state_id/validators", expectedStatus: http.StatusOK},
		// {method: "POST", target: "/eth/v1/beacon/states/:state_id/validators", expectedStatus: http.StatusOK},
		// {method: "GET", target: "/eth/v1/beacon/states/:state_id/validators/:validator_id", expectedStatus: http.StatusOK},
		{method: "GET", target: "/eth/v1/beacon/states/:state_id/validators/validator_balances", expectedStatus: http.StatusOK},
		{method: "POST", target: "/eth/v1/beacon/states/:state_id/validators/validator_balances", expectedStatus: http.StatusOK},
		{method: "GET", target: "/eth/v1/beacon/states/:state_id/committees", expectedStatus: http.StatusOK},
		{method: "GET", target: "/eth/v1/beacon/states/:state_id/sync_committees", expectedStatus: http.StatusOK},
		// {method: "GET", target: "/eth/v1/beacon/states/:state_id/randao", expectedStatus: http.StatusOK},
		// {method: "GET", target: "/eth/v1/beacon/headers", expectedStatus: http.StatusOK},
		// {method: "GET", target: "/eth/v1/beacon/headers/:block_id", expectedStatus: http.StatusOK},
		{method: "POST", target: "/eth/v1/beacon/blocks/blinded_blocks", expectedStatus: http.StatusNotImplemented},
		{method: "POST", target: "/eth/v2/beacon/blocks/blinded_blocks", expectedStatus: http.StatusNotImplemented},
		{method: "POST", target: "/eth/v1/beacon/blocks", expectedStatus: http.StatusNotImplemented},
		{method: "POST", target: "/eth/v2/beacon/blocks", expectedStatus: http.StatusNotImplemented},
		// {method: "GET", target: "/eth/v2/beacon/blocks/:block_id", expectedStatus: http.StatusOK},
		// {method: "GET", target: "/eth/v1/beacon/blocks/:block_id/root", expectedStatus: http.StatusOK},
		// {method: "GET", target: "/eth/v1/beacon/blocks/:block_id/attestations", expectedStatus: http.StatusOK},
		// {method: "GET", target: "/eth/v1/beacon/blob_sidecars/:block_id", expectedStatus: http.StatusOK},
		// {method: "POST", target: "/eth/v1/beacon/rewards/sync_committee/:block_id", expectedStatus: http.StatusOK},
		{method: "GET", target: "/eth/v1/beacon/deposit_snapshot", expectedStatus: http.StatusOK},
		// {method: "GET", target: "/eth/v1/beacon/rewards/block/:block_id", expectedStatus: http.StatusOK},
		{method: "POST", target: "/eth/v1/beacon/rewards/attestation/:epoch", expectedStatus: http.StatusOK},
		// {method: "GET", target: "/eth/v1/beacon/blinded_blocks/:block_id", expectedStatus: http.StatusOK},
		{method: "GET", target: "/eth/v1/beacon/light_client/bootstrap/:block_root", expectedStatus: http.StatusNotImplemented},
		{method: "GET", target: "/eth/v1/beacon/light_client/updates", expectedStatus: http.StatusNotImplemented},
		{method: "GET", target: "/eth/v1/beacon/light_client/finality_update", expectedStatus: http.StatusNotImplemented},
		{method: "GET", target: "/eth/v1/beacon/light_client/optimistic_update", expectedStatus: http.StatusNotImplemented},
		// {method: "GET", target: "/eth/v1/beacon/pool/attestations", expectedStatus: http.StatusOK},
		{method: "POST", target: "/eth/v1/beacon/pool/attestations", expectedStatus: http.StatusOK},
		{method: "GET", target: "/eth/v1/beacon/pool/attester_slashings", expectedStatus: http.StatusNotImplemented},
		{method: "POST", target: "/eth/v1/beacon/pool/attester_slashings", expectedStatus: http.StatusNotImplemented},
		{method: "GET", target: "/eth/v1/beacon/pool/proposer_slashings", expectedStatus: http.StatusNotImplemented},
		{method: "POST", target: "/eth/v1/beacon/pool/proposer_slashings", expectedStatus: http.StatusNotImplemented},
		{method: "POST", target: "/eth/v1/beacon/pool/sync_committees", expectedStatus: http.StatusOK},
		{method: "GET", target: "/eth/v1/beacon/pool/voluntary_exits", expectedStatus: http.StatusOK},
		{method: "POST", target: "/eth/v1/beacon/pool/voluntary_exits", expectedStatus: http.StatusOK},
		{method: "GET", target: "/eth/v1/beacon/pool/bls_to_execution_changes", expectedStatus: http.StatusOK},
		{method: "POST", target: "/eth/v1/beacon/pool/bls_to_execution_changes", expectedStatus: http.StatusOK},
		{method: "GET", target: "/eth/v1/builder/states/{state_id}/expected_withdrawals", expectedStatus: http.StatusNotImplemented},
		{method: "GET", target: "/eth/v1/config/fork_schedule", expectedStatus: http.StatusNotImplemented},
		{method: "GET", target: "/eth/v1/config/spec", expectedStatus: http.StatusOK},
		{method: "GET", target: "/eth/v1/config/deposit_contract", expectedStatus: http.StatusNotImplemented},
		{method: "GET", target: "/eth/v2/debug/beacon/states/:state_id", expectedStatus: http.StatusNotImplemented},
		{method: "GET", target: "/eth/v2/debug/beacon/states/heads", expectedStatus: http.StatusNotImplemented},
		{method: "GET", target: "/eth/v1/debug/fork_choice", expectedStatus: http.StatusNotImplemented},
		{method: "GET", target: "/eth/v1/events?topics=head&topics=proposer_slashing", expectedStatus: http.StatusOK},
		{method: "GET", target: "/eth/v1/node/identity", expectedStatus: http.StatusNotImplemented},
		{method: "GET", target: "/eth/v1/node/peers", expectedStatus: http.StatusNotImplemented},
		{method: "GET", target: "/eth/v1/node/peers/:peer_id", expectedStatus: http.StatusNotImplemented},
		{method: "GET", target: "/eth/v1/node/peers/peer_count", expectedStatus: http.StatusNotImplemented},
		{method: "GET", target: "/eth/v1/node/version", expectedStatus: http.StatusNotImplemented},
		{method: "GET", target: "/eth/v1/node/syncing", expectedStatus: http.StatusNotImplemented},
		{method: "GET", target: "/eth/v1/node/health", expectedStatus: http.StatusNotImplemented},
		{method: "POST", target: "/eth/v1/validator/duties/attester/:epoch", expectedStatus: http.StatusNotImplemented},
		{method: "GET", target: "/eth/v1/validator/duties/proposer/:epoch", expectedStatus: http.StatusOK},
		{method: "POST", target: "/eth/v1/validator/duties/sync/:epoch", expectedStatus: http.StatusNotImplemented},
		{method: "GET", target: "/eth/v3/validator/blocks/:slot", expectedStatus: http.StatusNotImplemented},
		{method: "GET", target: "/eth/v1/validator/attestation_data", expectedStatus: http.StatusNotImplemented},
		{method: "GET", target: "/eth/v1/validator/aggregate_attestation", expectedStatus: http.StatusNotImplemented},
		{method: "POST", target: "/eth/v1/validator/aggregate_and_proofs", expectedStatus: http.StatusNotImplemented},
		{method: "POST", target: "/eth/v1/validator/beacon_committee_subscriptions", expectedStatus: http.StatusNotImplemented},
		{method: "POST", target: "/eth/v1/validator/sync_committee_subscriptions", expectedStatus: http.StatusNotImplemented},
		{method: "POST", target: "/eth/v1/validator/beacon_committee_selections", expectedStatus: http.StatusNotImplemented},
		{method: "GET", target: "/eth/v1/validator/sync_committee_contribution", expectedStatus: http.StatusNotImplemented},
		{method: "POST", target: "/eth/v1/validator/sync_committee_selections", expectedStatus: http.StatusNotImplemented},
		{method: "POST", target: "/eth/v1/validator/contribution_and_proofs", expectedStatus: http.StatusNotImplemented},
		{method: "POST", target: "/eth/v1/validator/prepare_beacon_proposer", expectedStatus: http.StatusNotImplemented},
		{method: "POST", target: "/eth/v1/validator/register_validator", expectedStatus: http.StatusNotImplemented},
		{method: "POST", target: "/eth/v1/validator/liveness/:epoch", expectedStatus: http.StatusNotImplemented},
	}

	for _, tt := range tests {
		tt.target = remapParams(tt.target)
		req := httptest.NewRequest(tt.method, tt.target, nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, tt.expectedStatus, rec.Code, "Expected status code %d but got %d for path %s", tt.expectedStatus, rec.Code, tt.target)
	}
}

func remapParams(url string) string {
	url = strings.Replace(url, ":epoch", "1", -1)
	url = strings.Replace(url, ":state_id", "head", -1)
	url = strings.Replace(url, ":block_id", "head", -1)
	url = strings.Replace(url, ":slot", "1", -1)
	url = strings.Replace(url, ":peer_id", "QmYyQSo1c1Ym7orWxLYvCrM2EmxFTANf8wXmmE7DWjhx5N", -1)
	url = strings.Replace(url, ":block_root", "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2", -1)
	url = strings.Replace(url, ":validator_id", "1", -1)
	return url

}
