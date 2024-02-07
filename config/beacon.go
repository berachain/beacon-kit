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
	"github.com/ethereum/go-ethereum/common"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	"github.com/prysmaticlabs/prysm/v4/runtime/version"
)

// Beacon is the configuration for the beacon chain.
type Beacon struct {
	// DenebForkEpoch is used to represent the assigned fork epoch for deneb.
	DenebForkEpoch primitives.Epoch

	// Validator is the configuration for the validator. Only utilized when
	// this node is in the active validator set.
	Validator Validator
}

// Validator is the configuration for the validator. Only utilized when
// this node is in the active validator set.
type Validator struct {
	// Suggested FeeRecipient is the address that will receive the transaction fees
	// produced by any blocks from this node.
	SuggestedFeeRecipient common.Address

	// Grafitti is the string that will be included in the graffiti field of the beacon block.
	Graffiti string

	// PrepareAllPayloads informs the engine to prepare a block on every slot.
	PrepareAllPayloads bool
}

// DefaultBeaconConfig returns the default fork configuration.
func DefaultBeaconConfig() Beacon {
	return Beacon{
		DenebForkEpoch: primitives.Epoch(4294967295), //nolint:gomnd // we want it disabled rn.
		Validator:      DefaultValidatorConfig(),
	}
}

// DefaultValidatorConfig returns the default validator configuration.
func DefaultValidatorConfig() Validator {
	return Validator{
		SuggestedFeeRecipient: common.Address{},
		Graffiti:              "",
		PrepareAllPayloads:    true,
	}
}

// ActiveForkVersion returns the active fork version for a given slot.
func (c Beacon) ActiveForkVersion(epoch primitives.Epoch) int {
	if epoch >= c.DenebForkEpoch {
		return version.Deneb
	}

	// In BeaconKit we assume the Capella fork is always active.
	return version.Capella
}
