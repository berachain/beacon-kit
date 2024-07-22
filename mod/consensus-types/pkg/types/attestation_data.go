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

package types

import "github.com/berachain/beacon-kit/mod/primitives/pkg/common"

// AttestationData represents an attestation data.
//
//go:generate go run github.com/ferranbt/fastssz/sszgen -path attestation_data.go -objs AttestationData -include ../../../primitives/pkg/bytes,../../../primitives/pkg/common -output attestation_data.ssz.go
type AttestationData struct {
	// Slot is the slot number of the attestation data.
	Slot uint64
	// Index is the index of the validator.
	Index uint64
	// BeaconBlockRoot is the root of the beacon block.
	BeaconBlockRoot common.Root
}

// New creates a new attestation data instance.
func (a *AttestationData) New(
	slot, index uint64,
	beaconBlockRoot common.Root,
) *AttestationData {
	a = &AttestationData{
		Slot:            slot,
		Index:           index,
		BeaconBlockRoot: beaconBlockRoot,
	}
	return a
}

// GetSlot returns the slot of the attestation data.
func (a *AttestationData) GetSlot() uint64 {
	return a.Slot
}

// GetIndex returns the index of the attestation data.
func (a *AttestationData) GetIndex() uint64 {
	return a.Index
}

// GetBeaconBlockRoot returns the beacon block root of the attestation data.
func (a *AttestationData) GetBeaconBlockRoot() common.Root {
	return a.BeaconBlockRoot
}

// SetSlot sets the slot of the attestation data.
func (a *AttestationData) SetSlot(slot uint64) {
	a.Slot = slot
}

// SetIndex sets the index of the attestation data.
func (a *AttestationData) SetIndex(index uint64) {
	a.Index = index
}

// SetBeaconBlockRoot sets the beacon block root of the attestation data.
func (a *AttestationData) SetBeaconBlockRoot(root common.Root) {
	a.BeaconBlockRoot = root
}
