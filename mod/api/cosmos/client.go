package cosmos

import (
	"github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ChainQuerier struct {
	ContextGetter func(height int64, prove bool) (sdk.Context, error)
	ABCI          types.ABCI
}
