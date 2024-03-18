package rpc

import (
	"context"
	"fmt"
	"github.com/berachain/beacon-kit/config"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc/credentials"
	"net"

	"google.golang.org/grpc"

	"github.com/berachain/beacon-kit/runtime/service"
)

type Service struct {
	service.BaseService

	cfg        *config.RPC
	grpcServer *grpc.Server
	listener   net.Listener
}

func (s *Service) Start(ctx context.Context) {
	address := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		s.BaseService.Logger().Error("failed to listen: ", err)
		panic(err)
	}
	s.listener = lis
	s.BaseService.Logger().Info(fmt.Sprintf("gRPC server listening on port %d", s.cfg.Port))

	opts := []grpc.ServerOption{
		grpc.StatsHandler(&ocgrpc.ServerHandler{}),
	}

	if s.cfg.CertFlag != "" && s.cfg.KeyFlag != "" {
		creds, err := credentials.NewServerTLSFromFile(s.cfg.CertFlag, s.cfg.KeyFlag)
		if err != nil {
			s.Logger().Error(fmt.Sprintf("Failed to create gRPC server with TLS: %v", err))
		}
		opts = append(opts, grpc.Creds(creds))
	} else {
		s.Logger().Warn("You are using an insecure gRPC server. If you are running your beacon node and " +
			"validator on the same machines, you can ignore this message. If you want to know " +
			"how to enable secure connections, see: https://docs.prylabs.network/docs/prysm-usage/secure-grpc")
	}

	s.grpcServer = grpc.NewServer(opts...)

	go func() {
		if s.listener != nil {
			if err := s.grpcServer.Serve(s.listener); err != nil {
				s.Logger().Error(fmt.Sprintf("Failed to serve gRPC server: %v", err))
			}
		}
	}()
}
