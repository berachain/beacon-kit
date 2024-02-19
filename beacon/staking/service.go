package staking

import (
	"github.com/itsdevbear/bolaris/runtime/service"
)

// Service represents the staking service.
type Service struct {
	service.BaseService

	// vcp is responsible for applying validator set changes.
	vcp ValsetChangeProvider
}

// NewService returns a new Staking Service.
func NewService(
	base service.BaseService,
	opts ...Option) *Service {
	s := &Service{
		BaseService: base,
	}
	for _, opt := range opts {
		if err := opt(s); err != nil {
			s.Logger().Error("Failed to apply option", "error", err)
		}
	}
	return s
}
