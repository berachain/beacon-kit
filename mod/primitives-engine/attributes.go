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

import (
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
)

// PayloadAttributes is the attributes of a block payload.
//

type PayloadAttributes[
	Withdrawal SSZMarshallable,
] struct {
	// version is the version of the payload attributes.
	version uint32 `json:"-"`
	// Timestamp is the timestamp at which the block will be built at.
	Timestamp math.U64 `json:"timestamp"`
	// PrevRandao is the previous Randao value from the beacon chain as
	// per EIP-4399.
	PrevRandao primitives.Bytes32 `json:"prevRandao"`
	// SuggestedFeeRecipient is the suggested fee recipient for the block. If
	// the execution client has a different fee recipient, it will typically
	// ignore this value.
	SuggestedFeeRecipient common.ExecutionAddress `json:"suggestedFeeRecipient"`
	// Withdrawals is the list of withdrawals to be included in the block as per
	// EIP-4895
	Withdrawals []Withdrawal `json:"withdrawals"`
	// ParentBeaconBlockRoot is the root of the parent beacon block. (The block
	// prior)
	// to the block currently being processed. This field was added for
	// EIP-4788.
	ParentBeaconBlockRoot primitives.Root `json:"parentBeaconBlockRoot"`
}

// NewPayloadAttributes creates a new PayloadAttributes.
func NewPayloadAttributes[
	Withdrawal SSZMarshallable,
](
	forkVersion uint32,
	timestamp uint64,
	prevRandao primitives.Bytes32,
	suggestedFeeReceipient common.ExecutionAddress,
	withdrawals []Withdrawal,
	parentBeaconBlockRoot primitives.Root,
) (*PayloadAttributes[Withdrawal], error) {
	p := &PayloadAttributes[Withdrawal]{
		version:               forkVersion,
		Timestamp:             math.U64(timestamp),
		PrevRandao:            prevRandao,
		SuggestedFeeRecipient: suggestedFeeReceipient,
		Withdrawals:           withdrawals,
		ParentBeaconBlockRoot: parentBeaconBlockRoot,
	}

	if err := p.Validate(); err != nil {
		return nil, err
	}

	return p, nil
}

// Version returns the version of the PayloadAttributes.
func (p *PayloadAttributes[Withdrawal]) Version() uint32 {
	return p.version
}

// Validate validates the PayloadAttributes.
func (p *PayloadAttributes[Withdrawal]) Validate() error {
	if p.Timestamp == 0 {
		return ErrInvalidTimestamp
	}

	// TODO: how to handle? PrevRandao is empty on block 1.
	// TODO: Fix is to seed PrevRandao with Eth1BlockHash
	// as per spec.
	// if p.PrevRandao == [32]byte{} {
	// 	return ErrEmptyPrevRandao
	// }
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
