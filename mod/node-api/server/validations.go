package server

import (
	"regexp"
	"strconv"

	"github.com/go-playground/validator/v10"
)

func ConstructValidator() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("state_id", ValidateStateId)
	validate.RegisterValidation("block_id", ValidateBlockId)
	validate.RegisterValidation("validator_id", ValidateValidatorId)
	validate.RegisterValidation("validator_status", ValidateValidatorStatus)
	validate.RegisterValidation("epoch", ValidateUint64)
	validate.RegisterValidation("slot", ValidateUint64)
	validate.RegisterValidation("committee_index", ValidateUint64)
	validate.RegisterValidation("hex", ValidateHex)
	return validate
}

func ValidateStateId(fl validator.FieldLevel) bool {
	allowedValues := map[string]bool{
		"head":      true,
		"genesis":   true,
		"finalized": true,
		"justified": true,
	}
	return validateStateBlockIds(fl, allowedValues)
}

func ValidateBlockId(fl validator.FieldLevel) bool {
	allowedValues := map[string]bool{
		"head":      true,
		"genesis":   true,
		"finalized": true,
	}
	return validateStateBlockIds(fl, allowedValues)
}

// Used to validate slot, validator index,
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

// hex encoded public key (any bytes48 with 0x prefix) or validator index (uint64)
func ValidateValidatorId(fl validator.FieldLevel) bool {
	if validateRegex(fl, `^0x[0-9a-fA-F]{1,96}$`) {
		return true
	}
	if ValidateUint64(fl) {
		return true
	}
	return false
}

func ValidateHex(fl validator.FieldLevel) bool {
	return validateRegex(fl, `^0x[0-9a-fA-F]+$`)
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

func validateAllowedStrings(fl validator.FieldLevel, allowedValues map[string]bool) bool {
	value := fl.Field().String()
	if value == "" {
		return true
	}
	return allowedValues[value]
}

func validateRegex(fl validator.FieldLevel, hexPattern string) bool {
	value := fl.Field().String()
	if value == "" {
		return true
	}
	matched, _ := regexp.MatchString(hexPattern, value)
	return matched
}

func fieldEmpty(fl validator.FieldLevel) bool {
	if value := fl.Field().String(); value == "" {
		return false
	}
	return true
}

func validateStateBlockIds(fl validator.FieldLevel, allowedValues map[string]bool) bool {
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
