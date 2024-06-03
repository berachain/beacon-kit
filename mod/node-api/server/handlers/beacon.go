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

	response "github.com/berachain/beacon-kit/mod/node-api/server/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	echo "github.com/labstack/echo/v4"
)

func (rh RouteHandlers) GetGenesis(c echo.Context) error {
	genesis, err := rh.Backend.GetGenesis(context.Background())
	if err != nil {
		return err
	}
	if len(genesis.GenesisValidatorsRoot) == 0 {
		return echo.NewHTTPError(
			http.StatusNotFound,
			"Chain genesis info is not yet known",
		)
	}
	return c.JSON(http.StatusOK,
		WrapData(response.GenesisData{
			GenesisTime:           genesis.GenesisTime, // stub
			GenesisValidatorsRoot: genesis.GenesisValidatorsRoot,
			GenesisForkVersion:    genesis.GenesisForkVersion, // stub
		}))
}

func (rh RouteHandlers) GetStateRoot(c echo.Context) error {
	params, err := BindAndValidate[response.StateIDRequest](c)
	if err != nil {
		return err
	}
	if params == nil {
		return echo.ErrInternalServerError
	}
	stateRoot, err := rh.Backend.GetStateRoot(
		context.Background(),
		params.StateID,
	)
	if err != nil {
		return err
	}
	if len(stateRoot) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "State not found")
	}
	return c.JSON(http.StatusOK, response.ValidatorResponse{
		ExecutionOptimistic: false, // stubbed
		Finalized:           false, // stubbed
		Data:                WrapData(response.RootData{Root: stateRoot}),
	})
}

func (rh RouteHandlers) GetStateValidators(c echo.Context) error {
	params, err := BindAndValidate[response.StateValidatorsGetRequest](c)
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
		context.Background(),
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
	return c.JSON(http.StatusOK, response.ValidatorResponse{
		ExecutionOptimistic: false, // stubbed
		Finalized:           false, // stubbed
		Data:                validators})
}

func (rh RouteHandlers) PostStateValidators(c echo.Context) error {
	params, err := BindAndValidate[response.StateValidatorsPostRequest](c)
	if err != nil {
		return err
	}
	if params == nil {
		return echo.ErrInternalServerError
	}
	validators, err := rh.Backend.GetStateValidators(
		context.Background(),
		params.StateID,
		params.IDs,
		params.Statuses,
	)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, response.ValidatorResponse{
		ExecutionOptimistic: false, // stubbed
		Finalized:           false, // stubbed
		Data:                validators})
}

func (rh RouteHandlers) GetStateValidatorBalances(c echo.Context) error {
	params, err := BindAndValidate[response.ValidatorBalancesGetRequest](c)
	if err != nil {
		return err
	}
	if params == nil {
		return echo.ErrInternalServerError
	}
	balances, err := rh.Backend.GetStateValidatorBalances(
		context.Background(),
		params.StateID,
		params.ID,
	)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, response.ValidatorResponse{
		ExecutionOptimistic: false, // stubbed
		Finalized:           false, // stubbed
		Data:                balances,
	})
}

func (rh RouteHandlers) PostStateValidatorBalances(c echo.Context) error {
	params := &response.ValidatorBalancesPostRequest{}
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
		context.Background(),
		params.StateID,
		params.IDs,
	)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, response.ValidatorResponse{
		ExecutionOptimistic: false, // stubbed
		Finalized:           false, // stubbed
		Data:                balances,
	})
}

func (rh RouteHandlers) GetStateCommittees(c echo.Context) error {
	params, err := BindAndValidate[response.StateCommitteesRequest](c)
	if err != nil {
		return err
	}
	if params == nil {
		return echo.ErrInternalServerError
	}
	committees, err := rh.Backend.GetStateCommittees(
		context.Background(),
		params.StateID,
		params.ComitteeIndex,
		params.Epoch,
		params.Slot,
	)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, response.ValidatorResponse{
		ExecutionOptimistic: false, // stubbed
		Finalized:           false, // stubbed
		Data:                committees,
	})
}

func (rh RouteHandlers) GetStateSyncCommittees(c echo.Context) error {
	params, err := BindAndValidate[response.StateSyncCommitteesRequest](c)
	if err != nil {
		return err
	}
	if params == nil {
		return echo.ErrInternalServerError
	}
	committees, err := rh.Backend.GetStateSyncCommittees(
		context.Background(),
		params.StateID,
		params.Epoch,
	)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, response.ValidatorResponse{
		ExecutionOptimistic: false, // stubbed
		Finalized:           false, // stubbed
		Data:                committees,
	})
}

func (rh RouteHandlers) GetBlockHeaders(c echo.Context) error {
	params, err := BindAndValidate[response.BlockHeadersRequest](c)
	if err != nil {
		return err
	}
	if params == nil {
		return echo.ErrInternalServerError
	}
	parentRoot := primitives.Root{}
	if params.ParentRoot != "" {
		err = parentRoot.UnmarshalText([]byte(params.ParentRoot))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest,
				"Invalid parent_root: "+params.ParentRoot)
		}
	}
	headers, err := rh.Backend.GetBlockHeaders(
		context.Background(),
		params.Slot,
		parentRoot,
	)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, response.ValidatorResponse{
		ExecutionOptimistic: false, // stubbed
		Finalized:           false, // stubbed
		Data:                headers,
	})
}

func (rh RouteHandlers) GetBlockHeader(c echo.Context) error {
	params, err := BindAndValidate[response.BlockIDRequest](c)
	if err != nil {
		return err
	}
	if params == nil {
		return echo.ErrInternalServerError
	}
	header, err := rh.Backend.GetBlockHeader(
		context.Background(),
		params.BlockID,
	)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, response.ValidatorResponse{
		ExecutionOptimistic: false, // stubbed
		Finalized:           false, // stubbed
		Data:                header,
	})
}

func (rh RouteHandlers) GetBlock(c echo.Context) error {
	params, err := BindAndValidate[response.BlockIDRequest](c)
	if err != nil {
		return err
	}
	if params == nil {
		return echo.ErrInternalServerError
	}
	block, err := rh.Backend.GetBlock(
		context.Background(),
		params.BlockID,
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response.BlockResponse{
		Version:             "deneb", // stubbed
		ExecutionOptimistic: false,   // stubbed
		Finalized:           false,   // stubbed
		Data: &response.MessageSignature{
			Message:   block,
			Signature: crypto.BLSSignature{},
		},
	})
}

func (rh RouteHandlers) GetBlockBlobSidecars(c echo.Context) error {
	params, err := BindAndValidate[response.BlobSidecarRequest](c)
	if err != nil {
		return err
	}
	if params == nil {
		return echo.ErrInternalServerError
	}
	sidecars, err := rh.Backend.GetBlockBlobSidecars(
		context.Background(),
		params.BlockID,
		params.Indices,
	)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, response.ValidatorResponse{
		ExecutionOptimistic: false, // stubbed
		Finalized:           false, // stubbed
		Data:                sidecars,
	})
}

func (rh RouteHandlers) GetPoolVoluntaryExits(c echo.Context) error {
	exits, err := rh.Backend.GetPoolVoluntaryExits(
		context.Background(),
	)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, WrapData(response.ValidatorResponse{
		ExecutionOptimistic: false, // stubbed
		Finalized:           false, // stubbed
		Data:                exits,
	}))
}

func (rh RouteHandlers) GetPoolBtsToExecutionChanges(c echo.Context) error {
	changes, err := rh.Backend.GetPoolBtsToExecutionChanges(
		context.Background(),
	)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, WrapData(changes))
}
