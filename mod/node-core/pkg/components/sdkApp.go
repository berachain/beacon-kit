package components

import (
	"cosmossdk.io/core/transaction"
	"cosmossdk.io/depinject"
	"cosmossdk.io/runtime/v2"
)

type SDKAppInput[T transaction.Tx] struct {
	depinject.In
	AppBuilder *runtime.AppBuilder[T]
}

func ProvideSDKApp[T transaction.Tx](
	in SDKAppInput[T],
) (*runtime.App[T], error) {
	return in.AppBuilder.Build()
}
