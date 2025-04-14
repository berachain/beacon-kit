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

package merkle

import (
	"fmt"

	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/version"
)

const (
	// ProposerIndexGIndexBlock is the generalized index of the proposer index in the beacon block.
	ProposerIndexGIndexBlock = 9

	// StateGIndexBlock is the generalized index of the beacon state in the beacon block.
	StateGIndexBlock = 11

	// ZeroValidatorPubkeyGIndexDenebState is the generalized index of the 0
	// validator's pubkey in the beacon state in the Deneb forks. To get the
	// GIndex of the pubkey of validator at index n, the formula is:
	// GIndex = ZeroValidatorPubkeyGIndexDenebState +
	//          (ValidatorPubkeyGIndexOffset * n)
	ZeroValidatorPubkeyGIndexDenebState = 439804651110400

	// ZeroValidatorPubkeyGIndexDenebBlock is the generalized index of the 0
	// validator's pubkey in the beacon block in the Deneb forks. This is
	// calculated by concatenating the (ZeroValidatorPubkeyGIndexDenebState,
	// StateGIndexDenebBlock) GIndices. To get the GIndex of the pubkey of
	// validator at index n, the formula is:
	// GIndex = ZeroValidatorPubkeyGIndexDenebBlock +
	//          (ValidatorPubkeyGIndexOffset * n)
	ZeroValidatorPubkeyGIndexDenebBlock = 3254554418216960

	// ValidatorPubkeyGIndexOffset is the offset of a validator pubkey GIndex.
	ValidatorPubkeyGIndexOffset = 8

	// ZeroValidatorPubkeyGIndexElectraState is the generalized index of the 0
	// validator's pubkey in the beacon state in the Deneb forks. To get the
	// GIndex of the pubkey of validator at index n, the formula is:
	// GIndex = ZeroValidatorPubkeyGIndexElectraState +
	//          (ValidatorPubkeyGIndexOffset * n)
	ZeroValidatorPubkeyGIndexElectraState = 721279627821056

	// ZeroValidatorPubkeyGIndexElectraBlock is the generalized index of the 0
	// validator's pubkey in the beacon block in the Deneb forks. This is
	// calculated by concatenating the (ZeroValidatorPubkeyGIndexDenebState,
	// StateGIndexDenebBlock) GIndices. To get the GIndex of the pubkey of
	// validator at index n, the formula is:
	// GIndex = ZeroValidatorPubkeyGIndexElectraBlock +
	//          (ValidatorPubkeyGIndexOffset * n)
	ZeroValidatorPubkeyGIndexElectraBlock = 6350779162034176
)

// GetZeroValidatorPubkeyGIndexState determines the generalized index of the 0
// validator's pubkey in the beacon state based on the fork version.
func GetZeroValidatorPubkeyGIndexState(forkVersion common.Version) (int, error) {
	if version.EqualsOrIsAfter(forkVersion, version.Electra()) {
		return ZeroValidatorPubkeyGIndexElectraState, nil
	} else if version.EqualsOrIsAfter(forkVersion, version.Deneb()) {
		return ZeroValidatorPubkeyGIndexDenebState, nil
	}
	return 0, fmt.Errorf("unsupported fork version: %s", forkVersion)
}

// GetZeroValidatorPubkeyGIndexBlock determines the generalized index of the 0
// validator's pubkey in the beacon block based on the fork version.
func GetZeroValidatorPubkeyGIndexBlock(forkVersion common.Version) (uint64, error) {
	if version.EqualsOrIsAfter(forkVersion, version.Electra()) {
		return ZeroValidatorPubkeyGIndexElectraBlock, nil
	} else if version.EqualsOrIsAfter(forkVersion, version.Deneb()) {
		return ZeroValidatorPubkeyGIndexDenebBlock, nil
	}
	return 0, fmt.Errorf("unsupported fork version: %s", forkVersion)
}
