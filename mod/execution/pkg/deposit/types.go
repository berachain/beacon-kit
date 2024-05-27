package deposit

import (
	"context"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
)

type StorageBackend[
	AvailabilityStoreT any,
	BeaconStateT any,
	BlobSidecarsT any,
	DepositStoreT DepositStore,
] interface {
	// DepositStore returns the deposit store for the given context.
	DepositStore(context.Context) DepositStoreT
}

// DepositContract is the ABI for the deposit contract.
type DepositContract interface {
	GetDeposits(
		ctx context.Context,
		blockNumber uint64,
	) ([]*types.Deposit, error)
}

// DepositStore defines the interface for managing deposit operations.
type DepositStore interface {
	// PruneToIndex prunes the deposit store up to the specified index.
	PruneToIndex(index uint64) error
	// EnqueueDeposits adds a list of deposits to the deposit store.
	EnqueueDeposits(deposits []*types.Deposit) error
}
