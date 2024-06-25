package components

import (
	"errors"

	"cosmossdk.io/core/transaction"
	"github.com/cosmos/cosmos-sdk/client"
)

type TxDecoder[T transaction.Tx] struct {
	txConfig client.TxConfig
}

// Decode implements transaction.Codec.
func (t *TxDecoder[T]) Decode(bz []byte) (T, error) {
	var out T
	tx, err := t.txConfig.TxDecoder()(bz)
	if err != nil {
		return out, err
	}

	var ok bool
	out, ok = tx.(T)
	if !ok {
		return out, errors.New("unexpected Tx type")
	}

	return out, nil
}

// DecodeJSON implements transaction.Codec.
func (t *TxDecoder[T]) DecodeJSON(bz []byte) (T, error) {
	var out T
	tx, err := t.txConfig.TxJSONDecoder()(bz)
	if err != nil {
		return out, err
	}

	var ok bool
	out, ok = tx.(T)
	if !ok {
		return out, errors.New("unexpected Tx type")
	}

	return out, nil
}

func ProvideTxDecoder[T transaction.Tx](txConfig client.TxConfig) *TxDecoder[T] {
	return &TxDecoder[T]{
		txConfig: txConfig,
	}
}
