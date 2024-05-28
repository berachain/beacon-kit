package handlers

import (
	"context"
	"net/http"

	types "github.com/berachain/beacon-kit/mod/node-api/server/types"
	echo "github.com/labstack/echo/v4"
)

func (rh RouteHandlers) GetGenesis(c echo.Context) (err error) {
	genesisRoot, err := rh.Backend.GetGenesis(context.TODO())
	if err != nil {
		return err
	}
	if len(genesisRoot) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "Chain genesis info is not yet known")
	}
	return c.JSON(http.StatusOK,
		WrapData(types.GenesisData{
			GenesisTime:           "1590832934", //stub
			GenesisValidatorsRoot: genesisRoot,
			GenesisForkVersion:    "0x00000000", //stub
		}))
}

func (rh RouteHandlers) GetStateRoot(c echo.Context) (err error) {
	params, err := BindAndValidate[types.StateIdRequest](c)
	if err != nil {
		return err
	}
	stateRoot, err := rh.Backend.GetStateRoot(c.(context.Context), params.StateId)
	if err != nil {
		return err
	}
	if len(stateRoot) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "State not found")
	}
	return c.JSON(http.StatusOK, types.ValidatorResponse{
		ExecutionOptimistic: false, // stubbed
		Finalized:           false, // stubbed
		Data:                WrapData(types.RootData{Root: stateRoot}),
	})
}

func (rh RouteHandlers) GetStateValidators(c echo.Context) (err error) {
	params, err := BindAndValidate[types.StateValidatorsGetRequest](c)
	if err != nil {
		return err
	}
	if len(params.Status) > 0 {
		return echo.ErrNotImplemented
	}
	validators, err := rh.Backend.GetStateValidators(context.TODO(), params.StateId, params.Id, params.Status)
	if len(validators) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "State not found")

	}
	return c.JSON(http.StatusOK, types.ValidatorResponse{
		ExecutionOptimistic: false, // stubbed
		Finalized:           false, // stubbed
		Data:                validators})
}

func (rh RouteHandlers) PostStateValidators(c echo.Context) (err error) {
	params, err := BindAndValidate[types.StateValidatorsPostRequest](c)
	if err != nil {
		return err
	}
	validators, err := rh.Backend.GetStateValidators(context.TODO(), params.StateId, params.Ids, params.Statuses)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, types.ValidatorResponse{
		ExecutionOptimistic: false, // stubbed
		Finalized:           false, // stubbed
		Data:                validators})
}

func (rh RouteHandlers) GetStateValidatorBalances(c echo.Context) (err error) {
	params, err := BindAndValidate[types.ValidatorBalancesGetRequest](c)
	if err != nil {
		return err
	}
	balances, err := rh.Backend.GetStateValidatorBalances(context.TODO(), params.StateId, params.Id)
	return c.JSON(http.StatusOK, types.ValidatorResponse{
		ExecutionOptimistic: false, // stubbed
		Finalized:           false, // stubbed
		Data:                balances,
	})
}

func (rh RouteHandlers) PostStateValidatorBalances(c echo.Context) (err error) {
	params := &types.ValidatorBalancesPostRequest{}
	if err := (&echo.DefaultBinder{}).BindBody(c, &params.Ids); err != nil {
		return err
	}
	if err := echo.PathParamsBinder(c).String("state_id", &params.StateId).BindError(); err != nil {
		return err
	}
	if err := c.Validate(params); err != nil {
		return err
	}
	balances, err := rh.Backend.GetStateValidatorBalances(context.TODO(), params.StateId, params.Ids)
	return c.JSON(http.StatusOK, types.ValidatorResponse{
		ExecutionOptimistic: false, // stubbed
		Finalized:           false, // stubbed
		Data:                balances,
	})
}

func (rh RouteHandlers) GetBlockRewards(c echo.Context) (err error) {
	params, err := BindAndValidate[types.BlockIdRequest](c)
	if err != nil {
		return err
	}
	rewards, err := rh.Backend.GetBlockRewards(context.TODO(), params.BlockId)
	return c.JSON(http.StatusOK, types.ValidatorResponse{
		ExecutionOptimistic: false, // stubbed
		Finalized:           false, // stubbed
		Data:                rewards,
	})
}
