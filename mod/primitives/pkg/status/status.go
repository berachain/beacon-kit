package status

type Status int

const (
	ServiceNotStarted Status = iota
	ServiceStarted
	ServiceStopped
	ServiceUnhealthy
	ServiceHealthy
)
