package cometbft

import (
	"encoding/json"

	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	v1 "github.com/cometbft/cometbft/api/cometbft/types/v1"
	cmttypes "github.com/cometbft/cometbft/types"
)

type ConsensusParams struct {
	*cmttypes.ConsensusParams
}

func DefaultConsensusParams() *ConsensusParams {
	params := cmttypes.DefaultConsensusParams()
	params.Validator.PubKeyTypes = []string{crypto.CometBLSType}
	return &ConsensusParams{
		ConsensusParams: params,
	}
}

func (cp *ConsensusParams) Default() *ConsensusParams {
	return DefaultConsensusParams()
}

func (cp *ConsensusParams) MarshalJSON() ([]byte, error) {
	return json.Marshal(cp.ConsensusParams)
}

type ConsensusParamsStore struct {
	chainSpec common.ChainSpec
}

func NewConsensusParamsStore(chainSpec common.ChainSpec) *ConsensusParamsStore {
	return &ConsensusParamsStore{
		chainSpec: chainSpec,
	}
}

func (cps *ConsensusParamsStore) Get(slot uint64) (*v1.ConsensusParams, error) {
	params, ok := cps.chainSpec.GetCometBFTConfigForSlot(math.U64(slot)).(*ConsensusParams)
	if !ok {
		return nil, errors.Newf("failed to get consensus params for slot %d", slot)
	}
	p := params.ConsensusParams.ToProto()
	return &p, nil
}
