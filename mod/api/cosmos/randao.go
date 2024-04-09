package cosmos

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/berachain/beacon-kit/mod/api/beaconnode"
)

func (c ChainQuerier) GetStateRandao(params beaconnode.GetStateRandaoParams) beaconnode.GetStateRandaoRes {
	stateId := params.StateID
	if stateId == "" {
		return &beaconnode.GetStateRandaoBadRequest{
			Code:        1,
			Message:     "state_id is required in URL params",
			Stacktraces: nil,
		}
	}

	stateIdAsInt, err := strconv.ParseUint(stateId, 10, 64)
	if err != nil {
		return &beaconnode.GetStateRandaoBadRequest{
			Code:        2,
			Message:     "state_id must be a number",
			Stacktraces: nil,
		}
	}

	ctx, err := c.ContextGetter(int64(stateIdAsInt), false)
	if err != nil {
		return &beaconnode.GetStateRandaoInternalServerError{
			Code:        3,
			Message:     err.Error(),
			Stacktraces: nil,
		}
	}

	randao, err := c.Service.BeaconState(ctx).GetRandaoMixAtIndex(stateIdAsInt)
	if err != nil {
		return &beaconnode.GetStateRandaoInternalServerError{
			Code:        4,
			Message:     err.Error(),
			Stacktraces: nil,
		}
	}

	resp := &beaconnode.GetStateRandaoOK{
		ExecutionOptimistic: false,
		Finalized:           true,
		Data: beaconnode.GetStateRandaoOKData{
			Randao: hexutil.Encode(randao[:]),
		},
	}

	return resp
}
