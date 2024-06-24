package components

import (
	"cosmossdk.io/core/transaction"
	"cosmossdk.io/server/v2/cometbft"
	"github.com/cosmos/cosmos-sdk/client"
)

var _ transaction.Codec[transaction.Tx] = &temporaryTxDecoder{}

type temporaryTxDecoder struct {
	txConfig client.TxConfig
}

// Decode implements transaction.Codec.
func (t *temporaryTxDecoder) Decode(bz []byte) (transaction.Tx, error) {
	return t.txConfig.TxDecoder()(bz)
}

// DecodeJSON implements transaction.Codec.
func (t *temporaryTxDecoder) DecodeJSON(bz []byte) (transaction.Tx, error) {
	return t.txConfig.TxJSONDecoder()(bz)
}

// Temp hood
func ProvideCometServer(clientCtx client.Context) *cometbft.CometBFTServer[transaction.Tx] {
	return cometbft.New(&temporaryTxDecoder{clientCtx.TxConfig})
}
