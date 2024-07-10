package types

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

type ForkData[T any] interface {
	// New creates a new fork data.
	New(
		currentVersion common.Version,
		genesisValidatorsRoot common.Root,
	) T
	// ComputeDomain computes the domain of the fork data.
	ComputeDomain(domainType common.DomainType) (common.Domain, error)
	// ComputeRandaoSigningRoot computes the randao signing root.
	ComputeRandaoSigningRoot(
		domainType common.DomainType,
		epoch math.Epoch,
	) (common.Root, error)
}
