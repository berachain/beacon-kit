package app

import (
	"fmt"

	"github.com/berachain/beacon-kit/mod/config"
	consensusengine "github.com/berachain/beacon-kit/mod/consensus/pkg/engine"
	"github.com/berachain/beacon-kit/mod/depinject"
	"github.com/berachain/beacon-kit/mod/log/pkg/phuslu"
	"github.com/berachain/beacon-kit/mod/node/pkg/app/components"
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
	return &Builder[StorageBackendT, StateProcessorT]{
		app: &App[StorageBackendT, StateProcessorT]{},
	}
}

func (b *Builder[
	StorageBackendT, StateProcessorT,
]) Build(
	logger *phuslu.Logger,
	appOpts *components.AppOptions,
	config *config.Config,
) (*App[StorageBackendT, StateProcessorT], error) {
	var err error
	// depinject components into the app.
	container := depinject.NewContainer()
	if err = container.Supply(
		logger,
		appOpts,
		config,
	); err != nil {
		return nil, err
	}
	if err = container.Provide(b.components...); err != nil {
		return nil, err
	}

	// Resolve dependencies and construct the app.
	if err = container.Inject(
		&b.app.backend,
		&b.app.stateProcessor,
		&b.app.services,
		&b.app.Client,
	); err != nil {
		return nil, err
	}
	fmt.Println("DEPINJECTED")
	b.app.Logger = logger

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
]) WithConsensusClient(
	consensusClient consensusengine.Client,
) *Builder[StorageBackendT, StateProcessorT] {
	b.app.Client = consensusClient
	return b
}

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
