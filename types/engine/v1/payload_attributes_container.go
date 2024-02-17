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

package enginev1

import (
	"errors"

	"github.com/itsdevbear/bolaris/types/consensus/version"
	"github.com/itsdevbear/bolaris/types/engine/interfaces"
	v1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
	"google.golang.org/protobuf/proto"
)

// WrappedPayloadAttributesV3 wraps the PayloadAttributesV3 from
// Prysmatic Labs' Engine API v1.
var _ interfaces.PayloadAttributer = (*PayloadAttributesContainer)(nil)

func NewPayloadAttributesContainer(
	v int,
	timestamp uint64, prevRandao []byte,
	suggestedFeeReceipient []byte,
	withdrawals []*v1.Withdrawal,
	parentBeaconBlockRoot []byte,
) (*PayloadAttributesContainer, error) {
	switch v {
	case version.Deneb:
		return &PayloadAttributesContainer{
			Attributes: &PayloadAttributesContainer_V3{
				V3: &v1.PayloadAttributesV3{
					Timestamp:             timestamp,
					PrevRandao:            prevRandao,
					SuggestedFeeRecipient: suggestedFeeReceipient,
					Withdrawals:           withdrawals,
					ParentBeaconBlockRoot: parentBeaconBlockRoot,
				},
			},
		}, nil
	case version.Capella:
		return &PayloadAttributesContainer{
			Attributes: &PayloadAttributesContainer_V2{
				V2: &v1.PayloadAttributesV2{
					Timestamp:             timestamp,
					PrevRandao:            prevRandao,
					SuggestedFeeRecipient: suggestedFeeReceipient,
					Withdrawals:           withdrawals,
				},
			},
		}, nil
	default:
		return nil, errors.New("invalid version")
	}
}

// Version returns the consensus version for the Capella upgrade.
func (w *PayloadAttributesContainer) Version() int {
	switch w.ToProto().(type) {
	case *v1.PayloadAttributesV3:
		return version.Deneb
	case *v1.PayloadAttributesV2:
		return version.Capella
	default:
		return 0
	}
}

// IsEmpty returns true if the WrappedPayloadAttributesV3 is empty.
func (w *PayloadAttributesContainer) IsEmpty() bool {
	switch w := w.GetAttributes().(type) {
	case *PayloadAttributesContainer_V3:
		return w.V3 == nil
	case *PayloadAttributesContainer_V2:
		return w.V2 == nil
	default:
		return true
	}
}

func (w *PayloadAttributesContainer) ToProto() proto.Message {
	switch p := w.GetAttributes().(type) {
	case *PayloadAttributesContainer_V3:
		return p.V3
	case *PayloadAttributesContainer_V2:
		return p.V2
	default:
		return nil
	}
}

// GetPrevRandao returns the previous RANDAO value.
func (w *PayloadAttributesContainer) GetPrevRandao() []byte {
	payload := w.ToProto()
	switch p := payload.(type) {
	case *v1.PayloadAttributesV3:
		return p.GetPrevRandao()
	case *v1.PayloadAttributesV2:
		return p.GetPrevRandao()
	default:
		return nil
	}
}

func (w *PayloadAttributesContainer) GetTimestamp() uint64 {
	payload := w.ToProto()
	switch p := payload.(type) {
	case *v1.PayloadAttributesV3:
		return p.GetTimestamp()
	case *v1.PayloadAttributesV2:
		return p.GetTimestamp()
	default:
		return 0
	}
}

func (w *PayloadAttributesContainer) GetSuggestedFeeRecipient() []byte {
	payload := w.ToProto()
	switch p := payload.(type) {
	case *v1.PayloadAttributesV3:
		return p.GetSuggestedFeeRecipient()
	case *v1.PayloadAttributesV2:
		return p.GetSuggestedFeeRecipient()
	default:
		return nil
	}
}

func (w *PayloadAttributesContainer) GetWithdrawals() []*v1.Withdrawal {
	payload := w.ToProto()
	switch p := payload.(type) {
	case *v1.PayloadAttributesV3:
		return p.GetWithdrawals()
	case *v1.PayloadAttributesV2:
		return p.GetWithdrawals()
	default:
		return nil
	}
}
