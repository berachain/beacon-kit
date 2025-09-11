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

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/node-api/handlers/proof/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/merkle"
)

// ProveProposerPubkeyInBlock generates a proof for the proposer pubkey in the
// beacon block.
func ProveProposerPubkeyInBlock(
	bbh *ctypes.BeaconBlockHeader,
	bsm types.BeaconStateMarshallable,
) ([]common.Root, common.Root, error) {
	return ProveValidatorPubkeyInBlock(bbh.GetProposerIndex(), bbh, bsm)
}

// ProveValidatorPubkeyInBlock generates a proof for a validator's pubkey in the
// beacon block. The proof is verified against the beacon block root as a sanity
// check and the beacon block root is returned alongside the proof.
func ProveValidatorPubkeyInBlock(
	validatorIndex math.U64,
	bbh *ctypes.BeaconBlockHeader,
	bsm types.BeaconStateMarshallable,
) ([]common.Root, common.Root, error) {
	forkVersion := bsm.GetForkVersion()

	// Calculate the validator-specific offset.
	validatorOffset := ValidatorGIndexOffset * validatorIndex

	// 1. Proof inside the state.
	pubkeyInStateProof, leaf, err := ProveValidatorPubkeyInState(
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

	// 3. Combine proofs: state-level hashes come first, followed by block-level hashes.
	//
	//nolint:gocritic // ok.
	combinedProof := append(pubkeyInStateProof, stateInBlockProof...)

	// 4. Verify the combined proof against the beacon block root.
	beaconRoot, err := verifyPubkeyInBlock(
		forkVersion, bbh, validatorOffset, combinedProof, leaf,
	)
	if err != nil {
		return nil, common.Root{}, err
	}

	return combinedProof, beaconRoot, nil
}

// ProveValidatorPubkeyInState generates a proof for a validator's pubkey
// in the beacon state. The validatorOffset must be computed as
// (ValidatorGIndexOffset * validatorIndex).
func ProveValidatorPubkeyInState(
	forkVersion common.Version,
	bsm types.BeaconStateMarshallable,
	validatorOffset math.U64,
) ([]common.Root, common.Root, error) {
	stateProofTree, err := bsm.GetTree()
	if err != nil {
		return nil, common.Root{}, err
	}

	// Determine the correct gIndex based on the fork version.
	gIndex := int(validatorOffset) // #nosec G115 -- max validator offset is 8 * (2^40 - 1).
	zeroValidatorPubkeyGIndexState, err := GetZeroValidatorPubkeyGIndexState(forkVersion)
	if err != nil {
		return nil, common.Root{}, err
	}
	gIndex += zeroValidatorPubkeyGIndexState

	valPubkeyInStateProof, err := stateProofTree.Prove(gIndex)
	if err != nil {
		return nil, common.Root{}, err
	}

	proof := make([]common.Root, len(valPubkeyInStateProof.Hashes))
	for i, hash := range valPubkeyInStateProof.Hashes {
		proof[i] = common.NewRootFromBytes(hash)
	}
	return proof, common.NewRootFromBytes(valPubkeyInStateProof.Leaf), nil
}

// verifyPubkeyInBlock verifies the validator pubkey in the beacon block,
// returning the beacon block root used to verify against.
//
// TODO: verifying the proof is not absolutely necessary.
func verifyPubkeyInBlock(
	forkVersion common.Version,
	bbh *ctypes.BeaconBlockHeader,
	valOffset math.U64,
	proof []common.Root,
	leaf common.Root,
) (common.Root, error) {
	zeroValidatorPubkeyGIndexBlock, err := GetZeroValidatorPubkeyGIndexBlock(forkVersion)
	if err != nil {
		return common.Root{}, err
	}

	beaconRoot := bbh.HashTreeRoot()
	if !merkle.VerifyProof(
		beaconRoot, leaf, zeroValidatorPubkeyGIndexBlock+valOffset.Unwrap(), proof,
	) {
		return common.Root{}, fmt.Errorf(
			"proposer pubkey proof failed to verify against beacon root: %s", beaconRoot,
		)
	}

	return beaconRoot, nil
}
