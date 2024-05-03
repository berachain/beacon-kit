package hex

import (
	"encoding/json"
	"reflect"
)

// has0xPrefix returns true if s has a 0x prefix.
func has0xPrefix[T []byte | string](s T) bool {
	return len(s) >= 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X')
}

// isQuotedString returns true if input has quotes.
func isQuotedString[T []byte | string](input T) bool {
	return len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"'
}

// validateText validates the input text for a hex string.
func validateText(input []byte, wantPrefix bool) ([]byte, error) {
	if len(input) == 0 {
		return nil, nil // empty strings are allowed
	}
	if has0xPrefix(input) {
		input = input[2:]
	} else if wantPrefix {
		return nil, ErrMissingPrefix
	}
	if len(input)%2 != 0 {
		return nil, ErrOddLength
	}
	return input, nil
}

// validateNumber checks the input text for a hex number.
func validateNumber[T []byte | string](input T) (raw T, err error) {
	if len(input) == 0 {
		return *new(T), nil // empty strings are allowed
	}
	if !has0xPrefix(input) {
		return *new(T), ErrMissingPrefix
	}
	input = input[2:]
	if len(input) == 0 {
		return *new(T), ErrEmptyNumber
	}
	if len(input) > 1 && input[0] == '0' {
		return *new(T), ErrLeadingZero
	}
	return input, nil
}

// wrapUnmarshalError wraps an error occuring during JSON unmarshaling.
func wrapUnmarshalError(err error, t reflect.Type) error {
	if err != nil {
		err = &json.UnmarshalTypeError{Value: err.Error(), Type: t}
	}

	return err
}

const badNibble = ^uint64(0)

func decodeNibble(in byte) uint64 {
	switch {
	case in >= '0' && in <= '9':
		return uint64(in - '0')
	case in >= 'A' && in <= 'F':
		return uint64(in - 'A' + 10)
	case in >= 'a' && in <= 'f':
		return uint64(in - 'a' + 10)
	default:
		return badNibble
	}
}
