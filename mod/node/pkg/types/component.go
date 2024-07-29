package types

import (
	"context"
)

type NodeComponent[AppT App] interface {
	Start(context.Context) error
}
