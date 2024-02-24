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

package engine

import (
	"errors"

	"github.com/itsdevbear/bolaris/types/consensus/version"
	enginev1 "github.com/itsdevbear/bolaris/types/engine/v1"
)

// NewPayloadAttributesContainer creates a new PayloadAttributesContainer.
func NewPayloadAttributesContainer(
	v int,
	timestamp uint64, prevRandao []byte,
	suggestedFeeReceipient []byte,
	withdrawals []*enginev1.Withdrawal,
	parentBeaconBlockRoot []byte,
) (PayloadAttributer, error) {
	switch v {
	case version.Deneb:
		return &enginev1.PayloadAttributesContainer{
			Attributes: &enginev1.PayloadAttributesContainer_V3{
				V3: &enginev1.PayloadAttributesV3{
					Timestamp:             timestamp,
					PrevRandao:            prevRandao,
					SuggestedFeeRecipient: suggestedFeeReceipient,
					Withdrawals:           withdrawals,
					ParentBeaconBlockRoot: parentBeaconBlockRoot,
				},
			},
		}, nil
	default:
		return nil, errors.New("invalid version")
	}
}

// EmptyPayloadAttributesWithVersion creates a new PayloadAttributesContainer
// with no attributes.
func EmptyPayloadAttributesWithVersion(
	v int,
) PayloadAttributer {
	switch v {
	case version.Deneb:
		return &enginev1.PayloadAttributesContainer{
			Attributes: &enginev1.PayloadAttributesContainer_V3{
				V3: nil,
			},
		}
	default:
		return nil
	}
}

// NewWithdrawal creates a new Withdrawal.
func NewWithdrawal(
	_ []byte, // validatorPubkey
	amount uint64,
) *enginev1.Withdrawal {
	// TODO: implement
	return &enginev1.Withdrawal{
		Amount: amount,
	}
}
