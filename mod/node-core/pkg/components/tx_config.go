// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

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
