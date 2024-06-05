package app

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
