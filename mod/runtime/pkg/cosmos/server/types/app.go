package types

import (
	"encoding/json"
	"io"

	cmtproto "github.com/cometbft/cometbft/api/cometbft/types/v1"
	cmttypes "github.com/cometbft/cometbft/types"
	dbm "github.com/cosmos/cosmos-db"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
)

type (
	// AppOptions defines an interface that is passed into an application
	// constructor, typically used to set BaseApp options that are either supplied
	// via config file or through CLI arguments/flags. The underlying implementation
	// is defined by the server package and is typically implemented via a Viper
	// literal defined on the server Context. Note, casting Get calls may not yield
	// the expected types and could result in type assertion errors. It is recommend
	// to either use the cast package or perform manual conversion for safety.
	AppOptions interface {
		Get(string) interface{}
	}

	// Application defines an application interface that wraps abci.Application.
	// The interface defines the necessary contracts to be implemented in order
	// to fully bootstrap and start an application.
	Application interface {
		ABCI

		// CommitMultiStore return the multistore instance
		CommitMultiStore() storetypes.CommitMultiStore

		// Close is called in start cmd to gracefully cleanup resources.
		// Must be safe to be called multiple times.
		Close() error
	}

	// AppCreator is a function that allows us to lazily initialize an
	// application using various configurations.
	AppCreator[T Application] func(log.Logger, dbm.DB, io.Writer, AppOptions) T

	// ExportedApp represents an exported app state, along with
	// validators, consensus params and latest app height.
	ExportedApp struct {
		// AppState is the application state as JSON.
		AppState json.RawMessage
		// Validators is the exported validator set.
		Validators []cmttypes.GenesisValidator
		// Height is the app's latest block height.
		Height int64
		// ConsensusParams are the exported consensus params for ABCI.
		ConsensusParams cmtproto.ConsensusParams
	}

	// AppExporter is a function that dumps all app state to
	// JSON-serializable structure and returns the current validator set.
	AppExporter func(
		logger log.Logger,
		db dbm.DB,
		traceWriter io.Writer,
		height int64,
		forZeroHeight bool,
		jailAllowedAddrs []string,
		opts AppOptions,
		modulesToExport []string,
	) (ExportedApp, error)
)
