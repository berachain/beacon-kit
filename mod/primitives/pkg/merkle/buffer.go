package merkle

// initialBufferSize is the initial size of the internal buffer.
//
// TODO: choose a more appropriate size?
const initialBufferSize = 16

// Buffer is a re-usable buffer for merkle tree hashing. Prevents
// unnecessary allocations and garbage collection of byte slices.
//
// NOTE: this buffer is ONLY meant to be used in a single thread.
type Buffer[RootT ~[32]byte] struct {
	internal []RootT

	// TODO: add a mutex for multi-thread safety.
}

// NewBuffer creates a new buffer with the given capacity.
func NewBuffer[RootT ~[32]byte]() *Buffer[RootT] {
	return &Buffer[RootT]{
		internal: make([]RootT, initialBufferSize),
	}
}

// Get returns a slice of the internal buffer of roots of the given size.
func (b *Buffer[RootT]) Get(size int) []RootT {
	if size > len(b.internal) {
		b.grow(size - len(b.internal))
	}

	return b.internal[:size]
}

// TODO: add a Put method to return the buffer back for multi-threaded re-use.

// grow resizes the internal buffer by the requested size.
func (b *Buffer[RootT]) grow(newSize int) {
	b.internal = append(b.internal, make([]RootT, newSize)...)
}
