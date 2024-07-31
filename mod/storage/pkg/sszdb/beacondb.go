package sszdb

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/schema"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core/state"
	fastssz "github.com/ferranbt/fastssz"
)

type BeaconStateDB[
	BeaconBlockHeaderT,
	Eth1DataT,
	ExecutionPayloadHeaderT,
	ForkT,
	ValidatorT state.Validator[WithdrawalCredentialsT],
	ValidatorsT ~[]ValidatorT,
	WithdrawalT state.Withdrawal[WithdrawalT],
	WithdrawalCredentialsT state.WithdrawalCredentials,
] struct {
	*Backend
	schemaRoot schema.SSZType
}

func NewBeaconStateDB[
	BeaconBlockHeaderT,
	Eth1DataT,
	ExecutionPayloadHeaderT,
	ForkT,
	ValidatorT state.Validator[WithdrawalCredentialsT],
	ValidatorsT ~[]ValidatorT,
	WithdrawalT state.Withdrawal[WithdrawalT],
	WithdrawalCredentialsT state.WithdrawalCredentials,
](
	backend *Backend,
	monolith fastssz.HashRoot,
) (*BeaconStateDB[
	BeaconBlockHeaderT,
	Eth1DataT,
	ExecutionPayloadHeaderT,
	ForkT,
	ValidatorT,
	ValidatorsT,
	WithdrawalT,
	WithdrawalCredentialsT,
], error) {
	schemaRoot, err := CreateSchema(monolith)
	if err != nil {
		return nil, err
	}
	db := &BeaconStateDB[
		BeaconBlockHeaderT,
		Eth1DataT,
		ExecutionPayloadHeaderT,
		ForkT,
		ValidatorT,
		ValidatorsT,
		WithdrawalT,
		WithdrawalCredentialsT,
	]{
		Backend:    backend,
		schemaRoot: schemaRoot,
	}
	return db, db.bootstrap(monolith)
}

func (db *BeaconStateDB[
	BeaconBlockHeaderT,
	Eth1DataT,
	ExecutionPayloadHeaderT,
	ForkT,
	ValidatorT,
	ValidatorsT,
	WithdrawalT,
	WithdrawalCredentialsT,
]) bootstrap(
	monolith fastssz.HashRoot,
) error {
	bootstrapped, err := db.Get([]byte("bootstrapped"))
	if err != nil {
		return err
	}
	if bootstrapped != nil {
		return nil
	}
	err = db.SaveMonolith(monolith)
	if err != nil {
		return err
	}
	return db.Set([]byte("bootstrapped"), []byte{1})
}
