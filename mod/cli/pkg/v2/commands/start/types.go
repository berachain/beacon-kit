package start

import "context"

type Node interface {
	Start(context.Context) error
}
