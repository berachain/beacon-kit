package flags

const (
	// Execution Client.
	RPCDialURL      = "polaris.execution-client.rpc-dial-url"
	RPCTimeout      = "polaris.execution-client.rpc-timeout"
	RPCRetries      = "polaris.execution-client.rpc-retries"
	JWTSecretPath   = "polaris.execution-client.jwt-secret-path" //nolint:gosec // false positive.
	RequiredChainID = "polaris.execution-client.required-chain-id"

	// Beacon Config.
	AltairForkEpoch    = "polaris.beacon-config.altair-fork-epoch"
	BellatrixForkEpoch = "polaris.beacon-config.bellatrix-fork-epoch"
	CapellaForkEpoch   = "polaris.beacon-config.capella-fork-epoch"
	DenebForkEpoch     = "polaris.beacon-config.deneb-fork-epoch"
)
