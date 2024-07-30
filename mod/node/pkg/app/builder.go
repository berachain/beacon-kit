package app

import (
	"github.com/berachain/beacon-kit/mod/depinject"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/service"
)

type Builder[
	StorageBackendT any,
	StateProcessorT any,
] struct {
	app *App[StorageBackendT, StateProcessorT]

	// components is a slice of components to be added to the app
	// These components will be depinjected into the app and added
	// to the app at runtime.
	components []any
}

func NewBuilder[
	StorageBackendT any,
	StateProcessorT any,
]() *Builder[StorageBackendT, StateProcessorT] {
	return &Builder[StorageBackendT, StateProcessorT]{}
}

func (b *Builder[
	StorageBackendT, StateProcessorT,
]) Build() (*App[StorageBackendT, StateProcessorT], error) {
	var err error
	// depinject components into the app.
	container := depinject.NewContainer()
	if err = container.Supply(
	// supplied deps
	); err != nil {
		return nil, err
	}
	if err = container.Provide(b.components...); err != nil {
		return nil, err
	}

	// Resolve dependencies and construct the app.
	var (
		logger          log.Logger[any]
		storageBackend  StorageBackendT
		stateProcessor  StateProcessorT
		serviceRegistry service.Registry
	)
	if err = container.Inject(
		&logger,
		&storageBackend,
		&stateProcessor,
		&serviceRegistry,
	); err != nil {
		return nil, err
	}
	b.app = &App[StorageBackendT, StateProcessorT]{
		logger:         logger,
		backend:        storageBackend,
		stateProcessor: stateProcessor,
		services:       serviceRegistry,
	}

	return b.app, nil
}

// BUILDER OPTIONS

func (b *Builder[
	StorageBackendT, StateProcessorT,
]) WithComponents(
	components ...any,
) *Builder[StorageBackendT, StateProcessorT] {
	b.components = append(b.components, components...)
	return b
}

// For these methods, the paradigm will be to pass a pointer to the eventually
// constructed component. This allows the app to force the inclusion of the backend
// and state processor in the app without being strict on the actual components
// included, think of it as a minimal set of required components.

func (b *Builder[
	StorageBackendT, StateProcessorT,
]) WithStorageBackend(
	storageBackend StorageBackendT,
) *Builder[StorageBackendT, StateProcessorT] {
	b.app.backend = storageBackend
	return b
}

func (b *Builder[
	StorageBackendT, StateProcessorT,
]) WithStateProcessor(
	stateProcessor StateProcessorT,
) *Builder[StorageBackendT, StateProcessorT] {
	b.app.stateProcessor = stateProcessor
	return b
}
