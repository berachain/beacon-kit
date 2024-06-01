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

package builder

import (
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// getPayloadAttributes returns the payload attributes for the given state and
// slot. The attribute is required to initiate a payload build process in the
// context of an `engine_forkchoiceUpdated` call.
func (pb *PayloadBuilder[BeaconStateT, ExecutionPayloadT]) getPayloadAttribute(
	st BeaconStateT,
	slot math.Slot,
	timestamp uint64,
	prevHeadRoot [32]byte,
) (engineprimitives.PayloadAttributer, error) {
	var (
		prevRandao [32]byte
	)

	// Get the expected withdrawals to include in this payload.
	withdrawals, err := st.ExpectedWithdrawals()
	if err != nil {
		pb.logger.Error(
			"Could not get expected withdrawals to get payload attribute",
			"error",
			err,
		)
		return nil, err
	}

	epoch := pb.chainSpec.SlotToEpoch(slot)

	// Get the previous randao mix.
	prevRandao, err = st.GetRandaoMixAtIndex(
		uint64(epoch) % pb.chainSpec.EpochsPerHistoricalVector(),
	)
	if err != nil {
		return nil, err
	}

	return engineprimitives.NewPayloadAttributes(
		pb.chainSpec.ActiveForkVersionForEpoch(epoch),
		timestamp,
		prevRandao,
		pb.cfg.SuggestedFeeRecipient,
		withdrawals,
		prevHeadRoot,
	)
}
