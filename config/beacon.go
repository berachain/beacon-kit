// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	"github.com/prysmaticlabs/prysm/v4/runtime/version"
)

// Beacon represents configuration options for the consensus layer.
type Beacon struct {
	// AltairForkEpoch is used to represent the assigned fork epoch for altair.
	AltairForkEpoch primitives.Epoch
	// BellatrixForkEpoch is used to represent the assigned fork epoch for bellatrix.
	BellatrixForkEpoch primitives.Epoch
	// CapellaForkEpoch is used to represent the assigned fork epoch for capella.
	CapellaForkEpoch primitives.Epoch
	// DenebForkEpoch is used to represent the assigned fork epoch for deneb.
	DenebForkEpoch primitives.Epoch
}

// DefaultBeaconConfig returns the default fork configuration.
func DefaultBeaconConfig() Beacon {
	return Beacon{
		AltairForkEpoch:    0,
		BellatrixForkEpoch: 0,
		CapellaForkEpoch:   0,
		DenebForkEpoch:     primitives.Epoch(4294967295), //nolint:gomnd // we want it disabled rn.
	}
}

// ActiveForkVersion returns the active fork version for a given slot.
func (c *Beacon) ActiveForkVersion(epoch primitives.Epoch) int {
	if epoch >= c.DenebForkEpoch {
		return version.Deneb
	}

	if epoch >= c.CapellaForkEpoch {
		return version.Capella
	}

	if epoch >= c.BellatrixForkEpoch {
		return version.Bellatrix
	}

	if epoch >= c.AltairForkEpoch {
		return version.Altair
	}

	return (version.Phase0)
}
