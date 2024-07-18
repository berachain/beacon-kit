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

package utils

import (
	"errors"
	"net/http"

	"github.com/berachain/beacon-kit/mod/node-api/server/types"
	"github.com/labstack/echo/v4"
)

// HTTPErrorHandler is a custom HTTP error handler for the API.
func HTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	var message any = http.StatusText(code)
	httpError := &echo.HTTPError{}
	if errors.As(err, &httpError) {
		code = httpError.Code
		message = httpError.Message
	}
	c.Logger().Error(err)
	response := &types.ErrorResponse{
		Code:    code,
		Message: message,
	}
	if jsonErr := c.JSON(code, response); jsonErr != nil {
		c.Logger().Error(jsonErr)
	}
}

// BindAndValidate binds the request to the given type and validates the request.
func BindAndValidate[T any](c echo.Context) (*T, error) {
	t := new(T)
	if err := c.Bind(t); err != nil {
		return nil, echo.ErrBadRequest
	}
	if err := c.Validate(t); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return t, nil
}

// WrapData wraps the given data in a DataResponse.
func WrapData(nested any) types.DataResponse {
	return types.DataResponse{Data: nested}
}

// UseMiddlewares adds middlewares to the echo instance.
func UseMiddlewares(e *echo.Echo, middlewares ...echo.MiddlewareFunc) {
	for _, middleware := range middlewares {
		e.Use(middleware)
	}
}
