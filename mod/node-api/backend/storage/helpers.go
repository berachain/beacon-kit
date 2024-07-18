package storage

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/hex"
)

const (
	StateIDGenesis   = "genesis"
	StateIDFinalized = "finalized"
	StateIDJustified = "justified"
	StateIDHead      = "head"
)

const (
	Head int64 = iota
	Genesis
)

func heightFromStateID(stateID string) (int64, error) {
	switch stateID {
	case StateIDFinalized, StateIDJustified, StateIDHead:
		return Head, nil
	case StateIDGenesis:
		return Genesis, nil
	default:
		slot, err := hex.String(stateID).ToUint64()
		if err != nil {
			return 0, err
		}
		return int64(slot), nil
	}
}
