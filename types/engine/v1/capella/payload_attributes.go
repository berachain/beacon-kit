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

package capella

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/itsdevbear/bolaris/types/consensus/version"
	"github.com/itsdevbear/bolaris/types/engine/interfaces"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
	"google.golang.org/protobuf/proto"
)

// WrappedPayloadAttributesV2 ensures compatibility with the
// engine.PayloadAttributer interface.
var _ interfaces.PayloadAttributer = (*WrappedPayloadAttributesV2)(nil)

// WrappedPayloadAttributesV2 wraps the PayloadAttributesV2
// from Prysmatic Labs' EngineAPI v1 protobuf definitions.
type WrappedPayloadAttributesV2 struct {
	*enginev1.PayloadAttributesV2
}

// NewWrappedExecutionPayloadDeneb creates a new WrappedPayloadAttributesV2.
func NewWrappedPayloadAttributerV2(
	timestamp uint64, prevRandao []byte,
	suggestedFeeReceipient common.Address,
	withdrawals []*enginev1.Withdrawal,
) *WrappedPayloadAttributesV2 {
	return &WrappedPayloadAttributesV2{
		PayloadAttributesV2: &enginev1.PayloadAttributesV2{
			Timestamp:             timestamp,
			PrevRandao:            prevRandao,
			SuggestedFeeRecipient: suggestedFeeReceipient.Bytes(),
			Withdrawals:           withdrawals,
		},
	}
}

// Version returns the consensus version for the Capella upgrade.
func (w *WrappedPayloadAttributesV2) Version() int {
	return version.Capella
}

// IsEmpty returns true if the WrappedPayloadAttributesV3 is empty.
func (w *WrappedPayloadAttributesV2) IsEmpty() bool {
	return w.PayloadAttributesV2 == nil
}

// ToProto returns the WrappedPayloadAttributesV3 as a proto.Message.
func (w *WrappedPayloadAttributesV2) ToProto() proto.Message {
	return w.PayloadAttributesV2
}
