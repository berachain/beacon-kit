package beacondb

import (
	"cosmossdk.io/core/store"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
)

func (s *Store[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT, ForkT, ValidatorT,
]) SetGenesisValidatorRoot(
	root common.Root,
) {
	if s.changeSet == nil {
		s.changeSet = store.NewChangeset()
	}
	s.changeSet.Add([]byte("genesis"), []byte("genesis_validator_root"), root[:], false)
}

func (s *Store[
	BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT, ForkT, ValidatorT,
]) GetGenesisValidatorRoot() (common.Root, error) {
	bz, err := s.genesisValidatorsRoot.Get()
	if err != nil {
		return common.Root{}, err
	}
	return common.Root(bz), nil
}
