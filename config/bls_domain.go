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
	byteslib "github.com/berachain/beacon-kit/lib/bytes"
	"github.com/berachain/beacon-kit/primitives"
)

// BLSDomains conforms to the BeaconKitConfig interface.
var _ BeaconKitConfig[BLSDomains] = &BLSDomains{}

const (
	defaultDomainBeaconProposer    = 0x00000000
	defaultDomainBeaconAttester    = 0x01000000
	defaultDomainRandao            = 0x02000000
	defaultDomainDeposit           = 0x03000000
	defaultDomainVoluntaryExit     = 0x04000000
	defaultDomainSelectionProof    = 0x05000000
	defaultDomainAggregateAndProof = 0x06000000
	defaultDomainApplicationMask   = 0x00000001
)

// DefaultBLSDomainsConfig returns the default BLS domain configuration.
func DefaultBLSDomainsConfig() BLSDomains {
	return BLSDomains{
		DomainBeaconProposer:    byteslib.Uint32ToBytes4(defaultDomainBeaconProposer),
		DomainBeaconAttester:    byteslib.Uint32ToBytes4(defaultDomainBeaconAttester),
		DomainRandao:            byteslib.Uint32ToBytes4(defaultDomainRandao),
		DomainDeposit:           byteslib.Uint32ToBytes4(defaultDomainDeposit),
		DomainVoluntaryExit:     byteslib.Uint32ToBytes4(defaultDomainVoluntaryExit),
		DomainSelectionProof:    byteslib.Uint32ToBytes4(defaultDomainSelectionProof),
		DomainAggregateAndProof: byteslib.Uint32ToBytes4(defaultDomainAggregateAndProof),
		DomainApplicationMask:   byteslib.Uint32ToBytes4(defaultDomainApplicationMask),
	}
}

// BLSDomains is the configuration for BLS domain values.
// Spec: https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#domain-types
//
//nolint:lll // The spec url is long.
type BLSDomains struct {
	// DomainBeaconProposer defines the BLS signature domain
	// for beacon proposal verification.
	DomainBeaconProposer primitives.SSZDomain
	// DomainBeaconAttester defines the BLS signature domain
	// for attestation verification.
	DomainBeaconAttester primitives.SSZDomain
	// DomainRandao defines the BLS signature domain
	// for randao verification.
	DomainRandao primitives.SSZDomain
	// DomainDeposit defines the BLS signature domain
	// for deposit verification.
	DomainDeposit primitives.SSZDomain
	// DomainVoluntaryExit defines the BLS signature domain
	// for exit verification.
	DomainVoluntaryExit primitives.SSZDomain
	// DomainSelectionProof defines the BLS signature domain
	// for selection proof.
	DomainSelectionProof primitives.SSZDomain
	// DomainAggregateAndProof defines the BLS signature domain
	// for aggregate and proof.
	DomainAggregateAndProof primitives.SSZDomain
	// DomainApplicationMask defines the BLS signature domain
	// for application mask.
	DomainApplicationMask primitives.SSZDomain
}

// Parse parses the configuration.
func (c BLSDomains) Parse(parser parser.AppOptionsParser) (*BLSDomains, error) {
	var err error
	if c.DomainBeaconProposer, err = parser.GetBytes4(
		flags.DomainBeaconProposer,
	); err != nil {
		return nil, err
	}
	if c.DomainBeaconAttester, err = parser.GetBytes4(
		flags.DomainBeaconAttester,
	); err != nil {
		return nil, err
	}
	if c.DomainRandao, err = parser.GetBytes4(
		flags.DomainRandao,
	); err != nil {
		return nil, err
	}
	if c.DomainDeposit, err = parser.GetBytes4(
		flags.DomainDeposit,
	); err != nil {
		return nil, err
	}
	if c.DomainVoluntaryExit, err = parser.GetBytes4(
		flags.DomainVoluntaryExit,
	); err != nil {
		return nil, err
	}
	if c.DomainSelectionProof, err = parser.GetBytes4(
		flags.DomainSelectionProof,
	); err != nil {
		return nil, err
	}
	if c.DomainAggregateAndProof, err = parser.GetBytes4(
		flags.DomainAggregateAndProof,
	); err != nil {
		return nil, err
	}
	if c.DomainApplicationMask, err = parser.GetBytes4(
		flags.DomainApplicationMask,
	); err != nil {
		return nil, err
	}
	return &c, nil
}

// Template returns the configuration template.
func (c BLSDomains) Template() string {
	//nolint:lll
	return `
[beacon-kit.beacon-config.bls-domains]
# DomainBeaconProposer defines the BLS signature domain for beacon proposal verification.
domain-beacon-proposer = "{{ .BeaconKit.Beacon.BLSDomains.DomainBeaconProposer }}"
# DomainBeaconAttester defines the BLS signature domain for attestation verification.
domain-beacon-attester = "{{ .BeaconKit.Beacon.BLSDomains.DomainBeaconAttester }}"
# DomainRandao defines the BLS signature domain for randao verification.
domain-randao = "{{ .BeaconKit.Beacon.BLSDomains.DomainRandao }}"
# DomainDeposit defines the BLS signature domain for deposit verification.
domain-deposit = "{{ .BeaconKit.Beacon.BLSDomains.DomainDeposit }}"
# DomainVoluntaryExit defines the BLS signature domain for exit verification.
domain-voluntary-exit = "{{ .BeaconKit.Beacon.BLSDomains.DomainVoluntaryExit }}"
# DomainSelectionProof defines the BLS signature domain for selection proof.
domain-selection-proof = "{{ .BeaconKit.Beacon.BLSDomains.DomainSelectionProof }}"
# DomainAggregateAndProof defines the BLS signature domain for aggregate and proof.
domain-aggregate-and-proof = "{{ .BeaconKit.Beacon.BLSDomains.DomainAggregateAndProof }}"
# DomainApplicationMask defines the BLS signature domain for application mask.
domain-application-mask = "{{ .BeaconKit.Beacon.BLSDomains.DomainApplicationMask }}
`
}
