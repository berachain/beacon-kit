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
	"encoding/hex"

	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/components/signer"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/config/flags"
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
	PrivKey string `optional:"true"`
}

// ProvideBlsSigner is a function that provides the module to the application.
func ProvideBlsSigner(in BlsSignerInput) (crypto.BLSSigner, error) {
	defer func() {
		if r := recover(); r != nil {
			// TODO: NewNodeBuilder called before depinject AppOpts set
			return
		}
	}()
	signerType := cast.ToString(in.AppOpts.Get(flags.SignerType))
	if signerType == "bls" {
		homeDir := cast.ToString(in.AppOpts.Get(clientFlags.FlagHome))
		return signer.NewBLSSigner(
			homeDir+"/config/priv_validator_key.json",
			homeDir+"/data/priv_validator_state.json",
		), nil
	} else if signerType == "legacy" {
		return provideLegacy(in.PrivKey)
	}
	// TODO: panicking here so we can recover bc depinject late AF
	panic("invalid signer type")
}

func provideLegacy(privKey string) (crypto.BLSSigner, error) {
	if privKey == "" {
		return nil, signer.ErrValidatorPrivateKeyRequired
	}
	privKeyBz, err := hex.DecodeString(privKey)
	if err != nil {
		return nil, err
	}
	if len(privKeyBz) != constants.BLSSecretKeyLength {
		return nil, signer.ErrInvalidValidatorPrivateKeyLength
	}

	// creates bls signer that signs with the specified private key
	return signer.NewLegacySigner(
		[constants.BLSSecretKeyLength]byte(privKeyBz),
	)
}
