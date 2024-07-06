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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
)

// ActiveForkVersionForSlot returns the active fork version for a given slot.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) ActiveForkVersionForSlot(
	slot SlotT,
) uint32 {
	return c.ActiveForkVersionForEpoch(c.SlotToEpoch(slot))
}

// ActiveForkVersionForEpoch returns the active fork version for a given epoch.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) ActiveForkVersionForEpoch(
	epoch EpochT,
) uint32 {
	if epoch >= c.Data.ElectraForkEpoch {
		return version.Electra
	}

	return version.Deneb
}

// SlotToEpoch converts a slot to an epoch.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) SlotToEpoch(slot SlotT) EpochT {
	//#nosec:G701 // realistically fine in practice.
	return EpochT(uint64(slot) / c.SlotsPerEpoch())
}

// WithinDAPeriod checks if the block epoch is within
// MIN_EPOCHS_FOR_BLOB_SIDECARS_REQUESTS
// of the given current epoch.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT, CometBFTConfigT,
]) WithinDAPeriod(
	block, current SlotT,
) bool {
	return c.SlotToEpoch(
		block,
	)+EpochT(
		c.MinEpochsForBlobsSidecarsRequest(),
	) >= c.SlotToEpoch(
		current,
	)
}
