package middleware

import "github.com/berachain/beacon-kit/mod/errors"

var (
	ErrInitGenesisTimeout = func(errTimeout error) error {
		return errors.Wrapf(errTimeout,
			"A timeout occurred while waiting for genesis data processing",
		)
	}

	ErrBuildBeaconBlockTimeout = func(errTimeout error) error {
		return errors.Wrapf(errTimeout,
			"A timeout occurred while waiting for a beacon block to be built",
		)
	}

	ErrBuildSidecarsTimeout = func(errTimeout error) error {
		return errors.Wrapf(errTimeout,
			"A timeout occurred while waiting for blob sidecars to be built",
		)
	}

	ErrVerifyBeaconBlockTimeout = func(errTimeout error) error {
		return errors.Wrapf(
			errTimeout,
			"A timeout occurred while waiting for a beacon block to be verified",
		)
	}

	ErrVerifySidecarsTimeout = func(errTimeout error) error {
		return errors.Wrapf(errTimeout,
			"A timeout occurred while waiting for blob sidecars to be verified",
		)
	}

	ErrFinalValidatorUpdatesTimeout = func(errTimeout error) error {
		return errors.Wrapf(errTimeout,
			"A timeout occurred while waiting for final validator updates",
		)
	}
)
