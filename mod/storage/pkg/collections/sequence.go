package collections

import (
	sdkcollections "cosmossdk.io/collections"
)

// DefaultSequenceStart defines the default starting number of a sequence.
const DefaultSequenceStart uint64 = 0

// Sequence builds on top of an Item, and represents a monotonically increasing number.
type Sequence struct {
	Item[uint64]
}

// NewSequence instantiates a new sequence given
// a Schema, a Prefix and humanized name for the sequence.
func NewSequence(
	storeKey, key []byte, storeAccessor StoreAccessor,
) Sequence {
	seq := Sequence{
		NewItem(
			storeKey,
			key,
			sdkcollections.Uint64Value,
			storeAccessor,
		),
	}
	// i worry this breaks cause the collection set is not
	// committed yet.
	seq.Set(DefaultSequenceStart)
	return seq
}

// Peek returns the current sequence value, if no number
// is set then the DefaultSequenceStart is returned.
// Errors on encoding issues.
func (s *Sequence) Peek() (uint64, error) {
	return s.Item.Get()
}

// Next returns the next sequence number, and sets the next expected sequence.
// Errors on encoding issues.
func (s *Sequence) Next() (uint64, error) {
	seq, err := s.Peek()
	if err != nil {
		return 0, err
	}
	return seq, s.Set(seq + 1)
}

// Set hard resets the sequence to the provided value.
// Errors on encoding issues.
func (s *Sequence) Set(value uint64) error {
	return s.Item.Set(value)
}
