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
	"fmt"

	"github.com/berachain/beacon-kit/primitives/math"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
)

// StateAtSlot returns the beacon state at a particular slot using query context,
// resolving an input slot of 0 to the latest slot.
//
// This returns the beacon state of the version that was committed to disk at the requested slot,
// which has the empty state root in the latest block header. Hence, the most recent state and
// block roots are not updated.
func (b *Backend) StateAtSlot(slot math.Slot) (*statedb.StateDB, math.Slot, error) {
	queryCtx, err := b.node.CreateQueryContext(int64(slot), false) // #nosec G115 -- not an issue in practice.
	if err != nil {
		return nil, slot, fmt.Errorf("CreateQueryContext failed: %w", err)
	}
	st := b.sb.StateFromContext(queryCtx)

	// If using height 0 for the query context, make sure to return the latest slot.
	if slot == 0 {
		slot, err = st.GetSlot()
		if err != nil {
			return st, slot, fmt.Errorf("GetSlot failed: %w", err)
		}
	}
	return st, slot, nil
}
