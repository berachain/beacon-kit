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
	"github.com/berachain/beacon-kit/config/version"
	"github.com/berachain/beacon-kit/io/cli/parser"
	"github.com/berachain/beacon-kit/primitives"
)

// Beacon conforms to the BeaconKitConfig interface.
var _ BeaconKitConfig[Beacon] = Beacon{}

// Beacon is the configuration for the beacon chain.
type Beacon struct {
	// Execution is the configuration for the execution service.
	Execution Execution
	// Forks is the configuration for the beacon chain forks.
	Forks Forks
	// Limits is the configuration for limits (max/min) on the beacon chain.
	Limits Limits
	// Validator is the configuration for the validator. Only utilized when
	// this node is in the active validator set.
	Validator Validator
}

// DefaultBeaconConfig returns the default fork configuration.
func DefaultBeaconConfig() Beacon {
	return Beacon{
		Execution: DefaultExecutionConfig(),
		Forks:     DefaultForksConfig(),
		Limits:    DefaultLimitsConfig(),
		Validator: DefaultValidatorConfig(),
	}
}

// ActiveForkVersion returns the active fork version for a given slot.
func (c Beacon) ActiveForkVersion(epoch primitives.Epoch) int {
	if epoch >= c.Forks.ElectraForkEpoch {
		return version.Electra
	}

	// In BeaconKit we assume the Deneb fork is always active.
	return version.Deneb
}

// Parse parses the configuration.
func (c Beacon) Parse(parser parser.AppOptionsParser) (*Beacon, error) {
	var (
		err       error
		forks     *Forks
		limits    *Limits
		validator *Validator
		execution *Execution
	)

	// Parse the forks configuration.
	if forks, err = c.Forks.Parse(parser); err != nil {
		return nil, err
	}
	c.Forks = *forks

	// Parse the limits configuration.
	if limits, err = c.Limits.Parse(parser); err != nil {
		return nil, err
	}
	c.Limits = *limits

	// Parse the validator configuration.
	if validator, err = c.Validator.Parse(parser); err != nil {
		return nil, err
	}
	c.Validator = *validator

	// Parse the execution configuration.
	if execution, err = c.Execution.Parse(parser); err != nil {
		return nil, err
	}
	c.Execution = *execution

	return &c, nil
}

// Template returns the configuration template.
func (c Beacon) Template() string {
	return c.Forks.Template() + c.Limits.Template() +
		c.Validator.Template() + c.Execution.Template()
}
