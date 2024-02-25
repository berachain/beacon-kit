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
	sdk "github.com/cosmos/cosmos-sdk/types"
	protov2 "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/protoadapt"
)

// // HasMsgs defines an interface a transaction must fulfill.
// HasMsgs interface {
// 	// GetMsgs gets the all the transaction's messages.
// 	GetMsgs() []Msg
// }

// // Tx defines an interface a transaction must fulfill.
// Tx interface {
// 	HasMsgs

// 	// GetMsgsV2 gets the transaction's messages as
// google.golang.org/protobuf/proto.Message's.
// 	GetMsgsV2() ([]protov2.Message, error)
// }

var _ sdk.Tx = (*BeaconKitBlockContainer)(nil)

// BeaconKitBlockContainer is a wrapper around a BeaconKitBlock that implements
// the
// sdk.Tx interface.

func (bc *BeaconKitBlockContainer) GetMsgs() []sdk.Msg {
	return []sdk.Msg{bc}
}

func (bc *BeaconKitBlockContainer) GetMsgsV2() ([]protov2.Message, error) {
	return []protov2.Message{protoadapt.MessageV2Of(bc)}, nil
}
