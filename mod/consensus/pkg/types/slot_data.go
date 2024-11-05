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
	"time"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// SlotData represents the data to be used to propose a block.
type SlotData[AttestationDataT, SlashingInfoT any] struct {
	// slot is the slot number of the incoming slot.
	slot math.Slot
	// attestationData is the attestation data of the incoming slot.
	attestationData []AttestationDataT
	// slashingInfo is the slashing info of the incoming slot.
	slashingInfo []SlashingInfoT
	// nextPayloadTimestamp is the timestamp proposed by
	// consensus for the next payload to be proposed. It is also
	// used to bound current payload upon validation
	nextPayloadTimestamp math.U64
}

// New creates a new SlotData instance.
func (b *SlotData[AttestationDataT, SlashingInfoT]) New(
	slot math.Slot,
	attestationData []AttestationDataT,
	slashingInfo []SlashingInfoT,
	nextPayloadTimestamp time.Time,
) *SlotData[AttestationDataT, SlashingInfoT] {
	b = &SlotData[AttestationDataT, SlashingInfoT]{
		slot:                 slot,
		attestationData:      attestationData,
		slashingInfo:         slashingInfo,
		nextPayloadTimestamp: math.U64(nextPayloadTimestamp.Unix()),
	}
	return b
}

// GetSlot retrieves the slot of the SlotData.
func (b *SlotData[AttestationDataT, SlashingInfoT]) GetSlot() math.Slot {
	return b.slot
}

// GetAttestationData retrieves the attestation data of the SlotData.
func (b *SlotData[
	AttestationDataT,
	SlashingInfoT,
]) GetAttestationData() []AttestationDataT {
	return b.attestationData
}

// GetSlashingInfo retrieves the slashing info of the SlotData.
func (b *SlotData[
	AttestationDataT,
	SlashingInfoT,
]) GetSlashingInfo() []SlashingInfoT {
	return b.slashingInfo
}

// GetNextPayloadTimestamp retrieves the proposed next payload timestamp.
func (b *SlotData[
	AttestationDataT,
	SlashingInfoT,
]) GetNextPayloadTimestamp() math.U64 {
	return b.nextPayloadTimestamp
}

// SetAttestationData sets the attestation data of the SlotData.
func (b *SlotData[AttestationDataT, SlashingInfoT]) SetAttestationData(
	attestationData []AttestationDataT,
) {
	b.attestationData = attestationData
}

// SetSlashingInfo sets the slashing info of the SlotData.
func (b *SlotData[AttestationDataT, SlashingInfoT]) SetSlashingInfo(
	slashingInfo []SlashingInfoT,
) {
	b.slashingInfo = slashingInfo
}
