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
	"github.com/berachain/beacon-kit/config"
	"github.com/berachain/beacon-kit/primitives"
)

// Domain is the domain used for signing.
type Domain [DomainLength]byte

// Bytes returns the byte representation of the Domain.
func (d *Domain) Bytes() []byte {
	return d[:]
}

// computeDomain returns the domain for the DomainType and fork version.
func computeDomain(
	domainType DomainType,
	forkVersion primitives.Version,
	genesisValidatorsRoot primitives.Root,
) (Domain, error) {
	forkDataRoot, err := computeForkDataRoot(forkVersion, genesisValidatorsRoot)
	if err != nil {
		return Domain{}, err
	}
	var bz []byte
	bz = append(bz, domainType[:]...)
	bz = append(
		bz,
		forkDataRoot[:(primitives.RootLength-DomainTypeLength)]...)
	return Domain(bz), nil
}

// GetDomain returns the domain for the DomainType and epoch.
func GetDomain(
	cfg *config.Config,
	genesisValidatorsRoot primitives.Root,
	domainType DomainType,
	epoch primitives.Epoch,
) (Domain, error) {
	return computeDomain(
		domainType,
		VersionFromUint32(cfg.Beacon.ActiveForkVersionByEpoch(epoch)),
		genesisValidatorsRoot,
	)
}
