package middleware

type OuterBlock interface {
	Height() uint64
	GetBeaconBlockBytes() ([]byte, error)
	GetSidecarsBytes() ([]byte, error)
}
