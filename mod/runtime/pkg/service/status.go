package service

// StatusEvent represents a service status event.
type StatusEvent struct {
	name    string
	healthy bool
}

// NewStatusEvent creates a new status service.
func NewStatusEvent(name string, healthy bool) *StatusEvent {
	return &StatusEvent{
		name:    name,
		healthy: healthy,
	}
}

// Name returns the name of the service.
func (s *StatusEvent) Name() string {
	return s.name
}

// IsHealthy returns the health status of the service.
func (s *StatusEvent) IsHealthy() bool {
	return s.healthy
}
