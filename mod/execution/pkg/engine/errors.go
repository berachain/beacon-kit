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

package engine

import "github.com/berachain/beacon-kit/mod/errors"

var (
	// ErrExecutionClientDisconnected represents an error when
	/// the execution client is disconnected.
	ErrExecutionClientDisconnected = errors.New(
		"execution client disconnected")

	// ErrAcceptedSyncingPayloadStatus represents an error when
	// the payload status is SYNCING or ACCEPTED.
	ErrAcceptedSyncingPayloadStatus = errors.New(
		"payload status is SYNCING or ACCEPTED")

	// ErrInvalidPayloadStatus represents an error when the
	// payload status is INVALID.
	ErrInvalidPayloadStatus = errors.New(
		"payload status is INVALID")

	// ErrBadBlockProduced represents an error when the beacon
	// chain has produced a bad block.
	ErrBadBlockProduced = errors.New(
		"proposer has produced a bad block, RIP walrus",
	)
)
