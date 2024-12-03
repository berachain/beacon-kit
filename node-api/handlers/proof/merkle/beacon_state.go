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

// ProveBeaconStateInBlock generates a proof for the beacon state in the
// beacon block. It uses the fastssz library to generate the proof.
func ProveBeaconStateInBlock(
	bbh types.BeaconBlockHeader, verifyProof bool,
) ([]common.Root, error) {
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
		proof[i] = common.NewRootFromBytes(hash)
	}

	if verifyProof {
		if err = verifyBeaconStateInBlock(
			bbh, proof, common.NewRootFromBytes(stateInBlockProof.Leaf),
		); err != nil {
			return nil, err
		}
	}

	return proof, nil
}

// verifyBeaconStateInBlock verifies the beacon state proof in the block.
//
// TODO: verifying the proof is not absolutely necessary.
func verifyBeaconStateInBlock(
	bbh types.BeaconBlockHeader, proof []common.Root, leaf common.Root,
) error {
	beaconRoot := bbh.HashTreeRoot()
	if beaconRootVerified, err := merkle.VerifyProof(
		StateGIndexDenebBlock, leaf, proof, beaconRoot,
	); err != nil {
		return err
	} else if !beaconRootVerified {
		return errors.New(
			"beacon state proof failed to verify against beacon root",
		)
	}
	return nil
}
