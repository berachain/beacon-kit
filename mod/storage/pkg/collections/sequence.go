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

package collections

import (
	sdkcollections "cosmossdk.io/collections"
)

// DefaultSequenceStart defines the default starting number of a sequence.
const DefaultSequenceStart uint64 = 0

// Sequence builds on top of an Item, and represents a monotonically increasing number.
type Sequence struct {
	ItemKeeper[uint64]
}

// NewSequence instantiates a new sequence given
// a Schema, a Prefix and humanized name for the sequence.
func NewSequence(
	storeKey, key []byte, storeAccessor StoreAccessor,
) Sequence {
	seq := Sequence{
		NewItemKeeper(
			storeKey,
			key,
			sdkcollections.Uint64Value,
			storeAccessor,
		),
	}
	// i worry this breaks cause the collection set is not
	// committed yet.
	seq.Set(DefaultSequenceStart)
	return seq
}

// Peek returns the current sequence value, if no number
// is set then the DefaultSequenceStart is returned.
// Errors on encoding issues.
func (s *Sequence) Peek() (uint64, error) {
	return s.ItemKeeper.Get()
}

// Next returns the next sequence number, and sets the next expected sequence.
// Errors on encoding issues.
func (s *Sequence) Next() (uint64, error) {
	seq, err := s.Peek()
	if err != nil {
		return 0, err
	}
	return seq, s.Set(seq + 1)
}

// Set hard resets the sequence to the provided value.
// Errors on encoding issues.
func (s *Sequence) Set(value uint64) error {
	return s.ItemKeeper.Set(value)
}
