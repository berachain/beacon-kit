// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package types

import (
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constraints"
)

// Compile-time check to ensure Versionable implements Versionable interface
var _ constraints.Versionable = Versionable{}

// Versionable is a concrete type that implements the Versionable interface.
// This is used instead of embedding the interface directly for SSZ generation compatibility.
// The fork version is excluded from SSZ encoding using the ssz:"-" tag.
type Versionable struct {
	forkVersion common.Version `ssz:"-" json:"-"`
}

// NewVersionable creates a new Versionable with the given fork version.
func NewVersionable(forkVersion common.Version) Versionable {
	return Versionable{forkVersion: forkVersion}
}

// GetForkVersion returns the fork version.
func (v Versionable) GetForkVersion() common.Version {
	return v.forkVersion
}

// MarshalJSON returns null to exclude from JSON output.
func (v Versionable) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
}

// UnmarshalJSON does nothing since version is set elsewhere.
func (v *Versionable) UnmarshalJSON([]byte) error {
	return nil
}
