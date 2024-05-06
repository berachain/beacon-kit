package p2p

import "github.com/berachain/beacon-kit/mod/errors"

var (
	// ErrNilBeaconBlockInRequest is an error for when
	// the beacon block in an abci request is nil.
	ErrNilBeaconBlockInRequest = errors.New("nil beacon block in abci request")

	// ErrNoBeaconBlockInRequest is an error for when
	// there is no beacon block in an abci request.
	ErrNoBeaconBlockInRequest = errors.New("no beacon block in abci request")

	// ErrBzIndexOutOfBounds is an error for when the index
	// is out of bounds.
	ErrBzIndexOutOfBounds = errors.New("bzIndex out of bounds")

	// ErrNilABCIRequest is an error for when the abci request
	// is nil.
	ErrNilABCIRequest = errors.New("nil abci request")
)
