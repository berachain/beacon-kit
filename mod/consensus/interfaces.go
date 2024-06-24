package consensus

import (
	"context"

	"github.com/gogo/protobuf/proto"
)

type Consensus interface {
	InitGenesis(ctx context.Context, genesisBz []byte)
	Prepare(ctx context.Context, msg proto.Message) (proto.Message, error)
}
