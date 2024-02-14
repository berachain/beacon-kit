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
	"github.com/itsdevbear/bolaris/config/flags"
	"github.com/itsdevbear/bolaris/config/parser"
	"github.com/itsdevbear/bolaris/types/consensus/primitives"
)

// Forks conforms to the BeaconKitConfig interface.
var _ BeaconKitConfig[Forks] = &Forks{}

// DefaultForksConfig returns the default forks configuration.
func DefaultForksConfig() Forks {
	return Forks{
		DenebForkEpoch: primitives.Epoch(4294967295), //nolint:gomnd // we want it disabled rn.
	}
}

// Config represents the configuration struct for the forks.
type Forks struct {
	// DenebForkEpoch is used to represent the assigned fork epoch for deneb.
	DenebForkEpoch primitives.Epoch
}

// Parse parses the configuration.
func (c Forks) Parse(parser parser.AppOptionsParser) (*Forks, error) {
	var err error
	if c.DenebForkEpoch, err = parser.GetEpoch(
		flags.DenebForkEpoch,
	); err != nil {
		return nil, err
	}

	return &c, nil
}

// Template returns the configuration template.
func (c Forks) Template() string {
	return `
[beacon-kit.beacon-config.forks]
# Deneb fork epoch
deneb-fork-epoch = {{.BeaconKit.Beacon.Forks.DenebForkEpoch}}
`
}
