package components

import (
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/beacon/blockchain"
	"github.com/berachain/beacon-kit/mod/beacon/validator"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	dastore "github.com/berachain/beacon-kit/mod/da/pkg/store"
	datypes "github.com/berachain/beacon-kit/mod/da/pkg/types"
	engineclient "github.com/berachain/beacon-kit/mod/execution/pkg/client"
	"github.com/berachain/beacon-kit/mod/execution/pkg/deposit"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/metrics"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/services/version"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/feed"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/service"
	depositdb "github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
	"github.com/berachain/beacon-kit/mod/storage/pkg/manager"
	sdkversion "github.com/cosmos/cosmos-sdk/version"
	"github.com/ethereum/go-ethereum/event"
)

// ServiceRegistryInput is the input for the service registry provider.
type ServiceRegistryInput struct {
	depinject.In
	Logger       log.Logger
	ChainService *blockchain.Service[
		*dastore.Store[*types.BeaconBlockBody],
		*types.BeaconBlock,
		*types.BeaconBlockBody,
		BeaconState,
		*datypes.BlobSidecars,
		*types.Deposit,
		*depositdb.KVStore[*types.Deposit],
	]
	DBManagerService *manager.DBManager[
		*types.BeaconBlock,
		*feed.Event[*types.BeaconBlock],
		event.Subscription,
	]
	DepositService *deposit.Service[
		*types.BeaconBlock,
		*types.BeaconBlockBody,
		*feed.Event[*types.BeaconBlock],
		*types.Deposit,
		*types.ExecutionPayload,
		event.Subscription,
		types.WithdrawalCredentials,
	]
	EngineClient     *engineclient.EngineClient[*types.ExecutionPayload]
	TelemetrySink    *metrics.TelemetrySink
	ValidatorService *validator.Service[
		*types.BeaconBlock,
		*types.BeaconBlockBody,
		BeaconState,
		*datypes.BlobSidecars,
		*depositdb.KVStore[*types.Deposit],
		*types.ForkData,
	]
}

// ProvideServiceRegistry is the depinject provider for the service registry.
func ProvideServiceRegistry(
	in ServiceRegistryInput,
) *service.Registry {
	return service.NewRegistry(
		service.WithLogger(in.Logger.With("service", "service-registry")),
		service.WithService(in.ValidatorService),
		service.WithService(in.ChainService),
		service.WithService(in.DepositService),
		service.WithService(in.EngineClient),
		service.WithService(version.NewReportingService(
			in.Logger,
			in.TelemetrySink,
			sdkversion.Version,
		)),
		service.WithService(in.DBManagerService),
	)
}
