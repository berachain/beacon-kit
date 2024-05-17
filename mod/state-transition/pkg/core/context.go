package core

import "context"

type Context interface {
	GetValidateResult() bool
	GetSkipPayloadIfExists() bool
	GetOptimisticEngine() bool
}

// Context is the context for the state transition.
type ctx struct {
	context.Context

	ValidateResult      bool
	SkipPayloadIfExists bool
	OptimisticEngine    bool
}

// NewContext creates a new context for the state transition.
func NewContext(
	stdctx context.Context,
	validateResult, skipIfPayloadIfExists, optimisticEngine bool,
) Context {
	return &ctx{
		Context:             stdctx,
		ValidateResult:      validateResult,
		SkipPayloadIfExists: skipIfPayloadIfExists,
		OptimisticEngine:    optimisticEngine,
	}
}

func (c *ctx) GetValidateResult() bool {
	return c.ValidateResult
}

func (c *ctx) GetSkipPayloadIfExists() bool {
	return c.SkipPayloadIfExists
}

func (c *ctx) GetOptimisticEngine() bool {
	return c.OptimisticEngine
}
