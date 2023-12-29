// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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

package execution

import (
	"context"
	"crypto/rand"

	payloadattribute "github.com/prysmaticlabs/prysm/v4/consensus-types/payload-attribute"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
)

func (s *EngineNotifier) getPayloadAttributes(_ context.Context,
	_ /*slot*/ primitives.Slot, timestamp uint64) (payloadattribute.Attributer, error) {
	// TODO: modularize andn make better.

	// TODO: this is a hack to fill the PrevRandao field. It is not verifiable or safe
	// or anything, it's literally just to get basic functionality.
	var random [32]byte
	if _, err := rand.Read(random[:]); err != nil {
		return nil, err
	}

	// TODO: need to support multiple fork versions.
	return payloadattribute.New(&enginev1.PayloadAttributesV2{
		Timestamp:             timestamp,
		SuggestedFeeRecipient: s.etherbase.Bytes(),
		Withdrawals:           nil,
		PrevRandao:            append([]byte{}, random[:]...),
	})
}
