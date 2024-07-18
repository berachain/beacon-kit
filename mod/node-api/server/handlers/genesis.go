package handlers

import (
	"context"
	"net/http"

	types "github.com/berachain/beacon-kit/mod/node-api/server/types"
	echo "github.com/labstack/echo/v4"
)

func (rh RouteHandlers[_]) GetGenesis(c echo.Context) error {
	genesisRoot, err := rh.Backend.GetGenesis(context.TODO())
	if err != nil {
		return err
	}
	if len(genesisRoot) == 0 {
		return echo.NewHTTPError(
			http.StatusNotFound,
			"Chain genesis info is not yet known",
		)
	}
	return c.JSON(http.StatusOK,
		WrapData(types.GenesisData{
			GenesisTime:           "1590832934", // stub
			GenesisValidatorsRoot: genesisRoot,
			GenesisForkVersion:    "0x00000000", // stub
		}))
}
