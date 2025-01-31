// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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

import (
	"context"

	"github.com/berachain/beacon-kit/primitives/math"
)

// Context is the context for the state transition.
type Context struct {
	context.Context

	MeterGas bool
	// OptimisticEngine indicates whether to optimistically assume
	// the execution client has the correct state certain errors
	// are returned by the execution engine.
	OptimisticEngine bool
	// VerifyPayload indicates whether to call NewPayload on the
	// execution client. This can be done when the node is not
	// syncing, and the payload is already known to the execution client.
	VerifyPayload bool
	// ValidateRandao indicates whether to validate the Randao mix.
	ValidateRandao bool
	// ValidateResult indicates whether to validate the result of
	// the state transition.
	ValidateResult bool
	// Address of current block proposer
	ProposerAddress []byte
	// ConsensusTime returns the timestamp of current consensus request.
	// It is used to build next payload and to validate currentpayload.
	ConsensusTime math.U64
}

func (c *Context) GetMeterGas() bool {
	return c.MeterGas
}

// GetOptimisticEngine returns whether to optimistically assume the execution
// client has the correct state when certain errors are returned by the
// execution engine.
func (c *Context) GetOptimisticEngine() bool {
	return c.OptimisticEngine
}

// GetVerifyPayload returns whether to call NewPayload on the
// execution client. This can be done when the node is not syncing, and the
// payload is already known to the execution client.
func (c *Context) GetVerifyPayload() bool {
	return c.VerifyPayload
}

// GetValidateRandao returns whether to validate the Randao mix.
func (c *Context) GetValidateRandao() bool {
	return c.ValidateRandao
}

// GetValidateResult returns whether to validate the result of the state
// transition.
func (c *Context) GetValidateResult() bool {
	return c.ValidateResult
}

// GetProposerAddress returns the address of the validator
// selected by consensus to propose the block.
func (c *Context) GetProposerAddress() []byte {
	return c.ProposerAddress
}

// GetConsensusTime returns the timestamp of current consensus request.
// It is used to build next payload and to validate currentpayload.
func (c *Context) GetConsensusTime() math.U64 {
	return c.ConsensusTime
}

// Unwrap returns the underlying standard context.
func (c *Context) Unwrap() context.Context {
	return c.Context
}
