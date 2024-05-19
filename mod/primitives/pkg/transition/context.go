package transition

import "context"

// Context is the context for the state transition.
type Context struct {
	context.Context

	// ValidateResult indicates whether to validate the result of
	// the state transition.
	ValidateResult bool

	// SkipPayloadIfExists indicates whether to skip verifying
	// the payload if it already exists on the execution client.
	SkipPayloadIfExists bool

	// OptimisticEngine indicates whether to optimistically assume
	// the execution client has the correct state certain errors
	// are returned by the execution engine.
	OptimisticEngine bool
}

// NewContext creates a new context for the state transition.
func NewContext(
	stdctx context.Context,
	validateResult, skipIfPayloadIfExists, optimisticEngine bool,
) *Context {
	return &Context{
		Context:             stdctx,
		ValidateResult:      validateResult,
		SkipPayloadIfExists: skipIfPayloadIfExists,
		OptimisticEngine:    optimisticEngine,
	}
}

// GetValidateResult returns whether to validate the result of the state
// transition.
func (c *Context) GetValidateResult() bool {
	return c.ValidateResult
}

// GetSkipPayloadIfExists returns whether to skip verifying the payload if it
// already exists on the execution client.
func (c *Context) GetSkipPayloadIfExists() bool {
	return c.SkipPayloadIfExists
}

// GetOptimisticEngine returns whether to optimistically assume the execution
// client has the correct state when certain errors are returned by the
// execution engine.
func (c *Context) GetOptimisticEngine() bool {
	return c.OptimisticEngine
}

// Unwrap returns the underlying standard context.
func (c *Context) Unwrap() context.Context {
	return c.Context
}
