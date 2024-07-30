package dispatcher

import (
	"context"
	"sync"

	"github.com/berachain/beacon-kit/mod/async/pkg/broker"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/service"
)

var _ service.Basic = &Dispatcher{}

const (
	name = "dispatcher"
)

type Dispatcher struct {
	mu             sync.RWMutex
	brokerRegistry *broker.Registry
}

func NewDispatcher(
	brokerRegistry *broker.Registry,
) *Dispatcher {
	return &Dispatcher{
		mu:             sync.RWMutex{},
		brokerRegistry: brokerRegistry,
	}
}

func (d *Dispatcher) Start(ctx context.Context) error {
	return nil
}

func (*Dispatcher) Name() string {
	return name
}
