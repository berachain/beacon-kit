package internal

import (
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/karalabe/ssz"
)

type BeaconBlockHeader struct {
	*types.BeaconBlockHeader
}

// SizeSSZ returns the size of the BeaconBlockHeader object in SSZ encoding.
func (b *BeaconBlockHeader) SizeSSZ() uint32 {
	return 112 // Total size: Slot (8) + ProposerIndex (8) + ParentBlockRoot (32) + StateRoot (32) + BodyRoot (32)
}

// DefineSSZ defines the SSZ encoding for the BeaconBlockHeader object.
func (b *BeaconBlockHeader) DefineSSZ(codec *ssz.Codec) {
	if b.BeaconBlockHeader == nil {
		b.BeaconBlockHeader = &types.BeaconBlockHeader{}
	}
	ssz.DefineUint64(codec, &b.Slot)
	ssz.DefineUint64(codec, &b.ProposerIndex)
	ssz.DefineStaticBytes(codec, &b.ParentBlockRoot)
	ssz.DefineStaticBytes(codec, &b.StateRoot)
	ssz.DefineStaticBytes(codec, &b.BodyRoot)
}

// MarshalSSZToBytes marshals the BeaconBlockHeader object to SSZ format.
func (b *BeaconBlockHeader) MarshalSSZTo(buf []byte) ([]byte, error) {
	return buf, ssz.EncodeToBytes(buf, b)
}

// MarshalSSZ marshals the BeaconBlockBody object to SSZ format.
func (b *BeaconBlockHeader) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, b.SizeSSZ())
	return buf, ssz.EncodeToBytes(buf, b)
}

// UnmarshalSSZ unmarshals the BeaconBlockBody object from SSZ format.
func (b *BeaconBlockHeader) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, b)
}

// HashTreeRoot computes the SSZ hash tree root of the BeaconBlockHeader object.
func (b *BeaconBlockHeader) HashTreeRoot() ([32]byte, error) {
	return ssz.HashSequential(b), nil
}

// GetSlot retrieves the slot of the BeaconBlockBase.
func (b *BeaconBlockHeader) GetSlot() math.Slot {
	return math.Slot(b.Slot)
}

// GetSlot retrieves the slot of the BeaconBlockBase.
func (b *BeaconBlockHeader) GetProposerIndex() math.ValidatorIndex {
	return math.ValidatorIndex(b.ProposerIndex)
}

// GetParentBlockRoot retrieves the parent block root of the BeaconBlockBase.
func (b *BeaconBlockHeader) GetParentBlockRoot() common.Root {
	return b.ParentBlockRoot
}

// GetStateRoot retrieves the state root of the BeaconBlockDeneb.
func (b *BeaconBlockHeader) GetStateRoot() common.Root {
	return b.StateRoot
}

// SetStateRoot sets the state root of the BeaconBlockHeader.
func (b *BeaconBlockHeader) SetStateRoot(stateRoot common.Root) {
	b.StateRoot = stateRoot
}
