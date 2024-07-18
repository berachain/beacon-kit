package server

import (
	"context"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/node-api/backend"
	"github.com/berachain/beacon-kit/mod/node-api/backend/storage"
	"github.com/berachain/beacon-kit/mod/node-api/server/handlers"
	"github.com/berachain/beacon-kit/mod/node-api/server/routes"
	"github.com/berachain/beacon-kit/mod/node-api/server/utils"
	nodetypes "github.com/berachain/beacon-kit/mod/node-core/pkg/types"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/service"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core/state"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var _ service.Basic = (*Server[nodetypes.Node, any])(nil)

// Server is the server for the node API.
type Server[
	NodeT nodetypes.Node,
	ValidatorT any,
] struct {
	*echo.Echo
	config  Config
	backend handlers.Backend[NodeT, ValidatorT]
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
	NodeT nodetypes.Node,
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
		Eth1DataT, ExecutionPayloadHeaderT, ForkT, NodeT, StateStoreT,
		ValidatorT, WithdrawalT, WithdrawalCredentialsT,
	],
	corsConfig middleware.CORSConfig,
	loggingConfig middleware.LoggerConfig,
) *Server[NodeT, ValidatorT] {
	e := echo.New()
	e.HTTPErrorHandler = utils.HTTPErrorHandler
	e.Validator = &utils.CustomValidator{
		Validator: utils.ConstructValidator(),
	}
	utils.UseMiddlewares(
		e,
		middleware.CORSWithConfig(corsConfig),
		middleware.LoggerWithConfig(loggingConfig))
	routes.Assign(
		e,
		handlers.New(backend),
	)
	return &Server[NodeT, ValidatorT]{
		Echo:    e,
		config:  config,
		backend: backend,
	}
}

// Start starts the node API server.
func (s *Server[NodeT, ValidatorT]) Start(_ context.Context) error {
	if !s.config.Enabled {
		return nil
	}
	go s.Echo.Start(s.config.Address)
	return nil
}

// Name returns the name of the service.
func (s *Server[NodeT, ValidatorT]) Name() string {
	return "node-api"
}

func (s *Server[NodeT, ValidatorT]) AttachNode(node NodeT) {
	s.backend.AttachNode(node)
}
