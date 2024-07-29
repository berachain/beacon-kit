package types

import "context"

type App interface {
	Start(context.Context) error
}
