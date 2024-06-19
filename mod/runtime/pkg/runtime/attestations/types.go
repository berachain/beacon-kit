package attestations

import "github.com/berachain/beacon-kit/mod/primitives/pkg/common"

type AttestationData interface {
	GetSlot() uint64
	GetIndex() uint64
	GetBeaconBlockRoot() common.Root

	SetSlot(uint64)
	SetIndex(uint64)
	SetBeaconBlockRoot(uint64)
}
