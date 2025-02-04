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

// ReadOnlyContext defines an interface for managing state transition context.
type ReadOnlyContext interface {
	ConsensusCtx() context.Context
	ConsensusTime() math.U64
	ProposerAddress() []byte
	VerifyPayload() bool
	VerifyRandao() bool
	VerifyResult() bool
	MeterGas() bool
	OptimisticEngine() bool
}

// Context is the context for the state transition.
type Context struct {
	// consensusCtx is the context passed by CometBFT callbacks
	// We pass it down to be able to cancel processing (although
	// currently CometBFT context is set to TODO)
	consensusCtx context.Context
	// consensusTime returns the timestamp of current consensus request.
	// It is used to build next payload and to validate currentpayload.
	consensusTime math.U64
	// Address of current block proposer
	proposerAddress []byte

	// verifyPayload indicates whether to call NewPayload on the
	// execution client. This can be done when the node is not
	// syncing, and the payload is already known to the execution client.
	verifyPayload bool
	// verifyRandao indicates whether to validate the Randao mix.
	verifyRandao bool
	// verifyResult indicates whether to validate the result of
	// the state transition.
	verifyResult bool

	// meterGas controls whether gas data related to the execution
	// layer payload should be meter or not. We currently meter only
	// finalized blocks.
	meterGas bool
	// optimisticEngine indicates whether to optimistically assume
	// the execution client has the correct state certain errors
	// are returned by the execution engine.
	optimisticEngine bool
}

func NewTransitionCtx(
	consensusCtx context.Context,
	time math.U64,
	address []byte,
) *Context {
	return &Context{
		consensusCtx:    consensusCtx,
		consensusTime:   time,
		proposerAddress: address,

		// by default we don't meter gas
		// (we care only about finalized blocks gas)
		meterGas: false,

		// by default we don't have optimistic engine
		// as it basically mute some checks
		optimisticEngine: false,

		// by default we keep all verification
		verifyPayload: true,
		verifyRandao:  true,
		verifyResult:  true,
	}
}

// Setters to control context attributes.
func (c *Context) WithMeterGas(meter bool) *Context {
	c.meterGas = meter
	return c
}

func (c *Context) WithOptimisticEngine(optimistic bool) *Context {
	c.optimisticEngine = optimistic
	return c
}

func (c *Context) WithVerifyPayload(verifyPayload bool) *Context {
	c.verifyPayload = verifyPayload
	return c
}

func (c *Context) WithVerifyRandao(verifyRandao bool) *Context {
	c.verifyRandao = verifyRandao
	return c
}

func (c *Context) WithVerifyResult(verifyResult bool) *Context {
	c.verifyResult = verifyResult
	return c
}

// Getters of context attributes.
func (c *Context) ConsensusCtx() context.Context {
	return c.consensusCtx
}

func (c *Context) ConsensusTime() math.U64 {
	return c.consensusTime
}

func (c *Context) ProposerAddress() []byte {
	return c.proposerAddress
}

func (c *Context) VerifyPayload() bool {
	return c.verifyPayload
}

func (c *Context) VerifyRandao() bool {
	return c.verifyRandao
}

func (c *Context) VerifyResult() bool {
	return c.verifyResult
}

func (c *Context) MeterGas() bool {
	return c.meterGas
}

func (c *Context) OptimisticEngine() bool {
	return c.optimisticEngine
}
