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

package enginetypes

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	enginev1 "github.com/itsdevbear/bolaris/engine/types/v1"
	"github.com/itsdevbear/bolaris/primitives"
)

//go:generate go run github.com/fjl/gencodec -type PayloadAttributes -field-override payloadAttributesJSONMarshaling -out attributes.json.go

//nolint:lll // struct tags.
type PayloadAttributes struct {
	version               int
	Timestamp             uint64                      `json:"timestamp"             gencodec:"required"`
	PrevRandao            [32]byte                    `json:"prevRandao"            gencodec:"required"`
	SuggestedFeeRecipient primitives.ExecutionAddress `json:"suggestedFeeRecipient" gencodec:"required"`
	Withdrawals           []*enginev1.Withdrawal      `json:"withdrawals"`
	ParentBeaconBlockRoot *[32]byte                   `json:"parentBeaconBlockRoot"`
}

// JSON type overrides for PayloadAttributes.
type payloadAttributesJSONMarshaling struct {
	Timestamp             hexutil.Uint64
	PrevRandao            hexutil.Bytes
	ParentBeaconBlockRoot *common.Hash
}

// NewPayloadAttributes creates a new PayloadAttributes.
func NewPayloadAttributes(
	forkVersion int,
	timestamp uint64, prevRandao [32]byte,
	suggestedFeeReceipient primitives.ExecutionAddress,
	withdrawals []*enginev1.Withdrawal,
	parentBeaconBlockRoot [32]byte,
) (*PayloadAttributes, error) {
	if withdrawals == nil {
		withdrawals = make([]*enginev1.Withdrawal, 0)
	}

	return &PayloadAttributes{
		version:               forkVersion,
		Timestamp:             timestamp,
		PrevRandao:            prevRandao,
		SuggestedFeeRecipient: suggestedFeeReceipient,
		Withdrawals:           withdrawals,
		ParentBeaconBlockRoot: &parentBeaconBlockRoot,
	}, nil
}

// GetTimestamp returns the timestamp of the PayloadAttributes.
func (p *PayloadAttributes) GetTimestamp() uint64 {
	return p.Timestamp
}

// GetSuggestedFeeRecipient returns the suggested fee recipient address of the
// PayloadAttributes.
//
//nolint:lll
func (p *PayloadAttributes) GetSuggestedFeeRecipient() primitives.ExecutionAddress {
	return p.SuggestedFeeRecipient
}

// GetWithdrawals returns the list of withdrawals in the PayloadAttributes.
func (p *PayloadAttributes) GetWithdrawals() []*enginev1.Withdrawal {
	return p.Withdrawals
}

// GetParentBeaconBlockRoot returns the parent beacon block root of the
// PayloadAttributes.
// If the parent beacon block root is nil, a zero-value [32]byte is returned.
func (p *PayloadAttributes) GetParentBeaconBlockRoot() [32]byte {
	if p.ParentBeaconBlockRoot == nil {
		return [32]byte{}
	}
	return *p.ParentBeaconBlockRoot
}

// Version returns the version of the PayloadAttributes.
func (p *PayloadAttributes) Version() int {
	return p.version
}

// GetPrevRandao returns the previous Randao value of the PayloadAttributes.
func (p *PayloadAttributes) GetPrevRandao() [32]byte {
	return p.PrevRandao
}
