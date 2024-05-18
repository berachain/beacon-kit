// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package core

import "context"

// Context defines an interface for managing state transition context.
type Context interface {
	context.Context

	// GetValidateResult returns whether to validate the result of the state
	// transition.
	GetValidateResult() bool

	// GetSkipPayloadIfExists returns whether to skip verifying the payload if
	// it already exists on the execution client.
	GetSkipPayloadIfExists() bool

	// GetOptimisticEngine returns whether to optimistically assume the
	// execution client has the correct state when certain errors are returned
	// by the execution engine.
	GetOptimisticEngine() bool

	// WithContext sets the underlying context.
	WithContext(ctx context.Context) Context

	// Unwrap returns the underlying standard context.
	Unwrap() context.Context
}

// stateTransitionContext is the context for the state transition.
type stateTransitionContext struct {
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
) Context {
	return &stateTransitionContext{
		Context:             stdctx,
		ValidateResult:      validateResult,
		SkipPayloadIfExists: skipIfPayloadIfExists,
		OptimisticEngine:    optimisticEngine,
	}
}

// WithContext sets the underlying context.
func (c *stateTransitionContext) WithContext(ctx context.Context) Context {
	c.Context = ctx
	return c
}

// GetValidateResult returns whether to validate the result of the state
// transition.
func (c *stateTransitionContext) GetValidateResult() bool {
	return c.ValidateResult
}

// GetSkipPayloadIfExists returns whether to skip verifying the payload if it
// already exists on the execution client.
func (c *stateTransitionContext) GetSkipPayloadIfExists() bool {
	return c.SkipPayloadIfExists
}

// GetOptimisticEngine returns whether to optimistically assume the execution
// client has the correct state when certain errors are returned by the
// execution engine.
func (c *stateTransitionContext) GetOptimisticEngine() bool {
	return c.OptimisticEngine
}

// Unwrap returns the underlying standard context.
func (c *stateTransitionContext) Unwrap() context.Context {
	return c.Context
}
