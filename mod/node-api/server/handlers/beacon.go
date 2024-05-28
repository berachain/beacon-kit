// SPDX-License-IDentifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package handlers

import (
	"context"
	"net/http"

	types "github.com/berachain/beacon-kit/mod/node-api/server/types"
	echo "github.com/labstack/echo/v4"
)

func (rh RouteHandlers) GetGenesis(c echo.Context) error {
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

func (rh RouteHandlers) GetStateRoot(c echo.Context) error {
	params, err := BindAndValidate[types.StateIDRequest](c)
	if err != nil {
		return err
	}
	if params == nil {
		return echo.ErrInternalServerError
	}
	stateRoot, err := rh.Backend.GetStateRoot(
		context.TODO(),
		params.StateID,
	)
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

func (rh RouteHandlers) GetStateValidators(c echo.Context) error {
	params, err := BindAndValidate[types.StateValidatorsGetRequest](c)
	if err != nil {
		return err
	}
	if params == nil {
		return echo.ErrInternalServerError
	}
	if len(params.Status) > 0 {
		return echo.ErrNotImplemented
	}
	validators, err := rh.Backend.GetStateValidators(
		context.TODO(),
		params.StateID,
		params.ID,
		params.Status,
	)
	if err != nil {
		return err
	}
	if len(validators) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "State not found")
	}
	return c.JSON(http.StatusOK, types.ValidatorResponse{
		ExecutionOptimistic: false, // stubbed
		Finalized:           false, // stubbed
		Data:                validators})
}

func (rh RouteHandlers) PostStateValidators(c echo.Context) error {
	params, err := BindAndValidate[types.StateValidatorsPostRequest](c)
	if err != nil {
		return err
	}
	if params == nil {
		return echo.ErrInternalServerError
	}
	validators, err := rh.Backend.GetStateValidators(
		context.TODO(),
		params.StateID,
		params.IDs,
		params.Statuses,
	)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, types.ValidatorResponse{
		ExecutionOptimistic: false, // stubbed
		Finalized:           false, // stubbed
		Data:                validators})
}

func (rh RouteHandlers) GetStateValidatorBalances(c echo.Context) error {
	params, err := BindAndValidate[types.ValidatorBalancesGetRequest](c)
	if err != nil {
		return err
	}
	if params == nil {
		return echo.ErrInternalServerError
	}
	balances, err := rh.Backend.GetStateValidatorBalances(
		context.TODO(),
		params.StateID,
		params.ID,
	)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, types.ValidatorResponse{
		ExecutionOptimistic: false, // stubbed
		Finalized:           false, // stubbed
		Data:                balances,
	})
}

func (rh RouteHandlers) PostStateValidatorBalances(c echo.Context) error {
	params := &types.ValidatorBalancesPostRequest{}
	if err := (&echo.DefaultBinder{}).BindBody(c, &params.IDs); err != nil {
		return err
	}
	pathParamErr := echo.PathParamsBinder(c).
		String("state_id", &params.StateID).
		BindError()
	if pathParamErr != nil {
		return pathParamErr
	}
	if err := c.Validate(params); err != nil {
		return err
	}
	balances, err := rh.Backend.GetStateValidatorBalances(
		context.TODO(),
		params.StateID,
		params.IDs,
	)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, types.ValidatorResponse{
		ExecutionOptimistic: false, // stubbed
		Finalized:           false, // stubbed
		Data:                balances,
	})
}

func (rh RouteHandlers) GetBlockRewards(c echo.Context) error {
	params, err := BindAndValidate[types.BlockIDRequest](c)
	if err != nil {
		return err
	}
	if params == nil {
		return echo.ErrInternalServerError
	}
	rewards, err := rh.Backend.GetBlockRewards(context.TODO(), params.BlockID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, types.ValidatorResponse{
		ExecutionOptimistic: false, // stubbed
		Finalized:           false, // stubbed
		Data:                rewards,
	})
}
