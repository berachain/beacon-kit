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

package engineprimitives

import (
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
)

// PayloadAttributer represents payload attributes of a block.
type PayloadAttributer interface {
	constraints.Versionable
	constraints.Nillable
	// GetSuggestedFeeRecipient returns the suggested fee recipient for the
	// block.
	GetSuggestedFeeRecipient() gethprimitives.ExecutionAddress
}

// PayloadAttributes is the attributes of a block payload.
type PayloadAttributes[
	WithdrawalT any,
] struct {
	// version is the version of the payload attributes.
	version uint32 `json:"-"`
	// Timestamp is the timestamp at which the block will be built at.
	Timestamp math.U64 `json:"timestamp"`
	// PrevRandao is the previous Randao value from the beacon chain as
	// per EIP-4399.
	PrevRandao common.Bytes32 `json:"prevRandao"`
	// SuggestedFeeRecipient is the suggested fee recipient for the block. If
	// the execution client has a different fee recipient, it will typically
	// ignore this value.
	SuggestedFeeRecipient gethprimitives.ExecutionAddress `json:"suggestedFeeRecipient"`
	// Withdrawals is the list of withdrawals to be included in the block as per
	// EIP-4895
	Withdrawals []WithdrawalT `json:"withdrawals"`
	// ParentBeaconBlockRoot is the root of the parent beacon block. (The block
	// prior)
	// to the block currently being processed. This field was added for
	// EIP-4788.
	ParentBeaconBlockRoot common.Root `json:"parentBeaconBlockRoot"`
}

// NewPayloadAttributes creates a new PayloadAttributes.
func NewPayloadAttributes[
	WithdrawalT any,
](
	forkVersion uint32,
	timestamp uint64,
	prevRandao common.Bytes32,
	suggestedFeeRecipient gethprimitives.ExecutionAddress,
	withdrawals []WithdrawalT,
	parentBeaconBlockRoot common.Root,
) (*PayloadAttributes[WithdrawalT], error) {
	p := &PayloadAttributes[WithdrawalT]{
		version:               forkVersion,
		Timestamp:             math.U64(timestamp),
		PrevRandao:            prevRandao,
		SuggestedFeeRecipient: suggestedFeeRecipient,
		Withdrawals:           withdrawals,
		ParentBeaconBlockRoot: parentBeaconBlockRoot,
	}

	if err := p.Validate(); err != nil {
		return nil, err
	}

	return p, nil
}

// New empty PayloadAttributes.
func (p *PayloadAttributes[WithdrawalT]) New(
	forkVersion uint32,
	timestamp uint64,
	prevRandao common.Bytes32,
	suggestedFeeRecipient gethprimitives.ExecutionAddress,
	withdrawals []WithdrawalT,
	parentBeaconBlockRoot common.Root,
) (*PayloadAttributes[WithdrawalT], error) {
	var err error
	p, err = NewPayloadAttributes(
		forkVersion,
		timestamp,
		prevRandao,
		suggestedFeeRecipient,
		withdrawals,
		parentBeaconBlockRoot,
	)
	return p, err
}

// IsNil returns true if the PayloadAttributes is nil.
func (p *PayloadAttributes[WithdrawalT]) IsNil() bool {
	return p == nil
}

// GetSuggestedFeeRecipient returns the suggested fee recipient.
func (
	p *PayloadAttributes[WithdrawalT],
) GetSuggestedFeeRecipient() gethprimitives.ExecutionAddress {
	return p.SuggestedFeeRecipient
}

// Version returns the version of the PayloadAttributes.
func (p *PayloadAttributes[WithdrawalT]) Version() uint32 {
	return p.version
}

// Validate validates the PayloadAttributes.
func (p *PayloadAttributes[WithdrawalT]) Validate() error {
	if p.Timestamp == 0 {
		return ErrInvalidTimestamp
	}

	if p.PrevRandao == [32]byte{} {
		return ErrEmptyPrevRandao
	}

	if p.Withdrawals == nil && p.version >= version.Capella {
		return ErrNilWithdrawals
	}

	// TODO: currently beaconBlockRoot is 0x000 on block 1, we need
	// to fix this, before uncommenting the line below.
	// if p.ParentBeaconBlockRoot == [32]byte{} {
	// 	return ErrInvalidParentBeaconBlockRoot
	// }

	return nil
}
