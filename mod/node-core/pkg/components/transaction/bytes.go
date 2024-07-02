package transaction

import (
	"cosmossdk.io/core/transaction"
	"github.com/cosmos/gogoproto/proto"
)

type BytesTx []byte

var _ transaction.Tx = (*BytesTx)(nil)

func NewTxFromBytes[T transaction.Tx](bz []byte) T {
	return any(BytesTx(bz)).(T)
}

func (tx BytesTx) Bytes() []byte {
	return tx
}

func (tx BytesTx) Hash() [32]byte {
	return [32]byte{}
}

func (tx BytesTx) GetGasLimit() (uint64, error) { return 0, nil }

func (tx BytesTx) GetMessages() ([]proto.Message, error) { return nil, nil }

func (tx BytesTx) GetSenders() ([][]byte, error) { return nil, nil }

var _ transaction.Codec[transaction.Tx] = (*BytesTxCodec[transaction.Tx])(nil)

type BytesTxCodec[T transaction.Tx] struct{}

func (cdc *BytesTxCodec[T]) Decode(bz []byte) (T, error) {
	return NewTxFromBytes[T](bz), nil
}

// DecodeJSON decodes the tx JSON bytes into a DecodedTx
func (cdc *BytesTxCodec[T]) DecodeJSON(bz []byte) (T, error) {
	return NewTxFromBytes[T](bz), nil
}
