// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package core

import "github.com/berachain/beacon-kit/errors"

var (
	// ErrValSetCapExceeded is returned when the number of genesis deposits
	// exceeds the validator set cap.
	ErrValSetCapExceeded = errors.New("validator set cap exceeded at genesis")

	// ErrBlockSlotTooLow is returned when the block slot is too low.
	ErrBlockSlotTooLow = errors.New("block slot too low")

	// ErrSlotMismatch is returned when the slot in a block header does not
	// match the expected value.
	ErrSlotMismatch = errors.New("slot mismatch")

	// ErrProposerMismatch is returned when block builder does not match
	// with the proposer reported by consensus.
	ErrProposerMismatch = errors.New("proposer key mismatch")

	ErrDepositsRootMismatch = errors.New("deposits root mismatch")

	// ErrDepositsLengthMismatch is returned when length of deposits
	// listed in block is different from deposits from store.
	ErrDepositsLengthMismatch = errors.New("deposits lengths mismatched")

	// ErrDepositMismatch is returned when a specific deposit listed in
	// block is different from the correspondent one from store.
	ErrDepositMismatch = errors.New("deposit mismatched")

	// ErrDepositIndexOutOfOrder is returned when deposits are not in
	// contiguous order.
	ErrDepositIndexOutOfOrder = errors.New("deposit index out of order")

	// ErrParentRootMismatch is returned when the parent root in an execution
	// payload does not match the expected value.
	ErrParentRootMismatch = errors.New("parent root mismatch")

	// ErrParentPayloadHashMismatch is returned when the parent hash of an
	// execution payload does not match the expected value.
	ErrParentPayloadHashMismatch = errors.New("payload parent hash mismatch")

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

	// ErrExceedMaximumWithdrawals is returned when the number of withdrawals
	// in a block exceeds the maximum allowed.
	ErrExceedMaximumWithdrawals = errors.New("exceeds maximum withdrawals")

	// ErrZeroWithdrawals is returned when the number of withdrawals in a
	// block is zero. At least the EVM inflation withdrawal is always expected.
	ErrZeroWithdrawals = errors.New("zero withdrawals")

	// ErrNumWithdrawalsMismatch is returned when the number of withdrawals
	// in a block does not match the expected value.
	ErrNumWithdrawalsMismatch = errors.New("number of withdrawals mismatch")

	// ErrFirstWithdrawalNotEVMInflation is returned when the first withdrawal
	// in a block is not the EVM inflation withdrawal.
	ErrFirstWithdrawalNotEVMInflation = errors.New(
		"first withdrawal is not the EVM inflation withdrawal",
	)

	// ErrWithdrawalMismatch is returned when the withdrawals in a payload do
	// not match the local state's expected value.
	ErrWithdrawalMismatch = errors.New(
		"withdrawal mismatch between local state and payload")
)
