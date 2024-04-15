package ssz

import "errors"

var (
	// ErrInvalidNilSlice is returned when the input slice is nil.
	ErrInvalidNilSlice = errors.New("invalid empty slice")
)
