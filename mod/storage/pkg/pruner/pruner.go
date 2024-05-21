package pruner

import "github.com/berachain/beacon-kit/mod/storage/pkg/interfaces"

// might be unneccessary
type BasePruner struct {
	cfg      *Config
	prunable interfaces.Prunable
}
