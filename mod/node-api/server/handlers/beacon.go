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
