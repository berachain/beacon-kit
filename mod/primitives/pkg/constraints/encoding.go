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
}

// JSONMarshallable is an interface that combines the json.Marshaler and
// json.Unmarshaler interfaces.
type JSONMarshallable interface {
	// MarshalJSON marshals the object into a JSON byte slice and returns it
	// along with any error.
	MarshalJSON() ([]byte, error)
	// UnmarshalJSON unmarshals the object from the provided JSON byte slice and
	// returns an error if the unmarshaling fails.
	UnmarshalJSON([]byte) error
}
