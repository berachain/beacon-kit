package components

import (
	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/mod/config"
	"github.com/berachain/beacon-kit/mod/node-api/backend"
	"github.com/berachain/beacon-kit/mod/node-api/backend/storage"
	"github.com/berachain/beacon-kit/mod/node-api/server"
	nodetypes "github.com/berachain/beacon-kit/mod/node-core/pkg/types"
	"github.com/labstack/echo/v4/middleware"
)

type NodeAPIBackendInput struct {
	depinject.In

	StorageBackend *StorageBackend
}

func ProvideNodeAPIBackend(in NodeAPIBackendInput) *NodeAPIBackend {
	var node nodetypes.Node
	storageBackend := storage.NewBackend[
		*AvailabilityStore,
		*BeaconBlock,
		*BeaconBlockBody,
		*BeaconBlockHeader,
		*BeaconState,
		*BeaconStateMarshallable,
		*BlobSidecars,
		*BlockStore,
		*Deposit,
		*DepositStore,
		*Eth1Data,
		*ExecutionPayloadHeader,
		*Fork,
		*KVStore,
		*Validator,
		*Withdrawal,
		WithdrawalCredentials,
	](
		node,
		in.StorageBackend,
	)
	return backend.New[
		*AvailabilityStore,
		*BeaconBlock,
		*BeaconBlockBody,
		*BeaconBlockHeader,
		*BeaconState,
		*BeaconStateMarshallable,
		*BlobSidecars,
		*BlockStore,
		*Deposit,
		*DepositStore,
		*Eth1Data,
		*ExecutionPayloadHeader,
		*Fork,
		*KVStore,
		*Validator,
		*Withdrawal,
		WithdrawalCredentials,
	](storageBackend)
}

type NodeAPIServerInput struct {
	depinject.In

	Config         *config.Config
	NodeAPIBackend *NodeAPIBackend
}

func ProvideNodeAPIServer(in NodeAPIServerInput) *NodeAPIServer {
	return server.New[
		*AvailabilityStore,
		*BeaconBlock,
		*BeaconBlockBody,
		*BeaconBlockHeader,
		*BeaconState,
		*BeaconStateMarshallable,
		*BlobSidecars,
		*BlockStore,
		*Deposit,
		*DepositStore,
		*Eth1Data,
		*ExecutionPayloadHeader,
		*Fork,
		*KVStore,
		*Validator,
		*Withdrawal,
		WithdrawalCredentials,
	](
		in.Config.NodeAPI,
		in.NodeAPIBackend,
		middleware.DefaultCORSConfig,
		middleware.DefaultLoggerConfig,
	)
}
