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
	"github.com/berachain/beacon-kit/mod/storage/pkg/interfaces"
)

type Pruner struct {
	Interval  time.Duration // Interval at which the pruner runs
	Ticker    *time.Ticker  // Ticker for the pruner
	Quit      chan struct{}
	chainSpec primitives.ChainSpec
	db        interfaces.DB
	logger    log.Logger
}

func NewPruner(interval time.Duration, db interfaces.DB, logger log.Logger) *Pruner {
	return &Pruner{
		Interval: interval,
		Ticker:   time.NewTicker(interval),
		Quit:     make(chan struct{}),
		db:       db,
		logger:   logger,
	}
}

func (p *Pruner) Start() {
	go func() {
		for {
			select {
			case <-p.Ticker.C:
				// Do the pruning
				db, ok := p.db.(*filedb.DB)
				if !ok {
					p.logger.Error("DB is not a *filedb.DB instance")
					continue
				}
				err := p.prune(db)
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

func (p *Pruner) prune(db *filedb.DB) error {
	p.logger.Info("Pruning blobs")
	highestIndex := db.GetHighestSlot()
	lowestIndex := db.GetLowestSlot()
	// Calculate the difference between the highest and lowest indices.
	diff := highestIndex - lowestIndex

	p.logger.Info("Highest index: ", highestIndex, " Lowest index: ", lowestIndex, " Difference: ", diff)
	// Get the minimum epochs for blobs sidecars request from the chain spec.
	minEpochs := p.chainSpec.MinEpochsForBlobsSidecarsRequest()

	p.logger.Info("Min epochs: ", minEpochs)

	rangeDB := &filedb.RangeDB{
		DB: p.db,
	}

	// Convert the minimum epochs to slots.
	minEpochsInSots := minEpochs * p.chainSpec.SlotsPerEpoch()
	p.logger.Info("minEpochsInSots: ", minEpochsInSots)

	minEpochsInSots = 100
	// If the difference is greater than the minimum epochs, prune the blobs.
	if diff > minEpochsInSots {
		err := rangeDB.DeleteRange(lowestIndex, highestIndex-minEpochs)
		if err != nil {
			return err
		}
	}
	return nil
}
