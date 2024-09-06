package middleware

import (
	"encoding/json"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"

	consruntimetypes "github.com/berachain/beacon-kit/mod/consensus/pkg/types"
	datypes "github.com/berachain/beacon-kit/mod/da/pkg/types"
)

type (
	BeaconBlockT     = *types.BeaconBlock
	BlobSidecarsT    = *datypes.BlobSidecars
	validatorUpdates = transition.ValidatorUpdates
	GenesisT         = json.Unmarshaler
	SlotDataT        = *consruntimetypes.SlotData[
		*types.AttestationData,
		*types.SlashingInfo,
	]
)
