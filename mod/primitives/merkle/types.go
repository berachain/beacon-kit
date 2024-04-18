package merkle

// U64 is an interface that wraps the uint64 type.
// It is used to prevent circular dependencies between
// the merkle package and the primitives package.
type U64[T ~uint64] interface {
	~uint64
	// NextPowerOfTwo returns the smallest power of two that is greater than or equal to T.
	NextPowerOfTwo() T
	// ILog2Ceil returns the ceiling of the binary logarithm of T as a uint8.
	ILog2Ceil() uint8
}
