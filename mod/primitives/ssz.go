package primitives

// SSZMarshallable is an interface that combines the ssz.Marshaler and ssz.Unmarshaler interfaces.
type SSZMarshallable interface {
	// MarshalSSZ marshals the object into a new byte slice and returns it along with any error.
	MarshalSSZ() ([]byte, error)
	// UnmarshalSSZ unmarshals the object from the provided byte slice and returns an error if the unmarshaling fails.
	UnmarshalSSZ([]byte) error
	// SizeSSZ returns the size in bytes that the object would take when marshaled.
	SizeSSZ() int
}
