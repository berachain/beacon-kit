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

// NewStatus creates a new status service.
func (s *StatusEvent) Type() string {
	return s.name
}

// IsHealthy returns the health status of the service.
func (s *StatusEvent) IsHealthy() bool {
	return s.healthy
}
