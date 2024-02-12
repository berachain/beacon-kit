package deneb

import (
	"github.com/itsdevbear/bolaris/types/consensus/version"
	"github.com/itsdevbear/bolaris/types/engine/interfaces"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
)

// WrappedPayloadAttributesV2 wraps the PayloadAttributesV2 from the Prysmatic Labs' Engine API v1.
var _ interfaces.PayloadAttributer = (*WrappedPayloadAttributesV3)(nil)

// WrappedPayloadAttributesV2 is a struct that embeds enginev1.PayloadAttributesV2
// to provide additional functionality required by the PayloadAttributer interface.
type WrappedPayloadAttributesV3 struct {
	enginev1.PayloadAttributesV3
}

// Version returns the consensus version for the Capella upgrade.
func (w *WrappedPayloadAttributesV3) Version() int {
	return version.Deneb
}
