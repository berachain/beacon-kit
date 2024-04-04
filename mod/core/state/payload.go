package state

import types0 "github.com/berachain/beacon-kit/mod/execution/types"

func (s *BeaconStateDeneb) GetExecutionPayload() (types0.ExecutionPayload, error) {
	return s.LatestExecutionPayload, nil
}
