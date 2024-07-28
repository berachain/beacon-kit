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

package types

import "github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/merkle"

const (
	// ZeroValidatorPubkeyGIndexDenebState is the generalized index of the 0
	// validator's pubkey in the beacon state in the Deneb fork. To get the
	// GIndex of the pubkey of validator at index n, the formula is:
	// GIndex = ZeroValidatorPubkeyGIndexDenebState +
	//          (ValidatorPubkeyGIndexOffset * n)
	ZeroValidatorPubkeyGIndexDenebState = 439804651110400

	// StateGIndexDenebBlock is the generalized index of the beacon state in
	// the beacon block in the Deneb fork.
	StateGIndexDenebBlock = 11

	// ValidatorPubkeyGIndexOffset is the offset of a validator pubkey GIndex.
	ValidatorPubkeyGIndexOffset = 8

	// MaxBlockRoots is the maximum number of block roots in the beacon block.
	MaxBlockRoots uint64 = 8192

	// MaxStateRoots is the maximum number of state roots in the beacon state.
	MaxStateRoots uint64 = 8192

	// MaxBalances is the maximum number of balances in the beacon state.
	MaxBalances uint64 = 1099511627776

	// MaxRandaoMixes is the maximum number of randao mixes in the beacon state.
	MaxRandaoMixes uint64 = 65536

	// MaxSlashings is the maximum number of slashings in the beacon state.
	MaxSlashings uint64 = 1099511627776
)

var (
	// ZeroValidatorPubkeyGIndexDenebBlock is the generalized index of the 0
	// validator's pubkey in the beacon block in the Deneb fork. To get the
	// GIndex of the pubkey of validator at index n, the formula is:
	// GIndex = ZeroValidatorPubkeyGIndexDenebBlock +
	//          (ValidatorPubkeyGIndexOffset * n)
	ZeroValidatorPubkeyGIndexDenebBlock = merkle.GeneralizedIndices{
		StateGIndexDenebBlock,
		ZeroValidatorPubkeyGIndexDenebState,
	}.Concat()
)
