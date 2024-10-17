package components

import (
	"context"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

type (
	// BeaconState is the interface for the beacon state. It
	// is a combination of the read-only and write-only beacon state types.
	BeaconState[
		T,
		BeaconBlockHeaderT,
		BeaconStateMarshallableT,
		// Eth1DataT,
		ExecutionPayloadHeaderT,
		ForkT,
		KVStoreT,
		ValidatorT,
		ValidatorsT,
		WithdrawalT,
		WithdrawalsT any,
	] interface {
		NewFromDB(
			bdb KVStoreT,
			cs common.ChainSpec,
		) T
		Copy() T
		Context() context.Context
		HashTreeRoot() common.Root
		GetMarshallable() (BeaconStateMarshallableT, error)

		ReadOnlyBeaconState[
			BeaconBlockHeaderT, ExecutionPayloadHeaderT,
			ForkT, ValidatorT, ValidatorsT, WithdrawalT, WithdrawalsT,
		]
		WriteOnlyBeaconState[
			BeaconBlockHeaderT, ExecutionPayloadHeaderT,
			ForkT, ValidatorT, WithdrawalsT,
		]
	}

	// BeaconStore is the interface for the beacon store.
	BeaconStore[
		T any,
		BeaconBlockHeaderT any,
		ExecutionPayloadHeaderT any,
		ForkT any,
		ValidatorT any,
		ValidatorsT any,
		WithdrawalT any,
	] interface {
		// Context returns the context of the key-value store.
		Context() context.Context
		// WithContext returns a new key-value store with the given context.
		WithContext(
			ctx context.Context,
		) T
		// Copy returns a copy of the key-value store.
		Copy() T
		// GetLatestExecutionPayloadHeader retrieves the latest execution
		// payload
		// header.
		GetLatestExecutionPayloadHeader() (
			ExecutionPayloadHeaderT, error,
		)
		// SetLatestExecutionPayloadHeader sets the latest execution payload
		// header.
		SetLatestExecutionPayloadHeader(
			payloadHeader ExecutionPayloadHeaderT,
		) error
		// GetEth1DepositIndex retrieves the eth1 deposit index.
		GetEth1DepositIndex() (uint64, error)
		// SetEth1DepositIndex sets the eth1 deposit index.
		SetEth1DepositIndex(
			index uint64,
		) error
		// GetBalance retrieves the balance of a validator.
		GetBalance(idx math.ValidatorIndex) (math.Gwei, error)
		// SetBalance sets the balance of a validator.
		SetBalance(idx math.ValidatorIndex, balance math.Gwei) error
		// GetSlot retrieves the current slot.
		GetSlot() (math.Slot, error)
		// SetSlot sets the current slot.
		SetSlot(slot math.Slot) error
		// GetFork retrieves the fork.
		GetFork() (ForkT, error)
		// SetFork sets the fork.
		SetFork(fork ForkT) error
		// GetGenesisValidatorsRoot retrieves the genesis validators root.
		GetGenesisValidatorsRoot() (common.Root, error)
		// SetGenesisValidatorsRoot sets the genesis validators root.
		SetGenesisValidatorsRoot(root common.Root) error
		// GetLatestBlockHeader retrieves the latest block header.
		GetLatestBlockHeader() (BeaconBlockHeaderT, error)
		// SetLatestBlockHeader sets the latest block header.
		SetLatestBlockHeader(header BeaconBlockHeaderT) error
		// GetBlockRootAtIndex retrieves the block root at the given index.
		GetBlockRootAtIndex(index uint64) (common.Root, error)
		// StateRootAtIndex retrieves the state root at the given index.
		StateRootAtIndex(index uint64) (common.Root, error)
		// // GetEth1Data retrieves the eth1 data.
		// GetEth1Data() (Eth1DataT, error)
		// SetEth1Data sets the eth1 data.
		// SetEth1Data(data Eth1DataT) error
		// GetValidators retrieves all validators.
		GetValidators() (ValidatorsT, error)
		// GetBalances retrieves all balances.
		GetBalances() ([]uint64, error)
		// GetNextWithdrawalIndex retrieves the next withdrawal index.
		GetNextWithdrawalIndex() (uint64, error)
		// SetNextWithdrawalIndex sets the next withdrawal index.
		SetNextWithdrawalIndex(index uint64) error
		// GetNextWithdrawalValidatorIndex retrieves the next withdrawal
		// validator
		// index.
		GetNextWithdrawalValidatorIndex() (math.ValidatorIndex, error)
		// SetNextWithdrawalValidatorIndex sets the next withdrawal validator
		// index.
		SetNextWithdrawalValidatorIndex(index math.ValidatorIndex) error
		// GetTotalSlashing retrieves the total slashing.
		GetTotalSlashing() (math.Gwei, error)
		// SetTotalSlashing sets the total slashing.
		SetTotalSlashing(total math.Gwei) error
		// GetRandaoMixAtIndex retrieves the randao mix at the given index.
		GetRandaoMixAtIndex(index uint64) (common.Bytes32, error)
		// GetSlashings retrieves all slashings.
		GetSlashings() ([]uint64, error)
		// SetSlashingAtIndex sets the slashing at the given index.
		SetSlashingAtIndex(index uint64, amount math.Gwei) error
		// GetSlashingAtIndex retrieves the slashing at the given index.
		GetSlashingAtIndex(index uint64) (math.Gwei, error)
		// GetTotalValidators retrieves the total validators.
		GetTotalValidators() (uint64, error)
		// GetTotalActiveBalances retrieves the total active balances.
		GetTotalActiveBalances(uint64) (math.Gwei, error)
		// ValidatorByIndex retrieves the validator at the given index.
		ValidatorByIndex(index math.ValidatorIndex) (ValidatorT, error)
		// UpdateBlockRootAtIndex updates the block root at the given index.
		UpdateBlockRootAtIndex(index uint64, root common.Root) error
		// UpdateStateRootAtIndex updates the state root at the given index.
		UpdateStateRootAtIndex(index uint64, root common.Root) error
		// UpdateRandaoMixAtIndex updates the randao mix at the given index.
		UpdateRandaoMixAtIndex(index uint64, mix common.Bytes32) error
		// UpdateValidatorAtIndex updates the validator at the given index.
		UpdateValidatorAtIndex(
			index math.ValidatorIndex,
			validator ValidatorT,
		) error
		// ValidatorIndexByPubkey retrieves the validator index by the given
		// pubkey.
		ValidatorIndexByPubkey(
			pubkey crypto.BLSPubkey,
		) (math.ValidatorIndex, error)
		// AddValidator adds a validator.
		AddValidator(val ValidatorT) error
		// AddValidatorBartio adds a validator to the Bartio chain.
		AddValidatorBartio(val ValidatorT) error
		// ValidatorIndexByCometBFTAddress retrieves the validator index by the
		// given comet BFT address.
		ValidatorIndexByCometBFTAddress(
			cometBFTAddress []byte,
		) (math.ValidatorIndex, error)
		// GetValidatorsByEffectiveBalance retrieves validators by effective
		// balance.
		GetValidatorsByEffectiveBalance() ([]ValidatorT, error)
	}

	// ReadOnlyBeaconState is the interface for a read-only beacon state.
	ReadOnlyBeaconState[
		BeaconBlockHeaderT, ExecutionPayloadHeaderT, ForkT,
		ValidatorT, ValidatorsT, WithdrawalT, WithdrawalsT any,
	] interface {
		ReadOnlyRandaoMixes
		ReadOnlyStateRoots
		ReadOnlyValidators[ValidatorT]
		ReadOnlyWithdrawals[WithdrawalT]
		// GetBalances retrieves all balances.
		GetBalances() ([]uint64, error)
		GetBalance(math.ValidatorIndex) (math.Gwei, error)
		GetSlot() (math.Slot, error)
		GetFork() (ForkT, error)
		GetEth1DepositIndex() (uint64, error)
		GetLatestExecutionPayloadHeader() (
			ExecutionPayloadHeaderT, error,
		)
		GetGenesisValidatorsRoot() (common.Root, error)
		GetBlockRootAtIndex(uint64) (common.Root, error)
		GetLatestBlockHeader() (BeaconBlockHeaderT, error)
		GetTotalActiveBalances(uint64) (math.Gwei, error)
		GetValidators() (ValidatorsT, error)
		GetSlashingAtIndex(uint64) (math.Gwei, error)
		GetTotalSlashing() (math.Gwei, error)
		GetNextWithdrawalIndex() (uint64, error)
		GetNextWithdrawalValidatorIndex() (math.ValidatorIndex, error)
		GetTotalValidators() (uint64, error)
		GetValidatorsByEffectiveBalance() ([]ValidatorT, error)
		GetWithdrawals() (WithdrawalsT, error)
		ValidatorIndexByCometBFTAddress(
			cometBFTAddress []byte,
		) (math.ValidatorIndex, error)
	}

	// WriteOnlyBeaconState is the interface for a write-only beacon state.
	WriteOnlyBeaconState[
		BeaconBlockHeaderT, ExecutionPayloadHeaderT,
		ForkT, ValidatorT, WithdrawalsT any,
	] interface {
		WriteOnlyRandaoMixes
		WriteOnlyStateRoots
		WriteOnlyValidators[ValidatorT]

		SetGenesisValidatorsRoot(root common.Root) error
		SetFork(ForkT) error
		SetSlot(math.Slot) error
		SetEth1DepositIndex(uint64) error
		SetLatestExecutionPayloadHeader(
			ExecutionPayloadHeaderT,
		) error
		UpdateBlockRootAtIndex(uint64, common.Root) error
		SetLatestBlockHeader(BeaconBlockHeaderT) error
		IncreaseBalance(math.ValidatorIndex, math.Gwei) error
		DecreaseBalance(math.ValidatorIndex, math.Gwei) error
		UpdateSlashingAtIndex(uint64, math.Gwei) error
		SetNextWithdrawalIndex(uint64) error
		SetNextWithdrawalValidatorIndex(math.ValidatorIndex) error
		SetTotalSlashing(math.Gwei) error
		SetWithdrawals(WithdrawalsT) error
	}

	// WriteOnlyStateRoots defines a struct which only has write access to state
	// roots methods.
	WriteOnlyStateRoots interface {
		UpdateStateRootAtIndex(uint64, common.Root) error
	}

	// ReadOnlyStateRoots defines a struct which only has read access to state
	// roots
	// methods.
	ReadOnlyStateRoots interface {
		StateRootAtIndex(uint64) (common.Root, error)
	}

	// WriteOnlyRandaoMixes defines a struct which only has write access to
	// randao
	// mixes methods.
	WriteOnlyRandaoMixes interface {
		UpdateRandaoMixAtIndex(uint64, common.Bytes32) error
	}

	// ReadOnlyRandaoMixes defines a struct which only has read access to randao
	// mixes methods.
	ReadOnlyRandaoMixes interface {
		GetRandaoMixAtIndex(uint64) (common.Bytes32, error)
	}

	// WriteOnlyValidators has write access to validator methods.
	WriteOnlyValidators[ValidatorT any] interface {
		UpdateValidatorAtIndex(
			math.ValidatorIndex,
			ValidatorT,
		) error

		AddValidator(ValidatorT) error
		AddValidatorBartio(ValidatorT) error
	}

	// ReadOnlyValidators has read access to validator methods.
	ReadOnlyValidators[ValidatorT any] interface {
		ValidatorIndexByPubkey(
			crypto.BLSPubkey,
		) (math.ValidatorIndex, error)

		ValidatorByIndex(
			math.ValidatorIndex,
		) (ValidatorT, error)
	}

	// ReadOnlyWithdrawals only has read access to withdrawal methods.
	ReadOnlyWithdrawals[WithdrawalT any] interface {
		ExpectedWithdrawals() ([]WithdrawalT, error)
	}
)
