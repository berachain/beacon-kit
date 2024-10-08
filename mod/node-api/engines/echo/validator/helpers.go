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
	"regexp"
	"strconv"

	"github.com/berachain/beacon-kit/mod/node-api/engines/echo/validator/constants"
)

func validateUint64Dec(value string) bool {
	if value == "" {
		return true
	}
	if _, err := strconv.ParseUint(value, 10, 64); err == nil {
		return true
	}
	return false
}

// validateRoot checks if the provided field is a valid root.
// It validates against a 32 byte hex-encoded root with "0x" prefix.
func validateRoot(value string) bool {
	valid, err := validateRegex(value, constants.RootRegex)
	if err != nil {
		return false
	}
	return valid
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

func validateRegex(value string, hexPattern string) (bool, error) {
	if value == "" {
		return true, nil
	}
	matched, err := regexp.MatchString(hexPattern, value)
	if err != nil {
		return false, err
	}
	return matched, nil
}

func validateStateBlockIDs(value string, allowedValues map[string]bool) bool {
	// Check if value is one of the allowed values
	if validateAllowedStrings(value, allowedValues) {
		return true
	}
	// Check if value is a slot (unsigned 64-bit integer)
	if validateUint64Dec(value) {
		return true
	}
	// Check if value is a hex-encoded 32 byte root with "0x" prefix
	if validateRoot(value) {
		return true
	}
	return false
}
