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
	"github.com/ethereum/go-ethereum/common/hexutil"
	enginev1 "github.com/itsdevbear/bolaris/engine/types/v1"
)

//go:generate go run github.com/fjl/gencodec -type PayloadAttributes -field-override payloadAttributesJSONMarshaling -out attributes.json.go

type PayloadAttributes struct {
	version int

	Timestamp             uint64
	PrevRandao            []byte
	SuggestedFeeRecipient []byte
	Withdrawals           []*enginev1.Withdrawal
	ParentBeaconBlockRoot []byte
}

// JSON type overrides for PayloadAttributes.
type payloadAttributesJSONMarshaling struct {
	Timestamp hexutil.Uint64
}

// NewPayloadAttributes creates a new PayloadAttributes.
func NewPayloadAttributes(
	forkVersion int,
	timestamp uint64, prevRandao []byte,
	suggestedFeeReceipient []byte,
	withdrawals []*enginev1.Withdrawal,
	parentBeaconBlockRoot [32]byte,
) (*PayloadAttributes, error) {
	return &PayloadAttributes{
		version:               forkVersion,
		Timestamp:             timestamp,
		PrevRandao:            prevRandao,
		SuggestedFeeRecipient: suggestedFeeReceipient,
		Withdrawals:           withdrawals,
		ParentBeaconBlockRoot: parentBeaconBlockRoot[:],
	}, nil
}

func NewEmptyPayloadAttributesWithVersion(version int) *PayloadAttributes {
	return &PayloadAttributes{
		version: version,
	}
}

func (p *PayloadAttributes) GetTimestamp() uint64 {
	return p.Timestamp
}

func (p *PayloadAttributes) GetSuggestedFeeRecipient() []byte {
	return p.SuggestedFeeRecipient
}
func (p *PayloadAttributes) GetWithdrawals() []*enginev1.Withdrawal {
	return p.Withdrawals
}

func (p *PayloadAttributes) GetParentBeaconBlockRoot() []byte {
	return p.ParentBeaconBlockRoot
}

func (p *PayloadAttributes) Version() int {
	return p.version
}

func (p *PayloadAttributes) GetPrevRandao() []byte {
	return p.PrevRandao
}

func (p *PayloadAttributes) IsEmpty() bool {
	return p == nil
}
