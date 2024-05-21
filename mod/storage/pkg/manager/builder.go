package manager

import "github.com/berachain/beacon-kit/mod/storage/pkg/interfaces"

type DBManagerBuilder struct {
	DBManager *DBManager
	pruners   []interfaces.Pruner
}

func NewDBManagerBuilder() *DBManagerBuilder {
	return &DBManagerBuilder{}
}

func (b *DBManagerBuilder) Build() *DBManager {
	return b.DBManager
}

func (b *DBManagerBuilder) WithDBManager(dbManager *DBManager) *DBManagerBuilder {
	b.DBManager = dbManager
	return b
}

func (b *DBManagerBuilder) WithPruner(p interfaces.Pruner) *DBManagerBuilder {
	b.pruners = append(b.pruners, p)
	return b
}
