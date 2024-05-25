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
	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/components/signer"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	clientFlags "github.com/cosmos/cosmos-sdk/client/flags"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cast"
)

// BlsSignerInput is the input for the dep inject framework.
type BlsSignerInput struct {
	depinject.In
	AppOpts servertypes.AppOptions
	PrivKey LegacyKey `optional:"true"`
}

// type alias to LegacyKey used for LegacySinger construction.
type LegacyKey = signer.LegacyKey

// ProvideBlsSigner is a function that provides the module to the application.
func ProvideBlsSigner(in BlsSignerInput) (crypto.BLSSigner, error) {
	if in.PrivKey == [constants.BLSSecretKeyLength]byte{} {
		// if no private key is provided, use privval signer
		homeDir := cast.ToString(in.AppOpts.Get(clientFlags.FlagHome))
		return signer.NewBLSSigner(
			homeDir+"/config/priv_validator_key.json",
			homeDir+"/data/priv_validator_state.json",
		), nil
	}
	return signer.NewLegacySigner(in.PrivKey)
}

func GetLegacyKey(privKey string) (LegacyKey, error) {
	return signer.LegacyKeyFromString(privKey)
}
