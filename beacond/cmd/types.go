package main

import (
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/log/pkg/phuslu"
	nodetypes "github.com/berachain/beacon-kit/mod/node-core/pkg/types"
)

type (
	Node             = nodetypes.Node
	ExecutionPayload = types.ExecutionPayload

	Logger       = phuslu.Logger
	LoggerConfig = phuslu.Config
)
