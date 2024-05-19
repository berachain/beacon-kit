package handlers

import (
	"context"
	"net/http"

	"github.com/berachain/beacon-kit/mod/node-api/server/types"
	echo "github.com/labstack/echo/v4"
)

func (rh RouteHandlers) GetBlockProposerDuties(c echo.Context) (err error) {
	params, err := BindAndValidate[types.EpochRequest](c)
	if err != nil {
		return err
	}
	rh.Backend.GetBlockProposerDuties(context.TODO(), params.Epoch)
	return c.String(http.StatusOK, "GetBlockProposerDuties")
}
