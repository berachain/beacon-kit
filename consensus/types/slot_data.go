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
	"time"

	deneb "github.com/berachain/beacon-kit/consensus-types/deneb"
	"github.com/berachain/beacon-kit/primitives/math"
)

// SlotData represents the data to be used to propose a block.
type SlotData struct {
	// slot is the slot number of the incoming slot.
	slot math.Slot
	// attestationData is the attestation data of the incoming slot.
	attestationData []*deneb.AttestationData
	// slashingInfo is the slashing info of the incoming slot.
	slashingInfo []*deneb.SlashingInfo

	// some consensus data useful to build and verify the block
	*commonConsensusData
}

// NewSlotData creates a new SlotData instance.
func NewSlotData(
	slot math.Slot,
	attestationData []*deneb.AttestationData,
	slashingInfo []*deneb.SlashingInfo,
	proposerAddress []byte,
	consensusTime time.Time,
) *SlotData {
	return &SlotData{
		slot:            slot,
		attestationData: attestationData,
		slashingInfo:    slashingInfo,
		commonConsensusData: &commonConsensusData{
			proposerAddress: proposerAddress,
			consensusTime:   math.U64(consensusTime.Unix()), // #nosec G115
		},
	}
}

// GetSlot retrieves the slot of the SlotData.
func (b *SlotData) GetSlot() math.Slot {
	return b.slot
}

// GetAttestationData retrieves the attestation data of the SlotData.
func (b *SlotData) GetAttestationData() []*deneb.AttestationData {
	return b.attestationData
}

// GetSlashingInfo retrieves the slashing info of the SlotData.
func (b *SlotData) GetSlashingInfo() []*deneb.SlashingInfo {
	return b.slashingInfo
}

// SetAttestationData sets the attestation data of the SlotData.
func (b *SlotData) SetAttestationData(
	attestationData []*deneb.AttestationData,
) {
	b.attestationData = attestationData
}

// SetSlashingInfo sets the slashing info of the SlotData.
func (b *SlotData) SetSlashingInfo(
	slashingInfo []*deneb.SlashingInfo,
) {
	b.slashingInfo = slashingInfo
}
