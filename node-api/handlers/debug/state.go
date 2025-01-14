package debug

import (
	beacontypes "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/node-api/handlers/utils"
)

func (h *Handler[ContextT]) GetState(c ContextT) (any, error) {
	req, err := utils.BindAndValidate[beacontypes.GetStateRequest](
		c, h.Logger(),
	)
	if err != nil {
		return nil, err
	}
	slot, err := utils.SlotFromStateID(req.StateID, h.backend)
	if err != nil {
		return nil, err
	}
	stateRoot, err := h.backend.StateAtSlot(slot)
	if err != nil {
		return nil, err
	}
	return beacontypes.StateResponse{
		ExecutionOptimistic: false, // stubbed
		Finalized:           false, // stubbed
		Data:                stateRoot,
	}, nil
}
