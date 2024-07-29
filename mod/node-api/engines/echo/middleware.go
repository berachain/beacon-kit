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

package echo

import (
	"net/http"

	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/node-api/handlers"
	"github.com/berachain/beacon-kit/mod/node-api/handlers/types"
	"github.com/labstack/echo/v4"
)

// ErrorResponse is a response that is returned when an error occurs.
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// responseMiddleware is a middleware that converts errors to an HTTP status 
// code and response.
func responseMiddleware(
	handler *handlers.Route[Context],
) echo.HandlerFunc {
	return func(c Context) error {
		data, err := handler.Handler(c)
		code, response := responseFromError(data, err)
		return c.JSON(code, response)
	}
}

// responseFromErr converts an error to an HTTP status code and response. If
// the error is nil, the response is returned as is.
func responseFromError(data any, err error) (int, any) {
	switch {
	case err == nil:
		return http.StatusOK, data
	case errors.Is(err, types.ErrNotFound):
		return http.StatusNotFound, ErrorResponse{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}
	case errors.Is(err, types.ErrInvalidRequest):
		return http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	case errors.Is(err, types.ErrNotImplemented):
		return http.StatusNotImplemented, ErrorResponse{
			Code:    http.StatusNotImplemented,
			Message: err.Error(),
		}
	default:
		return http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
}
