package merkle

import "encoding/binary"

// HashFn is the generic hash function signature.
type HashFn func(input []byte) [32]byte

// Hasher describes an interface through which we can
// perform hash operations on byte arrays,indices,etc.
type Hasher2 interface {
	Hash(a []byte) [32]byte
	Combi(a [32]byte, b [32]byte) [32]byte
	MixIn(a [32]byte, i uint64) [32]byte
}

// HasherFunc defines a structure to hold a hash function and can be used for multiple rounds of
// hashing.
type HasherFunc[T ~[32]byte] struct {
	b        [64]byte
	hashFunc HashFn
}

// NewHasherFunc is the constructor for the object
// that fulfills the Hasher interface.
func NewHasherFunc[T ~[32]byte](h HashFn) *HasherFunc[T] {
	return &HasherFunc[T]{
		b:        [64]byte{},
		hashFunc: h,
	}
}

// Hash utilizes the provided hash function for
// the object.
func (h *HasherFunc[T]) Hash(a []byte) T {
	return T(h.hashFunc(a))
}

// Combi appends the two inputs and hashes them.
func (h *HasherFunc[T]) Combi(a, b T) T {
	copy(h.b[:32], a[:])
	copy(h.b[32:], b[:])
	return h.Hash(h.b[:])
}

// MixIn works like Combi, but using an integer as the second input.
func (h *HasherFunc[T]) MixIn(a T, i uint64) T {
	copy(h.b[:32], a[:])
	copy(h.b[32:], make([]byte, 32))
	binary.LittleEndian.PutUint64(h.b[32:], i)
	return h.Hash(h.b[:])
}
