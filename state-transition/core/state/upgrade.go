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

package state

import (
	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/primitives/math"
)

// IsPostFork2 returns true if the chain is post-upgrade (Fork2 on Boonet).
//
// TODO: Jank. Refactor into better fork version management.
func IsPostFork2(chainID uint64, slot math.Slot) bool {
	switch chainID {
	case spec.BartioChainID:
		return false
	case spec.BoonetEth1ChainID:
		if slot < math.U64(spec.BoonetFork2Height) {
			return false
		}

		return true
	default:
		return true
	}
}

// IsPostFork3 returns true if the chain is post-upgrade (Fork3 on Boonet).
//
// TODO: Jank. Refactor into better fork version management.
func IsPostFork3(chainID uint64, slot math.Slot) bool {
	switch chainID {
	case spec.BartioChainID:
		return false
	case spec.BoonetEth1ChainID:
		if slot < math.U64(spec.BoonetFork3Height) {
			return false
		}

		return true
	default:
		return true
	}
}
