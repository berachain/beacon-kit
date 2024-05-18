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

package core

import "github.com/berachain/beacon-kit/mod/errors"

var (
	// ErrParentRootMismatch is returned when the parent root in an execution
	// payload does not match the expected value.
	ErrParentRootMismatch = errors.New("parent root mismatch")

	// ErrRandaoMixMismatch is returned when the randao mix in an execution
	// payload does not match the expected value.
	ErrRandaoMixMismatch = errors.New("randao mix mismatch")

	// ErrExceedBlockBlobLimit is returned when the block exceeds the blob
	// limit.
	ErrExceedBlockBlobLimit = errors.New("block exceeds blob limit")

	// ErrSlashedProposer is returned when a block is processed in which
	// the proposer is slashed.
	ErrSlashedProposer = errors.New(
		"attempted to process a block with a slashed proposer")

	// ErrStateRootMismatch is returned when the state root in a block header
	// does not match the expected value.
	ErrStateRootMismatch = errors.New("state root mismatch")
)
