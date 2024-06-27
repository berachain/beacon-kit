package serializer

import "fmt"

// MarshalVectorFixed converts a slice of basic values into a byte slice.
func MarshalVectorFixed[T interface{ MarshalSSZ() ([]byte, error) }](
	out []byte, v []T,
) ([]byte, error) {
	// From the Spec:
	// fixed_parts = [
	// 		serialize(element)
	// 			if not is_variable_size(element)
	//			else None for element in value,
	// 		]
	// VectorBasic has all fixed types, so we simply
	// serialize each element and pack them together.
	for _, val := range v {
		bytes, err := val.MarshalSSZ()
		if err != nil {
			return out, err
		}
		out = append(out, bytes...)
	}
	return out, nil
}

// UnmarshalVectorFixed converts a byte slice into a slice of basic values.
func UnmarshalVectorFixed[
	T interface {
		NewFromSSZ([]byte) (T, error)
		SizeSSZ() int
	},
](
	buf []byte,
) ([]T, error) {
	var (
		err error
		t   T
	)
	elementSize := t.SizeSSZ()
	if len(buf)%elementSize != 0 {
		return nil, fmt.Errorf(
			"invalid buffer length %d for element size %d",
			len(buf),
			elementSize,
		)
	}

	result := make([]T, 0, len(buf)/elementSize)
	for i := 0; i < len(buf); i += elementSize {
		if t, err = t.NewFromSSZ(buf[i : i+elementSize]); err != nil {
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}
