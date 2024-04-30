// SPDX-License-Identifier: MIT
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

package math

import (
	"math/big"

	"github.com/cockroachdb/errors"
)

var (
	// ErrUnexpectedInputLengthBase is the base error for unexpected input
	// length errors.
	ErrUnexpectedInputLengthBase = errors.New("unexpected input length")

	// ErrNilBigInt is returned when a nil big.Int is provided to a.
	ErrNilBigInt = errors.New("big.Int is nil")

	// ErrNegativeBigInt is returned when a negative big.Int is provided to a
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
