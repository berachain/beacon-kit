package hex

import (
	"encoding/hex"
	"reflect"

	"github.com/berachain/beacon-kit/mod/errors"
)

var bytesT = reflect.TypeOf(Bytes(nil))

// Bytes marshals/unmarshals as a JSON string with 0x prefix.
// The empty slice marshals as "0x".
type Bytes []byte

// MarshalText implements encoding.TextMarshaler
func (b Bytes) MarshalText() ([]byte, error) {
	result := make([]byte, len(b)*2+2)
	copy(result, `0x`)
	hex.Encode(result[2:], b)
	return result, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *Bytes) UnmarshalJSON(input []byte) error {
	if !isQuotedString(input) {
		return wrapUnmarshalError(ErrNonQuotedString, bytesT)
	}

	return wrapUnmarshalError(b.UnmarshalText(input[1:len(input)-1]), bytesT)
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (b *Bytes) UnmarshalText(input []byte) error {
	raw, err := validateText(input, true)
	if err != nil {
		return err
	}
	dec := make([]byte, len(raw)/2)
	if _, err = hex.Decode(dec, raw); err != nil {
		return err
	}
	*b = dec

	return nil
}

// String returns the hex encoding of b.
func (b Bytes) String() String {
	return StrFromBytes(b)
}

// UnmarshalFixedJSON decodes the input as a string with 0x prefix. The length of out
// determines the required input length. This function is commonly used to implement the
// UnmarshalJSON method for fixed-size types.
func UnmarshalFixedJSON(typ reflect.Type, input, out []byte) error {
	if !isQuotedString(input) {
		return wrapUnmarshalError(ErrNonQuotedString, bytesT)
	}
	return wrapUnmarshalError(
		UnmarshalFixedText(typ.String(), input[1:len(input)-1], out), typ,
	)
}

// UnmarshalFixedText decodes the input as a string with 0x prefix. The length of out
// determines the required input length. This function is commonly used to implement the
// UnmarshalText method for fixed-size types.
func UnmarshalFixedText(typname string, input, out []byte) error {
	raw, err := validateText(input, true)
	if err != nil {
		return err
	}
	if len(raw)/2 != len(out) {
		return errors.Newf(
			"hex string has length %d, want %d for %s",
			len(raw), len(out)*2, typname,
		)
	}
	// Pre-verify syntax before modifying out.
	for _, b := range raw {
		if decodeNibble(b) == badNibble {
			return ErrInvalidString
		}
	}
	if _, err = hex.Decode(out, raw); err != nil {
		return err
	}

	return nil
}
