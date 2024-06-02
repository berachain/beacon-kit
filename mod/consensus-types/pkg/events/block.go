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

package events

import (
	"context"
)

// Block represents a generic block in the beacon chain.
type Block[BeaconBlockT any] struct {
	// ctx is the context associated with the block.
	ctx context.Context
	// block is the actual beacon block.
	block BeaconBlockT
}

// NewBlock creates a new Block with the given context and beacon block.
func NewBlock[
	BeaconBlockT any,
](ctx context.Context, block BeaconBlockT) Block[BeaconBlockT] {
	return Block[BeaconBlockT]{
		ctx:   ctx,
		block: block,
	}
}

// Context returns the context associated with the block.
func (e Block[BeaconBlockT]) Context() context.Context {
	return e.ctx
}

// Block returns the beacon block.
func (e Block[BeaconBlockT]) Block() BeaconBlockT {
	return e.block
}
