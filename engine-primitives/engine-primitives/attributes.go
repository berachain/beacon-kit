// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
)

// PayloadAttributes is the attributes of a block payload.
type PayloadAttributes struct {
	// Timestamp is the timestamp at which the block will be built at.
	Timestamp math.U64 `json:"timestamp"`
	// PrevRandao is the previous Randao value from the beacon chain as
	// per EIP-4399.
	PrevRandao common.Bytes32 `json:"prevRandao"`
	// SuggestedFeeRecipient is the suggested fee recipient for the block. If
	// the execution client has a different fee recipient, it will typically
	// ignore this value.
	SuggestedFeeRecipient common.ExecutionAddress `json:"suggestedFeeRecipient"`
	// Withdrawals is the list of withdrawals to be included in the block as per
	// EIP-4895
	Withdrawals Withdrawals `json:"withdrawals"`
	// ParentBeaconBlockRoot is the root of the parent beacon block. (The block
	// prior)
	// to the block currently being processed. This field was added for
	// EIP-4788.
	ParentBeaconBlockRoot common.Root `json:"parentBeaconBlockRoot"`

	// PrevProposerPubKey carries the public key of previous block proposed
	// Needed to enshrine distributeFor in EVM
	PrevProposerPubKey crypto.BLSPubkey `json:"prevProposerPubKey"`
}

// NewPayloadAttributes creates a new PayloadAttributes and validates it for
// the given fork version.
func NewPayloadAttributes(
	forkVersion common.Version,
	timestamp math.U64,
	prevRandao common.Bytes32,
	suggestedFeeRecipient common.ExecutionAddress,
	withdrawals Withdrawals,
	parentBeaconBlockRoot common.Root,
	prevProposerPubKey crypto.BLSPubkey,
) (*PayloadAttributes, error) {
	pa := &PayloadAttributes{
		Timestamp:             timestamp,
		PrevRandao:            prevRandao,
		SuggestedFeeRecipient: suggestedFeeRecipient,
		Withdrawals:           withdrawals,
		ParentBeaconBlockRoot: parentBeaconBlockRoot,
		PrevProposerPubKey:    prevProposerPubKey,
	}

	if err := pa.Validate(forkVersion); err != nil {
		return nil, err
	}

	return pa, nil
}

// GetSuggestedFeeRecipient returns the suggested fee recipient.
func (p *PayloadAttributes) GetSuggestedFeeRecipient() common.ExecutionAddress {
	return p.SuggestedFeeRecipient
}

// Validate validates the PayloadAttributes for the given fork version.
func (p *PayloadAttributes) Validate(forkVersion common.Version) error {
	if p.Timestamp == 0 {
		return ErrInvalidTimestamp
	}

	if p.PrevRandao == [32]byte{} {
		return ErrEmptyPrevRandao
	}

	// For any fork version Capella onwards, withdrawals are required.
	if p.Withdrawals == nil && version.EqualsOrIsAfter(forkVersion, version.Capella()) {
		return ErrNilWithdrawals
	}

	emptyBLSPubKey := crypto.BLSPubkey{}
	if version.IsBefore(forkVersion, version.Electra1()) {
		if p.PrevProposerPubKey != emptyBLSPubKey {
			return ErrNonEmptyPrevProposerPubKey
		}
	} else {
		if p.PrevProposerPubKey == emptyBLSPubKey {
			return ErrEmptyPrevProposerPubKey
		}
	}

	return nil
}
