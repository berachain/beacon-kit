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

package da

// BlobProcessor is the interface for the blobs processor.
type BlobProcessor[
	AvailabilityStoreT,
	ConsensusSidecarsT, BlobSidecarsT any,
] interface {
	// ProcessSidecars processes the blobs and ensures they match the local
	// state.
	ProcessSidecars(
		avs AvailabilityStoreT,
		sidecars BlobSidecarsT,
	) error
	// VerifySidecars verifies the blobs and ensures they match the local state.
	VerifySidecars(sidecars ConsensusSidecarsT) error
}

type ConsensusSidecars[BlobSidecarsT any, BeaconBlockHeaderT any] interface {
	GetSidecars() BlobSidecarsT
	GetHeader() BeaconBlockHeaderT
}

// BlobSidecar is the interface for the blob sidecar.
type BlobSidecar interface {
	// Len returns the length of the sidecar.
	Len() int
	// IsNil checks if the sidecar is nil.
	IsNil() bool
}
