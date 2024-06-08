package components

// DefaultComponents is the default set of components
// that are provided by beacon-kit.
func DefaultComponents() []any {
	return []any{
		ProvideAvailibilityStore,
		ProvideBlsSigner,
		ProvideTrustedSetup,
		ProvideDepositStore,
		ProvideConfig,
		ProvideEngineClient,
		ProvideJWTSecret,
		ProvideBlobProofVerifier,
		ProvideTelemetrySink,
		ProvideExecutionEngine,
		ProvideBeaconDepositContract,
		ProvideLocalBuilder,
		ProvideStateProcessor,
	}
}
