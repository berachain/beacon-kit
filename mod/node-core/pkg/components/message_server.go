package components

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/cometbft"
)

// MessageServerInput is the input for the dep inject framework.
type MessageServerInput struct {
	depinject.In
	Environment appmodule.Environment
	ChainSpec   common.ChainSpec
}

// ProvideMessageServer is a function that provides the message server to the application.
func ProvideMessageServer(
	in MessageServerInput,
) *cometbft.MsgServer {
	return cometbft.NewMsgServer(
		in.Environment.EventService,
		in.ChainSpec,
	)
}
