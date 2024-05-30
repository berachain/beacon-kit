// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package components

import (
	"errors"

	txsigning "cosmossdk.io/x/tx/signing"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	signingtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
)

// ProvideNoopTxConfig returns a no-op TxConfig.
func ProvideNoopTxConfig() client.TxConfig {
	return NoOpTxConfig{}
}

// NoOpTxConfig is a no-op implementation of the TxConfig interface.
type NoOpTxConfig struct{}

// TxEncoder returns a no-op TxEncoder.
func (NoOpTxConfig) TxEncoder() sdk.TxEncoder {
	return func(sdk.Tx) ([]byte, error) {
		return nil, nil
	}
}

// TxDecoder returns a no-op TxDecoder.
func (NoOpTxConfig) TxDecoder() sdk.TxDecoder {
	return func([]byte) (sdk.Tx, error) {
		return nil, errors.New("skip decoding")
	}
}

// TxJSONEncoder returns a no-op TxJSONEncoder.
func (NoOpTxConfig) TxJSONEncoder() sdk.TxEncoder {
	return func(sdk.Tx) ([]byte, error) {
		return nil, errors.New("skip decoding")
	}
}

// TxJSONDecoder returns a no-op TxJSONDecoder.
func (NoOpTxConfig) TxJSONDecoder() sdk.TxDecoder {
	return func([]byte) (sdk.Tx, error) {
		return nil, nil
	}
}

// MarshalSignatureJSON returns a no-op MarshalSignatureJSON.
func (NoOpTxConfig) MarshalSignatureJSON(
	[]signingtypes.SignatureV2,
) ([]byte, error) {
	return make([]byte, 0), nil
}

// UnmarshalSignatureJSON returns a no-op UnmarshalSignatureJSON.
func (NoOpTxConfig) UnmarshalSignatureJSON(
	[]byte,
) ([]signingtypes.SignatureV2, error) {
	return nil, nil
}

// NewTxBuilder returns a no-op TxBuilder.
func (NoOpTxConfig) NewTxBuilder() client.TxBuilder {
	return nil
}

// WrapTxBuilder returns a no-op WrapTxBuilder.
func (NoOpTxConfig) WrapTxBuilder(sdk.Tx) (client.TxBuilder, error) {
	//nolint:nilnil // its whatever.
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
