package capella

import (
	"github.com/itsdevbear/bolaris/types/consensus/version"
	"github.com/itsdevbear/bolaris/types/engine/interfaces"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
)

// WrappedPayloadAttributesV2 wraps the PayloadAttributesV2 from the Prysmatic Labs' Engine API v1.
var _ interfaces.PayloadAttributer = (*WrappedPayloadAttributesV2)(nil)

// WrappedPayloadAttributesV2 is a struct that embeds enginev1.PayloadAttributesV2
// to provide additional functionality required by the PayloadAttributer interface.
type WrappedPayloadAttributesV2 struct {
	enginev1.PayloadAttributesV2
}

// Version returns the consensus version for the Capella upgrade.
func (w *WrappedPayloadAttributesV2) Version() int {
	return version.Capella
}
