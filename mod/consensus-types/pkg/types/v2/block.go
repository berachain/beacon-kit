package v2

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/karalabe/ssz"
)

// BeaconBlockHeaderBase represents the base of a beacon block header.
type BeaconBlock struct {
	// Slot represents the position of the block in the chain.
	Slot math.Slot
	// ProposerIndex is the index of the validator who proposed the block.
	ProposerIndex math.ValidatorIndex
	// ParentBlockRoot is the hash of the parent block
	ParentBlockRoot common.Root
	// StateRoot is the hash of the state at the block.
	StateRoot common.Root
	// Body is the body of the BeaconBlockDeneb, containing the block's
	// operations.
	Body *BeaconBlockBody
}

// DefineSSZ defines the SSZ encoding for the BeaconBlock object.
func (b *BeaconBlock) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineUint64(codec, &b.Slot)
	ssz.DefineUint64(codec, &b.ProposerIndex)
	ssz.DefineStaticBytes(codec, &b.ParentBlockRoot)
	ssz.DefineStaticBytes(codec, &b.StateRoot)
	ssz.DefineDynamicObjectContent(codec, &b.Body)
}

// SizeSSZ returns the size of the BeaconBlock object in SSZ encoding.
func (b *BeaconBlock) SizeSSZ(isFixed bool) uint32 {
	return 4 + // Slot
		8 + // ProposerIndex
		32 + // ParentBlockRoot
		32 + // StateRoot
		b.Body.SizeSSZ(isFixed) // Body
}

// Version returns the version of the BeaconBlock.
func (b *BeaconBlock) Version() uint32 {
	return b.Body.version
}

// IsNil checks if the BeaconBlock instance is nil.
func (b *BeaconBlock) IsNil() bool {
	return b == nil
}

// SetStateRoot sets the state root of the BeaconBlock.
func (b *BeaconBlock) SetStateRoot(root common.Root) {
	b.StateRoot = root
}

// GetBody retrieves the body of the BeaconBlock.
func (b *BeaconBlock) GetBody() *BeaconBlockBody {
	return b.Body
}

// GetSlot retrieves the slot of the BeaconBlock.
func (b *BeaconBlock) GetSlot() math.Slot {
	return b.Slot
}

// GetParentBlockRoot
func (b *BeaconBlock) GetParentBlockRoot() common.Root {
	return b.ParentBlockRoot
}

// GetStateRoot
func (b *BeaconBlock) GetStateRoot() common.Root {
	return b.StateRoot
}

// MarshalSSZ serializes the BeaconBlock into a writer.
func (b *BeaconBlock) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, b.SizeSSZ(false))
	return b.MarshalSSZTo(buf)
}

// MarshalSSZTo serializes the BeaconBlock into a writer.
func (b *BeaconBlock) MarshalSSZTo(buf []byte) ([]byte, error) {
	return buf, ssz.EncodeToBytes(buf, b)
}

// UnmarshalSSZ deserializes the BeaconBlock from a byte slice.
func (b *BeaconBlock) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, b)
}

// HashTreeRoot returns the hash tree root of the BeaconBlock.
func (b *BeaconBlock) HashTreeRoot() ([32]byte, error) {
	return ssz.HashSequential(b), nil
}

// GetHeader builds a BeaconBlockHeader from the BeaconBlock.
func (b *BeaconBlock) GetHeader() *BeaconBlockHeader {
	return &BeaconBlockHeader{
		Slot:            b.Slot,
		ProposerIndex:   b.ProposerIndex,
		ParentBlockRoot: b.ParentBlockRoot,
		StateRoot:       b.StateRoot,
		BodyRoot:        common.Root(b.GetBody().HashTreeRoot()),
	}
}
