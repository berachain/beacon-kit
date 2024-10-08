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

package validator

import (
	"fmt"
	"net/http"

	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// TODO: these validators need to be un-janked to 1) not use `FieldLevel` for
// repeated `.Field().String()` calls and 2) strongly type the allowed IDs,
// putting validation logic on each type.

// CustomValidator is a custom validator for the API.
type CustomValidator struct {
	Validator *validator.Validate
}

// Validate validates the given interface.
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		var validationErrors validator.ValidationErrors
		hasValidationErrors := errors.As(err, &validationErrors)
		if !hasValidationErrors || len(validationErrors) == 0 {
			return nil
		}
		firstError := validationErrors[0]
		field := firstError.Field()
		value := firstError.Value()
		return echo.NewHTTPError(http.StatusBadRequest,
			fmt.Sprintf("Invalid %s: %s", field, value))
	}
	return nil
}

// ConstructValidator constructs a new validator.
func ConstructValidator() *validator.Validate {
	validators := map[string](func(fl validator.FieldLevel) bool){
		"state_id":         ValidateStateID,
		"block_id":         ValidateBlockID,
		"execution_id":     ValidateExecutionID,
		"validator_id":     ValidateValidatorID,
		"epoch":            ValidateUint64,
		"slot":             ValidateUint64,
		"validator_status": ValidateValidatorStatus,
	}
	validate := validator.New()
	for tag, fn := range validators {
		err := validate.RegisterValidation(tag, fn)
		if err != nil {
			panic(err)
		}
	}
	return validate
}
