package flags

const (
	// Block Store Service Config.
	blockStoreServiceRoot          = beaconKitRoot + "block-store-service."
	BlockStoreServiceEnabled       = blockStoreServiceRoot + "enabled"
	BlockStoreServicePrunerEnabled = blockStoreServiceRoot +
		"pruner-enabled"
	BlockStoreServiceAvailabilityWindow = blockStoreServiceRoot +
		"availability-window"
)
