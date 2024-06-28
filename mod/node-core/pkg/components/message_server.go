package components

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/comet"
)

// MessageServerInput is the input for the dep inject framework.
type MessageServerInput struct {
	depinject.In
	Environment appmodule.Environment
}

// ProvideMessageServer is a function that provides the message server to the application.
func ProvideMessageServer(
	in MessageServerInput,
) *comet.MsgServer {
	return comet.NewMsgServer(in.Environment.EventService)
}
