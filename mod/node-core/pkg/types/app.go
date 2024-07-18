package types

import (
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Application is the interface that wraps the methods of the cosmos-sdk Application.
// It also adds a few methods for creating query contexts.
type Application interface {
	servertypes.Application

	CreateQueryContext(height int64, prove bool) (sdk.Context, error)
}
