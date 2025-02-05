package main

import (
	"flag"
	"fmt"
)

var (
	debug bool
)

func main() {
	// Define the flags
	seed := flag.Int64("seed", 0, "Seed")
	sel := flag.Uint("sel", 4, "Selector")
	valLen := flag.Uint("valLen", 5, "Validator length")
	debugFlag := flag.Bool("debug", false, "Debug")
	hardcoded := flag.Bool("hardcoded", false, "Use hardcoded values")

	// Parse the flags
	flag.Parse()

	// Set global debug value
	debug = *debugFlag

	// Use hardcoded values when `--hardcoded` is set
	if *hardcoded {
		seedValue := int64(123)
		seed = &seedValue
		selValue := uint(0)
		sel = &selValue
		valLenValue := uint(5)
		valLen = &valLenValue
		debug = true
	}

	// Seed the random number generator
	initializeRNG(*seed)

	// Select the proof to generate based on the inputted selector
	if *sel == 0 {
		_beaconBlockRoot, _proposerIndex, _proposerPubkey, _proposerPubkeyProof := generateValidatorPubkeyProof(*valLen)
		encodeValidatorPubkeyParams(_beaconBlockRoot, _proposerIndex, _proposerPubkey[:], _proposerPubkeyProof)
	} else if *sel == 1 {
		_beaconBlockRoot, _proposerIndex, _proposerPubkeyProof := generateProposerIndexProof(*valLen)
		encodeProposerIndexParams(_beaconBlockRoot, _proposerIndex, _proposerPubkeyProof)
	} else {
		fmt.Println("Invalid selector")
	}
}
