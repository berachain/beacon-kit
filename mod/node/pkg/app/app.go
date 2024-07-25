package app

import (
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/service"
)

type App[
	StorageBackendT any,
	StateProcessorT any,
] struct {
	logger log.Logger[any]

	// Services contained within the service registry
	// are the services used to push the chain and
	// related components forward.
	services service.Registry
	// The backend is the central data access layer for
	// the application.
	backend StorageBackendT
	// The state processor is the component that is
	// responsible for transitioning the state of the
	// chain.
	stateProcessor StateProcessorT
}
