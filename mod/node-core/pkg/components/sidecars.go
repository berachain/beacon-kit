package components

import (
	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	dablob "github.com/berachain/beacon-kit/mod/da/pkg/blob"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/metrics"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
)

type SidecarFactoryInput struct {
	depinject.In
	ChainSpec     common.ChainSpec
	TelemetrySink *metrics.TelemetrySink
}

func ProvideSidecarFactory(in SidecarFactoryInput) *SidecarFactory {
	return dablob.NewSidecarFactory[
		*BeaconBlock,
		*BeaconBlockBody,
		*BeaconBlockHeader,
		*BlobsBundle,
		*BlobSidecar,
		*BlobSidecars,
		*Deposit,
		*Eth1Data,
		*ExecutionPayload,
	](
		in.ChainSpec,
		types.KZGPositionDeneb,
		in.TelemetrySink,
	)
}
