package main

import (
	"fmt"

	"github.com/berachain/beacon-kit/mod/node-api/handlers/proof/merkle"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// Generate a proof with some random values to be validated by the verifyValidatorPubkeyInBeaconBlock function in BeaconVerifier.sol
func generateValidatorPubkeyProof(validatorCount uint) (common.Root, math.U64, [48]byte, []common.Root) {
	// Generate beacon state and block header
	beaconState := rndBeaconState(validatorCount)
	beaconBlockHeader := rndBeaconBlockHeaderFromState(beaconState)

	// Get the proof of the proposer pubkey in the beacon state.
	proposerOffset := merkle.ValidatorPubkeyGIndexOffset * beaconBlockHeader.GetProposerIndex()
	valPubkeyInStateProof, _, err := merkle.ProveProposerPubkeyInState(
		beaconState, proposerOffset,
	)
	if err != nil {
		panic(err)
	}

	// Then get the proof of the beacon state in the beacon block.
	stateInBlockProof, err := merkle.ProveBeaconStateInBlock(beaconBlockHeader, true)
	if err != nil {
		panic(err)
	}

	combinedProof := append(valPubkeyInStateProof, stateInBlockProof...)

	beaconBlockRoot := beaconBlockHeader.HashTreeRoot()
	validatorIndex := beaconBlockHeader.GetProposerIndex()
	validatorPubkey := beaconState.Validators[beaconBlockHeader.GetProposerIndex()].Pubkey

	if debug {
		fmt.Println("Proof")
		for _, root := range combinedProof {
			fmt.Println("  ", root)
		}
		fmt.Println("Beacon block root", beaconBlockRoot)
		fmt.Println("Validator index", validatorIndex)
		fmt.Println("Validator pubkey", validatorPubkey)
	}

	return beaconBlockRoot, validatorIndex, validatorPubkey, combinedProof
}

// Generate a proof with some random values to be validated by the verifyProposerIndexInBeaconBlock function in BeaconVerifier.sol
func generateProposerIndexProof(validatorCount uint) (common.Root, math.U64, []common.Root) {
	// Generate beacon state and block header
	beaconState := rndBeaconState(validatorCount)
	beaconBlockHeader := rndBeaconBlockHeaderFromState(beaconState)

	// Get the proof of the proposer index in the beacon state.
	proposerIndexProof, _, err := merkle.ProveProposerIndexInBlock(beaconBlockHeader)
	if err != nil {
		panic(err)
	}

	beaconBlockRoot := beaconBlockHeader.HashTreeRoot()
	validatorIndex := beaconBlockHeader.GetProposerIndex()

	if debug {
		fmt.Println("Proof")
		for _, root := range proposerIndexProof {
			fmt.Println("  ", root)
		}
		fmt.Println("Beacon block root", beaconBlockRoot)
		fmt.Println("Validator index", validatorIndex)
	}

	return beaconBlockRoot, validatorIndex, proposerIndexProof
}
