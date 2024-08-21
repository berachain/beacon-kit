package dispatcher

import (
	"github.com/berachain/beacon-kit/mod/async/pkg/broker"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/async"
)

// RegisterEvent registers an event with the given eventID with the dispatcher.
func RegisterEvent[eventT async.BaseEvent](dispatcher *Dispatcher, eventID string) {
	dispatcher.RegisterBrokers(broker.New[eventT](eventID))
}
