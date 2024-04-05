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

package types

import (
	"github.com/cockroachdb/errors"
)

var (
	// errInvalidIndex          = errors.New("index out of bounds").
	errInvalidInclusionProof = errors.New(
		"invalid KZG commitment inclusion proof",
	)
	// errNilBlockHeader = errors.New("received nil beacon block header").
)

var (
	// ErrInvalidExecutionValue is an error for when the
	// execution value is invalid.
	ErrInvalidExecutionValue = errors.New("invalid execution value")

	// ErrForkVersionNotSupported is an error for when the fork
	// version is not supported.
	ErrForkVersionNotSupported = errors.New("fork version not supported")

	// ErrNilBlk is an error for when the beacon block is nil.
	ErrNilBlk = errors.New("nil beacon block")

	// ErrNilPayload is an error for when there is no payload
	// in a beacon block.
	ErrNilPayload = errors.New("nil payload in beacon block")

	// ErrNilBlkBody is an error for when the block body is nil.
	ErrNilBlkBody = errors.New("nil block body")

	// ErrNilBlockHeader is returned when a block header from a block is nil.
	ErrNilBlockHeader = errors.New("nil block header")

	// ErrNilDeposits is an error for when the deposits are nil.
	ErrNilDeposits = errors.New("nil deposits")

	// ErrNilWithdrawal is an error for when the deposit is nil.
	ErrNilWithdrawal = errors.New("nil withdrawal")

	// ErrNilWithdrawals is an error for when the deposits are nil.
	ErrNilWithdrawals = errors.New("nil withdrawals")

	// ErrNilBlobsBundle is an error for when the blobs bundle is nil.
	ErrNilBlobsBundle = errors.New("nil blobs bundle")

	// ErrInvalidWithdrawalCredentials is an error for when the.
	ErrInvalidWithdrawalCredentials = errors.New(
		"invalid withdrawal credentials",
	)

	// ErrDepositMessage is an error for when the deposit signature doesn't
	// match.
	ErrDepositMessage = errors.New("invalid deposit message")
)
