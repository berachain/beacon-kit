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

package chain

import "github.com/berachain/beacon-kit/errors"

var (
	// ErrInsufficientMaxWithdrawalsPerPayload is returned when the max
	// withdrawals per payload less than 2. Must allow at least one for the EVM
	// inflation withdrawal, and at least one more for a validator withdrawal
	// per block.
	ErrInsufficientMaxWithdrawalsPerPayload = errors.New(
		"max withdrawals per payload must be greater than 1")

	// ErrInvalidValidatorSetCap is returned when the validator set cap is
	// greater than the validator registry limit.
	ErrInvalidValidatorSetCap = errors.New(
		"validator set cap must be less than the validator registry limit",
	)

	// ErrInvalidMinEpochsToInactivityPenalty is returned when the minimum epochs
	// to inactivity penalty is zero.
	ErrInvalidMinEpochsToInactivityPenalty = errors.New(
		"minimum epochs to inactivity penalty must be greater than zero",
	)

	// ErrInvalidSlotsPerEpoch is returned when slots per epoch is zero.
	ErrInvalidSlotsPerEpoch = errors.New(
		"slots per epoch must be greater than zero",
	)

	// ErrInvalidFieldElementsPerBlob is returned when field elements per blob is zero.
	ErrInvalidFieldElementsPerBlob = errors.New(
		"field elements per blob must be greater than zero",
	)

	// ErrInvalidBytesPerBlob is returned when bytes per blob is zero.
	ErrInvalidBytesPerBlob = errors.New(
		"bytes per blob must be greater than zero",
	)

	// ErrExcessiveMaxWithdrawalsPerPayload is returned when the max withdrawals per payload
	// is greater than a reasonable limit that could cause resource exhaustion.
	ErrExcessiveMaxWithdrawalsPerPayload = errors.New(
		"max withdrawals per payload must not exceed system limits")
)
