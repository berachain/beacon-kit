// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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

package execution

import (
	"context"

	"github.com/itsdevbear/bolaris/runtime/dispatch"
	"github.com/itsdevbear/bolaris/types"
)

// forkchoiceStoreProvider is an interface that wraps the basic ForkChoiceStore method.
type forkchoiceStoreProvider interface {
	// ForkChoiceStore returns a fork choice store in the provided context.
	ForkChoiceStore(ctx context.Context) types.ForkChoiceStore
}

// GrandCentralDispatch is an interface that wraps the basic GetQueue method.
// It is used to retrieve a dispatch queue by its ID.
type GrandCentralDispatch interface {
	// GetQueue returns a queue with the provided ID.
	GetQueue(id string) dispatch.Queue
}

// NotifyForkchoiceUpdateArg is the argument for the forkchoice
// update notification `notifyForkchoiceUpdate`.
type NotifyForkchoiceUpdateArg struct {
	// headHash is the hash of the head block we are building ontop of.=
	headHash []byte
	// safeHash is the hash of the last safe block.
	safeHash []byte
	// finalHash is the hash of the last finalized block.
	finalHash []byte
}

// NewNotifyForkchoiceUpdateArg creates a new NotifyForkchoiceUpdateArg.
func NewNotifyForkchoiceUpdateArg(
	headHash []byte,
	safeHash []byte,
	finalHash []byte,
) *NotifyForkchoiceUpdateArg {
	return &NotifyForkchoiceUpdateArg{
		headHash:  headHash,
		safeHash:  safeHash,
		finalHash: finalHash,
	}
}
