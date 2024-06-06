package components

import (
	"cosmossdk.io/core/log"
	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	dablob "github.com/berachain/beacon-kit/mod/da/pkg/blob"
	"github.com/berachain/beacon-kit/mod/da/pkg/kzg"
	dastore "github.com/berachain/beacon-kit/mod/da/pkg/store"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/components/metrics"
	"github.com/berachain/beacon-kit/mod/primitives"
)

type BlobProcessorInput struct {
	depinject.In
	Logger            log.Logger
	ChainSpec         primitives.ChainSpec
	BlobProofVerifier kzg.BlobProofVerifier
	TelemetrySink     *metrics.TelemetrySink
}

func ProvideBlobProcessor(
	in BlobProcessorInput,
) *dablob.Processor[
	*dastore.Store[types.BeaconBlockBody],
	types.BeaconBlockBody,
] {
	return dablob.NewProcessor[
		*dastore.Store[types.BeaconBlockBody],
		types.BeaconBlockBody,
	](
		in.Logger.With("service", "blob-processor"),
		in.ChainSpec,
		dablob.NewVerifier(in.BlobProofVerifier, in.TelemetrySink),
		types.BlockBodyKZGOffset,
		in.TelemetrySink,
	)
}
