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

// SlotData represents the data to be used to propose a block.
type SlotData struct {
	// Slot is the slot number of the incoming slot.
	Slot uint64
	// AttestationData is the attestation data of the incoming slot.
	AttestationData []*AttestationData
	// SlashingInfo is the slashing info of the incoming slot.
	SlashingInfo []*SlashingInfo
}

// GetSlot retrieves the slot of the SlotData.
func (b *SlotData) GetSlot() uint64 {
	return b.Slot
}

// GetAttestationData retrieves the attestation data of the SlotData.
func (b *SlotData) GetAttestationData() []*AttestationData {
	return b.AttestationData
}

// GetSlashingInfo retrieves the slashing info of the SlotData.
func (b *SlotData) GetSlashingInfo() []*SlashingInfo {
	return b.SlashingInfo
}

// SetSlot sets the slot of the SlotData.
func (b *SlotData) SetSlot(slot uint64) {
	b.Slot = slot
}

// SetAttestationData sets the attestation data of the SlotData.
func (b *SlotData) SetAttestationData(attestationData []*AttestationData) {
	b.AttestationData = attestationData
}

// SetSlashingInfo sets the slashing info of the SlotData.
func (b *SlotData) SetSlashingInfo(slashingInfo []*SlashingInfo) {
	b.SlashingInfo = slashingInfo
}
