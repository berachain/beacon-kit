package storage

import (
	"context"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/storage"
	nodetypes "github.com/berachain/beacon-kit/mod/node-core/pkg/types"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core/state"
)

func NewBackend[
	AvailabilityStoreT AvailabilityStore[
		BeaconBlockBodyT, BlobSidecarsT,
	],
	BeaconBlockT any,
	BeaconBlockBodyT types.RawBeaconBlockBody,
	BeaconBlockHeaderT core.BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconStateT core.BeaconState[
		BeaconStateT, BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
		ForkT, StateStoreT, ValidatorT, WithdrawalT,
	],
	BeaconStateMarshallableT state.BeaconStateMarshallable[
		BeaconStateMarshallableT, BeaconBlockHeaderT, Eth1DataT,
		ExecutionPayloadHeaderT, ForkT, ValidatorT,
	],
	BlobSidecarsT any,
	BlockStoreT BlockStore[BeaconBlockT],
	DepositT Deposit,
	DepositStoreT DepositStore[DepositT],
	Eth1DataT,
	ExecutionPayloadHeaderT,
	ForkT any,
	StateStoreT state.KVStore[
		StateStoreT, BeaconBlockHeaderT, Eth1DataT,
		ExecutionPayloadHeaderT, ForkT, ValidatorT,
	],
	ValidatorT Validator[WithdrawalCredentialsT],
	WithdrawalT Withdrawal[WithdrawalT],
	WithdrawalCredentialsT WithdrawalCredentials,
](
	node nodetypes.Node,
	sb *storage.Backend[
		AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
		BeaconStateT, BeaconStateMarshallableT, BlobSidecarsT, BlockStoreT,
		DepositT, DepositStoreT, Eth1DataT, ExecutionPayloadHeaderT, ForkT,
		StateStoreT, ValidatorT, WithdrawalT, WithdrawalCredentialsT,
	],
) Backend[
	AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
	BeaconStateT, BlobSidecarsT, BlockStoreT, DepositT, DepositStoreT,
	Eth1DataT, ExecutionPayloadHeaderT, ForkT, StateStoreT, ValidatorT,
	WithdrawalT, WithdrawalCredentialsT,
] {
	return &backend[
		AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
		BeaconStateT, BeaconStateMarshallableT, BlobSidecarsT, BlockStoreT,
		DepositT, DepositStoreT, Eth1DataT, ExecutionPayloadHeaderT, ForkT,
		StateStoreT, ValidatorT, WithdrawalT, WithdrawalCredentialsT,
	]{
		node:    node,
		Backend: sb,
	}
}

type backend[
	AvailabilityStoreT AvailabilityStore[
		BeaconBlockBodyT, BlobSidecarsT,
	],
	BeaconBlockT any,
	BeaconBlockBodyT types.RawBeaconBlockBody,
	BeaconBlockHeaderT core.BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconStateT core.BeaconState[
		BeaconStateT, BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
		ForkT, StateStoreT, ValidatorT, WithdrawalT,
	],
	BeaconStateMarshallableT state.BeaconStateMarshallable[
		BeaconStateMarshallableT, BeaconBlockHeaderT, Eth1DataT,
		ExecutionPayloadHeaderT, ForkT, ValidatorT,
	],
	BlobSidecarsT any,
	BlockStoreT BlockStore[BeaconBlockT],
	DepositT Deposit,
	DepositStoreT DepositStore[DepositT],
	Eth1DataT,
	ExecutionPayloadHeaderT,
	ForkT any,
	StateStoreT state.KVStore[
		StateStoreT, BeaconBlockHeaderT, Eth1DataT,
		ExecutionPayloadHeaderT, ForkT, ValidatorT,
	],
	ValidatorT Validator[WithdrawalCredentialsT],
	WithdrawalT Withdrawal[WithdrawalT],
	WithdrawalCredentialsT WithdrawalCredentials,
] struct {
	node nodetypes.Node
	*storage.Backend[
		AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
		BeaconStateT, BeaconStateMarshallableT, BlobSidecarsT, BlockStoreT,
		DepositT, DepositStoreT, Eth1DataT, ExecutionPayloadHeaderT, ForkT,
		StateStoreT, ValidatorT, WithdrawalT, WithdrawalCredentialsT,
	]
}

// StateFromContext returns a state from the context.
// It wraps the StorageBackend.StateFromContext method to allow for state ID
// querying.
func (b *backend[
	_, _, _, _, BeaconStateT, _, _, _, _, _, _, _, _, _, _, _, _,
]) StateFromContext(
	ctx context.Context,
	stateID string,
) (BeaconStateT, error) {
	var state BeaconStateT
	height, err := heightFromStateID(stateID)
	if err != nil {
		return state, err
	}
	queryCtx, err := b.node.CreateQueryContext(height, false)
	if err != nil {
		panic(err)
	}
	return b.Backend.StateFromContext(queryCtx), nil
}
