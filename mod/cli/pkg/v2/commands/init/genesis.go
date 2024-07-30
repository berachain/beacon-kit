package initcli

import (
	"encoding/json"
	"os"
)

type Genesis[
	ConsensusParamsT interface{ json.Marshaler },
	GenesisStateT interface{ json.Marshaler },
] struct {
	State           GenesisStateT    `json:"state"`
	ConsensusParams ConsensusParamsT `json:"consensus_params"`
}

func (g *Genesis[_, _]) Save(filePath string) error {
	bz, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, bz, 0644)
}
