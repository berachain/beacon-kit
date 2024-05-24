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
	"fmt"

	"github.com/berachain/beacon-kit/mod/beacon/dispatcher"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core/state"
)

// Compile time assertion to ensure Handler implements dispatcher.Handler
// var _ dispatcher.Handler = (*Handler)(nil)

const ServiceName = "deposit_handler"

type Handler[BeaconStateT state.BeaconState] struct {
	ch     <-chan dispatcher.Event
	cs     primitives.ChainSpec
	signer crypto.BLSSigner
	logger log.Logger[any]
}

func NewHandler[BeaconStateT state.BeaconState](
	cs primitives.ChainSpec,
	signer crypto.BLSSigner,
	logger log.Logger[any],
	ch <-chan dispatcher.Event,
) *Handler[BeaconStateT] {
	return &Handler[BeaconStateT]{
		ch:     ch,
		cs:     cs,
		signer: signer,
		logger: logger,
	}
}

func (h *Handler[BeaconStateT]) Start(ctx context.Context) error {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-h.ch:
				depositEvent, ok := event.(*Event[BeaconStateT])
				if !ok {
					panic(ErrInvalidEventType)
				}

				st := depositEvent.st

				// Prune deposits.
				// TODO: This should be moved into a go-routine in the background.
				// Watching for logs should be completely decoupled as well.
				idx, err := st.GetEth1DepositIndex()
				if err != nil {
					return
				}

				fmt.Println(idx)

				// // TODO: pruner shouldn't be in main block processing thread.
				// if err = s.PruneDepositEvents(ctx, idx); err != nil {
				// 	return err
				// }

				// var lph engineprimitives.ExecutionPayloadHeader
				// lph, err = st.GetLatestExecutionPayloadHeader()
				// if err != nil {
				// 	return  err
				// }

				// // Process the logs from the previous blocks execution payload.
				// // TODO: This should be moved out of the main block processing flow.
				// // TODO: eth1FollowDistance should be done actually proper
				// eth1FollowDistance := math.U64(1)
				// if err = s.retrieveDepositsFromBlock(
				// 	ctx, lph.GetNumber()-eth1FollowDistance,
				// ); err != nil {
				// 	// s.logger.Error("failed to process logs", "error", err)
				// 	return err
				// } 

			}
		}
	}()

	return nil
}

func (h *Handler[BeaconStateT]) Name() string {
	return ServiceName
}

func (h *Handler[BeaconStateT]) Status() error {
	return nil
}

func (h *Handler[BeaconStateT]) WaitForHealthy(ctx context.Context) {}
