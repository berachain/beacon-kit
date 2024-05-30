package components

import (
	"errors"

	txsigning "cosmossdk.io/x/tx/signing"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	signingtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
)

// ProvideNoopTxConfig returns a no-op TxConfig.
func ProvideNoopTxConfig() (client.TxConfig, runtime.BaseAppOption) {
	return NoOpTxConfig{}, func(app *baseapp.BaseApp) {
		app.SetTxDecoder(NoOpTxConfig{}.TxDecoder())
		app.SetTxEncoder(NoOpTxConfig{}.TxEncoder())
	}
}

// TxDecoder returns a no-op TxDecoder.
func (NoOpTxConfig) TxDecoder() sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, error) {
		return fakeTx{}, errors.New("skip decoding")
	}
}

// TxJSONEncoder returns a no-op TxJSONEncoder.
func (NoOpTxConfig) TxJSONEncoder() sdk.TxEncoder {
	return func(tx sdk.Tx) ([]byte, error) {
		return nil, errors.New("skip decoding")
	}
}

// TxJSONDecoder returns a no-op TxJSONDecoder.
func (NoOpTxConfig) TxJSONDecoder() sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, error) {
		return fakeTx{}, nil
	}
}

// MarshalSignatureJSON returns a no-op MarshalSignatureJSON.
func (NoOpTxConfig) MarshalSignatureJSON([]signingtypes.SignatureV2) ([]byte, error) {
	return make([]byte, 0), nil
}

// UnmarshalSignatureJSON returns a no-op UnmarshalSignatureJSON.
func (NoOpTxConfig) UnmarshalSignatureJSON([]byte) ([]signingtypes.SignatureV2, error) {
	return nil, nil
}

// NewTxBuilder returns a no-op TxBuilder.
func (NoOpTxConfig) NewTxBuilder() client.TxBuilder {
	return nil
}

// WrapTxBuilder returns a no-op WrapTxBuilder.
func (NoOpTxConfig) WrapTxBuilder(sdk.Tx) (client.TxBuilder, error) {
	return nil, nil
}

// SignModeHandler returns a no-op SignModeHandler.
func (NoOpTxConfig) SignModeHandler() *txsigning.HandlerMap {
	return &txsigning.HandlerMap{}
}

// SigningContext returns a no-op SigningContext.
func (NoOpTxConfig) SigningContext() *txsigning.Context {
	return &txsigning.Context{}
}
