package beacon

import (
	echo "github.com/labstack/echo/v4"
)

// Handler is the handler for the beacon API.
type Handler[ValidatorT any] struct {
	backend Backend[ValidatorT]
}

// NewHandler creates a new handler for the beacon API.
func NewHandler[ValidatorT any](
	backend Backend[ValidatorT],
) Handler[ValidatorT] {
	return Handler[ValidatorT]{
		backend: backend,
	}
}

// NotImplemented is a placeholder for the beacon API.
func (h Handler[_]) NotImplemented(_ echo.Context) error {
	return echo.ErrNotImplemented
}
