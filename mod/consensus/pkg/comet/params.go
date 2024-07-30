package cometbft

import (
	"encoding/json"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
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
