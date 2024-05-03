package hex

import "errors"

var (
	ErrEmptyString     = errors.New("empty hex string")
	ErrMissingPrefix   = errors.New("hex string without 0x prefix")
	ErrOddLength       = errors.New("hex string of odd length")
	ErrNonQuotedString = errors.New("non-quoted hex string")
	ErrInvalidString   = errors.New("invalid hex string")

	ErrLeadingZero = errors.New("hex number with leading zero digits")
	ErrEmptyNumber = errors.New("hex string \"0x\"")
	ErrUint64Range = errors.New("hex number > 64 bits")
	ErrBig256Range = errors.New("hex number > 256 bits")
)
