// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
	ssztypes "github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/types"
)

// Fork as defined in the Ethereum 2.0 specification:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#fork
//
//go:generate go run github.com/ferranbt/fastssz/sszgen -path fork.go -objs Fork -include ../../../primitives/pkg/bytes,../../../primitives/pkg/math,../../../primitives/pkg/common -output fork.ssz.go
//nolint:lll
type Fork struct {
	// PreviousVersion is the last version before the fork.
	PreviousVersion common.Version
	// CurrentVersion is the first version after the fork.
	CurrentVersion common.Version
	// Epoch is the epoch at which the fork occurred.
	Epoch math.Epoch
}

// New creates a new fork.
func (f *Fork) New(
	previousVersion common.Version,
	currentVersion common.Version,
	epoch math.Epoch,
) *Fork {
	return &Fork{
		PreviousVersion: previousVersion,
		CurrentVersion:  currentVersion,
		Epoch:           epoch,
	}
}

func (*Fork) Schema() *ssz.Schema[*Fork] {
	s := &ssz.Schema[*Fork]{}
	s.DefineField(
		"previous_version",
		func(f *Fork) ssztypes.MinimalSSZType { return f.PreviousVersion },
	)
	s.DefineField(
		"current_version",
		func(f *Fork) ssztypes.MinimalSSZType { return f.CurrentVersion },
	)
	s.DefineField(
		"epoch",
		func(f *Fork) ssztypes.MinimalSSZType { return ssz.U64(f.Epoch) },
	)
	return s
}

func (f *Fork) Default() *Fork {
	if f == nil {
		return &Fork{}
	}
	return f
}
