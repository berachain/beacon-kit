package chainspec

import "github.com/berachain/beacon-kit/mod/chain-spec/pkg/chain"

type BartioChainSpec[
	DomainTypeT ~[4]byte,
	EpochT ~uint64,
	ExecutionAddressT ~[20]byte,
	SlotT ~uint64,
	CometBFTConfigT any,
] struct {
	// BGTContractAddress
	BGTContractAddress ExecutionAddressT `mapstructure:"bgt-contract-address"`
	// SpecData is the underlying data structure for chain-specific parameters.
	chain.SpecData[
		DomainTypeT,
		EpochT,
		ExecutionAddressT,
		SlotT,
		CometBFTConfigT,
	]
}
