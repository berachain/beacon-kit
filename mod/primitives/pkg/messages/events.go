package messages

// events.
const (
	BeaconBlockFinalizedEvent = "beacon-block-finalized"
)

// events, topologically sorted.
const (
	// genesis data events
	GenesisDataReceived  = "genesis-data-received"
	GenesisDataProcessed = "genesis-data-processed"
	// pre proposal events
	NewSlot          = "new-slot"
	BuiltBeaconBlock = "built-beacon-block"
	BuiltSidecars    = "built-sidecars"
	// proposal processing events
	BeaconBlockReceived = "beacon-block-received"
	SidecarsReceived    = "sidecars-received"
	// finalize block events
	FinalBeaconBlockReceived       = "final-beacon-block-received"
	FinalBlobSidecarsReceived      = "final-blob-sidecars-received"
	FinalValidatorUpdatesProcessed = "final-validator-updates"
	FinalSidecarsProcessed         = "final-sidecars-processed"
)
