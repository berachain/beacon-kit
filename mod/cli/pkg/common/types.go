package common

import (
	"io"

	"cosmossdk.io/log"
	dbm "github.com/cosmos/cosmos-db"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
)

type (
	// Application defines an application interface that wraps abci.Application.
	Application = servertypes.Application
	// AppCreator is a function that allows us to lazily initialize an
	// application using various configurations.
	AppCreator[T Application] func(log.Logger, dbm.DB, io.Writer, AppOptions) T

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
)
