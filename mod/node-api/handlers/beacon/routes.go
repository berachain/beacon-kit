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

package beacon

import (
	"net/http"

	"github.com/berachain/beacon-kit/mod/node-api/handlers"
)

//nolint:funlen // routes are long
func (h *Handler[_, ContextT, _, _]) RegisterRoutes() {
	h.routes.Routes = []handlers.Route[ContextT]{
		{
			Method:  http.MethodGet,
			Path:    "/eth/v1/beacon/genesis",
			Handler: h.GetGenesis,
		},
		{
			Method:  http.MethodGet,
			Path:    "/eth/v1/beacon/states/:state_id/root",
			Handler: h.GetStateRoot,
		},
		{
			Method:  http.MethodGet,
			Path:    "/eth/v1/beacon/states/:state_id/fork",
			Handler: h.GetStateFork,
		},
		{
			Method:  http.MethodGet,
			Path:    "/eth/v1/beacon/states/:state_id/finality_checkpoints",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodGet,
			Path:    "/eth/v1/beacon/states/:state_id/validators",
			Handler: h.GetStateValidators,
		},
		{
			Method:  http.MethodPost,
			Path:    "/eth/v1/beacon/states/:state_id/validators",
			Handler: h.PostStateValidators,
		},
		{
			Method:  http.MethodGet,
			Path:    "/eth/v1/beacon/states/:state_id/validators/:validator_id",
			Handler: h.GetStateValidator,
		},
		{
			Method:  http.MethodGet,
			Path:    "/eth/v1/beacon/states/:state_id/validator_balances",
			Handler: h.GetStateValidatorBalances,
		},
		{
			Method:  http.MethodPost,
			Path:    "/eth/v1/beacon/states/:state_id/validator_balances",
			Handler: h.PostStateValidatorBalances,
		},
		{
			Method:  http.MethodGet,
			Path:    "/eth/v1/beacon/states/:state_id/committees",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodGet,
			Path:    "/eth/v1/beacon/states/:state_id/sync_committees",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodGet,
			Path:    "/eth/v1/beacon/states/:state_id/randao",
			Handler: h.GetRandao,
		},
		{
			Method:  http.MethodGet,
			Path:    "/eth/v1/beacon/headers",
			Handler: h.GetBlockHeaders,
		},
		{
			Method:  http.MethodGet,
			Path:    "/eth/v1/beacon/headers/:block_id",
			Handler: h.GetBlockHeaderByID,
		},
		{
			Method:  http.MethodPost,
			Path:    "/eth/v1/beacon/blocks/blinded_blocks",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodPost,
			Path:    "eth/v2/beacon/blocks/blinded_blocks",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodPost,
			Path:    "/eth/v1/beacon/blocks",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodPost,
			Path:    "eth/v2/beacon/blocks",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodGet,
			Path:    "eth/v2/beacon/blocks/:block_id",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodGet,
			Path:    "/eth/v1/beacon/blocks/:block_id/root",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodGet,
			Path:    "/eth/v1/beacon/blocks/:block_id/attestations",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodGet,
			Path:    "/eth/v1/beacon/blob_sidecars/:block_id",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodPost,
			Path:    "/eth/v1/beacon/rewards/sync_committee/:block_id",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodGet,
			Path:    "/eth/v1/beacon/deposit_snapshot",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodPost,
			Path:    "/eth/v1/beacon/rewards/attestation/:epoch",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodGet,
			Path:    "/eth/v1/beacon/blinded_blocks/:block_id",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodGet,
			Path:    "/eth/v1/beacon/light_client/bootstrap/:block_root",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodGet,
			Path:    "/eth/v1/beacon/light_client/updates",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodGet,
			Path:    "/eth/v1/beacon/light_client/finality_update",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodGet,
			Path:    "/eth/v1/beacon/light_client/optimistic_update",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodGet,
			Path:    "/eth/v1/beacon/pool/attestations",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodPost,
			Path:    "/eth/v1/beacon/pool/attestations",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodGet,
			Path:    "/eth/v1/beacon/pool/attester_slashings",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodPost,
			Path:    "/eth/v1/beacon/pool/attester_slashings",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodGet,
			Path:    "/eth/v1/beacon/pool/proposer_slashings",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodPost,
			Path:    "/eth/v1/beacon/pool/proposer_slashings",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodPost,
			Path:    "/eth/v1/beacon/pool/sync_committees",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodGet,
			Path:    "/eth/v1/beacon/pool/voluntary_exits",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodPost,
			Path:    "/eth/v1/beacon/pool/voluntary_exits",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodGet,
			Path:    "/eth/v1/beacon/pool/bls_to_execution_changes",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodPost,
			Path:    "/eth/v1/beacon/pool/bls_to_execution_changes",
			Handler: h.NotImplemented,
		},
	}
}
