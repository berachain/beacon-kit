package types

import "context"

type Node interface {
	Start(context.Context) error
}
