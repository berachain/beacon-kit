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

package deposit

import "errors"

var (
	// ErrInvalidPubKeyLength is returned when the public key is invalid.
	ErrInvalidPubKeyLength = errors.New(
		"invalid public key length",
	)

	// ErrInvalidWithdrawalCredentialsLength is returned when the withdrawal
	// credentials are invalid.
	ErrInvalidWithdrawalCredentialsLength = errors.New(
		"invalid withdrawal credentials length",
	)

	// ErrInvalidAmount is returned when the deposit amount is invalid.
	ErrInvalidAmount = errors.New(
		"invalid amount",
	)

	// ErrInvalidSignatureLength is returned when the signature is invalid.
	ErrInvalidSignatureLength = errors.New(
		"invalid signature length",
	)

	// ErrInvalidVersionLength is returned when the deposit version is invalid.
	ErrInvalidVersionLength = errors.New(
		"invalid version",
	)

	// ErrInvalidRootLength is returned when the deposit root is invalid.
	ErrInvalidRootLength = errors.New(
		"invalid root length",
	)

	// ErrDepositTransactionFailed is returned when the deposit transaction
	// fails.
	ErrDepositTransactionFailed = errors.New(
		"deposit transaction failed",
	)

	// ErrPrivateKeyRequired is returned when the broadcast flag is set but a
	// private key is not provided.
	ErrPrivateKeyRequired = errors.New(
		"private key required",
	)

	// ErrValidatorPrivateKeyRequired is returned when the validator private key
	// is required but not provided.
	ErrValidatorPrivateKeyRequired = errors.New(
		"validator private key required",
	)

	// ErrInvalidValidatorPrivateKeyLength is returned when the validator
	// private
	// key has an invalid length.
	ErrInvalidValidatorPrivateKeyLength = errors.New(
		"invalid validator private key length",
	)
)
