package constraints

// SSZMarshallable is an interface that combines the ssz.Marshaler and
// ssz.Unmarshaler interfaces.
type SSZMarshallable interface {
	// MarshalSSZTo marshals the object into the provided byte slice and returns
	// it along with any error.
	MarshalSSZTo([]byte) ([]byte, error)
	// MarshalSSZ marshals the object into a new byte slice and returns it along
	// with any error.
	MarshalSSZ() ([]byte, error)
	// UnmarshalSSZ unmarshals the object from the provided byte slice and
	// returns an error if the unmarshaling fails.
	UnmarshalSSZ([]byte) error
	// SizeSSZ returns the size in bytes that the object would take when
	// marshaled.
	SizeSSZ() int
	// HashTreeRoot returns the hash tree root of the object.
	HashTreeRoot() ([32]byte, error)
}

// NewFromSSZable is an interface that combines the SSZMarshallable interface
// with a NewFromSSZ method that creates a new object from an ssz byte slice.
type NewFromSSZable[T any] interface {
	SSZMarshallable
	NewFromSSZ(bz []byte, forkVersion uint32) (T, error)
}
