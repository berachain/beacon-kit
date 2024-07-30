package types

import "context"

type Broker interface {
	Name() string
	Start(context.Context) error
}
