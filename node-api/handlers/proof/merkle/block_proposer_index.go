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
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/node-api/handlers/proof/types"
	"github.com/berachain/beacon-kit/primitives/pkg/common"
	"github.com/berachain/beacon-kit/primitives/pkg/encoding/ssz/merkle"
)

// ProveProposerIndexInBlock generates a proof for the proposer index in the
// beacon block. The proof is then verified against the beacon block root as a
// sanity check. Returns the proof along with the beacon block root. It uses
// the fastssz library to generate the proof.
func ProveProposerIndexInBlock[
	BeaconBlockHeaderT types.BeaconBlockHeader,
](bbh BeaconBlockHeaderT) ([]common.Root, common.Root, error) {
	blockProofTree, err := bbh.GetTree()
	if err != nil {
		return nil, common.Root{}, err
	}

	proposerIndexProof, err := blockProofTree.Prove(
		ProposerIndexGIndexDenebBlock,
	)
	if err != nil {
		return nil, common.Root{}, err
	}

	proof := make([]common.Root, len(proposerIndexProof.Hashes))
	for i, hash := range proposerIndexProof.Hashes {
		proof[i] = common.NewRootFromBytes(hash)
	}

	beaconRoot, err := verifyProposerIndexInBlock(
		bbh, proof, common.NewRootFromBytes(proposerIndexProof.Leaf),
	)
	if err != nil {
		return nil, common.Root{}, err
	}

	return proof, beaconRoot, nil
}

// verifyProposerIndexInBlock verifies the proposer index proof in the block.
//
// TODO: verifying the proof is not absolutely necessary.
func verifyProposerIndexInBlock(
	bbh types.BeaconBlockHeader, proof []common.Root, leaf common.Root,
) (common.Root, error) {
	beaconRoot := bbh.HashTreeRoot()
	if beaconRootVerified, err := merkle.VerifyProof(
		ProposerIndexGIndexDenebBlock, leaf, proof, beaconRoot,
	); err != nil {
		return common.Root{}, err
	} else if !beaconRootVerified {
		return common.Root{}, errors.New(
			"proposer index proof failed to verify against beacon root",
		)
	}

	return beaconRoot, nil
}
