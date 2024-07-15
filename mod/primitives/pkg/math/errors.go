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

package math

import (
	"math/big"

	"github.com/berachain/beacon-kit/mod/errors"
)

var (
	// ErrUnexpectedInputLengthBase is the base error for unexpected input
	// length errors.
	ErrUnexpectedInputLengthBase = errors.New("unexpected input length")

	// ErrNilBigInt is returned when a nil big.Int is provided to a.
	ErrNilBigInt = errors.New("big.Int is nil")

	// ErrNegativeBigIntBase is returned when a negative big.Int is provided to a
	// function that requires a positive big.Int.
	ErrNegativeBigIntBase = errors.New("big.Int is negative")
)

// ErrUnexpectedInputLength returns an error indicating that the input length.
func ErrUnexpectedInputLength(expected, actual int) error {
	return errors.Wrapf(
		ErrUnexpectedInputLengthBase,
		"expected %d, got %d", expected, actual,
	)
}

// ErrNegativeBigInt returns an error indicating that a negative big.Int was
// provided.
func ErrNegativeBigInt(actual *big.Int) error {
	return errors.Wrapf(
		ErrNegativeBigIntBase, "big.Int is negative: got %s", actual.String())
}
