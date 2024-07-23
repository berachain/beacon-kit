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

package utils

import (
	"strconv"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// ValidatorIndexByID parses a validator index from a string.
// The string can be either a validator index or a validator pubkey.
func ValidatorIndexByID[
	BeaconStateT interface {
		ValidatorIndexByPubkey(key crypto.BLSPubkey) (math.U64, error)
	},
](st BeaconStateT, keyOrIndex string) (math.U64, error) {
	index, err := strconv.ParseUint(keyOrIndex, 10, 64)
	if err == nil {
		return math.U64(index), nil
	}
	var key crypto.BLSPubkey
	if err = key.UnmarshalText([]byte(keyOrIndex)); err != nil {
		return math.U64(0), err
	}
	return st.ValidatorIndexByPubkey(key)
}
