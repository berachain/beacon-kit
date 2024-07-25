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

package p2p

import (
	"context"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/encoding"
)

// NoopBlobHandler is a gossip handler that simply returns the
// ssz marshalled data as a "reference" to the object it receives.
type NoopBlobHandler[
	BlobT constraints.SSZMarshallable, _ encoding.ABCIRequest,
] struct {
	NoopGossipHandler[BlobT, []byte]
}

func NewNoopBlobHandler[
	BlobT constraints.SSZMarshallable, ReqT encoding.ABCIRequest,
]() NoopBlobHandler[BlobT, ReqT] {
	return NoopBlobHandler[BlobT, ReqT]{
		NoopGossipHandler: NoopGossipHandler[BlobT, []byte]{},
	}
}

// Publish takes a Blob and returns the ssz marshalled data.
func (n NoopBlobHandler[BlobT, _]) Publish(
	_ context.Context,
	data BlobT,
) ([]byte, error) {
	return data.MarshalSSZ()
}

// Request takes an ABCI Request and returns a Blob.
func (n NoopBlobHandler[BlobT, ReqT]) Request(
	_ context.Context,
	req ReqT,
) (BlobT, error) {
	return encoding.UnmarshalBlobSidecarsFromABCIRequest[BlobT](
		req,
		1,
	)
}
