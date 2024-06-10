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

import (
	echo "github.com/labstack/echo/v4"
)

type Handlers interface {
	NotImplemented(c echo.Context) error
	GetGenesis(c echo.Context) error
	GetStateRoot(c echo.Context) error
	GetStateValidators(c echo.Context) error
	PostStateValidators(c echo.Context) error
	GetStateValidatorBalances(c echo.Context) error
	PostStateValidatorBalances(c echo.Context) error
	GetBlockRewards(c echo.Context) error
}

func UseMiddlewares(e *echo.Echo, middlewares ...echo.MiddlewareFunc) {
	for _, middleware := range middlewares {
		e.Use(middleware)
	}
}

func AssignRoutes(e *echo.Echo, handler Handlers) {
	assignBeaconRoutes(e, handler)
	assignBuilderRoutes(e, handler)
	assignConfigRoutes(e, handler)
	assignDebugRoutes(e, handler)
	assignEventsRoutes(e, handler)
	aasignNodeRoutes(e, handler)
	assignValidatorRoutes(e, handler)
	assignRewardsRoutes(e, handler)
}

func assignBeaconRoutes(e *echo.Echo, h Handlers) {
	e.GET("/eth/v1/beacon/genesis",
		h.GetGenesis)
	e.GET("/eth/v1/beacon/states/:state_id/root",
		h.GetStateRoot)
	e.GET("/eth/v1/beacon/states/:state_id/fork",
		h.NotImplemented)
	e.GET("/eth/v1/beacon/states/:state_id/finality_checkpoints",
		h.NotImplemented)
	e.GET("/eth/v1/beacon/states/:state_id/validators",
		h.GetStateValidators)
	e.POST("/eth/v1/beacon/states/:state_id/validators",
		h.PostStateValidators)
	e.GET("/eth/v1/beacon/states/:state_id/validators/:validator_id",
		h.NotImplemented)
	e.GET("/eth/v1/beacon/states/:state_id/validator_balances",
		h.GetStateValidatorBalances)
	e.POST("/eth/v1/beacon/states/:state_id/validator_balances",
		h.PostStateValidatorBalances)
	e.GET("/eth/v1/beacon/states/:state_id/committees",
		h.NotImplemented)
	e.GET("/eth/v1/beacon/states/:state_id/sync_committees",
		h.NotImplemented)
	e.GET("/eth/v1/beacon/states/:state_id/randao",
		h.NotImplemented)
	e.GET("/eth/v1/beacon/headers",
		h.NotImplemented)
	e.GET("/eth/v1/beacon/headers/:block_id",
		h.NotImplemented)
	e.POST("/eth/v1/beacon/blocks/blinded_blocks",
		h.NotImplemented)
	e.POST("/eth/v2/beacon/blocks/blinded_blocks",
		h.NotImplemented)
	e.POST("/eth/v1/beacon/blocks",
		h.NotImplemented)
	e.POST("/eth/v2/beacon/blocks",
		h.NotImplemented)
	e.GET("/eth/v2/beacon/blocks/:block_id",
		h.NotImplemented)
	e.GET("/eth/v1/beacon/blocks/:block_id/root",
		h.NotImplemented)
	e.GET("/eth/v1/beacon/blocks/:block_id/attestations",
		h.NotImplemented)
	e.GET("/eth/v1/beacon/blob_sidecars/:block_id",
		h.NotImplemented)
	e.POST("/eth/v1/beacon/rewards/sync_committee/:block_id",
		h.NotImplemented)
	e.GET("/eth/v1/beacon/deposit_snapshot",
		h.NotImplemented)
	e.POST("/eth/v1/beacon/rewards/attestation/:epoch",
		h.NotImplemented)
	e.GET("/eth/v1/beacon/blinded_blocks/:block_id",
		h.NotImplemented)
	e.GET("/eth/v1/beacon/light_client/bootstrap/:block_root",
		h.NotImplemented)
	e.GET("/eth/v1/beacon/light_client/updates",
		h.NotImplemented)
	e.GET("/eth/v1/beacon/light_client/finality_update",
		h.NotImplemented)
	e.GET("/eth/v1/beacon/light_client/optimistic_update",
		h.NotImplemented)
	e.GET("/eth/v1/beacon/pool/attestations",
		h.NotImplemented)
	e.POST("/eth/v1/beacon/pool/attestations",
		h.NotImplemented)
	e.GET("/eth/v1/beacon/pool/attester_slashings",
		h.NotImplemented)
	e.POST("/eth/v1/beacon/pool/attester_slashings",
		h.NotImplemented)
	e.GET("/eth/v1/beacon/pool/proposer_slashings",
		h.NotImplemented)
	e.POST("/eth/v1/beacon/pool/proposer_slashings",
		h.NotImplemented)
	e.POST("/eth/v1/beacon/pool/sync_committees",
		h.NotImplemented)
	e.GET("/eth/v1/beacon/pool/voluntary_exits",
		h.NotImplemented)
	e.POST("/eth/v1/beacon/pool/voluntary_exits",
		h.NotImplemented)
	e.GET("/eth/v1/beacon/pool/bls_to_execution_changes",
		h.NotImplemented)
	e.POST("/eth/v1/beacon/pool/bls_to_execution_changes",
		h.NotImplemented)
}

func assignBuilderRoutes(e *echo.Echo, h Handlers) {
	e.GET("/eth/v1/builder/states/:state_id/expected_withdrawals",
		h.NotImplemented)
}

func assignConfigRoutes(e *echo.Echo, h Handlers) {
	e.GET("/eth/v1/config/fork_schedule",
		h.NotImplemented)
	e.GET("/eth/v1/config/spec",
		h.NotImplemented)
	e.GET("/eth/v1/config/deposit_contract",
		h.NotImplemented)
}

func assignDebugRoutes(e *echo.Echo, h Handlers) {
	e.GET("/eth/v2/debug/beacon/states/:state_id",
		h.NotImplemented)
	e.GET("/eth/v2/debug/beacon/states/heads",
		h.NotImplemented)
	e.GET("/eth/v1/debug/fork_choice",
		h.NotImplemented)
}

func assignEventsRoutes(e *echo.Echo, h Handlers) {
	e.GET("/eth/v1/events",
		h.NotImplemented)
}

func aasignNodeRoutes(e *echo.Echo, h Handlers) {
	e.GET("/eth/v1/node/identity",
		h.NotImplemented)
	e.GET("/eth/v1/node/peers",
		h.NotImplemented)
	e.GET("/eth/v1/node/peers/:peer_id",
		h.NotImplemented)
	e.GET("/eth/v1/node/peers/peer_count",
		h.NotImplemented)
	e.GET("/eth/v1/node/version",
		h.NotImplemented)
	e.GET("/eth/v1/node/syncing",
		h.NotImplemented)
	e.GET("/eth/v1/node/health",
		h.NotImplemented)
}

func assignValidatorRoutes(e *echo.Echo, h Handlers) {
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

func assignRewardsRoutes(e *echo.Echo, h Handlers) {
	e.POST("/eth/v1/beacon/rewards/sync_committee/:block_id",
		h.NotImplemented)
	e.GET("/eth/v1/beacon/rewards/blocks/:block_id",
		h.GetBlockRewards)
	e.POST("/eth/v1/beacon/rewards/attestations/:epoch",
		h.NotImplemented)
}
