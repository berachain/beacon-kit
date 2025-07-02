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
	"github.com/berachain/beacon-kit/primitives/math"
)

// SlashingInfoSize is the size of the SlashingInfo object in SSZ encoding.
const SlashingInfoSize = 16 // 8 bytes for Slot + 8 bytes for Index

// Compile-time assertion to ensure SlashingInfoSize matches the SizeSSZ method.
var _ = [1]struct{}{}[16-SlashingInfoSize]

//go:generate sszgen --path . --include ../../primitives/math --objs SlashingInfo --output slashing_info_sszgen.go

// SlashingInfo represents a slashing info.
type SlashingInfo struct {
	// Slot is the slot number of the slashing info.
	Slot math.Slot
	// ValidatorIndex is the validator index of the slashing info.
	Index math.U64
}

func (*SlashingInfo) ValidateAfterDecodingSSZ() error { return nil }

/* -------------------------------------------------------------------------- */
/*                             Getters and Setters                            */
/* -------------------------------------------------------------------------- */

// GetSlot returns the slot of the slashing info.
func (s *SlashingInfo) GetSlot() math.Slot {
	return s.Slot
}

// GetIndex returns the index of the slashing info.
func (s *SlashingInfo) GetIndex() math.U64 {
	return s.Index
}

// SetSlot sets the slot of the slashing info.
func (s *SlashingInfo) SetSlot(slot math.Slot) {
	s.Slot = slot
}

// SetIndex sets the index of the slashing info.
func (s *SlashingInfo) SetIndex(index math.U64) {
	s.Index = index
}
