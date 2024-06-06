package components

import (
	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/mod/beacon/blockchain"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	datypes "github.com/berachain/beacon-kit/mod/da/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	execution "github.com/berachain/beacon-kit/mod/execution/pkg/engine"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
)

type StateProcessorInput struct {
	depinject.In
	ChainSpec       primitives.ChainSpec
	ExecutionEngine *execution.Engine[*types.ExecutionPayload]
	Signer          crypto.BLSSigner
}

// ProvideStateProcessor provides the state processor to the depinject
// framework.
func ProvideStateProcessor(
	in StateProcessorInput,
) blockchain.StateProcessor[
	*types.BeaconBlock,
	BeaconState,
	*datypes.BlobSidecars,
	*transition.Context,
	*types.Deposit,
] {
	return core.NewStateProcessor[
		*types.BeaconBlock,
		types.BeaconBlockBody,
		*types.BeaconBlockHeader,
		BeaconState,
		*datypes.BlobSidecars,
		*transition.Context,
		*types.Deposit,
		*types.ExecutionPayload,
		*types.ExecutionPayloadHeader,
		*types.Fork,
		*types.ForkData,
		*types.Validator,
		*engineprimitives.Withdrawal,
		types.WithdrawalCredentials,
	](
		in.ChainSpec,
		in.ExecutionEngine,
		in.Signer,
	)
}
