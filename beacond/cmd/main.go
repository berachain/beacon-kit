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

package main

import (
	"os"

	app "github.com/berachain/beacon-kit/beacond/app"
	nodebuilder "github.com/berachain/beacon-kit/mod/node-builder"
)

func main() {
	nb := nodebuilder.NewNodeBuilder[app.BeaconApp]().
		WithAppInfo(
			&nodebuilder.AppInfo[app.BeaconApp]{
				Name:            "beacond",
				Description:     "beacond is a beacon node for any beacon-kit chain",
				Creator:         app.NewBeaconKitAppWithDefaultBaseAppOptions,
				DepInjectConfig: app.Config(),
			},
		)
	if err := nb.RunNode(); err != nil {
		os.Exit(1)
	}
}
