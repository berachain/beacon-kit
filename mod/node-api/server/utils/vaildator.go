package utils

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

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
