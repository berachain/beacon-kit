package types

import "github.com/cosmos/gogoproto/proto"

// COSMOS TYPES

type (
	Msg      = proto.Message
	Identity = []byte
)

type Tx interface {
	// Hash returns the unique identifier for the Tx.
	Hash() [32]byte
	// GetMessages returns the list of state transitions of the Tx.
	GetMessages() ([]Msg, error)
	// GetSenders returns the tx state transition sender.
	GetSenders() ([]Identity, error) // TODO reduce this to a single identity if accepted
	// GetGasLimit returns the gas limit of the tx. Must return math.MaxUint64 for infinite gas
	// txs.
	GetGasLimit() (uint64, error)
	// Bytes returns the encoded version of this tx. Note: this is ideally cached
	// from the first instance of the decoding of the tx.
	Bytes() []byte
}

type tx[T Tx] []byte

func NewTx[T Tx](b []byte) T {
	return Tx(tx[T](b)).(T)
}

func (t tx[T]) Hash() [32]byte {
	return [32]byte(t)
}

func (t tx[T]) GetMessages() ([]Msg, error) {
	return nil, nil
}

func (t tx[T]) GetSenders() ([]Identity, error) {
	return nil, nil
}

func (t tx[T]) GetGasLimit() (uint64, error) {
	return 0, nil
}

func (t tx[T]) Bytes() []byte {
	return t
}
