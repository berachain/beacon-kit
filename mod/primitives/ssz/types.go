package ssz

type U64[T ~uint64] interface {
	~uint64
	Unwrap() uint64
	NextPowerOfTwo() T
	ILog2Ceil() uint8
}

type Basic interface{}

type BasicVecList[B Basic, RootT ~[32]byte] []B

type Composite[RootT ~[32]byte] interface {
	SizeSSZ() int
	HashTreeRoot() (RootT, error)
}

type Container interface {
	Marshallable
}
