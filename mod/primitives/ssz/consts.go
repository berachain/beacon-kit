package ssz

const (
	BytesPerLengthOffset = 4
	MaximumLength        = 1 << (8 * BytesPerLengthOffset)
	BitsPerByte          = 8
)
