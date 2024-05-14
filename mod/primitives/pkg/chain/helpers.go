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

package chain

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
)

// ActiveForkVersion returns the active fork version for a given slot.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT,
]) ActiveForkVersionForSlot(
	slot SlotT,
) uint32 {
	return c.ActiveForkVersionForEpoch(c.SlotToEpoch(slot))
}

// ActiveForkVersionBySlot returns the active fork version for a given epoch.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT,
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
	DomainTypeT, EpochT, ExecutionAddressT, SlotT,
]) SlotToEpoch(slot SlotT) EpochT {
	//#nosec:G701 // realistically fine in practice.
	return EpochT(uint64(slot) / c.SlotsPerEpoch())
}

// WithinDAPeriod checks if the block epoch is within
// MIN_EPOCHS_FOR_BLOB_SIDECARS_REQUESTS
// of the given current epoch.
func (c chainSpec[
	DomainTypeT, EpochT, ExecutionAddressT, SlotT,
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
