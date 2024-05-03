package hex

import (
	"encoding/hex"
	"errors"
)

var (
	ErrEmptyString   = errors.New("empty hex string")
	ErrMissingPrefix = errors.New("hex string without 0x prefix")
)

// String represents a hex string with 0x prefix.
type String string

// FromBytes creates a hex string with 0x prefix.
func StrFromBytes(b []byte) String {
	enc := make([]byte, len(b)*2+2)
	copy(enc, "0x")
	hex.Encode(enc[2:], b)
	return String(enc)
}

// Has0xPrefix returns true if s has a 0x prefix.
func (s String) Has0xPrefix() bool {
	return has0xPrefix[string](string(s))
}

// IsEmpty returns true if s is empty.
func (s String) IsEmpty() bool {
	return len(s) == 0
}

// ToBytes decodes a hex string with 0x prefix.
func (s String) ToBytes() ([]byte, error) {
	if s.IsEmpty() {
		return nil, ErrEmptyString
	} else if s.Has0xPrefix() {
		return nil, ErrMissingPrefix
	}
	return hex.DecodeString(string(s[2:]))
}

// MustToBytes decodes a hex string with 0x prefix.
// It panics for invalid input.
func (s String) MustToBytes() []byte {
	b, err := s.ToBytes()
	if err != nil {
		panic(err)
	}
	return b
}
