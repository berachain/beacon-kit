package state

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// Marshallable represents an interface for a beacon state
// with generic types.
type Marshallable[
	T any,
	BeaconBlockHeaderT,
	Eth1DataT,
	ExecutionPayloadHeaderT,
	ForkT,
	ValidatorT any,
] interface {
	constraints.SSZMarshallable
	// New returns a new instance of the BeaconStateMarshallable.
	New(
		forkVersion uint32,
		genesisValidatorsRoot common.Bytes32,
		slot math.U64,
		fork ForkT,
		latestBlockHeader BeaconBlockHeaderT,
		blockRoots []common.Bytes32,
		stateRoots []common.Bytes32,
		eth1Data Eth1DataT,
		eth1DepositIndex uint64,
		latestExecutionPayloadHeader ExecutionPayloadHeaderT,
		validators []ValidatorT,
		balances []uint64,
		randaoMixes []common.Bytes32,
		nextWithdrawalIndex uint64,
		nextWithdrawalValidatorIndex math.U64,
		slashings []uint64, totalSlashing math.U64,
	) (T, error)
}
