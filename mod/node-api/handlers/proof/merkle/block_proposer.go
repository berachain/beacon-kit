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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// ProveProposerInBlock generates a proof for the proposer pubkey in the
// beacon block. The proof is then verified against the beacon block root as a
// sanity check. Returns the proof along with the beacon block root. It uses
// the fastssz library to generate the proof.
func ProveProposerInBlock[
	BeaconBlockHeaderT types.BeaconBlockHeader,
	BeaconStateMarshallableT types.BeaconStateMarshallable,
	ExecutionPayloadHeaderT types.ExecutionPayloadHeader,
	ValidatorT any,
](
	bbh BeaconBlockHeaderT,
	bs types.BeaconState[
		BeaconStateMarshallableT, ExecutionPayloadHeaderT, ValidatorT,
	],
) ([]common.Root, common.Root, error) {
	// Get the proof of the proposer pubkey in the beacon state.
	proposerOffset := ValidatorPubkeyGIndexOffset * bbh.GetProposerIndex()
	valPubkeyInStateProof, leaf, err := ProveProposerPubkeyInState(
		bs, proposerOffset,
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
		bbh, proposerOffset, combinedProof, leaf,
	)
	if err != nil {
		return nil, common.Root{}, err
	}

	return combinedProof, beaconRoot, nil
}

// ProveProposerPubkeyInState generates a proof for the proposer pubkey
// in the beacon state. It uses the fastssz library to generate the proof.
func ProveProposerPubkeyInState[
	BeaconStateMarshallableT types.BeaconStateMarshallable,
	ExecutionPayloadHeaderT types.ExecutionPayloadHeader,
	ValidatorT any,
](
	bs types.BeaconState[
		BeaconStateMarshallableT, ExecutionPayloadHeaderT, ValidatorT,
	],
	proposerOffset math.U64,
) ([]common.Root, common.Root, error) {
	bsm, err := bs.GetMarshallable()
	if err != nil {
		return nil, common.Root{}, err
	}
	stateProofTree, err := bsm.GetTree()
	if err != nil {
		return nil, common.Root{}, err
	}

	//#nosec:G701 // max proposer offset is 8 * (2^40 - 1).
	gIndex := ZeroValidatorPubkeyGIndexDenebState + int(proposerOffset)
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
	bbh types.BeaconBlockHeader,
	valOffset math.U64,
	proof []common.Root,
	leaf common.Root,
) (common.Root, error) {
	beaconRoot := bbh.HashTreeRoot()
	if beaconRootVerified, err := merkle.VerifyProof(
		merkle.GeneralizedIndex(ZeroValidatorPubkeyGIndexDenebBlock+valOffset),
		leaf, proof, beaconRoot,
	); err != nil {
		return common.Root{}, err
	} else if !beaconRootVerified {
		return common.Root{}, errors.New(
			"proof failed to verify against beacon root",
		)
	}

	return beaconRoot, nil
}
