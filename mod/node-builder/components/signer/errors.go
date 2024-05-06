package signer

import "errors"

// ErrInvalidSignature is returned when a signature is invalid.
var ErrInvalidSignature = errors.New("invalid signature")
