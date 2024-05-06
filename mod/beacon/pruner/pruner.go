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

package pruner

import (
	"context"
	"time"

	"github.com/berachain/beacon-kit/mod/node-builder/service"
	"github.com/berachain/beacon-kit/mod/storage/pkg/prunedb"
)

// Service is the blobs pruning service.
type Service struct {
	service.BaseService
	db     *prunedb.DB
	ticker *time.Ticker

	windowSize      uint64
	lastPrunedIndex uint64
}

func (s *Service) Start(ctx context.Context) {
	go func() {
		defer s.ticker.Stop()

		for {
			select {
			case <-s.ticker.C:
				// Do the pruning
				if err := s.prune(); err != nil {
					s.Logger().Error("Error pruning: ", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (s *Service) prune() error {
	// Increment the start window and end window
	// if the latest index is greater than the end window.
	latest := s.db.GetLatestIndex()
	left := max(0, latest-s.windowSize)

	// If the last pruned index is greater than the left side
	// of the window, we can skip.
	if left <= s.lastPrunedIndex {
		return nil
	}

	// Prune the range and set the previously pruned index.
	s.Logger().Info("🙈 da monkeys are pruning da blobs 🙉")
	if err := s.db.DeleteRange(s.lastPrunedIndex, left); err != nil {
		return err
	}
	s.lastPrunedIndex = left - 1

	return nil
}
