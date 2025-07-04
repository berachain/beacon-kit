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
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/node-api/handlers/proof/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/merkle"
)

// ProveProposerPubkeyInBlock generates a proof for the proposer pubkey in the
// beacon block. The proof is then verified against the beacon block root as a
// sanity check. Returns the proof along with the beacon block root. It uses
// the fastssz library to generate the proof.
func ProveProposerPubkeyInBlock(
	bbh *ctypes.BeaconBlockHeader,
	bsm types.BeaconStateMarshallable,
) ([]common.Root, common.Root, error) {
	forkVersion := bsm.GetForkVersion()

	// Get the proof of the proposer pubkey in the beacon state.
	proposerOffset := ValidatorGIndexOffset * bbh.GetProposerIndex()
	valPubkeyInStateProof, leaf, err := ProveProposerPubkeyInState(
		forkVersion, bsm, proposerOffset,
	)
	if err != nil {
		return nil, common.Root{}, err
	}

	// Then get the proof of the beacon state in the beacon block.
	stateInBlockProof, err := ProveBeaconStateInBlock(bbh, false)
	if err != nil {
		return nil, common.Root{}, err
	}

	// Sanity check that the combined proof verifies against our beacon root.
	//
	//nolint:gocritic // ok.
	combinedProof := append(valPubkeyInStateProof, stateInBlockProof...)
	beaconRoot, err := verifyProposerInBlock(
		forkVersion, bbh, proposerOffset, combinedProof, leaf,
	)
	if err != nil {
		return nil, common.Root{}, err
	}

	return combinedProof, beaconRoot, nil
}

// ProveProposerPubkeyInState generates a proof for the proposer pubkey
// in the beacon state. It uses the fastssz library to generate the proof.
func ProveProposerPubkeyInState(
	forkVersion common.Version,
	bsm types.BeaconStateMarshallable,
	proposerOffset math.U64,
) ([]common.Root, common.Root, error) {
	stateProofTree, err := bsm.GetTree()
	if err != nil {
		return nil, common.Root{}, err
	}

	// Determine the correct gIndex based on the fork version.
	gIndex := int(proposerOffset) // #nosec G115 -- max proposer offset is 8 * (2^40 - 1).
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

// verifyProposerInBlock verifies the proposer pubkey in the beacon block,
// returning the beacon block root used to verify against.
//
// TODO: verifying the proof is not absolutely necessary.
func verifyProposerInBlock(
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

	beaconRootBytes, err := bbh.HashTreeRoot()
	if err != nil {
		return common.Root{}, err
	}
	beaconRoot := common.NewRootFromBytes(beaconRootBytes[:])
	if !merkle.VerifyProof(
		beaconRoot, leaf, zeroValidatorPubkeyGIndexBlock+valOffset.Unwrap(), proof,
	) {
		return common.Root{}, errors.Wrapf(
			errors.New("proposer pubkey proof failed to verify against beacon root"),
			"beacon root: 0x%s", beaconRoot,
		)
	}

	return beaconRoot, nil
}
