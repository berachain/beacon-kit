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
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/berachain/beacon-kit/mod/node-api/handlers/utils"
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

func ConstructValidator() *validator.Validate {
	validators := map[string](func(fl validator.FieldLevel) bool){
		"state_id":         ValidateStateID,
		"block_id":         ValidateBlockID,
		"execution_id":     ValidateExecutionID,
		"validator_id":     ValidateValidatorID,
		"validator_status": ValidateValidatorStatus,
		"epoch":            ValidateUint64,
		"slot":             ValidateUint64,
		"committee_index":  ValidateUint64,
		"hex":              ValidateHex,
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

func ValidateStateID(fl validator.FieldLevel) bool {
	allowedValues := map[string]bool{
		"head":      true,
		"genesis":   true,
		"finalized": true,
		"justified": true,
	}
	return validateStateBlockIDs(fl, allowedValues)
}

func ValidateBlockID(fl validator.FieldLevel) bool {
	allowedValues := map[string]bool{
		"head":      true,
		"genesis":   true,
		"finalized": true,
	}
	return validateStateBlockIDs(fl, allowedValues)
}

func ValidateExecutionID(fl validator.FieldLevel) bool {
	allowedValues := map[string]bool{
		utils.StateIDHead:      true,
		utils.StateIDGenesis:   true,
		utils.StateIDFinalized: true,
		utils.StateIDJustified: true,
	}

	if utils.IsExecutionNumberPrefix(fl.Field().String()) {
		return true
	}

	return validateStateBlockIDs(fl, allowedValues)
}

func ValidateUint64(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true
	}
	if _, err := strconv.ParseUint(value, 10, 64); err == nil {
		return true
	}
	return false
}

// ValidateValidatorID checks if the provided field is a valid
// validator identifier. It validates against a hex-encoded public key
// or a numeric validator index.
func ValidateValidatorID(fl validator.FieldLevel) bool {
	valid, err := validateRegex(fl, `^0x[0-9a-fA-F]{1,96}$`)
	if err != nil {
		return false
	}
	if valid {
		return true
	}
	if ValidateUint64(fl) {
		return true
	}
	return false
}

func ValidateHex(fl validator.FieldLevel) bool {
	valid, err := validateRegex(fl, `^0x[0-9a-fA-F]+$`)
	if err != nil {
		return false
	}
	return valid
}

func ValidateValidatorStatus(fl validator.FieldLevel) bool {
	// Eth Beacon Node API specs: https://hackmd.io/ofFJ5gOmQpu1jjHilHbdQQ
	allowedStatuses := map[string]bool{
		"pending_initialized": true,
		"pending_queued":      true,
		"active_ongoing":      true,
		"active_exiting":      true,
		"active_slashed":      true,
		"exited_unslashed":    true,
		"exited_slashed":      true,
		"withdrawal_possible": true,
		"withdrawal_done":     true,
	}
	return validateAllowedStrings(fl, allowedStatuses)
}

func validateAllowedStrings(
	fl validator.FieldLevel,
	allowedValues map[string]bool,
) bool {
	value := fl.Field().String()
	if value == "" {
		return true
	}
	return allowedValues[value]
}

func validateRegex(fl validator.FieldLevel, hexPattern string) (
	bool, error) {
	value := fl.Field().String()
	if value == "" {
		return true, nil
	}
	matched, err := regexp.MatchString(hexPattern, value)
	if err != nil {
		return false, err
	}
	return matched, nil
}

func validateStateBlockIDs(
	fl validator.FieldLevel,
	allowedValues map[string]bool,
) bool {
	// Check if value is one of the allowed values
	if validateAllowedStrings(fl, allowedValues) {
		return true
	}
	// Check if value is a slot (unsigned 64-bit integer)
	if ValidateUint64(fl) {
		return true
	}
	// Check if value is a hex-encoded state root with "0x" prefix
	if ValidateHex(fl) {
		return true
	}
	return false
}
