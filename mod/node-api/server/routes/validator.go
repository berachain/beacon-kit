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
