package blob

// Processor is the interface for the blobs processor.
type Processor[
	AvailabilityStoreT any,
	BeaconBlockBodyT any,
	BlobSidecarsT any,
	ExecutionPayloadT any,
] interface {
	// ProcessSidecars processes the blobs and ensures they match the local
	// state.
	ProcessSidecars(
		avs AvailabilityStoreT,
		sidecars BlobSidecarsT,
	) error
	// VerifySidecars verifies the blobs and ensures they match the local state.
	VerifySidecars(
		sidecars BlobSidecarsT,
	) error
}
