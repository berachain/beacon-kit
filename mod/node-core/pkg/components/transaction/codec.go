package transaction

import "cosmossdk.io/core/transaction"

var _ transaction.Codec[transaction.Tx] = (*SSZTxCodec[transaction.Tx])(nil)

type SSZTxCodec[T transaction.Tx] struct{}

func (cdc *SSZTxCodec[T]) Decode(bz []byte) (T, error) {
	return NewTxFromSSZ[T](bz), nil
}

// DecodeJSON decodes the tx JSON bytes into a DecodedTx
func (cdc *SSZTxCodec[T]) DecodeJSON(bz []byte) (T, error) {
	return NewTxFromSSZ[T](bz), nil
}
