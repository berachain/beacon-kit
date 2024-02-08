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
	"github.com/itsdevbear/bolaris/config/flags"
	"github.com/itsdevbear/bolaris/config/parser"
	"github.com/itsdevbear/bolaris/third_party/go-ethereum/common"
)

// Validator conforms to the BeaconKitConfig interface.
var _ BeaconKitConfig[Validator] = &Validator{}

// DefaultValidatorConfig returns the default validator configuration.
func DefaultValidatorConfig() Validator {
	return Validator{
		SuggestedFeeRecipient: common.Address{},
		Graffiti:              "",
		PrepareAllPayloads:    true,
	}
}

// Config represents the configuration struct for the validator.
type Validator struct {
	// Suggested FeeRecipient is the address that will receive the transaction fees
	// produced by any blocks from this node.
	SuggestedFeeRecipient common.Address

	// Grafitti is the string that will be included in the graffiti field of the beacon block.
	Graffiti string

	// PrepareAllPayloads informs the engine to prepare a block on every slot.
	PrepareAllPayloads bool
}

// Parse parses the configuration.
func (c Validator) Parse(parser parser.AppOptionsParser) (*Validator, error) {
	var err error
	if c.SuggestedFeeRecipient, err = parser.GetExecutionAddress(
		flags.SuggestedFeeRecipient,
	); err != nil {
		return nil, err
	}
	if c.Graffiti, err = parser.GetString(flags.Graffiti); err != nil {
		return nil, err
	}

	if c.PrepareAllPayloads, err = parser.GetBool(
		flags.PrepareAllPayloads,
	); err != nil {
		return nil, err
	}

	return &c, nil
}

// Template returns the configuration template.
func (c Validator) Template() string {
	return `
[beacon-kit.beacon-config.validator]
# Post bellatrix, this address will receive the transaction fees produced by any blocks 
# from this node.
suggested-fee-recipient = "{{.BeaconKit.Beacon.Validator.SuggestedFeeRecipient}}"

# Graffiti string that will be included in the graffiti field of the beacon block.
graffiti = "{{.BeaconKit.Beacon.Validator.Graffiti}}"

# Prepare all payloads informs the engine to prepare a block on every slot.
prepare-all-payloads = {{.BeaconKit.Beacon.Validator.PrepareAllPayloads}}
`
}
