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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package merkle

import (
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/node-api/handlers/proof/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/merkle"
)

// ProveValidatorBalanceInState generates a proof for a validator's balance in the
// beacon state. The validatorOffset must be computed as (ValidatorGIndexOffset * validatorIndex).
func ProveValidatorBalanceInState(
	forkVersion common.Version,
	bsm types.BeaconStateMarshallable,
	validatorOffset math.U64,
) ([]common.Root, common.Root, error) {
	stateProofTree, err := bsm.GetTree()
	if err != nil {
		return nil, common.Root{}, err
	}

	zeroBalanceGIndexState, err := GetZeroValidatorBalanceGIndexState(forkVersion)
	if err != nil {
		return nil, common.Root{}, err
	}

	gIndex := zeroBalanceGIndexState + int(validatorOffset) // #nosec G115 -- bounded
	balProof, err := stateProofTree.Prove(gIndex)
	if err != nil {
		return nil, common.Root{}, err
	}

	proof := make([]common.Root, len(balProof.Hashes))
	for i, hash := range balProof.Hashes {
		proof[i] = common.NewRootFromBytes(hash)
	}
	return proof, common.NewRootFromBytes(balProof.Leaf), nil
}

// ProveValidatorBalanceInBlock generates a proof for a validator's balance in the
// beacon block. The proof is verified against the beacon block root and returned
// alongside the beacon block root used for verification.
func ProveValidatorBalanceInBlock(
	validatorIndex math.U64,
	bbh *ctypes.BeaconBlockHeader,
	bsm types.BeaconStateMarshallable,
) ([]common.Root, common.Root, error) {
	forkVersion := bsm.GetForkVersion()
	validatorOffset := ValidatorGIndexOffset * validatorIndex

	balInStateProof, leaf, err := ProveValidatorBalanceInState(forkVersion, bsm, validatorOffset)
	if err != nil {
		return nil, common.Root{}, err
	}

	stateInBlockProof, err := ProveBeaconStateInBlock(bbh, false)
	if err != nil {
		return nil, common.Root{}, err
	}

	combinedProof := append(balInStateProof, stateInBlockProof...)
	beaconRoot, err := verifyValidatorBalanceInBlock(forkVersion, bbh, validatorOffset, combinedProof, leaf)
	if err != nil {
		return nil, common.Root{}, err
	}
	return combinedProof, beaconRoot, nil
}

// verifyValidatorBalanceInBlock verifies the provided Merkle proof of a validator's
// balance inside the beacon block and returns the beacon block root that the proof
// was verified against.
func verifyValidatorBalanceInBlock(
	forkVersion common.Version,
	bbh *ctypes.BeaconBlockHeader,
	validatorOffset math.U64,
	proof []common.Root,
	leaf common.Root,
) (common.Root, error) {
	zeroBalanceGIndexBlock, err := GetZeroValidatorBalanceGIndexBlock(forkVersion)
	if err != nil {
		return common.Root{}, err
	}

	beaconRoot := bbh.HashTreeRoot()
	if !merkle.VerifyProof(beaconRoot, leaf, zeroBalanceGIndexBlock+validatorOffset.Unwrap(), proof) {
		return common.Root{}, errors.Wrapf(
			errors.New("validator balance proof failed to verify against beacon root"),
			"beacon root: 0x%s", beaconRoot,
		)
	}

	return beaconRoot, nil
}
