package rpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/berachain/beacon-kit/beacon/rpc/beacon"
	"github.com/berachain/beacon-kit/config"
	"github.com/gorilla/mux"
	"net/http"
	"time"

	"github.com/berachain/beacon-kit/runtime/service"
)

type Service struct {
	service.BaseService

	cfg    *config.RPC
	Router *mux.Router
	server *http.Server
}

func (s *Service) Start(ctx context.Context) {
	logger := s.Logger().With("module", "rpc")
	address := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	s.Router = newRouter()

	s.initializeBeaconServerRoutes(&beacon.Server{})

	s.server = &http.Server{
		Addr:              address,
		Handler:           s.Router,
		ReadHeaderTimeout: time.Second,
	}

	go func() {
		logger.With("address", address).Info("Starting gRPC gateway")
		if err := s.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Error(fmt.Sprintf("Failed to start RPC server: %v", err))
			return
		}
	}()
}

func (s *Service) initializeBeaconServerRoutes(beaconServer *beacon.Server) {
	s.Router.HandleFunc("/eth/v1/beacon/states/{state_id}/randao", beaconServer.GetRandao).Methods(http.MethodGet)
}
