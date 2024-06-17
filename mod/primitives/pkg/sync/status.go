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

package sync

// ClientStatus represents the synchronization status of a CL.
type ClientStatus uint8

// Constants representing the possible states of ELStatus.
const (
	// StatusDisconnected indicates that the client is disconnected.
	StatusDisconnected ClientStatus = iota
	// StatusNotSynced indicates that the client is not synced.
	StatusNotSynced
	// StatusSynced indicates that the client is synced.
	StatusSynced
)

// Status represents the synchronization status of both CL and EL.
type Status struct {
	// clStatus represents the synchronization status of a CL.
	clStatus ClientStatus
	// elStatus represents the synchronization status of an EL.
	elStatus ClientStatus
}

// Healthy returns true if both CL and EL are synced.
func (s *Status) Healthy() bool {
	return s.clStatus == StatusSynced &&
		s.elStatus == StatusSynced
}

// CLStatus returns the synchronization status of a CL.
func (s *Status) CLStatus() ClientStatus {
	return s.clStatus
}

// ELStatus returns the synchronization status of an EL.
func (s *Status) ELStatus() ClientStatus {
	return s.elStatus
}
