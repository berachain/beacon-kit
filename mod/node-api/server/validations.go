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

package server

import (
	"regexp"
	"strconv"

	"github.com/go-playground/validator/v10"
)

func ConstructValidator() *validator.Validate {
	validators := map[string](func(fl validator.FieldLevel) bool){
		"state_id":         ValidateStateID,
		"block_id":         ValidateBlockID,
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

// hex encoded public key (any bytes48 with 0x prefix) or validator index.
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

func fieldEmpty(fl validator.FieldLevel) bool {
	if value := fl.Field().String(); value == "" {
		return false
	}
	return true
}

func validateStateBlockIDs(
	fl validator.FieldLevel,
	allowedValues map[string]bool,
) bool {
	if fieldEmpty(fl) {
		return true
	}
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
