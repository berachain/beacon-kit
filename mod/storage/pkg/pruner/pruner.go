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
	"time"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/storage/pkg/filedb"
)

type Pruner struct {
	Interval time.Duration // Interval at which the pruner runs
	Ticker   *time.Ticker  // Ticker for the pruner
	Quit     chan struct{}
	// chainSpec       primitives.ChainSpec
	logger          log.Logger
	DB              *filedb.PrunerDB
	startWindow     uint64
	endWindow       uint64
	lastPrunedIndex uint64
}

// No other methods needs to be implemented

func NewPruner(
	interval time.Duration,
	db *filedb.PrunerDB,
	logger log.Logger,
	chainSpec primitives.ChainSpec) *Pruner {
	return &Pruner{
		Interval:    interval,
		Ticker:      time.NewTicker(interval),
		Quit:        make(chan struct{}),
		DB:          db,
		logger:      logger,
		startWindow: 0,
		// Setting endWindow to DA window (in slots)
		endWindow: chainSpec.SlotsPerEpoch() *
			chainSpec.MinEpochsForBlobsSidecarsRequest(),
		lastPrunedIndex: 0,
	}
}

func (p *Pruner) Start() {
	go func() {
		for {
			select {
			case <-p.Ticker.C:
				// Do the pruning
				err := p.prune(p.DB)
				if err != nil {
					p.logger.Error("Error pruning: ", err)
				}
			case <-p.Quit:
				p.Ticker.Stop()
				return
			}
		}
	}()
}

func (p *Pruner) Stop() {
	close(p.Quit)
}

func (p *Pruner) prune(db *filedb.PrunerDB) error {
	p.logger.Info("Pruning blobs")

	// Increment the start window and end window
	// if the latest index is greater than the end window.
	if db.LatestIndex > p.endWindow {
		difference := db.LatestIndex - p.endWindow
		p.startWindow += difference
		p.endWindow = db.LatestIndex
	}

	// If the difference is greater than the minimum epochs, prune the blobs.
	if p.startWindow <= p.lastPrunedIndex {
		return nil
	}
	err := db.DeleteRange(p.lastPrunedIndex, p.startWindow)
	if err != nil {
		return err
	}

	p.lastPrunedIndex = p.startWindow - 1
	return nil
}
