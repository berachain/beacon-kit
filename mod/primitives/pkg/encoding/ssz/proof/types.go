package proof

import "github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/types"

type SSZType interface {
	Type() types.Type
	SizeSSZ() int
	IsFixed() bool
	HashTreeRoot() ([32]byte, error)
	MarshalSSZ() ([]byte, error)
}

type Elements interface {
	SSZType
	ElementType() types.Type
	ElementAtIndex(i uint64) SSZType
}
type Container[T SSZType] interface {
	SSZType
	GetFieldByName(name string) T
	GetFieldIndex(name string) uint64
}
