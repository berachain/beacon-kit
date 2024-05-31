// SPDX-License-Identifier: MIT
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

	"github.com/berachain/beacon-kit/mod/node-api/server/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	echo "github.com/labstack/echo/v4"
)

func (rh RouteHandlers) GetBlockPropserDuties(c echo.Context) error {
	params, err := BindAndValidate[types.EpochRequest](c)
	if err != nil {
		return err
	}
	if params == nil {
		return echo.ErrInternalServerError
	}

	duties, err := rh.Backend.GetBlockPropserDuties(
		context.TODO(),
		params.Epoch,
	)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, types.BlockProposerDutiesResponse{
		DependentRoot:       primitives.Root{0x1}, // stubbed
		ExecutionOptimistic: false,                // stubbed
		Data:                duties,
	})
}
