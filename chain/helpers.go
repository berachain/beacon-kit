// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package chain

import (
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
)

// ActiveForkVersionForTimestamp returns the active fork version for a given timestamp.
// Note that ActiveForkVersionForTimestamp will NOT check for
//   - ElectraDisableWithdrawalsForkTime
//   - ElectraEnableWithdrawalsForkTime
//
// As there is no logic that relies on the above forks that goes through ActiveForkVersionForTimestamp.
// It also gives flexibility for these forks to occur at any point after Electra, e.g. after Electra1.
func (s spec) ActiveForkVersionForTimestamp(timestamp math.U64) common.Version {
	time := timestamp.Unwrap()
	if time >= s.ElectraForkTime() {
		return version.Electra()
	}
	if time >= s.Deneb1ForkTime() {
		return version.Deneb1()
	}
	return version.Deneb()
}

// WithdrawalsDisabled is a switch that can be used to freeze withdrawals in an emergency scenario.
// An exception is made for the EVM inflation withdrawal which is always active.
func (s spec) WithdrawalsDisabled(timestamp math.U64) bool {
	time := timestamp.Unwrap()
	return time >= s.ElectraDisableWithdrawalsForkTime() && time < s.ElectraEnableWithdrawalsForkTime()
}

// GenesisForkVersion returns the fork version at genesis.
func (s spec) GenesisForkVersion() common.Version {
	return s.ActiveForkVersionForTimestamp(math.U64(s.GenesisTime()))
}

// SlotToEpoch converts a slot to an epoch.
func (s spec) SlotToEpoch(slot math.Slot) math.Epoch {
	return math.Epoch(slot.Unwrap() / s.SlotsPerEpoch())
}

// WithinDAPeriod checks if the block epoch is within MIN_EPOCHS_FOR_BLOB_SIDECARS_REQUESTS
// of the given current epoch.
func (s spec) WithinDAPeriod(block, current math.Slot) bool {
	return s.SlotToEpoch(block)+s.MinEpochsForBlobsSidecarsRequest() >= s.SlotToEpoch(current)
}
