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

package signer

import "errors"

// ErrInvalidSignature is returned when a signature is invalid.
var (
	// ErrInvalidSignature is returned when a signature is invalid.
	ErrInvalidSignature = errors.New("invalid BLS signature")

	// ErrValidatorPrivateKeyRequired is returned when the validator private key
	// is required but not provided.
	ErrValidatorPrivateKeyRequired = errors.New(
		"validator private key required",
	)
	// ErrInvalidValidatorPrivateKeyLength is returned when the validator
	// private key has an invalid length.
	ErrInvalidValidatorPrivateKeyLength = errors.New(
		"invalid validator private key length",
	)
)
