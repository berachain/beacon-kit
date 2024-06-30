package types

type Type uint8

const (
	Basic Type = iota
	Composite
)

// SSZType is the interface for all SSZ types.
type SSZType[T any] interface {
	// NewFromSSZ creates a new composite type from an SSZ byte slice.
	NewFromSSZ([]byte) (T, error)
	// MarshalSSZ serializes the composite type to an SSZ byte slice.
	MarshalSSZ() ([]byte, error)
	// HashTreeRoot returns the hash tree root of the composite type.
	HashTreeRoot() ([32]byte, error)
	// MarshalSSZ marshals the type into SSZ format.
	IsFixed() bool
	// SizeSSZ returns the size of the type in bytes.
	SizeSSZ() int
	// TODO: Enable these
	//
	// ChunkCount returns the number of chunks required to store the type.
	// ChunkCount() uint64
	// Type returns the type of the SSZ object.
	// Type() Type
}

// SSZEnumerable is the interface for all SSZ enumerable types must implement.
type SSZEnumerable[
	SelfT SSZType[SelfT],
	ElementT SSZType[ElementT],
] interface {
	SSZType[SelfT]
	// N returns the N value as defined in the SSZ specification.
	N() uint64
	// Elements returns the elements of the enumerable type.
	Elements() []SSZType[ElementT]
}
