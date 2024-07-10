package storage

import "context"

type Backend[
	AvailabilityStoreT any,
	BeaconStateT any,
	DepositStoreT any,
] interface {
	// AvailabilityStore returns the availability store.
	AvailabilityStore(ctx context.Context) AvailabilityStoreT
	// StateFromContext returns the beacon state from the given context.
	StateFromContext(ctx context.Context) BeaconStateT
	// DepositStore returns the deposit store.
	DepositStore(ctx context.Context) DepositStoreT
}
