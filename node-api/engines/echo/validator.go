// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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

	"github.com/berachain/beacon-kit/node-api/handlers/utils"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
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
		"timestamp_id":     ValidateTimestampID,
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

func ValidateStateID(fl validator.FieldLevel) bool {
	allowedValues := map[string]bool{
		"head":      true,
		"genesis":   true,
		"finalized": true,
		"justified": true,
	}
	return validateStateBlockIDs(fl.Field().String(), allowedValues)
}

func ValidateBlockID(fl validator.FieldLevel) bool {
	allowedValues := map[string]bool{
		"head":      true,
		"genesis":   true,
		"finalized": true,
	}
	return validateStateBlockIDs(fl.Field().String(), allowedValues)
}

func ValidateTimestampID(fl validator.FieldLevel) bool {
	allowedValues := map[string]bool{
		utils.StateIDHead:      true,
		utils.StateIDGenesis:   true,
		utils.StateIDFinalized: true,
		utils.StateIDJustified: true,
	}

	value := fl.Field().String()
	if utils.IsTimestampIDPrefix(value) {
		return ValidateUint64Dec(value[1:])
	}

	return validateStateBlockIDs(value, allowedValues)
}

func ValidateUint64Dec(value string) bool {
	if value == "" {
		return true
	}
	_, err := math.U64FromString(value)
	return err == nil
}

func ValidateUint64(fl validator.FieldLevel) bool {
	return ValidateUint64Dec(fl.Field().String())
}

// ValidateValidatorID checks if the provided field is a valid
// validator identifier. It validates against a hex-encoded public key
// or a numeric validator index.
func ValidateValidatorID(fl validator.FieldLevel) bool {
	var key crypto.BLSPubkey
	err := key.UnmarshalText([]byte(fl.Field().String()))
	if err == nil {
		return true
	}
	if ValidateUint64(fl) {
		return true
	}
	return false
}

// ValidateRoot checks if the provided field is a valid root.
// It validates against a 32 byte hex-encoded root with "0x" prefix.
func ValidateRoot(value string) bool {
	_, err := common.NewRootFromHex(value)
	return err == nil
}

func ValidateValidatorStatus(fl validator.FieldLevel) bool {
	// Eth Beacon Node API specs: https://hackmd.io/ofFJ5gOmQpu1jjHilHbdQQ
	allowedStatuses := map[string]bool{
		constants.ValidatorStatusActiveExiting:      true,
		constants.ValidatorStatusActiveOngoing:      true,
		constants.ValidatorStatusActiveSlashed:      true,
		constants.ValidatorStatusExitedSlashed:      true,
		constants.ValidatorStatusExitedUnslashed:    true,
		constants.ValidatorStatusPendingInitialized: true,
		constants.ValidatorStatusPendingQueued:      true,
		constants.ValidatorStatusWithdrawalDone:     true,
		constants.ValidatorStatusWithdrawalPossible: true,
	}
	return validateAllowedStrings(fl.Field().String(), allowedStatuses)
}

func validateAllowedStrings(
	value string,
	allowedValues map[string]bool,
) bool {
	if value == "" {
		return true
	}
	return allowedValues[value]
}

func validateStateBlockIDs(value string, allowedValues map[string]bool) bool {
	// Check if value is one of the allowed values
	if validateAllowedStrings(value, allowedValues) {
		return true
	}
	// Check if value is a slot (unsigned 64-bit integer)
	if ValidateUint64Dec(value) {
		return true
	}
	// Check if value is a hex-encoded 32 byte root with "0x" prefix
	if ValidateRoot(value) {
		return true
	}
	return false
}
