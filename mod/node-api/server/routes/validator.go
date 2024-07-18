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

package routes

import (
	"github.com/berachain/beacon-kit/mod/node-api/server/handlers"
	"github.com/labstack/echo/v4"
)

func assignValidatorRoutes[ValidatorT any](
	e *echo.Echo,
	h handlers.RouteHandlers[ValidatorT],
) {
	e.POST("/eth/v1/validator/duties/attester/:epoch",
		h.NotImplemented)
	e.GET("/eth/v1/validator/duties/proposer/:epoch",
		h.NotImplemented)
	e.POST("/eth/v1/validator/duties/sync/:epoch",
		h.NotImplemented)
	e.GET("/eth/v3/validator/blocks/:slot",
		h.NotImplemented)
	e.GET("/eth/v1/validator/attestation_data",
		h.NotImplemented)
	e.GET("/eth/v1/validator/aggregate_attestation",
		h.NotImplemented)
	e.POST("/eth/v1/validator/aggregate_and_proofs",
		h.NotImplemented)
	e.POST("/eth/v1/validator/beacon_committee_subscriptions",
		h.NotImplemented)
	e.POST("/eth/v1/validator/sync_committee_subscriptions",
		h.NotImplemented)
	e.POST("/eth/v1/validator/beacon_committee_selections",
		h.NotImplemented)
	e.GET("/eth/v1/validator/sync_committee_contribution",
		h.NotImplemented)
	e.POST("/eth/v1/validator/sync_committee_selections",
		h.NotImplemented)
	e.POST("/eth/v1/validator/contribution_and_proofs",
		h.NotImplemented)
	e.POST("/eth/v1/validator/prepare_beacon_proposer",
		h.NotImplemented)
	e.POST("/eth/v1/validator/register_validator",
		h.NotImplemented)
	e.POST("/eth/v1/validator/liveness/:epoch",
		h.NotImplemented)
}
