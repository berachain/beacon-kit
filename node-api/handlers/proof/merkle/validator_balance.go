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
	"github.com/pkg/errors"
)

// bytesPerBalance is the number of bytes in a single balance (uint64).
const bytesPerBalance uint64 = 8

// ProveBalanceInState generates a proof for a validator's balance in the beacon state.
func ProveBalanceInState(
	forkVersion common.Version,
	bsm types.BeaconStateMarshallable,
	validatorIndex math.U64,
) ([]common.Root, common.Root, error) {
	stateProofTree, err := bsm.GetTree()
	if err != nil {
		return nil, common.Root{}, err
	}

	// Determine the starting generalized index for the 0-th validator's
	// balance for this fork.
	zeroBalanceGIndexState, err := GetZeroValidatorBalanceGIndexState(forkVersion)
	if err != nil {
		return nil, common.Root{}, err
	}

	// Since balances are packed 4 per leaf, calculate the leaf offset
	leafOffset := validatorIndex / BalancesPerLeaf

	// Calculate the generalized index for the target validator's balance leaf.
	// The offset multiplication is bounded by the number of validators, so
	// converting to int is safe on 64-bit architectures.
	gIndex := zeroBalanceGIndexState + int(leafOffset) // #nosec G115

	balanceProof, err := stateProofTree.Prove(gIndex)
	if err != nil {
		return nil, common.Root{}, err
	}

	proof := make([]common.Root, len(balanceProof.Hashes))
	for i, hash := range balanceProof.Hashes {
		proof[i] = common.NewRootFromBytes(hash)
	}

	// The leaf contains 4 packed uint64 balances
	return proof, common.NewRootFromBytes(balanceProof.Leaf), nil
}

// ProveBalanceInBlock generates a proof for a validator's balance in the beacon block.
// Returns the proof, the leaf containing the packed balances, and the beacon block root.
func ProveBalanceInBlock(
	validatorIndex math.U64,
	bbh *ctypes.BeaconBlockHeader,
	bsm types.BeaconStateMarshallable,
	allBalances []uint64,
) ([]common.Root, common.Root, common.Root, error) {
	forkVersion := bsm.GetForkVersion()

	// 1. Proof inside the state.
	balanceInStateProof, leaf, err := ProveBalanceInState(
		forkVersion, bsm, validatorIndex,
	)
	if err != nil {
		return nil, common.Root{}, common.Root{}, err
	}

	// 2. Build the balance leaf and assert that it matches the proof's leaf.
	builtLeaf := buildBalanceLeaf(allBalances, validatorIndex)
	if !leaf.Equals(builtLeaf) {
		return nil, common.Root{}, common.Root{}, fmt.Errorf(
			"balance leaf mismatch -- proof tree leaf: 0x%s, built leaf: 0x%s", leaf, builtLeaf,
		)
	}

	// 3. Proof of the state inside the block.
	stateInBlockProof, err := ProveBeaconStateInBlock(bbh, false)
	if err != nil {
		return nil, common.Root{}, common.Root{}, err
	}

	// 4. Combine proofs: state-level hashes come first, followed by block-level
	// hashes (same order as ProveProposerPubkeyInBlock).
	//
	//nolint:gocritic // ok.
	combinedProof := append(balanceInStateProof, stateInBlockProof...)

	// 5. Verify the combined proof against the beacon block root.
	// Since balances are packed 4 per leaf, calculate the leaf offset.
	leafOffset := validatorIndex / BalancesPerLeaf
	beaconRoot, err := verifyBalanceInBlock(
		forkVersion, bbh, leafOffset.Unwrap(), combinedProof, leaf,
	)
	if err != nil {
		return nil, common.Root{}, common.Root{}, err
	}

	return combinedProof, leaf, beaconRoot, nil
}

// buildBalanceLeaf constructs the 32-byte leaf containing the packed balances
// for the group of validators that includes `validatorIndex`. Balances are
// packed 4 per leaf (little-endian uint64s).
func buildBalanceLeaf(allBalances []uint64, validatorIndex math.U64) common.Root {
	var leafBytes common.Root

	// Determine which leaf the validator belongs to and the starting index in
	// the balances slice.
	leafIndex := validatorIndex / BalancesPerLeaf
	startIdx := leafIndex * BalancesPerLeaf

	// Pack up to 4 balances (little-endian) into the 32-byte array.
	for i := range uint64(BalancesPerLeaf) {
		idx := startIdx.Unwrap() + i
		if idx >= uint64(len(allBalances)) {
			break
		}
		bal := allBalances[idx]
		for j := range bytesPerBalance {
			leafBytes[i*bytesPerBalance+j] = byte(bal >> (j * bytesPerBalance))
		}
	}

	return leafBytes
}

// verifyBalanceInBlock verifies the provided Merkle proof of a
// validator's balance inside the beacon block and returns the
// beacon block root that the proof was verified against.
//
// NOTE: Proof verification is not strictly necessary for operation, but we do
// it as a sanity check to avoid propagating malformed proofs downstream.
func verifyBalanceInBlock(
	forkVersion common.Version,
	bbh *ctypes.BeaconBlockHeader,
	balanceOffset uint64,
	proof []common.Root,
	leaf common.Root,
) (common.Root, error) {
	zeroBalanceGIndexBlock, err := GetZeroValidatorBalanceGIndexBlock(forkVersion)
	if err != nil {
		return common.Root{}, err
	}

	beaconRoot := bbh.HashTreeRoot()
	if !merkle.VerifyProof(
		beaconRoot,
		leaf,
		zeroBalanceGIndexBlock+balanceOffset,
		proof,
	) {
		return common.Root{}, errors.Wrapf(
			errors.New("balance proof failed to verify against beacon root"),
			"beacon root: 0x%s", beaconRoot,
		)
	}

	return beaconRoot, nil
}
