package hex

import (
	"encoding/hex"
	"strconv"
)

// String represents a hex string with 0x prefix.
type String string

// StrFromBytes creates a hex string with 0x prefix.
func StrFromBytes(b []byte) String {
	enc := make([]byte, len(b)*2+2)
	copy(enc, "0x")
	hex.Encode(enc[2:], b)
	return String(enc)
}

// StrFromUint64 encodes i as a hex string with 0x prefix.
func StrFromUint64(i uint64) String {
	enc := make([]byte, 2, 10)
	copy(enc, "0x")
	return String(strconv.AppendUint(enc, i, 16))
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

// ToUint64 decodes a hex string with 0x prefix.
func (s String) ToUint64() (uint64, error) {
	raw, err := validateNumber(string(s))
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(raw, 16, 64)
}

// MustToUint64 decodes a hex string with 0x prefix.
// It panics for invalid input.
func (s String) MustToUint64() uint64 {
	i, err := s.ToUint64()
	if err != nil {
		panic(err)
	}
	return i
}

// Unwrap returns the string value.
func (s String) Unwrap() string {
	return string(s)
}
