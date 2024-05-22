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

package validator

import "github.com/berachain/beacon-kit/mod/errors"

var (
	// ErrForkVersionNotSupported is an error for when the fork
	// version is not supported.
	ErrForkVersionNotSupported = errors.New("fork version not supported")

	// ErrNilPayload is an error for when there is no payload
	// in a beacon block.
	ErrNilPayload = errors.New("nil payload in beacon block")

	// ErrNilBlkBody is an error for when the block body is nil.
	ErrNilBlkBody = errors.New("nil block body")

	// ErrNilWithdrawal is an error for when the deposit is nil.
	ErrNilWithdrawal = errors.New("nil withdrawal")

	// ErrNilWithdrawals is an error for when the deposits are nil.
	ErrNilWithdrawals = errors.New("nil withdrawals")

	// ErrNilBlobsBundle is an error for when the blobs bundle is nil.
	ErrNilBlobsBundle = errors.New("nil blobs bundle")
)
