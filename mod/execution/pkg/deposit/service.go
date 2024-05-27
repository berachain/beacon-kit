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

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/ethereum/go-ethereum/event"
)

type Service[DepositStoreT Store] struct {
	feed *event.FeedOf[types.BlockEvent]
	dc   Contract
	sb   StorageBackend[
		any, any, any, DepositStoreT,
	]

	logger log.Logger[any]
}

func NewService[DepositStoreT Store](
	feed *event.FeedOf[types.BlockEvent],
	logger log.Logger[any],
	sb StorageBackend[
		any, any, any, DepositStoreT,
	],
	dc Contract,
) *Service[DepositStoreT] {
	return &Service[DepositStoreT]{
		feed:   feed,
		logger: logger,
		sb:     sb,
		dc:     dc,
	}
}

func (s *Service[DepositStoreT]) Start(ctx context.Context) error {
	ch := make(chan types.BlockEvent)
	feed := s.feed.Subscribe(ch)
	go func() {
		for {
			select {
			case <-ctx.Done():
				feed.Unsubscribe()
			case event := <-ch:
				if err := s.handleDepositEvent(event); err != nil {
					s.logger.Error("failed to handle deposit event", "err", err)
				}
			}
		}
	}()

	return nil
}

func (s *Service[DepositStoreT]) Name() string {
	return "deposit-handler"
}

func (s *Service[DepositStoreT]) Status() error {
	return nil
}

func (s *Service[DepositStoreT]) WaitForHealthy(_ context.Context) {}
