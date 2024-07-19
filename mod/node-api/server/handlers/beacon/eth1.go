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

package beacon

import (
	"context"
	"net/http"

	types "github.com/berachain/beacon-kit/mod/node-api/server/types"
	"github.com/berachain/beacon-kit/mod/node-api/server/utils"
	echo "github.com/labstack/echo/v4"
)

func (h *Handler[_]) GetStateRoot(c echo.Context) error {
	params, err := utils.BindAndValidate[types.StateIDRequest](c)
	if err != nil {
		return err
	}
	if params == nil {
		return echo.ErrInternalServerError
	}
	stateRoot, err := h.backend.GetStateRoot(
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
		Data:                utils.WrapData(types.RootData{Root: stateRoot}),
	})
}
