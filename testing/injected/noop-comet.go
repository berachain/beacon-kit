package injected

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type NoopComet struct{}

func (n *NoopComet) Start(ctx context.Context) error {
	// TODO implement me
	panic("implement me")
}

func (n *NoopComet) Stop() error {
	// TODO implement me
	panic("implement me")
}

func (n *NoopComet) Name() string {
	// TODO implement me
	panic("implement me")
}

func (n *NoopComet) CreateQueryContext(height int64, prove bool) (sdk.Context, error) {
	// TODO implement me
	panic("implement me")
}

func (n *NoopComet) LastBlockHeight() int64 {
	// TODO implement me
	panic("implement me")
}
