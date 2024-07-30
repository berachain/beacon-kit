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

// ProveExecutionNumberInBlock generates a proof for the block number of the
// execution payload in the beacon block. The proof is then verified against
// the beacon block root as a sanity check. Returns the proof along with the
// beacon block root. It uses the fastssz library to generate the proof.
func ProveExecutionNumberInBlock[
	BeaconBlockHeaderT types.BeaconBlockHeader,
	BeaconStateMarshallableT types.BeaconStateMarshallable,
	ExecutionPayloadHeaderT types.ExecutionPayloadHeader,
	ValidatorT any,
](
	bbh BeaconBlockHeaderT,
	bs types.BeaconState[
		BeaconBlockHeaderT, BeaconStateMarshallableT, ExecutionPayloadHeaderT,
		ValidatorT,
	],
) ([]common.Root, common.Root, error) {
	// Get the proof of the execution number in the beacon state.
	numberInStateProof, leaf, err := ProveExecutionNumberInState(bs)
	if err != nil {
		return nil, common.Root{}, err
	}

	// Then get the proof of the beacon state in the beacon block.
	stateInBlockProof, err := ProveStateInBlock(bbh)
	if err != nil {
		return nil, common.Root{}, err
	}

	// Sanity check that the combined proof verifies against our beacon root.
	//
	//nolint:gocritic // ok.
	combinedProof := append(numberInStateProof, stateInBlockProof...)
	beaconRoot, err := verifyExecutionNumberInBlock(bbh, combinedProof, leaf)
	if err != nil {
		return nil, common.Root{}, err
	}

	return combinedProof, beaconRoot, nil
}

// ProveExecutionNumberInState generates a proof for the block number of the
// execution payload in the beacon state. It uses the fastssz library.
func ProveExecutionNumberInState[
	BeaconBlockHeaderT types.BeaconBlockHeader,
	BeaconStateMarshallableT types.BeaconStateMarshallable,
	ExecutionPayloadHeaderT types.ExecutionPayloadHeader,
	ValidatorT any,
](
	bs types.BeaconState[
		BeaconBlockHeaderT, BeaconStateMarshallableT, ExecutionPayloadHeaderT,
		ValidatorT,
	],
) ([]common.Root, common.Root, error) {
	bsm, err := bs.GetMarshallable()
	if err != nil {
		return nil, common.Root{}, err
	}
	stateProofTree, err := bsm.GetTree()
	if err != nil {
		return nil, common.Root{}, err
	}

	numberInStateProof, err := stateProofTree.Prove(
		ExecutionPayloadNumberGIndexDenebState,
	)
	if err != nil {
		return nil, common.Root{}, err
	}

	proof := make([]common.Root, len(numberInStateProof.Hashes))
	for i, hash := range numberInStateProof.Hashes {
		proof[i] = common.Root(hash)
	}
	return proof, common.Root(numberInStateProof.Leaf), nil
}

// verifyExecutionNumberInBlock verifies the execution number in the beacon
// block, returning the beacon block root used to verify against.
//
// TODO: verifying the proof is not absolutely necessary.
func verifyExecutionNumberInBlock(
	bbh types.BeaconBlockHeader,
	proof []common.Root,
	leaf common.Root,
) (common.Root, error) {
	beaconRoot, err := bbh.HashTreeRoot()
	if err != nil {
		return common.Root{}, err
	}

	var beaconRootVerified bool
	if beaconRootVerified, err = merkle.VerifyProof(
		ExecutionPayloadNumberGIndexDenebBlock,
		leaf, proof, beaconRoot,
	); err != nil {
		return common.Root{}, err
	} else if !beaconRootVerified {
		return common.Root{}, errors.Newf(
			"proof failed to verify against beacon root: 0x%x", beaconRoot[:],
		)
	}

	return beaconRoot, nil
}
