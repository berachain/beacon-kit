package manager

import "github.com/berachain/beacon-kit/mod/storage/pkg/interfaces"

// TODO - Make DBManager implement Service interface and register it with the
// registry
type DBManager struct {
	Pruners []interfaces.Pruner
}

func NewDBManager(pruners []interfaces.Pruner) *DBManager {
	return &DBManager{Pruners: pruners}
}
