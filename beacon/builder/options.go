package builder

import (
	"github.com/itsdevbear/bolaris/config"
	"github.com/itsdevbear/bolaris/runtime/service"
)

// WithBaseService sets the base service.
func WithBaseService(svc service.BaseService) service.Option[Service] {
	return func(s *Service) error {
		s.BaseService = svc
		return nil
	}
}

// WithBuilderConfig sets the builder config.
func WithBuilderConfig(cfg *config.Builder) service.Option[Service] {
	return func(s *Service) error {
		s.cfg = cfg
		return nil
	}
}

// WithLocalBuilder sets the local builder.
func WithLocalBuilder(builder PayloadBuilder) service.Option[Service] {
	return func(s *Service) error {
		s.localBuilder = builder
		return nil
	}
}

// WithRemoteBuilders sets the remote builders.
func WithRemoteBuilders(builders ...PayloadBuilder) service.Option[Service] {
	return func(s *Service) error {
		s.remoteBuilders = append(s.remoteBuilders, builders...)
		return nil
	}
}
