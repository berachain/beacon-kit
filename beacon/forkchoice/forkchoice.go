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

package forkchoice

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

// ForkChoicer represents the full fork choice interface composed of all the
// sub-interfaces.
type ForkChoicer interface {
	Reader
	Writer

	WithContext(ctx context.Context) ForkChoicer
}

type Reader interface {
	JustifiedPayloadBlockHash() common.Hash
	FinalizedPayloadBlockHash() common.Hash

	// TODO: eventually deprecate this.
	HeadBeaconBlock() [32]byte
}

type Writer interface {
	BlockProcessor

	// Eventually deprecate this
	UpdateHeadBeaconBlock([32]byte)
}

// BlockProcessor processes the block that's used for accounting fork choice.
type BlockProcessor interface {
	// InsertNode inserts a new node into the forkchoice.
	InsertNode(common.Hash) error
}
