package ssz

import (
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/merkleizer"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/types"
)

// Vector conforms to the SSZEenumerable interface.
var _ types.SSZEnumerable[types.BaseSSZType] = (*Container)(nil)

type Container struct {
	elements []types.BaseSSZType
}

// SizeSSZ returns the size of the container in bytes.
func (c *Container) SizeSSZ() int {
	size := 0
	for _, element := range c.elements {
		size += element.SizeSSZ()
	}
	return size
}

// ContainerFromElements creates a new Container from elements.
func ContainerFromElements(elements ...types.BaseSSZType) *Container {
	return &Container{
		elements: elements,
	}
}

// IsFixed returns true if the container is fixed size.
func (c *Container) IsFixed() bool {
	for _, element := range c.elements {
		if !element.IsFixed() {
			return false
		}
	}
	return true
}

// N returns the N value as defined in the SSZ specification.
func (c *Container) N() uint64 {
	return uint64(len(c.elements))
}

// Type returns the type of the container.
func (*Container) Type() types.Type {
	return types.Composite
}

// ChunkCount returns the number of chunks in the container.
func (c *Container) ChunkCount() uint64 {
	return c.N()
}

// Elements returns the elements of the container.
func (c *Container) Elements() []types.BaseSSZType {
	return c.elements
}

/* -------------------------------------------------------------------------- */
/*                                Merkleization                               */
/* -------------------------------------------------------------------------- */

// HashTreeRoot returns the hash tree root of the container.
func (c *Container) HashTreeRootWith(
	merkleizer VectorMerkleizer[[32]byte, types.BaseSSZType],
) ([32]byte, error) {
	return merkleizer.MerkleizeVectorCompositeOrContainer(c.elements)
}

// HashTreeRoot returns the hash tree root of the container.
func (c *Container) HashTreeRoot() ([32]byte, error) {
	return c.HashTreeRootWith(merkleizer.New[[32]byte, types.BaseSSZType]())
}

/* -------------------------------------------------------------------------- */
/*                                Serialization                               */
/* -------------------------------------------------------------------------- */

// MarshalSSZToBytes marshals the VectorBasic into SSZ format.
func (c *Container) MarshalSSZTo(_ []byte) ([]byte, error) {
	return nil, errors.New("not implemented yet")
}

// MarshalSSZ marshals the VectorBasic into SSZ format.
func (c *Container) MarshalSSZ() ([]byte, error) {
	return c.MarshalSSZTo(make([]byte, 0, c.SizeSSZ()))
}

// NewFromSSZ creates a new Container from SSZ format.
func (c *Container) NewFromSSZ(_ []byte) (*Container, error) {
	return nil, errors.New("not implemented yet")
}
