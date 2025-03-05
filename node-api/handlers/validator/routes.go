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
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package validator

import (
	"net/http"

	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/node-api/handlers"
)

func (h *Handler) RegisterRoutes(logger log.Logger) {
	h.SetLogger(logger)
	h.BaseHandler.AddRoutes([]*handlers.Route{
		{
			Method:  http.MethodPost,
			Path:    "/eth/v1/validator/duties/attester/:epoch",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodGet,
			Path:    "/eth/v1/validator/duties/proposer/:epoch",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodPost,
			Path:    "/eth/v1/validator/duties/sync/:epoch",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodGet,
			Path:    "/eth/v3/validator/blocks/:slot",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodPost,
			Path:    "/eth/v1/validator/attestation_data",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodGet,
			Path:    "/eth/v1/validator/aggregate_attestation",
			Handler: h.Deprecated,
		},
		{
			Method:  http.MethodGet,
			Path:    "/eth/v2/validator/aggregate_attestation",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodPost,
			Path:    "/eth/v1/validator/aggregate_and_proofs",
			Handler: h.Deprecated,
		},
		{
			Method:  http.MethodPost,
			Path:    "/eth/v2/validator/aggregate_and_proofs",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodPost,
			Path:    "/eth/v1/validator/beacon_committee_subscriptions",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodPost,
			Path:    "/eth/v1/validator/sync_committee_subscriptions",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodPost,
			Path:    "/eth/v1/validator/beacon_committee_selections",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodGet,
			Path:    "/eth/v1/validator/sync_committee_contribution",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodPost,
			Path:    "/eth/v1/validator/contribution_and_proofs",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodPost,
			Path:    "/eth/v1/validator/prepare_beacon_proposer",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodPost,
			Path:    "/eth/v1/validator/register_validator",
			Handler: h.NotImplemented,
		},
		{
			Method:  http.MethodPost,
			Path:    "/eth/v1/validator/liveness/:epoch",
			Handler: h.NotImplemented,
		},
	})
}
