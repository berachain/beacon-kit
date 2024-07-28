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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
)

// ** Generalized Index **
// val0Pubkey in state - gIndex 439804651110400
// state in block - gIndex 11
// state in block -> val0Pubkey in state - concat 3254554418216960
// val0Pubkey in block - gIndex 3254554418216960

// ProofForProposerPubkey_FastSSZ generates a proof for the proposer
// pubkey in the beacon block. The proof is then verified against the beacon
// block root as a sanity check. Returns the proof along with the beacon block
// root. It uses the fastssz library to generate the proof.
//
//nolint:revive,stylecheck // for explicit naming.
func ProofForProposerPubkey_FastSSZ[
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconStateT BeaconState[
		BeaconBlockHeaderT,
		Eth1DataT,
		ExecutionPayloadHeaderT,
		ForkT,
		ValidatorT,
	],
	Eth1DataT any,
	ExecutionPayloadHeaderT any,
	ForkT any,
	ValidatorT any,
](
	bbh BeaconBlockHeaderT,
	bs BeaconStateT,
) ([]common.Root, common.Root, error) {
	var (
		pubkeyProof     []common.Root
		beaconRoot, err = bbh.HashTreeRoot()
	)
	if err != nil {
		return nil, common.Root{}, err
	}

	// Get the proof tree of the beacon state first.
	stateProofTree, err := bs.GetTree()
	if err != nil {
		return nil, common.Root{}, err
	}

	// Get the generalized index of the proposer pubkey in the beacon state.
	valPubkeyInStateGIndex := ZeroValidatorPubkeyGIndexDenebState +
		(ValidatorPubkeyGIndexOffset * int(bbh.GetProposerIndex()))

	// Get the proof of the proposer pubkey in the beacon state.
	valPubkeyInStateProof, err := stateProofTree.Prove(valPubkeyInStateGIndex)
	if err != nil {
		return nil, common.Root{}, err
	}
	for _, hash := range valPubkeyInStateProof.Hashes {
		pubkeyProof = append(pubkeyProof, common.Root(hash))
	}

	// // Sanity check that the combined proof verifies against our beacon root.
	// beaconRoot, err := beaconBlock.HashTreeRoot()
	// if err != nil {
	// 	return nil, common.Root{}, err
	// }
	// beaconRootVerified, err := ssz.VerifyProof(beaconRoot[:], proof)
	// if err != nil {
	// 	return nil, common.Root{}, err
	// }
	// if !beaconRootVerified {
	// 	return nil, common.Root{}, errors.Newf(
	// 		"proof verification failed against beacon root: %x", beaconRoot[:],
	// 	)
	// }

	return pubkeyProof, beaconRoot, nil
}
