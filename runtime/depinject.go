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

package runtime

import (
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/config"
)

//nolint:gochecknoinits // GRRRR fix later.
func init() {
	if err := depinject.Inject(
		depinject.Provide(ProvideRuntime),
	); err != nil {
		panic(err)
	}
}

// DepInjectInput is the input for the dep inject framework.
type DepInjectInput struct {
	depinject.In

	Config *config.Config
	Bsp    BeaconStorageBackend
	Vcp    ValsetChangeProvider
	Logger log.Logger
}

// DepInjectOutput is the output for the dep inject framework.
type DepInjectOutput struct {
	depinject.Out

	Runtime *BeaconKitRuntime
}

// ProvideModule is a function that provides the module to the application.
func ProvideRuntime(in DepInjectInput) DepInjectOutput {
	r, err := NewDefaultBeaconKitRuntime(
		in.Config,
		in.Bsp,
		in.Vcp,
		in.Logger,
	)
	if err != nil {
		panic(err)
	}

	return DepInjectOutput{
		Runtime: r,
	}
}
