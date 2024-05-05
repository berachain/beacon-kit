package pruner

import (
	"fmt"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/storage/pkg/filedb"
	"github.com/berachain/beacon-kit/mod/storage/pkg/interfaces"
	"time"
)

type Pruner struct {
	Interval  time.Duration // Interval at which the pruner runs
	Ticker    *time.Ticker  // Ticker for the pruner
	Quit      chan struct{}
	chainSpec primitives.ChainSpec
	db        interfaces.DB
}

func NewPruner(interval time.Duration, db interfaces.DB) *Pruner {
	return &Pruner{
		Interval: interval,
		Ticker:   time.NewTicker(interval),
		Quit:     make(chan struct{}),
		db:       db,
	}
}

func (p *Pruner) Start() error {
	for {
		select {
		case <-p.Ticker.C:
			// Do the pruning
			db, ok := p.db.(*filedb.DB)
			if !ok {
				return fmt.Errorf("DB is not a *filedb.DB instance")
			}
			err := p.prune(db)
			if err != nil {
				return err
			}
		case <-p.Quit:
			p.Ticker.Stop()
			return nil
		}
	}
}

func (p *Pruner) Stop() {
	close(p.Quit)
}

func (p *Pruner) prune(db *filedb.DB) error {

	fmt.Println("Pruning blobs")

	highestIndex := db.GetHighestSlot()

	fmt.Println("Highest index: ", highestIndex)

	lowestIndex := db.GetLowestSlot()

	fmt.Println("lowestIndex index: ", lowestIndex)

	// Calculate the difference between the highest and lowest indices.
	diff := highestIndex - lowestIndex

	fmt.Println("diff : ", diff)

	// Get the minimum epochs for blobs sidecars request from the chain spec.
	minEpochs := p.chainSpec.MinEpochsForBlobsSidecarsRequest()

	fmt.Println("minEpochs : ", minEpochs)

	rangeDB := &filedb.RangeDB{
		DB: p.db,
	}

	// Convert the minimum epochs to slots.
	minEpochsInSots := minEpochs * p.chainSpec.SlotsPerEpoch()
	fmt.Println("minEpochsInSots : ", minEpochsInSots)

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
