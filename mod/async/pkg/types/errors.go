package types

import "github.com/berachain/beacon-kit/mod/errors"

var (
	// ErrBrokerAlreadyExists defines an error for when a service already
	// exists.
	ErrBrokerAlreadyExists = func(brokerName string) error {
		return errors.Newf("broker already exists: %v", brokerName)
	}

	// ErrInputIsNotPointer defines an error for when the input must
	// be of pointer type.
	ErrInputIsNotPointer = func(valueType any) error {
		return errors.Newf(
			"input must be of pointer type, received value type instead: %T",
			valueType,
		)
	}

	// errUnknownService defines is returned when an unknown service is seen.
	ErrUnknownService = func(brokerType any) error {
		return errors.Newf("unknown service: %T", brokerType)
	}
)
