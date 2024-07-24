package app

import "context"

type Application interface {
	Start(context.Context) error
}
