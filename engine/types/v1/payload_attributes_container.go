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
	"github.com/itsdevbear/bolaris/types/consensus/version"
	"google.golang.org/protobuf/proto"
)

// Version returns the consensus version of the PayloadAttributesContainer.
func (w *PayloadAttributesContainer) Version() int {
	switch w.ToProto().(type) {
	case *PayloadAttributesV3:
		return version.Deneb
	default:
		return 0
	}
}

// IsEmpty returns true if the WrappedPayloadAttributesV3 is empty.
func (w *PayloadAttributesContainer) IsEmpty() bool {
	switch w := w.GetAttributes().(type) {
	case *PayloadAttributesContainer_V3:
		return w.V3 == nil
	default:
		return true
	}
}

// ToProto converts the PayloadAttributesContainer to its protobuf
// representation.
func (w *PayloadAttributesContainer) ToProto() proto.Message {
	switch p := w.GetAttributes().(type) {
	case *PayloadAttributesContainer_V3:
		return p.V3
	default:
		return nil
	}
}

// GetPrevRandao retrieves the previous RANDAO value from the
// PayloadAttributesContainer.
func (w *PayloadAttributesContainer) GetPrevRandao() []byte {
	payload := w.ToProto()
	switch p := payload.(type) {
	case *PayloadAttributesV3:
		return p.GetPrevRandao()
	case *PayloadAttributesV2:
		return p.GetPrevRandao()
	default:
		return nil
	}
}

// GetTimestamp fetches the timestamp from the PayloadAttributesContainer.
func (w *PayloadAttributesContainer) GetTimestamp() uint64 {
	payload := w.ToProto()
	switch p := payload.(type) {
	case *PayloadAttributesV3:
		return p.GetTimestamp()
	case *PayloadAttributesV2:
		return p.GetTimestamp()
	default:
		return 0
	}
}

// GetSuggestedFeeRecipient retrieves the suggested fee recipient address
// from the PayloadAttributesContainer.
func (w *PayloadAttributesContainer) GetSuggestedFeeRecipient() []byte {
	payload := w.ToProto()
	switch p := payload.(type) {
	case *PayloadAttributesV3:
		return p.GetSuggestedFeeRecipient()
	case *PayloadAttributesV2:
		return p.GetSuggestedFeeRecipient()
	default:
		return nil
	}
}

// GetWithdrawals fetches the list of withdrawals from the
// PayloadAttributesContainer.
func (w *PayloadAttributesContainer) GetWithdrawals() []*Withdrawal {
	payload := w.ToProto()
	switch p := payload.(type) {
	case *PayloadAttributesV3:
		return p.GetWithdrawals()
	case *PayloadAttributesV2:
		return p.GetWithdrawals()
	default:
		return nil
	}
}
