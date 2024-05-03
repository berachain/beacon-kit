package hex

import (
	"encoding/json"
	"reflect"

	"github.com/holiman/uint256"
)

var u256T = reflect.TypeOf((*uint256.Int)(nil))

// U256 marshals/unmarshals as a JSON string with 0x prefix.
// The zero value marshals as "0x0".
type U256 uint256.Int

// MarshalText implements encoding.TextMarshaler
func (b U256) MarshalText() ([]byte, error) {
	u256 := (*uint256.Int)(&b)
	return []byte(u256.Hex()), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *U256) UnmarshalJSON(input []byte) error {
	// The uint256.Int.UnmarshalJSON method accepts "dec", "0xhex"; we must be
	// more strict, hence we check string and invoke SetFromHex directly.
	if !isQuotedString(input) {
		return ErrNonQuotedString
	}
	// The hex decoder needs to accept empty string ("") as '0', which uint256.Int
	// would reject.
	if len(input) == 2 {
		(*uint256.Int)(b).Clear()
		return nil
	}
	err := (*uint256.Int)(b).SetFromHex(string(input[1 : len(input)-1]))
	if err != nil {
		return &json.UnmarshalTypeError{Value: err.Error(), Type: u256T}
	}
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler
func (b *U256) UnmarshalText(input []byte) error {
	// The uint256.Int.UnmarshalText method accepts "dec", "0xhex"; we must be
	// more strict, hence we check string and invoke SetFromHex directly.
	return (*uint256.Int)(b).SetFromHex(string(input))
}

// String returns the hex encoding of b.
func (b *U256) String() string {
	return (*uint256.Int)(b).Hex()
}
