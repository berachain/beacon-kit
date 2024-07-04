package components

import (
	"fmt"

	"cosmossdk.io/core/transaction"
	"cosmossdk.io/depinject"
	"cosmossdk.io/runtime/v2"
	bkappmanager "github.com/berachain/beacon-kit/mod/node-core/pkg/components/appmanager"
)

type SDKAppInput[T transaction.Tx] struct {
	depinject.In
	AppBuilder *runtime.AppBuilder[T]
	Middleware *ABCIMiddleware
}

func ProvideSDKApp[T transaction.Tx](
	in SDKAppInput[T],
) (*runtime.App[T], error) {
	fmt.Println("PP BUILDER", in.AppBuilder)
	app, err := in.AppBuilder.Build()
	if err != nil {
		return nil, err
	}
	// set app manager
	appManager := bkappmanager.NewAppManager(
		app.GetAppManager(),
		in.Middleware,
	)
	app.AppManager = appManager
	return app, nil
}
