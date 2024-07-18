package beacon

import (
	"context"
	"net/http"

	types "github.com/berachain/beacon-kit/mod/node-api/server/types"
	"github.com/berachain/beacon-kit/mod/node-api/server/utils"
	echo "github.com/labstack/echo/v4"
)

func (h Handler[_]) GetGenesis(c echo.Context) error {
	genesisRoot, err := h.backend.GetGenesis(context.TODO())
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
		utils.WrapData(types.GenesisData{
			GenesisTime:           "1590832934", // stub
			GenesisValidatorsRoot: genesisRoot,
			GenesisForkVersion:    "0x00000000", // stub
		}))
}
