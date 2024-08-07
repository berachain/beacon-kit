// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

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
	// Persist indicates whether to persist the state changes.
	PersistState bool
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

func (c *Context) Persist() bool {
	return c.PersistState
}
