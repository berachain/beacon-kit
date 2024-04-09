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
	"github.com/berachain/beacon-kit/mod/forks/version"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

//nolint:lll // struct tags.
//go:generate go run github.com/fjl/gencodec -type PayloadAttributes -field-override payloadAttributesJSONMarshaling -out attributes.json.go
type PayloadAttributes struct {
	// version is the version of the payload attributes.
	version uint32
	// Timestamp is the timestamp at which the block will be built at.
	Timestamp uint64 `json:"timestamp"             gencodec:"required"`
	// PrevRandao is the previous Randao value from the beacon chain as
	// per EIP-4399.
	PrevRandao [32]byte `json:"prevRandao"            gencodec:"required"`
	// SuggestedFeeRecipient is the suggested fee recipient for the block. If
	// the execution client has a different fee recipient, it will typically
	// ignore this value.
	SuggestedFeeRecipient primitives.ExecutionAddress `json:"suggestedFeeRecipient" gencodec:"required"`
	// Withdrawals is the list of withdrawals to be included in the block as per
	// EIP-4895
	Withdrawals []*primitives.Withdrawal `json:"withdrawals"`
	// ParentBeaconBlockRoot is the root of the parent beacon block. (The block
	// prior)
	// to the block currently being processed. This field was added in EIP-4788.
	ParentBeaconBlockRoot [32]byte `json:"parentBeaconBlockRoot"`
}

// JSON type overrides for PayloadAttributes.
type payloadAttributesJSONMarshaling struct {
	Timestamp             hexutil.Uint64
	PrevRandao            hexutil.Bytes
	ParentBeaconBlockRoot hexutil.Bytes
}

// NewPayloadAttributes creates a new PayloadAttributes.
func NewPayloadAttributes(
	forkVersion uint32,
	timestamp uint64, prevRandao [32]byte,
	suggestedFeeReceipient primitives.ExecutionAddress,
	withdrawals []*primitives.Withdrawal,
	parentBeaconBlockRoot [32]byte,
) (*PayloadAttributes, error) {
	p := &PayloadAttributes{
		version:               forkVersion,
		Timestamp:             timestamp,
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

// Validate validates the PayloadAttributes.
func (p *PayloadAttributes) Validate() error {
	if p.Timestamp == 0 {
		return ErrInvalidTimestamp
	}

	// TODO: how to handle? PrevRandao is empty on block 1.
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

// Version returns the version of the PayloadAttributes.
func (p *PayloadAttributes) Version() uint32 {
	return p.version
}
