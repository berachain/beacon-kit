package flags

const (
	// Builder Config.
	builderRoot              = beaconKitRoot + "payload-builder."
	SuggestedFeeRecipient    = builderRoot + "suggested-fee-recipient"
	LocalBuilderEnabled      = builderRoot + "local-builder-enabled"
	LocalBuildPayloadTimeout = builderRoot + "local-build-payload-timeout"
)
