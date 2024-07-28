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

package merkle

import (
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/node-api/handlers/proof/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/merkle"
)

// ProveProposerPubkey_FastSSZ generates a proof for the proposer pubkey in the
// beacon block. The proof is then verified against the beacon block root as a
// sanity check. Returns the proof along with the beacon block root. It uses
// the fastssz library to generate the proof.
//
//nolint:revive,stylecheck // for explicit naming.
func ProveProposerPubkey_FastSSZ[
	BeaconBlockHeaderT types.BeaconBlockHeader,
	BeaconStateT types.BeaconState[BeaconStateMarshallableT, ValidatorT],
	BeaconStateMarshallableT types.BeaconStateMarshallable,
	ValidatorT any,
](bbh BeaconBlockHeaderT, bs BeaconStateT) ([]common.Root, common.Root, error) {
	// Get the proof of the proposer pubkey in the beacon state.
	proposerOffset := ValidatorPubkeyGIndexOffset * int(bbh.GetProposerIndex())
	valPubkeyInStateProof, leaf, err := ProveProposerPubkeyInState_FastSSZ(
		bs, proposerOffset,
	)
	if err != nil {
		return nil, common.Root{}, err
	}

	// Then get the proof of the beacon state in the beacon block.
	stateInBlockProof, err := ProveStateInBlock_FastSSZ(bbh)
	if err != nil {
		return nil, common.Root{}, err
	}

	// Sanity check that the combined proof verifies against our beacon root.
	combinedProof := append(valPubkeyInStateProof, stateInBlockProof...)
	beaconRoot, err := verifyProposerInBlock(
		bbh, proposerOffset, combinedProof, leaf,
	)
	if err != nil {
		return nil, common.Root{}, err
	}

	return combinedProof, beaconRoot, nil
}

// ProveProposerPubkeyInState_FastSSZ generates a proof for the proposer pubkey
// in the beacon state. It uses the fastssz library to generate the proof.
func ProveProposerPubkeyInState_FastSSZ[
	BeaconStateT types.BeaconState[BeaconStateMarshallableT, ValidatorT],
	BeaconStateMarshallableT types.BeaconStateMarshallable,
	ValidatorT any,
](bs BeaconStateT, proposerOffset int) ([]common.Root, common.Root, error) {
	bsm, err := bs.GetMarshallable()
	if err != nil {
		return nil, common.Root{}, err
	}
	stateProofTree, err := bsm.GetTree()
	if err != nil {
		return nil, common.Root{}, err
	}

	gIndex := ZeroValidatorPubkeyGIndexDenebState + proposerOffset
	valPubkeyInStateProof, err := stateProofTree.Prove(gIndex)
	if err != nil {
		return nil, common.Root{}, err
	}

	proof := make([]common.Root, len(valPubkeyInStateProof.Hashes))
	for i, hash := range valPubkeyInStateProof.Hashes {
		proof[i] = common.Root(hash)
	}
	return proof, common.Root(valPubkeyInStateProof.Leaf), nil
}

// ProveStateInBlock_FastSSZ generates a proof for the beacon state in the
// beacon block. It uses the fastssz library to generate the proof.
func ProveStateInBlock_FastSSZ[
	BeaconBlockHeaderT types.BeaconBlockHeader,
](bbh BeaconBlockHeaderT) ([]common.Root, error) {
	blockProofTree, err := bbh.GetTree()
	if err != nil {
		return nil, err
	}

	stateInBlockProof, err := blockProofTree.Prove(StateGIndexDenebBlock)
	if err != nil {
		return nil, err
	}

	proof := make([]common.Root, len(stateInBlockProof.Hashes))
	for i, hash := range stateInBlockProof.Hashes {
		proof[i] = common.Root(hash)
	}
	return proof, nil
}

// verifyProposerInBlock verifies the proposer pubkey in the beacon block,
// returning the beacon block root used to verify against.
func verifyProposerInBlock[BeaconBlockHeaderT types.BeaconBlockHeader](
	bbh BeaconBlockHeaderT, offset int, proof []common.Root, leaf common.Root,
) (common.Root, error) {
	beaconRoot, err := bbh.HashTreeRoot()
	if err != nil {
		return common.Root{}, err
	}

	if beaconRootVerified, err := merkle.VerifyProof(
		merkle.GeneralizedIndex(ZeroValidatorPubkeyGIndexDenebBlock+offset),
		leaf, proof, beaconRoot,
	); err != nil {
		return common.Root{}, err
	} else if !beaconRootVerified {
		return common.Root{}, errors.Newf(
			"proof verification failed against beacon root: %x", beaconRoot[:],
		)
	}

	return beaconRoot, nil
}
