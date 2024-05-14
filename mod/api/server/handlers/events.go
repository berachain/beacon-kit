package handlers

import (
	"net/http"
	"github.com/berachain/beacon-kit/mod/api/server/types"
	echo "github.com/labstack/echo/v4"
)

func (rh RouteHandlers) GetEvents(c echo.Context) (err error) {
	params, err := BindAndValidate[types.EventsRequest](c)
	if err != nil {
		return err
	}
	rh.Backend.SubscribeEvents(params.Topics)
	return c.String(http.StatusOK, "GetEvents")
}
