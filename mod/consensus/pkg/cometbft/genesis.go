package cometbft

import (
	"time"

	cmttypes "github.com/cometbft/cometbft/types"
)

type Genesis struct {
	*cmttypes.GenesisDoc `mapstructure:",squash"`
}

func NewGenesis(chainID string, appState []byte, consensusParams *cmttypes.ConsensusParams) *Genesis {
	return &Genesis{
		GenesisDoc: &cmttypes.GenesisDoc{
			GenesisTime:     time.Now(),
			ChainID:         chainID,
			InitialHeight:   1,
			ConsensusParams: consensusParams,
			AppState:        appState,
		},
	}
}

// func (g *Genesis) UnmarshalJSON(data []byte) error {
// 	return cmtjson.Unmarshal(data, g.GenesisDoc)
// }

func (g *Genesis) Export(path string) error {
	if err := g.ValidateAndComplete(); err != nil {
		return err
	}
	return g.SaveAs(path)
}

func (g *Genesis) ToGenesisDoc() (*cmttypes.GenesisDoc, error) {
	return g.GenesisDoc, nil
}
