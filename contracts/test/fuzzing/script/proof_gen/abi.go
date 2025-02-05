package main

import (
	"fmt"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Encode the parameters for the verifyValidatorPubkeyInBeaconBlock function in BeaconVerifier.sol
func encodeValidatorPubkeyParams(_beaconBlockRoot [32]byte, _proposerIndex math.U64, _proposerPubkey []byte, _proposerPubkeyProof []common.Root) {
	// Define the ABI types
	beaconBlockRootType, err := abi.NewType("bytes32", "", nil)
	if err != nil {
		fmt.Println("Error defining beaconBlockRootType:", err)
		return
	}

	proposerIndexType, err := abi.NewType("uint64", "", nil)
	if err != nil {
		fmt.Println("Error defining proposerIndexType:", err)
		return
	}

	proposerPubkeyType, err := abi.NewType("bytes", "", nil)
	if err != nil {
		fmt.Println("Error defining proposerPubkeyType:", err)
		return
	}

	proposerPubkeyProofType, err := abi.NewType("bytes32[]", "", nil)
	if err != nil {
		fmt.Println("Error defining proposerPubkeyProofType:", err)
		return
	}

	// Create the arguments slice
	arguments := abi.Arguments{
		{
			Name:    "beaconBlockRoot",
			Type:    beaconBlockRootType,
			Indexed: false,
		},
		{
			Name:    "proposerIndex",
			Type:    proposerIndexType,
			Indexed: false,
		},
		{
			Name:    "proposerPubkey",
			Type:    proposerPubkeyType,
			Indexed: false,
		},
		{
			Name:    "proposerPubkeyProof",
			Type:    proposerPubkeyProofType,
			Indexed: false,
		},
	}

	// Prepare the values in the order of arguments
	values := []interface{}{
		_beaconBlockRoot,
		_proposerIndex,
		_proposerPubkey,
		_proposerPubkeyProof,
	}

	// Pack the values
	packed, err := arguments.Pack(values...)
	if err != nil {
		fmt.Println("Error packing values:", err)
		return
	}

	// Output the packed bytes in hex format
	if !debug {
		fmt.Print(hexutil.Encode(packed))
	}
}

// Encode the parameters for the verifyProposerIndexInBeaconBlock function in BeaconVerifier.sol
func encodeProposerIndexParams(beaconBlockRoot common.Root, proposerIndex math.U64, proposerIndexProof []common.Root) {
	// Define the ABI types
	beaconBlockRootType, err := abi.NewType("bytes32", "", nil)
	if err != nil {
		fmt.Println("Error defining beaconBlockRootType:", err)
		return
	}

	proposerIndexType, err := abi.NewType("uint64", "", nil)
	if err != nil {
		fmt.Println("Error defining proposerIndexType:", err)
		return
	}

	proposersIndexProofType, err := abi.NewType("bytes32[]", "", nil)
	if err != nil {
		fmt.Println("Error defining proposersIndexProofType:", err)
		return
	}

	// Create the arguments slice
	arguments := abi.Arguments{
		{
			Name:    "beaconBlockRoot",
			Type:    beaconBlockRootType,
			Indexed: false,
		},
		{
			Name:    "proposerIndex",
			Type:    proposerIndexType,
			Indexed: false,
		},
		{
			Name:    "proposersIndexProof",
			Type:    proposersIndexProofType,
			Indexed: false,
		},
	}

	// Prepare the values in the order of arguments
	values := []interface{}{
		beaconBlockRoot,
		proposerIndex,
		proposerIndexProof,
	}

	// Pack the values
	packed, err := arguments.Pack(values...)
	if err != nil {
		fmt.Println("Error packing values:", err)
		return
	}

	// Output the packed bytes in hex format
	if !debug {
		fmt.Print(hexutil.Encode(packed))
	}
}
