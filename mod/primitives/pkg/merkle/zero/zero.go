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

package zero

import "github.com/berachain/beacon-kit/mod/primitives/pkg/crypto/sha256"

// NumZeroHashes is the number of pre-computed zero-hashes.
const NumZeroHashes = 64

// Hashes is a pre-computed list of zero-hashes for each depth level.
//
//nolint:gochecknoglobals // saves recomputing.
var Hashes [NumZeroHashes + 1][32]byte

// InitZeroHashes the zero-hashes pre-computed data
// with the given hash-function.
func InitZeroHashes(zeroHashesLevels int) {
	v := [64]byte{}
	for i := range zeroHashesLevels {
		copy(v[:32], Hashes[i][:])
		copy(v[32:], Hashes[i][:])
		Hashes[i+1] = sha256.Hash(v[:])
	}
}

//nolint:gochecknoinits // saves recomputing.
func init() {
	InitZeroHashes(NumZeroHashes)
}
