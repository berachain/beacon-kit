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

package core

import (
	"cosmossdk.io/log"

	"github.com/itsdevbear/bolaris/config"
	"github.com/itsdevbear/bolaris/types/state"
	payloadattribute "github.com/prysmaticlabs/prysm/v4/consensus-types/payload-attribute"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
	"github.com/prysmaticlabs/prysm/v4/runtime/version"
)

// BuildPayloadAttributes builds payload attributes for a given state.
func BuildPayloadAttributes(
	beaconConfig *config.Beacon,
	st state.BeaconState,
	logger log.Logger,
	prevRando []byte,
	headRoot []byte,
	t uint64,
) payloadattribute.Attributer {
	emptyAttri := payloadattribute.EmptyWithVersion(st.Version())
	var attr payloadattribute.Attributer
	switch st.Version() {
	case version.Deneb:
		withdrawals, err := st.ExpectedWithdrawals()
		if err != nil {
			logger.Error(
				"Could not get expected withdrawals to get payload attribute", "error", err)
			return emptyAttri
		}
		attr, err = payloadattribute.New(&enginev1.PayloadAttributesV3{
			Timestamp:             t,
			PrevRandao:            prevRando,
			SuggestedFeeRecipient: beaconConfig.Validator.SuggestedFeeRecipient[:],
			Withdrawals:           withdrawals,
			ParentBeaconBlockRoot: headRoot,
		})
		if err != nil {
			logger.Error("Could not get payload attribute", "error", err)
			return emptyAttri
		}
	case version.Capella:
		withdrawals, err := st.ExpectedWithdrawals()
		if err != nil {
			logger.Error(
				"Could not get expected withdrawals to get payload attribute", "error", err)
			return emptyAttri
		}
		attr, err = payloadattribute.New(&enginev1.PayloadAttributesV2{
			Timestamp:             t,
			PrevRandao:            prevRando,
			SuggestedFeeRecipient: beaconConfig.Validator.SuggestedFeeRecipient[:],
			Withdrawals:           withdrawals,
		})
		if err != nil {
			logger.Error(
				"Could not get payload attribute", "error", err)
			return emptyAttri
		}
	case version.Bellatrix:
		var err error
		attr, err = payloadattribute.New(&enginev1.PayloadAttributes{
			Timestamp:             t,
			PrevRandao:            prevRando,
			SuggestedFeeRecipient: beaconConfig.Validator.SuggestedFeeRecipient[:],
		})
		if err != nil {
			logger.Error("Could not get payload attribute", "error", err)
			return emptyAttri
		}
	default:
		logger.Error(
			"Could not get payload attribute due to unknown state version", "version", st.Version())
		return emptyAttri
	}

	return attr
}
