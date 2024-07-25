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

package types

import (
	ssz "github.com/ferranbt/fastssz"

	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
)

// GetProofForProposer_FastSSZ generates a proof for the proposer validator in
// the beacon block. The proof is then verified against the beacon block root
// as a sanity check. Returns the proof along with the beacon block root. It
// uses the fastssz library to generate the proof.
//
//nolint:revive,stylecheck // for explicit naming.
func GetProofForProposer_FastSSZ(
	beaconBlock *BeaconBlockForValidator,
) ([]common.Root, common.Root, error) {
	// Get the proof tree to generate the proof.
	proofTree, err := beaconBlock.GetTree()
	if err != nil {
		return nil, common.Root{}, err
	}

	// Get the generalized index of the proposer validator in the tree.
	gIndex := int(ZeroValidatorGIndexDenebPlus + beaconBlock.ProposerIndex)

	// Get the proof of the proposer validator in the tree.
	proof, err := proofTree.Prove(gIndex)
	if err != nil {
		return nil, common.Root{}, err
	}
	validatorProof := make([]common.Root, len(proof.Hashes))
	for i, hash := range proof.Hashes {
		validatorProof[i] = common.Root(hash)
	}

	// Sanity check that the proof verifies against our beacon root.
	beaconRoot, err := beaconBlock.HashTreeRoot()
	if err != nil {
		return nil, common.Root{}, err
	}
	beaconRootVerified, err := ssz.VerifyProof(beaconRoot[:], proof)
	if err != nil {
		return nil, common.Root{}, err
	}
	if !beaconRootVerified {
		return nil, common.Root{}, errors.Newf(
			"proof verification failed against beacon root: %x", beaconRoot[:],
		)
	}

	return validatorProof, beaconRoot, nil
}
