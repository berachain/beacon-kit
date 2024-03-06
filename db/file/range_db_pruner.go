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

package file

import "context"

// RangeDBPruner prunes old indexes from the range DB.
type RangeDBPruner[T numeric] struct {
	db             *RangeDB[T]
	notifyNewIndex chan T
	pruneWindow    T
}

func NewRangeDBPruner[T numeric](
	db *RangeDB[T],
	pruneWindow T,
	notifyNewIndex chan T,
) *RangeDBPruner[T] {
	return &RangeDBPruner[T]{
		db:             db,
		notifyNewIndex: notifyNewIndex,
		pruneWindow:    pruneWindow,
	}
}

// Start starts a goroutine that listens for new indices to prune.
func (p *RangeDBPruner[T]) Start(ctx context.Context) {
	go p.runLoop(ctx)
}

// NotifyNewIndex notifies the pruner of a new index.
func (p *RangeDBPruner[T]) runLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case index := <-p.notifyNewIndex:
			// Do nothing until the prune window is reached.
			if p.pruneWindow > index {
				continue
			}

			// Prune all indexes below the most recently seen
			// index minus the prune window.
			if err := p.db.DeleteRange(0, index-p.pruneWindow); err != nil {
				return
			}
		}
	}
}
