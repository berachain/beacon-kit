package runtime

import "context"

// Basic is the minimal interface for a service.
type BasicService interface {
	// Start spawns any goroutines required by the service.
	Start(ctx context.Context) error
	// Name returns the name of the service.
	Name() string
}
