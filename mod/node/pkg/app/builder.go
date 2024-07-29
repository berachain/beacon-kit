package app

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

func (b *Builder[
	StorageBackendT, StateProcessorT,
]) Build() *App[StorageBackendT, StateProcessorT] {
	// depinject components into the app.
	return b.app
}

// For these methods, the paradigm will be to pass a pointer to the eventually
// constructed component. This allows the app to force the inclusion of the backend
// and state processor in the app without being strict on the actual components
// included, think of it as a minimal set of required components.

func (b *Builder[
	StorageBackendT, _,
]) WithStorageBackend(storageBackend StorageBackendT) {
	b.app.backend = storageBackend
}

func (b *Builder[
	_, StateProcessorT,
]) WithStateProcessor(stateProcessor StateProcessorT) {
	b.app.stateProcessor = stateProcessor
}
