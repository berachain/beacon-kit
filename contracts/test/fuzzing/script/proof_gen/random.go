package main

import (
	"fmt"
	"math/rand"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/node-api/handlers/proof/merkle/mock"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/holiman/uint256"
)

const (
	MAX_BLOCK_ROOTS  = 8192
	MAX_STATE_ROOTS  = 8192
	MAX_EXTRA_DATA   = 32
	MAX_VALIDATORS   = 10000
	MAX_RANDDAOMIXES = 65536
)

var (
	rng *rand.Rand
)

// Initialize the random number generator with a seed
func initializeRNG(seed int64) {
	source := rand.NewSource(seed)
	rng = rand.New(source)
}

// Generate a random BeaconState
func rndBeaconState(validatorCount uint) *mock.BeaconState {
	execPayloadHeader := rndExecutionPayloadHeader()
	validators := rndValidators(validatorCount)

	var (
		beaconStateMarshallable = &mock.BeaconStateMarshallable{}
		err                     error
	)
	beaconStateMarshallable, err = beaconStateMarshallable.New(
		0,                                    // unused
		rndRoot(),                            // genesisValidatorsRoot (common.Root)
		rndU64(),                             // slot (math.Slot)
		rndFork(),                            // fork (ForkT)
		rndBeaconBlockHeader(),               // latestBlockHeader (BeaconBlockHeaderT)
		rndSlice(MAX_BLOCK_ROOTS, rndRoot),   // blockRoots ([]common.Root)
		rndSlice(MAX_STATE_ROOTS, rndRoot),   // stateRoots ([]common.Root)
		rndEth1Data(),                        // eth1Data (Eth1DataT)
		rng.Uint64(),                         // eth1DepositIndex (uint64)
		execPayloadHeader,                    // latestExecutionPayloadHeader (ExecutionPayloadHeaderT)
		validators,                           // validators ([]ValidatorT)
		rndSlice(MAX_VALIDATORS, rng.Uint64), // balances ([]uint64)
		rndSlice(MAX_RANDDAOMIXES, rndB32),   // randaoMixes ([]common.Bytes32)
		rng.Uint64(),                         // nextWithdrawalIndex (uint64)
		rndU64(),                             // nextWithdrawalValidatorIndex (math.ValidatorIndex)
		rndSlice(MAX_VALIDATORS, rndU64),     // slashings ([]math.Gwei)
		rndU64(),                             // totalSlashing (math.Gwei)
	)

	if err != nil {
		panic(err)
	}

	return &mock.BeaconState{BeaconStateMarshallable: beaconStateMarshallable}
}

// Generate a random Fork
func rndFork() *types.Fork {
	return (&types.Fork{}).New(
		rndB4(),  // PreviousVersion
		rndB4(),  // CurrentVersion
		rndU64(), // Epoch
	)
}

// Generate a random BeaconBlockHeader
func rndBeaconBlockHeader() *types.BeaconBlockHeader {
	return (&types.BeaconBlockHeader{}).New(
		rndU64(),  // slot
		rndU64(),  // proposerIndex
		rndRoot(), // parentBlockRoot
		rndRoot(), // stateRoot
		rndRoot(), // bodyRoot
	)
}

// Generate random Eth1Data
func rndEth1Data() *types.Eth1Data {
	return (&types.Eth1Data{}).New(
		rndRoot(),          // DepositRoot
		rndU64(),           // DepositCount
		rndExecutionHash(), // BlockHash
	)
}

// Generate a random BeaconBlockHeader from a BeaconState
func rndBeaconBlockHeaderFromState(beaconState *mock.BeaconState) *types.BeaconBlockHeader {
	// Pick a random validator as proposer
	validatorIndex := math.U64(rng.Intn(len(beaconState.Validators)))

	// Generate the beacon block header
	stateRoot := beaconState.BeaconStateMarshallable.HashTreeRoot()
	return (&types.BeaconBlockHeader{}).New(
		beaconState.Slot, // slot
		validatorIndex,   // proposerIndex
		rndRoot(),        // parentBlockRoot
		stateRoot,        // stateRoot
		rndRoot(),        // bodyRoot
	)
}

// Generate a random execution payload header
func rndExecutionPayloadHeader() *types.ExecutionPayloadHeader {
	execPayloadHeader := &types.ExecutionPayloadHeader{
		ParentHash:       rndExecutionHash(),
		FeeRecipient:     rndExecutionAddress(),
		StateRoot:        rndB32(),
		ReceiptsRoot:     rndB32(),
		LogsBloom:        rndB256(),
		Random:           rndB32(),
		Number:           rndU64(),
		GasLimit:         rndU64(),
		GasUsed:          rndU64(),
		Timestamp:        rndU64(),
		ExtraData:        rndBytes(MAX_EXTRA_DATA),
		BaseFeePerGas:    rndInt256(),
		BlockHash:        rndExecutionHash(),
		TransactionsRoot: rndRoot(),
		WithdrawalsRoot:  rndRoot(),
		BlobGasUsed:      rndU64(),
		ExcessBlobGas:    rndU64()}
	return execPayloadHeader
}

// Generate a random validator slice
func rndValidators(validatorCount uint) types.Validators {
	validators := make(types.Validators, validatorCount)
	for i := range validators {
		validators[i] = &types.Validator{
			Pubkey:                     rndValidatorPubkey(),
			EffectiveBalance:           rndU64(),
			Slashed:                    rndBool(),
			ActivationEligibilityEpoch: rndU64(),
			ActivationEpoch:            rndU64(),
			ExitEpoch:                  rndU64(),
			WithdrawableEpoch:          rndU64(),
		}
	}
	return validators
}

// Generate a random root
func rndRoot() common.Root {
	var root common.Root
	mustRndRead(root[:])
	return root
}

// Generate a random ExecutionHash
func rndExecutionHash() common.ExecutionHash {
	var hash common.ExecutionHash
	mustRndRead(hash[:])
	return hash
}

// Generate a random ExecutionAddress
func rndExecutionAddress() common.ExecutionAddress {
	return common.NewExecutionAddressFromHex(
		rndAddress(),
	)
}

// Generate a random validator pubkey
func rndValidatorPubkey() [48]byte {
	var validatorPubkey [48]byte
	mustRndRead(validatorPubkey[:])
	return validatorPubkey
}

// Generate a random address
func rndAddress() string {
	address := make([]byte, 20)
	mustRndRead(address)
	return fmt.Sprintf("0x%x", address)
}

// Generate a random B4
func rndB4() bytes.B4 {
	var b4 bytes.B4
	mustRndRead(b4[:])
	return b4
}

// Generate a random B32
func rndB32() bytes.B32 {
	var b32 bytes.B32
	mustRndRead(b32[:])
	return b32
}

// Generate a random B256
func rndB256() bytes.B256 {
	var b256 bytes.B256
	mustRndRead(b256[:])
	return b256
}

// Generate random bytes of length between 0 and max
func rndBytes(max int) []byte {
	count := rng.Intn(max) + 1
	bytes := make([]byte, count)
	mustRndRead(bytes)
	return bytes
}

// Generate a random uint64
func rndU64() math.U64 {
	return math.U64(rng.Uint64())
}

// Generate random uint256.Int
func rndInt256() *uint256.Int {
	// Generate 32 random bytes (256 bits)
	randomBytes := make([]byte, 32)
	mustRndRead(randomBytes)

	// Convert the random bytes to a uint256.Int
	randomUint256 := new(uint256.Int).SetBytes(randomBytes)

	return randomUint256
}

// Generate a random boolean
func rndBool() bool {
	return rng.Intn(2) == 1
}

// Fill `target` with random bytes
func mustRndRead[T []byte](target T) T {
	_, err := rng.Read(target)
	if err != nil {
		panic(err)
	}
	return target
}

// Define a function type for random value generators
type RandomGenerator[T any] func() T

// Generic function to generate a slice of random elements
func rndSlice[T any](max int, generator RandomGenerator[T]) []T {
	count := rng.Intn(max) + 1
	slice := make([]T, count)
	for i := 0; i < count; i++ {
		slice[i] = generator()
	}
	return slice
}
