package vm

import (
	"github.com/ava-labs/avalanchego/utils/logging"
	"github.com/ava-labs/avalanchego/version"
	"github.com/ava-labs/avalanchego/vms"

	"github.com/berachain/beacon-kit/mod/consensus/pkg/miniavalanche/middleware"
)

var (
	_ vms.Factory = (*Factory)(nil)

	vmVersion = &version.Semantic{
		Major: 0,
		Minor: 0,
		Patch: 0,
	}
)

// Entry point for node to build the VM
type Factory struct {
	Config     Config
	Middleware middleware.VMMiddleware
}

func (f *Factory) New(logging.Logger) (interface{}, error) {
	return &VM{
		validators: f.Config.Validators,
		middleware: f.Middleware,
	}, nil
}
