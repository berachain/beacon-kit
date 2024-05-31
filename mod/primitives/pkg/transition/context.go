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

package transition

import "context"

// Context is the context for the state transition.
type Context struct {
	context.Context
	// OptimisticEngine indicates whether to optimistically assume
	// the execution client has the correct state certain errors
	// are returned by the execution engine.
	OptimisticEngine bool
	// SkipPayloadVerification indicates whether to skip calling NewPayload
	// on the execution client. This can be done when the node is not
	// syncing, and the payload is already known to the execution client.
	SkipPayloadVerification bool
	// SkipValidateRandao indicates whether to skip validating the Randao mix.
	SkipValidateRandao bool
	// SkipValidateResult indicates whether to validate the result of
	// the state transition.
	SkipValidateResult bool
}

// GetOptimisticEngine returns whether to optimistically assume the execution
// client has the correct state when certain errors are returned by the
// execution engine.
func (c *Context) GetOptimisticEngine() bool {
	return c.OptimisticEngine
}

// GetSkipPayloadVerification returns whether to skip calling NewPayload on the
// execution client. This can be done when the node is not syncing, and the
// payload is already known to the execution client.
func (c *Context) GetSkipPayloadVerification() bool {
	return c.SkipPayloadVerification
}

// GetSkipValidateRandao returns whether to skip validating the Randao mix.
func (c *Context) GetSkipValidateRandao() bool {
	return c.SkipValidateRandao
}

// GetSkipValidateResult returns whether to validate the result of the state
// transition.
func (c *Context) GetSkipValidateResult() bool {
	return c.SkipValidateResult
}

// Unwrap returns the underlying standard context.
func (c *Context) Unwrap() context.Context {
	return c.Context
}
