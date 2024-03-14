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

package config

import (
	"cosmossdk.io/depinject"
	"github.com/cosmos/cosmos-sdk/client/flags"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cast"
)

// DepInjectInput is the input for the dep inject framework.
type DepInjectInput struct {
	depinject.In

	AppOpts servertypes.AppOptions
}

// DepInjectOutput is the output for the dep inject framework.
type DepInjectOutput struct {
	depinject.Out

	Network Network
}

// ProvideNetworkCfg get the CometBFT network config
// from the app options and inject it into the dep inject output.
func ProvideNetworkCfg(in DepInjectInput) DepInjectOutput {
	chainID := cast.ToString(in.AppOpts.Get(flags.FlagChainID))

	return DepInjectOutput{
		Network: Network{
			ChainID: chainID,
		},
	}
}
