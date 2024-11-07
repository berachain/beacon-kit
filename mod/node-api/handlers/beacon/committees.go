package beacon

import (
	beacontypes "github.com/berachain/beacon-kit/mod/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/mod/node-api/handlers/types"
	"github.com/berachain/beacon-kit/mod/node-api/handlers/utils"
)

func (h *Handler[_, ContextT, _, _, _]) GetStateCommittees(
	c ContextT,
) (any, error) {
	req, err := utils.BindAndValidate[beacontypes.GetStateCommitteesRequest](
		c, h.Logger(),
	)
	if err != nil {
		return nil, err
	}
	slot, err := utils.SlotFromStateID(req.StateID, h.backend)
	if err != nil {
		return nil, err
	}
	// TODO : To be implemented
	//committees, err := h.backend.CommitteesByStateID(slot)
	//if err != nil {
	//	return nil, err
	//}

	// Stub the response
	return types.Wrap(beacontypes.CommitteeResponseData{
		Index:      0,
		Slot:       slot,
		Validators: []uint64{0},
	}), nil
}
