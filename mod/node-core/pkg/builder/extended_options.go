package builder

import (
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/comet"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// we can decouple from cosmos-sdk/types by having these options accept
// the middleware type instead, and just extract the handers. not really an
// issue rn tho

//nolint:lll // todo fix xD
func WithCometParamStore(chainSpec primitives.ChainSpec) func(bApp *baseapp.BaseApp) {
	return func(bApp *baseapp.BaseApp) {
		bApp.SetParamStore(comet.NewConsensusParamsStore(chainSpec))
	}
}

func WithPrepareProposal(handler sdk.PrepareProposalHandler) func(bApp *baseapp.BaseApp) {
	return func(bApp *baseapp.BaseApp) {
		bApp.SetPrepareProposal(handler)
	}
}

func WithProcessProposal(handler sdk.ProcessProposalHandler) func(bApp *baseapp.BaseApp) {
	return func(bApp *baseapp.BaseApp) {
		bApp.SetProcessProposal(handler)
	}
}

func WithPreBlocker(preBlocker sdk.PreBlocker) func(bApp *baseapp.BaseApp) {
	return func(bApp *baseapp.BaseApp) {
		bApp.SetPreBlocker(preBlocker)
	}
}
