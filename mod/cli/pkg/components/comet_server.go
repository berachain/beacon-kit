package components

import (
	"errors"

	"cosmossdk.io/core/transaction"
	serverv2 "cosmossdk.io/server/v2"
	"cosmossdk.io/server/v2/cometbft"
	"github.com/cosmos/cosmos-sdk/client"
)

type temporaryTxDecoder[T transaction.Tx] struct {
	txConfig client.TxConfig
}

// Decode implements transaction.Codec.
func (t *temporaryTxDecoder[T]) Decode(bz []byte) (T, error) {
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
func (t *temporaryTxDecoder[T]) DecodeJSON(bz []byte) (T, error) {
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

// Temp hood
func ProvideCometServer[NodeT serverv2.AppI[T], T transaction.Tx](
	clientCtx client.Context) *cometbft.CometBFTServer[NodeT, T] {
	return cometbft.New[NodeT, T](&temporaryTxDecoder[T]{clientCtx.TxConfig})
}
