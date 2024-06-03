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
	// ErrBlockSlotTooLow is returned when the block slot is too low.
	ErrBlockSlotTooLow = errors.New("block slot too low")

	// ErrBeaconStateOutOfSync is returned when the state is either too far
	// behind
	// or too far ahead of the head and we must abort the state transition.
	ErrBeaconStateOutOfSync = errors.New("state is out of sync with head")

	// ErrSlotMismatch is returned when the slot in a block header does not
	// match the expected value.
	ErrSlotMismatch = errors.New("slot mismatch")

	// ErrParentRootMismatch is returned when the parent root in an execution
	// payload does not match the expected value.
	ErrParentRootMismatch = errors.New("parent root mismatch")

	// ErrRandaoMixMismatch is returned when the randao mix in an execution
	// payload does not match the expected value.
	ErrRandaoMixMismatch = errors.New("randao mix mismatch")

	// ErrExceedsBlockDepositLimit is returned when the block exceeds the
	// deposit limit.
	ErrExceedsBlockDepositLimit = errors.New("block exceeds deposit limit")

	// ErrRewardsLengthMismatch is returned when the length of the rewards
	// in a block does not match the expected value.
	ErrRewardsLengthMismatch = errors.New("rewards length mismatch")

	// ErrPenaltiesLengthMismatch is returned when the length of the penalties
	// in a block does not match the expected value.
	ErrPenaltiesLengthMismatch = errors.New("penalties length mismatch")

	// ErrExceedsBlockBlobLimit is returned when the block exceeds the blob
	// limit.
	ErrExceedsBlockBlobLimit = errors.New("block exceeds blob limit")

	// ErrSlashedProposer is returned when a block is processed in which
	// the proposer is slashed.
	ErrSlashedProposer = errors.New(
		"attempted to process a block with a slashed proposer")

	// ErrStateRootMismatch is returned when the state root in a block header
	// does not match the expected value.
	ErrStateRootMismatch = errors.New("state root mismatch")

	// ErrInvalidSignature is returned when the signature is invalid.
	ErrInvalidSignature = errors.New("invalid signature")

	// ErrXorInvalid is returned when the XOR operation is invalid.
	ErrXorInvalid = errors.New("xor invalid")
)
