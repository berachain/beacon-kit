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

// CLStatus represents the synchronization status of a CL.
type CLStatus uint8

// ELStatus represents the synchronization status of an EL.
type ELStatus uint8

// Constants representing the possible states of CLStatus.
const (
	// CLStatusNotSynced indicates that the CL is not synced.
	CLStatusNotSynced CLStatus = iota
	// CLStatusSynced indicates that the CL is synced.
	CLStatusSynced
)

// Constants representing the possible states of ELStatus.
const (
	// ELStatusDisconnected indicates that the EL is disconnected.
	ELStatusDisconnected ELStatus = iota
	// ELStatusNotSynced indicates that the EL is not synced.
	ELStatusNotSynced
	// ELStatusSynced indicates that the EL is synced.
	ELStatusSynced
)

// Status represents the synchronization status of both CL and EL.
type Status struct {
	// clStatus represents the synchronization status of a CL.
	clStatus CLStatus
	// elStatus represents the synchronization status of an EL.
	elStatus ELStatus
}

// Healthy returns true if both CL and EL are synced.
func (s *Status) Healthy() bool {
	return s.clStatus == CLStatusSynced &&
		s.elStatus == ELStatusSynced
}

// CLStatus returns the synchronization status of a CL.
func (s *Status) CLStatus() CLStatus {
	return s.clStatus
}

// ELStatus returns the synchronization status of an EL.
func (s *Status) ELStatus() ELStatus {
	return s.elStatus
}
