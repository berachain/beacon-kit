package components

import (
	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/mod/async/pkg/broker"
	"github.com/berachain/beacon-kit/mod/async/pkg/dispatcher"
)

type BrokerRegistryInput struct {
	depinject.In

	SidecarsBroker        *SidecarsBroker
	BlockBroker           *BlockBroker
	GenesisBroker         *GenesisBroker
	SlotBroker            *SlotBroker
	StatusBroker          *StatusBroker
	ValidatorUpdateBroker *ValidatorUpdateBroker
}

func ProvideBrokerRegistry(in BrokerRegistryInput) *BrokerRegistry {
	brokerRegistry := broker.NewRegistry()
	brokerRegistry.RegisterBroker(in.SidecarsBroker)
	brokerRegistry.RegisterBroker(in.BlockBroker)
	brokerRegistry.RegisterBroker(in.GenesisBroker)
	brokerRegistry.RegisterBroker(in.SlotBroker)
	brokerRegistry.RegisterBroker(in.StatusBroker)
	brokerRegistry.RegisterBroker(in.ValidatorUpdateBroker)
	return brokerRegistry
}

type DispatcherInput struct {
	depinject.In
	BrokerRegistry *BrokerRegistry
}

func ProvideDispatcher(input DispatcherInput) *Dispatcher {
	return dispatcher.NewDispatcher(input.BrokerRegistry)
}
