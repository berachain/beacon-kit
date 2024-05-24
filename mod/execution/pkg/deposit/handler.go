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

	"github.com/berachain/beacon-kit/mod/beacon/blockchain"
	"github.com/berachain/beacon-kit/mod/beacon/dispatcher"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

const ServiceName = "deposit_handler"

type Handler[
	DepositStoreT DepositStore,
] struct {
	ch <-chan dispatcher.Event

	cs     primitives.ChainSpec
	signer crypto.BLSSigner
	sb     StorageBackend[
		any, any, any, DepositStoreT,
	]
	dc     DepositContract
	logger log.Logger[any]
}

func NewHandler[
	DepositStoreT DepositStore,
](
	sb StorageBackend[
		any, any, any, DepositStoreT,
	],
	dc DepositContract,
	ch <-chan dispatcher.Event,
	logger log.Logger[any],
) *Handler[DepositStoreT] {
	return &Handler[DepositStoreT]{
		ch:     ch,
		sb:     sb,
		dc:     dc,
		logger: logger,
	}
}

func (h *Handler[DepositStoreT]) Start(ctx context.Context) error {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case e := <-h.ch:
				event, ok := e.(*blockchain.Event)
				if !ok {
					panic(ErrInvalidEventType)
				}

				if err := h.Process(event.Ctx, event.Slot); err != nil {
					panic(err)
				}
			}
		}
	}()

	return nil
}

func (h *Handler[DepositStoreT]) Name() string {
	return ServiceName
}

func (h *Handler[DepositStoreT]) Status() error {
	return nil
}

func (h *Handler[DepositStoreT]) WaitForHealthy(ctx context.Context) {}

func (h *Handler[DepositStoreT]) Process(
	ctx context.Context,
	slot math.U64,
) error {
	h.logger.Info("💵 processing deposits 💵", "slot", slot)
	deposits, err := h.dc.GetDeposits(ctx, slot.Unwrap())
	if err != nil {
		return err
	}

	if err := h.sb.DepositStore(ctx).EnqueueDeposits(
		deposits,
	); err != nil {
		return err
	}
	return nil
}
