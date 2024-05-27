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

	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/storage/pkg/interfaces"
)

// Pruner is a struct that holds the prunable interface and a notifier channel.
type Pruner struct {
	prunable      interfaces.Prunable
	pruneRequests chan uint64
	logger        log.Logger[any]
	name          string
}

func NewPruner(logger log.Logger[any], prunable interfaces.Prunable) *Pruner {
	return &Pruner{
		logger:        logger,
		prunable:      prunable,
		pruneRequests: make(chan uint64),
	}
}

// Start starts the pruner by listening for new indexes to prune.
func (p *Pruner) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case index := <-p.pruneRequests:
				if err := p.prunable.Prune(index); err != nil {
					p.logger.Error("Error pruning", "error", err)
				}
			}
		}
	}()
}

// Notify sends a new index to the pruner through the notifier channel.
func (p *Pruner) Notify(index uint64) {
	p.pruneRequests <- index
}

// WithLogger sets the logger for the pruner.
func (p *Pruner) WithLogger(l log.Logger[any]) {
	p.logger = l
}

// Name returns the name of the pruner.
func (p *Pruner) Name() string {
	return p.name
}
