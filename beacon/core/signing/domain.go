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

package signing

import (
	"github.com/berachain/beacon-kit/primitives"
	"github.com/prysmaticlabs/prysm/v5/config/params"
)

// Domain is the domain used for signing.
type Domain [DomainLength]byte

// Bytes returns the byte representation of the Domain.
func (d *Domain) Bytes() []byte {
	return d[:]
}

// ComputeDomain returns the domain for the DomainType and fork version.
// def compute_domain(domain_type: DomainType, fork_version: Version=None, genesis_validators_root: Root=None) -> Domain:
//
//	"""
//	Return the domain for the domain_type and fork_version.
//	"""
//	if fork_version is None:
//	    fork_version = GENESIS_FORK_VERSION
//	if genesis_validators_root is None:
//	    genesis_validators_root = Root()  # all bytes zero by default
//	fork_data_root = compute_fork_data_root(fork_version, genesis_validators_root)
//	return Domain(domain_type + fork_data_root[:28])
func ComputeDomain(
	domainType DomainType,
	forkVersion Version,
	genesisValidatorsRoot primitives.HashRoot,
) Domain {
	if forkVersion == nil {
		forkVersion = params.BeaconConfig().GenesisForkVersion
	}
	if genesisValidatorsRoot == nil {
		genesisValidatorsRoot = params.BeaconConfig().ZeroHash[:]
	}
	var forkBytes [ForkVersionByteLength]byte
	copy(forkBytes[:], forkVersion)

	forkDataRoot, err := computeForkDataRoot(forkBytes[:], genesisValidatorsRoot)
	if err != nil {
		return nil, err
	}

	return domain(domainType, forkDataRoot[:]), nil

}

// This returns the bls domain given by the domain type and fork data root.
func domain(domainType [DomainByteLength]byte, forkDataRoot []byte) []byte {
	var b []byte
	b = append(b, domainType[:4]...)
	b = append(b, forkDataRoot[:28]...)
	return b
}
