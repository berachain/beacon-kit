package randao

import (
	"github.com/itsdevbear/bolaris/beacon/state"
	"github.com/itsdevbear/bolaris/types/consensus/primitives"
)

// This is the internal representation of the randao reveal
// Although it is 32 bytes now, it can change
// We use the same size as Ed25519 sig
// TODO: update to 96 bytes when moving to BLS
type RandaoReveal = [32]byte

// This is the external representation of the randao random number
// We fix this to 32 bytes
type RandaoRandom = [32]byte

type RandaoCompatible interface {
	// Given a slot, derive the current epoch number
	// This is used for generating the provably random number
	// TODO: add other necessary args
	GetEpochBySlot(primitives.Slot) primitives.Epoch

	// Get the historical randao reveals
	// For now track everything, but we probably want to only track
	// up to some number based on config
	// TODO: is this the right input struct arg?
	GetRandaoReveals(state.ReadOnlyBeaconState) []RandaoReveal

	// This is done by signing the epoch number
	// TODO: add other necessary args
	GenerateRandaoReveal(state.ReadOnlyBeaconState) (RandaoReveal, error)

	// Update the state with the calculated reveal
	SetRandaoReveal(state.ReadOnlyBeaconState, RandaoReveal) error

	// Given a randao reveal and the pub key of a validator, validatate the reveal
	// It should be the signature over the epoch
	VerifyRandaoReveal(state.ReadOnlyBeaconState, RandaoReveal, PubKey) bool

	// This mixes in the randao reveal and updates the internal state
	MixAndUpdateRandaoReveal(*state.ReadOnlyBeaconState, RandaoReveal) error
}

/*

How it works:
1. Each epoch, validators compute their randaoReveal which is a sig over the epoch number
	Uses: GetEpochBySlot, GenerateRandaoReveal, SetRandaoReveal
2. On each block the proposer proposes its randaoReveal, which can be verified (but not predicted) by other validators
	Uses: VerifyRandaoReveal
3. When the block is processed, the valid randaoReveal will be mixed in with the chain's randao value to get a new value
4. The updating is done by hashing the randaoReveal and XOR-ing it with the current randao value
	Uses: MixAndUpdateRandaoReveal
5. We might want to track historical randao mixes for various reasons
	Uses: GetRandaoReveals
*/
