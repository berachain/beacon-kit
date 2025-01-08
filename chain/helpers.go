// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
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
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
)

// ActiveForkVersionForSlot returns the active fork version for a given slot.
func (s spec) ActiveForkVersionForSlot(slot math.Slot) uint32 {
	return s.ActiveForkVersionForEpoch(s.SlotToEpoch(slot))
}

// ActiveForkVersionForEpoch returns the active fork version for a given epoch.
func (s spec) ActiveForkVersionForEpoch(epoch math.Epoch) uint32 {
	if epoch >= s.ElectraForkEpoch() {
		return version.Electra
	} else if epoch >= s.DenebPlusForkEpoch() {
		return version.DenebPlus
	}

	return version.Deneb
}

// SlotToEpoch converts a slot to an epoch.
func (s spec) SlotToEpoch(slot math.Slot) math.Epoch {
	//#nosec:G701 // realistically fine in practice.
	return math.Epoch(uint64(slot) / s.SlotsPerEpoch())
}

// WithinDAPeriod checks if the block epoch is within MIN_EPOCHS_FOR_BLOB_SIDECARS_REQUESTS
// of the given current epoch.
func (s spec) WithinDAPeriod(block, current math.Slot) bool {
	return (s.SlotToEpoch(block) + s.MinEpochsForBlobsSidecarsRequest()) >= s.SlotToEpoch(current)
}
