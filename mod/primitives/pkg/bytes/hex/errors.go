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

package hex

import (
	"encoding/json"
	"errors"
	"reflect"
)

var (
	ErrEmptyString     = errors.New("empty hex string")
	ErrMissingPrefix   = errors.New("hex string without 0x prefix")
	ErrOddLength       = errors.New("hex string of odd length")
	ErrNonQuotedString = errors.New("non-quoted hex string")
	ErrInvalidString   = errors.New("invalid hex string")

	ErrLeadingZero = errors.New("hex number with leading zero digits")
	ErrEmptyNumber = errors.New("hex string \"0x\"")
	ErrUint64Range = errors.New("hex number > 64 bits")
	ErrBig256Range = errors.New("hex number > 256 bits")

	ErrInvalidBigWordSize = errors.New("weird big.Word size")
)

// WrapUnmarshalError wraps an error occurring during JSON unmarshaling.
func WrapUnmarshalError(err error, t reflect.Type) error {
	if err != nil {
		err = &json.UnmarshalTypeError{Value: err.Error(), Type: t}
	}

	return err
}
