package broker

import (
	"context"
	"reflect"

	"github.com/berachain/beacon-kit/mod/async/pkg/types"
)

// grand hack central for now
type Registry struct {
	brokers map[string]types.Broker
}

func NewRegistry() *Registry {
	return &Registry{
		brokers: make(map[string]types.Broker),
	}
}

func (r *Registry) RegisterBroker(broker types.Broker) {
	r.brokers[broker.Name()] = broker
}

func (r *Registry) FetchBroker(broker types.Broker) error {
	brokerType := reflect.TypeOf(broker)
	if brokerType.Kind() != reflect.Ptr {
		return types.ErrInputIsNotPointer(brokerType)
	}

	elem := reflect.ValueOf(broker).Elem()

	typeName := ""
	for name, b := range r.brokers {
		bType := reflect.TypeOf(b)
		if bType.AssignableTo(bType.Elem()) {
			typeName = name
			break
		}
	}

	if typeName == "" {
		return types.ErrUnknownService(brokerType)
	}

	if registeredBroker, ok := r.brokers[typeName]; ok {
		elem.Set(reflect.ValueOf(registeredBroker))
	}

	return nil
}

func (s *Registry) StartAll(ctx context.Context) {
	for _, broker := range s.brokers {
		broker.Start(ctx)
	}
}
