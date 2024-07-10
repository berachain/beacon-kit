package state

import (
	"context"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// BeaconState is the interface for the beacon state. It
// is a combination of the read-only and write-only beacon state types.
type BeaconState[
	T any,
	BeaconBlockHeaderT any,
	Eth1DataT,
	ExecutionPayloadHeaderT,
	ForkT,
	KVStoreT,
	ValidatorT,
	WithdrawalT any,
] interface {
	NewFromDB(
		bdb KVStoreT,
		cs common.ChainSpec,
	) T
	Copy() T
	Save()
	Context() context.Context
	HashTreeRoot() ([32]byte, error)
	ReadOnlyBeaconState[
		BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
		ValidatorT, WithdrawalT,
	]
	WriteOnlyBeaconState[
		BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
		ForkT, ValidatorT,
	]
}

// ReadOnlyBeaconState is the interface for a read-only beacon state.
type ReadOnlyBeaconState[
	BeaconBlockHeaderT any,
	Eth1DataT,
	ExecutionPayloadHeaderT,
	ValidatorT,
	WithdrawalT any,
] interface {
	ReadOnlyEth1Data[Eth1DataT, ExecutionPayloadHeaderT]
	ReadOnlyRandaoMixes
	ReadOnlyStateRoots
	ReadOnlyValidators[ValidatorT]
	ReadOnlyWithdrawals[WithdrawalT]

	GetBalance(math.ValidatorIndex) (math.Gwei, error)
	GetSlot() (math.Slot, error)
	GetGenesisValidatorsRoot() (common.Root, error)
	GetBlockRootAtIndex(uint64) (common.Root, error)
	GetLatestBlockHeader() (BeaconBlockHeaderT, error)
	GetTotalActiveBalances(uint64) (math.Gwei, error)
	GetValidators() ([]ValidatorT, error)
	GetTotalSlashing() (math.Gwei, error)
	GetNextWithdrawalIndex() (uint64, error)
	GetNextWithdrawalValidatorIndex() (math.ValidatorIndex, error)
	GetTotalValidators() (uint64, error)
	GetValidatorsByEffectiveBalance() ([]ValidatorT, error)
	ValidatorIndexByCometBFTAddress(
		cometBFTAddress []byte,
	) (math.ValidatorIndex, error)
}

// WriteOnlyBeaconState is the interface for a write-only beacon state.
type WriteOnlyBeaconState[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT, ForkT, ValidatorT any,
] interface {
	WriteOnlyEth1Data[Eth1DataT, ExecutionPayloadHeaderT]
	WriteOnlyRandaoMixes
	WriteOnlyStateRoots
	WriteOnlyValidators[ValidatorT]

	SetGenesisValidatorsRoot(root common.Root) error
	SetFork(ForkT) error
	SetSlot(math.Slot) error
	UpdateBlockRootAtIndex(uint64, common.Root) error
	SetLatestBlockHeader(BeaconBlockHeaderT) error
	IncreaseBalance(math.ValidatorIndex, math.Gwei) error
	DecreaseBalance(math.ValidatorIndex, math.Gwei) error
	UpdateSlashingAtIndex(uint64, math.Gwei) error
	SetNextWithdrawalIndex(uint64) error
	SetNextWithdrawalValidatorIndex(math.ValidatorIndex) error
	RemoveValidatorAtIndex(math.ValidatorIndex) error
	SetTotalSlashing(math.Gwei) error
}

// WriteOnlyStateRoots defines a struct which only has write access to state
// roots methods.
type WriteOnlyStateRoots interface {
	UpdateStateRootAtIndex(uint64, common.Root) error
}

// ReadOnlyStateRoots defines a struct which only has read access to state roots
// methods.
type ReadOnlyStateRoots interface {
	StateRootAtIndex(uint64) (common.Root, error)
}

// WriteOnlyRandaoMixes defines a struct which only has write access to randao
// mixes methods.
type WriteOnlyRandaoMixes interface {
	UpdateRandaoMixAtIndex(uint64, common.Bytes32) error
}

// ReadOnlyRandaoMixes defines a struct which only has read access to randao
// mixes methods.
type ReadOnlyRandaoMixes interface {
	GetRandaoMixAtIndex(uint64) (common.Bytes32, error)
}

// WriteOnlyValidators has write access to validator methods.
type WriteOnlyValidators[ValidatorT any] interface {
	UpdateValidatorAtIndex(
		math.ValidatorIndex,
		ValidatorT,
	) error

	AddValidator(ValidatorT) error
	AddValidatorBartio(ValidatorT) error
}

// ReadOnlyValidators has read access to validator methods.
type ReadOnlyValidators[ValidatorT any] interface {
	ValidatorIndexByPubkey(
		crypto.BLSPubkey,
	) (math.ValidatorIndex, error)

	ValidatorByIndex(
		math.ValidatorIndex,
	) (ValidatorT, error)
}

// WriteOnlyEth1Data has write access to eth1 data.
type WriteOnlyEth1Data[Eth1DataT, ExecutionPayloadHeaderT any] interface {
	SetEth1Data(Eth1DataT) error
	SetEth1DepositIndex(uint64) error
	SetLatestExecutionPayloadHeader(
		ExecutionPayloadHeaderT,
	) error
}

// ReadOnlyEth1Data has read access to eth1 data.
type ReadOnlyEth1Data[Eth1DataT, ExecutionPayloadHeaderT any] interface {
	GetEth1Data() (Eth1DataT, error)
	GetEth1DepositIndex() (uint64, error)
	GetLatestExecutionPayloadHeader() (
		ExecutionPayloadHeaderT, error,
	)
}

// ReadOnlyWithdrawals only has read access to withdrawal methods.
type ReadOnlyWithdrawals[WithdrawalT any] interface {
	ExpectedWithdrawals() ([]WithdrawalT, error)
}
