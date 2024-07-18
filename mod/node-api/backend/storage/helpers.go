package storage

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/hex"
)

const (
	Finalized int64 = iota
	Genesis
)

func heightFromStateID(stateID string) (int64, error) {
	switch stateID {
	case "finalized", "justified", "head":
		return Finalized, nil
	case "genesis":
		return Genesis, nil
	default:
		slot, err := hex.String(stateID).ToUint64()
		if err != nil {
			return 0, err
		}
		return int64(slot), nil
	}
}
