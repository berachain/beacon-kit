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

package engineprimitives

import "github.com/berachain/beacon-kit/mod/errors"

var (
	// ErrInvalidTimestamp indicates that the provided timestamp is not valid.
	ErrInvalidTimestamp = errors.New("invalid timestamp")

	// ErrInvalidRandao indicates that the provided RANDAO value is not valid.
	ErrInvalidRandao = errors.New("invalid randao")

	// ErrNilWithdrawals indicates that the withdrawals are in a
	// Capella versioned payload.
	ErrNilWithdrawals = errors.New("nil withdrawals post capella")

	// ErrEmptyPrevRandao indicates that the previous RANDAO value is empty.
	ErrEmptyPrevRandao = errors.New("empty randao")

	// ErrFailedToUnmarshalTx indicates that the transaction could not be
	// unmarshaled.
	ErrFailedToUnmarshalTx = errors.New("failed to unmarshal transaction")

	// ErrInvalidVersionedHash indicates that the versioned hash is invalid.
	ErrInvalidVersionedHash = errors.New("invalid versioned hash")

	// ErrMismatchedNumVersionedHashes indicates that the number of blobs in the
	// payload does not match the expected number.
	ErrMismatchedNumVersionedHashes = errors.New(
		"mismatch in number of versioned hashes",
	)

	// ErrPayloadBlockHashMismatch represents an error when the
	// block hash in the payload does not match from the assembled
	// block.
	ErrPayloadBlockHashMismatch = errors.New(
		"block hash in payload does not match assembled block",
	)
)
