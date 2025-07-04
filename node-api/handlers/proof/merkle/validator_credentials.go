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

// ProveWithdrawalCredentialsInState generates a proof for a validator's
// withdrawal credentials in the beacon state. The validatorOffset must be
// computed as (ValidatorWithdrawalCredentialsGIndexOffset * validatorIndex).
func ProveWithdrawalCredentialsInState(
	forkVersion common.Version,
	bsm types.BeaconStateMarshallable,
	validatorOffset math.U64,
) ([]common.Root, common.Root, error) {
	stateProofTree, err := bsm.GetTree()
	if err != nil {
		return nil, common.Root{}, err
	}

	// Determine the starting generalized index for the 0-th validator's
	// withdrawal credentials for this fork.
	zeroWithdrawalGIndexState, err := GetZeroValidatorCredentialsGIndexState(forkVersion)
	if err != nil {
		return nil, common.Root{}, err
	}

	// Calculate the generalized index for the target validator. The offset
	// multiplication is bounded by (2^40-1)*8 < 2^43 < 2^63, so converting to
	// int is safe on 64-bit architectures.
	gIndex := zeroWithdrawalGIndexState + int(validatorOffset) // #nosec G115

	withdrawalProof, err := stateProofTree.Prove(gIndex)
	if err != nil {
		return nil, common.Root{}, err
	}

	proof := make([]common.Root, len(withdrawalProof.Hashes))
	for i, hash := range withdrawalProof.Hashes {
		proof[i] = common.NewRootFromBytes(hash)
	}
	return proof, common.NewRootFromBytes(withdrawalProof.Leaf), nil
}

// ProveWithdrawalCredentialsInBlock generates a proof for a validator's
// withdrawal credentials in the beacon block. The proof is verified against
// the beacon block root as a sanity check and the "correct" beacon block root
// is returned alongside the proof.
func ProveWithdrawalCredentialsInBlock(
	validatorIndex math.U64,
	bbh *ctypes.BeaconBlockHeader,
	bsm types.BeaconStateMarshallable,
) ([]common.Root, common.Root, error) {
	forkVersion := bsm.GetForkVersion()

	// Calculate the validator-specific offset.
	validatorOffset := ValidatorGIndexOffset * validatorIndex

	// 1. Proof inside the state.
	withdrawalInStateProof, leaf, err := ProveWithdrawalCredentialsInState(
		forkVersion, bsm, validatorOffset,
	)
	if err != nil {
		return nil, common.Root{}, err
	}

	// 2. Proof of the state inside the block.
	stateInBlockProof, err := ProveBeaconStateInBlock(bbh, false)
	if err != nil {
		return nil, common.Root{}, err
	}

	// 3. Combine proofs: state-level hashes come first, followed by block-level
	// hashes (same order as ProveProposerPubkeyInBlock).
	//
	//nolint:gocritic // ok.
	combinedProof := append(withdrawalInStateProof, stateInBlockProof...)

	// 4. Verify the combined proof against the beacon block root.
	beaconRoot, err := verifyWithdrawalCredentialsInBlock(
		forkVersion, bbh, validatorOffset, combinedProof, leaf,
	)
	if err != nil {
		return nil, common.Root{}, err
	}

	return combinedProof, beaconRoot, nil
}

// verifyWithdrawalCredentialsInBlock verifies the provided Merkle proof of a
// validator's withdrawal credentials inside the beacon block and returns the
// beacon block root that the proof was verified against.
//
// NOTE: Proof verification is not strictly necessary for operation, but we do
// it as a sanity check to avoid propagating malformed proofs downstream.
func verifyWithdrawalCredentialsInBlock(
	forkVersion common.Version,
	bbh *ctypes.BeaconBlockHeader,
	validatorOffset math.U64,
	proof []common.Root,
	leaf common.Root,
) (common.Root, error) {
	zeroWithdrawalGIndexBlock, err := GetZeroValidatorCredentialsGIndexBlock(forkVersion)
	if err != nil {
		return common.Root{}, err
	}

	beaconRootBytes, err := bbh.HashTreeRoot()
	if err != nil {
		return common.Root{}, err
	}
	beaconRoot := common.NewRootFromBytes(beaconRootBytes[:])
	if !merkle.VerifyProof(
		beaconRoot,
		leaf,
		zeroWithdrawalGIndexBlock+validatorOffset.Unwrap(),
		proof,
	) {
		return common.Root{}, errors.Wrapf(
			errors.New("withdrawal credentials proof failed to verify against beacon root"),
			"beacon root: 0x%s", beaconRoot,
		)
	}

	return beaconRoot, nil
}
