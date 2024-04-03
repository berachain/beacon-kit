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
	"os"

	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/mod/node-builder/components/signer"
	"github.com/berachain/beacon-kit/mod/primitives/constants"
	"github.com/cometbft/cometbft/p2p"
	"github.com/cosmos/cosmos-sdk/client/flags"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cast"
)

// BlsSignerInput is the input for the dep inject framework.
type BlsSignerInput struct {
	depinject.In

	AppOpts servertypes.AppOptions
}

// BlsSignerOutput is the output for the dep inject framework.
type BlsSignerOutput struct {
	depinject.Out

	BlsSigner *signer.BLSSigner
}

func ProvideBlsSigner(in BlsSignerInput) BlsSignerOutput {
	homeDir := cast.ToString(in.AppOpts.Get(flags.FlagHome))

	key, err := NewSignerFromFile(
		homeDir + "/config/priv_validator_key.json",
	)
	if err != nil {
		os.Exit(1)
	}

	return BlsSignerOutput{
		BlsSigner: key,
	}
}

// NewSignerFromFile creates a new Signer instance given a
// file path to a CometBFT node key.
//
// TODO: In order to make RANDAO signing compatible with the BLS12-381
// we need to extend the PrivVal interface to allow signing arbitrary things.
// @tac0turtle.
func NewSignerFromFile(filePath string) (*signer.BLSSigner, error) {
	key, err := p2p.LoadNodeKey(filePath)
	if err != nil {
		return nil, err
	}

	return signer.NewBLSSigner(
		[constants.BLSSecretKeyLength]byte(key.PrivKey.Bytes()))
}
