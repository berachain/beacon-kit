package runtime

import (
	"context"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/log"
	"cosmossdk.io/core/store"
)

// NewEnvironment creates a new environment for the application
// For setting custom services that aren't set by default, use the EnvOption
// Note: Depinject always provide an environment with all services (mandatory and optional)
func NewEnvironment(
	kvService store.KVStoreService,
	logger log.Logger,
	opts ...EnvOption,
) appmodule.Environment {
	env := appmodule.Environment{
		Logger:             logger,
		EventService:       nil,
		HeaderService:      nil,
		BranchService:      nil,
		GasService:         nil,
		TransactionService: nil,
		KVStoreService:     kvService,
		MsgRouterService:   nil,
		QueryRouterService: nil,
		MemStoreService:    failingMemStore{},
	}

	for _, opt := range opts {
		opt(&env)
	}

	return env
}

type EnvOption func(*appmodule.Environment)

func EnvWithMemStoreService(memStoreService store.MemoryStoreService) EnvOption {
	return func(env *appmodule.Environment) {
		env.MemStoreService = memStoreService
	}
}

// failingMemStore is a memstore that panics when accessed
// this is to ensure all fields are set by in environment
type failingMemStore struct {
	store.MemoryStoreService
}

func (failingMemStore) OpenMemoryStore(context.Context) store.KVStore {
	panic("memory store not set")
}
