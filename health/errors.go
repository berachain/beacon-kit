package health

import "errors"

// ErrHealthCheckTimeout is returned when a health check times out.
var ErrHealthCheckTimeout = errors.New("health check timed out")
