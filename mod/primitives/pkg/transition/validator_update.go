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

package transition

import (
	"sort"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// ValidatorUpdates is a list of validator updates.
type ValidatorUpdates []*ValidatorUpdate

// ValidatorUpdate is a struct that holds the validator update.
type ValidatorUpdate struct {
	// Pubkey is the public key of the validator.
	Pubkey crypto.BLSPubkey
	// EffectiveBalance is the effective balance of the validator.
	EffectiveBalance math.Gwei
}

// RemoveDuplicates removes duplicate validator updates. We
// iterate through the list backwards since we want the last
// update to be the one that is kept.
func (vu ValidatorUpdates) RemoveDuplicates() ValidatorUpdates {
	duplicateCheck := make(map[crypto.BLSPubkey]struct{})
	j := len(vu) - 1
	for i := j; i >= 0; i-- {
		if _, exists := duplicateCheck[vu[i].Pubkey]; !exists {
			duplicateCheck[vu[i].Pubkey] = struct{}{}
			vu[j] = vu[i]
			j--
		}
	}

	return vu[j+1:]
}

// Sort sorts the validator updates.
func (vu ValidatorUpdates) Sort() ValidatorUpdates {
	sort.SliceStable(vu, func(i, j int) bool {
		return string((vu)[i].Pubkey[:]) < string((vu)[j].Pubkey[:])
	})
	return vu
}
