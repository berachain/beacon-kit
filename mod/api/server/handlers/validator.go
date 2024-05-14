package handlers

import (
	"net/http"
	"github.com/berachain/beacon-kit/mod/api/server/types"
	echo "github.com/labstack/echo/v4"
)

func (rh RouteHandlers) GetBlockProposerDuties(c echo.Context) (err error) {
	params, err := BindAndValidate[types.EpochRequest](c)
	if err != nil {
		return err
	}
	rh.Backend.GetBlockProposerDuties(params.Epoch)
	return c.String(http.StatusOK, "GetBlockProposerDuties")
}
