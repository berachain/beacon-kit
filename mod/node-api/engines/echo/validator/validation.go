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
	"github.com/berachain/beacon-kit/mod/node-api/engines/echo/validator/constants"
	"github.com/berachain/beacon-kit/mod/node-api/handlers/utils"
	"github.com/go-playground/validator/v10"
)

func ValidateStateID(fl validator.FieldLevel) bool {
	allowedValues := map[string]bool{
		constants.StateIDHead:      true,
		constants.StateIDGenesis:   true,
		constants.StateIDFinalized: true,
		constants.StateIDJustified: true,
	}
	return validateStateBlockIDs(fl.Field().String(), allowedValues)
}

func ValidateBlockID(fl validator.FieldLevel) bool {
	allowedValues := map[string]bool{
		constants.StateIDHead:      true,
		constants.StateIDGenesis:   true,
		constants.StateIDFinalized: true,
	}
	return validateStateBlockIDs(fl.Field().String(), allowedValues)
}

func ValidateExecutionID(fl validator.FieldLevel) bool {
	allowedValues := map[string]bool{
		constants.StateIDHead:      true,
		constants.StateIDGenesis:   true,
		constants.StateIDFinalized: true,
		constants.StateIDJustified: true,
	}

	value := fl.Field().String()
	if utils.IsExecutionNumberPrefix(value) {
		return validateUint64Dec(value[1:])
	}

	return validateStateBlockIDs(value, allowedValues)
}

// ValidateValidatorID checks if the provided field is a valid
// validator identifier. It validates against a hex-encoded public key
// or a numeric validator index.
func ValidateValidatorID(fl validator.FieldLevel) bool {
	valid, err := validateRegex(fl.Field().String(), constants.ValidatorIDRegex)
	if err != nil {
		return false
	}
	if valid {
		return true
	}
	if validateUint64(fl) {
		return true
	}
	return false
}

func ValidateValidatorStatus(fl validator.FieldLevel) bool {
	// Eth Beacon Node API specs: https://hackmd.io/ofFJ5gOmQpu1jjHilHbdQQ
	allowedStatuses := map[string]bool{
		constants.StatusPendingInitialized: true,
		constants.StatusPendingQueued:      true,
		constants.StatusActiveOngoing:      true,
		constants.StatusActiveExiting:      true,
		constants.StatusActiveSlashed:      true,
		constants.StatusExitedUnslashed:    true,
		constants.StatusExitedSlashed:      true,
		constants.StatusWithdrawalPossible: true,
		constants.StatusWithdrawalDone:     true,
	}
	return validateAllowedStrings(fl.Field().String(), allowedStatuses)
}
