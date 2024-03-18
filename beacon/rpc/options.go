package rpc

import (
	"github.com/berachain/beacon-kit/config"
	"github.com/berachain/beacon-kit/runtime/service"
)

// WithBaseService sets the base service for the RPC service.
func WithBaseService(bs service.BaseService) service.Option[Service] {
	return func(s *Service) error {
		s.BaseService = bs

		return nil
	}
}

// WithConfig sets the configuration for the RPC service.
func WithConfig(cfg *config.RPC) service.Option[Service] {
	return func(s *Service) error {
		s.cfg = cfg

		return nil
	}
}
