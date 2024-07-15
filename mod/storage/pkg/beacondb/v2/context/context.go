package context

import (
	"context"

	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/v2/changeset"
)

// Context is a wrapper around the context.Context that holds a changeset
type Context struct {
	context.Context
	Changeset *changeset.Changeset
}

// New initializes a new Context with an empty changeset
func New(ctx context.Context) *Context {
	return &Context{
		Context:   ctx,
		Changeset: changeset.New(),
	}
}

func (c *Context) Copy() *Context {
	return &Context{
		Context:   c.Context,
		Changeset: c.Changeset.Copy(),
	}
}
