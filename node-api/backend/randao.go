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

package backend

import (
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
)

func (b *Backend) RandaoAtEpoch(slot math.Slot, epoch math.Epoch) (common.Bytes32, error) {
	// Get the state at the given slot.
	st, resolvedSlot, err := b.StateAtSlot(slot)
	if err != nil {
		return common.Bytes32{}, errors.Wrapf(err, "failed to get state from slot %d", slot)
	}

	// Infer the epoch if not provided.
	if epoch == 0 {
		epoch = b.cs.SlotToEpoch(resolvedSlot)
	}

	index := epoch.Unwrap() % b.cs.EpochsPerHistoricalVector()
	return st.GetRandaoMixAtIndex(index)
}
