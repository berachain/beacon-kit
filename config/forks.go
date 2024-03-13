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
	"github.com/berachain/beacon-kit/config/flags"
	"github.com/berachain/beacon-kit/io/cli/parser"
	"github.com/berachain/beacon-kit/primitives"
)

const (
	defaultElectraForkEpoch   = 9999999999999999
	defaultDenebForkEpoch     = 0 // Deneb is supported from the genesis.
	defaultGenesisForkVersion = 0
	defaultDenebForkVersion   = 4
)

// Forks conforms to the BeaconKitConfig interface.
var _ BeaconKitConfig[Forks] = &Forks{}

// DefaultForksConfig returns the default forks configuration.
func DefaultForksConfig() Forks {
	return Forks{
		SlotsPerEpoch: primitives.Slot(1),
		ElectraForkEpoch: primitives.Epoch(
			defaultElectraForkEpoch,
		),
		GenesisForkVersion: defaultGenesisForkVersion,
		DenebForkVersion:   defaultDenebForkVersion,
		DenebForkEpoch:     primitives.Epoch(defaultDenebForkEpoch),
	}
}

// Config represents the configuration struct for the forks.
type Forks struct {
	SlotsPerEpoch primitives.Slot
	// ElectraForkEpoch is used to represent the assigned fork epoch for
	// electra.
	ElectraForkEpoch primitives.Epoch
	// GenesisForkVersion represents the genesis fork version.
	GenesisForkVersion uint32
	// DenebForkVersion represents the Deneb fork version.
	// We skip Altair, Bellatrix, and Capella.
	DenebForkVersion uint32
	// DenebForkEpoch is used to represent
	// the assigned fork epoch for Deneb.
	DenebForkEpoch primitives.Epoch
}

// Parse parses the configuration.
func (c Forks) Parse(parser parser.AppOptionsParser) (*Forks, error) {
	var err error
	if c.SlotsPerEpoch, err = parser.GetUint64(
		flags.SlotsPerEpoch,
	); err != nil {
		return nil, err
	}

	if c.ElectraForkEpoch, err = parser.GetEpoch(
		flags.ElectraForkEpoch,
	); err != nil {
		return nil, err
	}

	if c.DenebForkEpoch, err = parser.GetEpoch(
		flags.DenebForkEpoch,
	); err != nil {
		return nil, err
	}

	if c.GenesisForkVersion, err = parser.GetUint32(
		flags.GenesisForkVersion,
	); err != nil {
		return nil, err
	}

	if c.DenebForkVersion, err = parser.GetUint32(
		flags.DenebForkVersion,
	); err != nil {
		return nil, err
	}

	return &c, nil
}

// Template returns the configuration template.
func (c Forks) Template() string {
	return `
[beacon-kit.beacon-config.forks]
# slots per epoch
slots-per-epoch = {{.BeaconKit.Beacon.Forks.SlotsPerEpoch}}
# Electra fork epoch
electra-fork-epoch = {{.BeaconKit.Beacon.Forks.ElectraForkEpoch}}
# Deneb fork epoch
deneb-fork-epoch = {{.BeaconKit.Beacon.Forks.DenebForkEpoch}}
# Genesis fork version
genesis-fork-version = {{.BeaconKit.Beacon.Forks.GenesisForkVersion}}
# Deneb fork version
deneb-fork-version = {{.BeaconKit.Beacon.Forks.DenebForkVersion}}
`
}

// ForkAtEpoch returns the fork version at the given epoch.
func (c Forks) ForkAtEpoch(epoch primitives.Epoch) uint32 {
	if epoch < c.DenebForkEpoch {
		return c.GenesisForkVersion
	}
	// Deneb is the latest supported fork.
	return c.DenebForkVersion
}
