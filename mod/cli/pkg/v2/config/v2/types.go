package config

import "context"

// Config is a generic type containing the app config
// information.
//
// Note: there is dynamic connasence between the config
// type implementation and the encoding tool used to
// encode/decode the configuration.
type Config[T any] interface {
	// Default returns the default configuration.
	Default() T
}

// Node is a generic type representing a node in the app.
type Node[ConfigT Config[ConfigT]] interface {
	// Start starts the node.
	Start(context.Context) error
	// DefaultConfig returns the default configuration.
	DefaultConfig() ConfigT
}
