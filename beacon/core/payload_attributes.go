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

	"github.com/itsdevbear/bolaris/beacon/state"
	"github.com/itsdevbear/bolaris/config"
	"github.com/itsdevbear/bolaris/types/consensus/version"
	"github.com/itsdevbear/bolaris/types/engine"
	capella "github.com/itsdevbear/bolaris/types/engine/v1/capella"
	deneb "github.com/itsdevbear/bolaris/types/engine/v1/deneb"
)

// BuildPayloadAttributes builds payload attributes for a given state.
func BuildPayloadAttributes(
	beaconConfig *config.Beacon,
	st state.BeaconState,
	logger log.Logger,
	prevRando []byte,
	parentBeaconBlockRoot []byte,
	timestamp uint64,
) engine.PayloadAttributer {
	// Since beacon-kit is always post capella we can assume calling ExpectedWithdrawals
	// is valid.
	withdrawals, err := st.ExpectedWithdrawals()
	if err != nil {
		logger.Error(
			"Could not get expected withdrawals to get payload attribute", "error", err)
		return nil
	}

	switch st.Version() {
	case version.Deneb:
		return deneb.NewWrappedPayloadAttributerV3(
			timestamp, prevRando, beaconConfig.Validator.SuggestedFeeRecipient, withdrawals,
			parentBeaconBlockRoot,
		)
	case version.Capella:
		return capella.NewWrappedPayloadAttributerV2(
			timestamp, prevRando, beaconConfig.Validator.SuggestedFeeRecipient, withdrawals,
		)
	default:
		logger.Error(
			"Could not get payload attribute due to unknown state version",
			"version", st.Version(),
		)
	}
	return nil
}
