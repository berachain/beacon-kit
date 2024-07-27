package depinject

import "errors"

var (
	ErrTargetMustBePointer = errors.New("target must be a pointer")
)
