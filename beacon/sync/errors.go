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

package sync

import "errors"

var (
	// ErrConsensusClientIsSyncing indicates that the consensus client
	// is still in the process of syncing.
	ErrConsensusClientIsSyncing = errors.New(
		"consensus client is still syncing",
	)

	// ErrExecutionClientIsSyncing indicates that the execution client
	// is still in the process of syncing.
	ErrExecutionClientIsSyncing = errors.New(
		"execution client is still syncing",
	)

	// ErrNotRunning indicates that the service is not currently running.
	ErrNotRunning = errors.New(
		"service is not running",
	)

	// ErrInsufficientCLPeers indicates that the execution client does not have
	// enough peers to be considered healthy.
	ErrInsufficientCLPeers = errors.New(
		"number of consensus client peers below configured threshold",
	)

	// ErrInsufficientExPeers indicates that the execution client does not have
	// enough peers to be considered healthy.
	ErrInsufficientELPeers = errors.New(
		"number of execution client peers below configured threshold",
	)
)
