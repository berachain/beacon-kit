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
	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/core"
	"github.com/berachain/beacon-kit/mod/node-builder/config"
	"github.com/berachain/beacon-kit/mod/runtime"
)

// RuntimeInjectInput is the input for the dep inject framework.
type RuntimeInjectInput struct {
	depinject.In

	BeaconCfg *config.Config
	Logger    log.Logger
	Signer    core.BLSSigner
	Bsp       runtime.BeaconStorageBackend
}

// RuntimeInjectOutput is the output for the dep inject framework.
type RuntimeInjectOutput struct {
	depinject.Out

	Runtime *runtime.BeaconKitRuntime
}

// ProvideRuntime is a function that provides the module to the application.
func ProvideRuntime(in RuntimeInjectInput) RuntimeInjectOutput {
	r, err := runtime.NewDefaultBeaconKitRuntime(
		in.BeaconCfg,
		in.Signer,
		in.Logger,
		in.Bsp,
	)
	if err != nil {
		in.Logger.Error(
			"failed to create beacon-kit runtime, exiting...", "error", err)
		os.Exit(1)
	}

	return RuntimeInjectOutput{
		Runtime: r,
	}
}
