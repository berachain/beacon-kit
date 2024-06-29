package components

import (
	"path/filepath"

	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	serverv2 "cosmossdk.io/server/v2"
	"cosmossdk.io/store/v2"
	"cosmossdk.io/store/v2/commitment/iavl"
	"cosmossdk.io/store/v2/db"
	"cosmossdk.io/store/v2/root"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cast"
)

type StoreConfigInput struct {
	depinject.In

	Logger  log.Logger
	AppOpts servertypes.AppOptions
}

func ProvideStoreOptions(in StoreConfigInput) *root.FactoryOptions {
	homeDir := cast.ToString(in.AppOpts.Get(serverv2.FlagHome))
	scRawDb, err := db.NewPebbleDB("application", filepath.Join(homeDir, "data"))
	if err != nil {
		panic(err)
	}
	return &root.FactoryOptions{
		Logger:  in.Logger,
		RootDir: homeDir,
		SSType:  0,
		SCType:  0,
		SCPruningOption: &store.PruningOption{
			KeepRecent: 0,
			Interval:   0,
		},
		IavlConfig: &iavl.Config{
			CacheSize:              100_000,
			SkipFastStorageUpgrade: true,
		},
		SCRawDB: scRawDb,
	}
}
