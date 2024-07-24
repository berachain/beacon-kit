package app

import (
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/service"
)

type App struct {
	logger log.Logger[any]
	service.Registry
}

