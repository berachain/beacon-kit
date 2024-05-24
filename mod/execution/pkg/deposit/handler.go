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

package deposit

import (
	"context"

	"github.com/berachain/beacon-kit/mod/beacon/dispatcher"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core/state"
)

// Compile time assertion to ensure Handler implements dispatcher.Handler
// var _ dispatcher.Handler = (*Handler)(nil)

type Handler[BeaconStateT state.BeaconState] struct {
	events chan dispatcher.Event

	cs     primitives.ChainSpec
	signer crypto.BLSSigner
}

func NewHandler[BeaconStateT state.BeaconState](
	cs primitives.ChainSpec,
	signer crypto.BLSSigner,
) *Handler[BeaconStateT] {
	return &Handler[BeaconStateT]{
		events: make(chan dispatcher.Event),
		cs:     cs,
		signer: signer,
	}
}

func (h *Handler[BeaconStateT]) Start(ctx context.Context) error {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-h.events:
				depositEvent, ok := event.(*Event[BeaconStateT])
				if !ok {
					panic(ErrInvalidEventType)
				}

				if err := h.processDeposits(
					depositEvent.st.(BeaconStateT),
					depositEvent.deposits,
				); err != nil {
					panic(ErrInvalidEventType)
				}
			}
		}
	}()

	return nil
}

func (h *Handler[BeaconStateT]) Notify(event dispatcher.Event) {
	h.events <- event
}
