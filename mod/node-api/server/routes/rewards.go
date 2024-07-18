package routes

import (
	"github.com/berachain/beacon-kit/mod/node-api/server/handlers"
	"github.com/labstack/echo/v4"
)

func assignRewardsRoutes[ValidatorT any](
	e *echo.Echo,
	h handlers.RouteHandlers[ValidatorT],
) {
	e.POST("/eth/v1/beacon/rewards/sync_committee/:block_id",
		h.NotImplemented)
	e.GET("/eth/v1/beacon/rewards/blocks/:block_id",
		h.GetBlockRewards)
	e.POST("/eth/v1/beacon/rewards/attestations/:epoch",
		h.NotImplemented)
}
