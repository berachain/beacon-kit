package client

import "github.com/itsdevbear/bolaris/types/consensus/primitives"

// beaconConfig is an interface for the beacon chain configuration
type beaconConfig interface {
	ActiveForkVersion(primitives.Epoch) int
}
