package server

import (
	"context"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/node-api/backend"
	"github.com/berachain/beacon-kit/mod/node-api/backend/storage"
	"github.com/berachain/beacon-kit/mod/node-api/server/handlers"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/service"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core/state"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var _ service.Basic = (*Server)(nil)

// Server is the server for the node API.
type Server struct {
	*echo.Echo
	config Config
}

// New creates a new node API server.
func New[
	AvailabilityStoreT storage.AvailabilityStore[
		BeaconBlockBodyT, BlobSidecarsT,
	],
	BeaconBlockT any,
	BeaconBlockBodyT types.RawBeaconBlockBody,
	BeaconBlockHeaderT core.BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconStateT core.BeaconState[
		BeaconStateT, BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
		ForkT, StateStoreT, ValidatorT, WithdrawalT,
	],
	BeaconStateMarshallableT state.BeaconStateMarshallable[
		BeaconStateMarshallableT, BeaconBlockHeaderT, Eth1DataT,
		ExecutionPayloadHeaderT, ForkT, ValidatorT,
	],
	BlobSidecarsT any,
	BlockStoreT storage.BlockStore[BeaconBlockT],
	DepositT storage.Deposit,
	DepositStoreT storage.DepositStore[DepositT],
	Eth1DataT,
	ExecutionPayloadHeaderT,
	ForkT any,
	StateStoreT state.KVStore[
		StateStoreT, BeaconBlockHeaderT, Eth1DataT,
		ExecutionPayloadHeaderT, ForkT, ValidatorT,
	],
	ValidatorT storage.Validator[WithdrawalCredentialsT],
	WithdrawalT storage.Withdrawal[WithdrawalT],
	WithdrawalCredentialsT storage.WithdrawalCredentials,
](
	config Config,
	backend *backend.Backend[
		AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
		BeaconStateT, BlobSidecarsT, BlockStoreT, DepositT, DepositStoreT,
		Eth1DataT, ExecutionPayloadHeaderT, ForkT, StateStoreT, ValidatorT,
		WithdrawalT, WithdrawalCredentialsT,
	],
	corsConfig middleware.CORSConfig,
	loggingConfig middleware.LoggerConfig,
) *Server {
	e := echo.New()
	e.HTTPErrorHandler = handlers.CustomHTTPErrorHandler
	e.Validator = &handlers.CustomValidator{
		Validator: ConstructValidator(),
	}
	UseMiddlewares(
		e,
		middleware.CORSWithConfig(corsConfig),
		middleware.LoggerWithConfig(loggingConfig))
	AssignRoutes(
		e,
		handlers.RouteHandlers[ValidatorT]{Backend: backend},
	)
	return &Server{
		Echo:   e,
		config: config,
	}
}

// Start starts the node API server.
func (s *Server) Start(_ context.Context) error {
	if !s.config.Enabled {
		return nil
	}
	go s.Echo.Start(s.config.Address)
	return nil
}

// Name returns the name of the service.
func (s *Server) Name() string {
	return "node-api"
}
