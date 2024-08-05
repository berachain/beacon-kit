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
)

// NoopGossipHandler is a gossip handler that simply returns the
// ssz marshalled data as a "reference" to the object it receives.
type NoopGossipHandler[
	DataT interface {
		constraints.Empty[DataT]
		constraints.SSZMarshallable
	}, BytesT ~[]byte,
] struct{}

// Publish creates a new NoopGossipHandler.
func (n NoopGossipHandler[DataT, BytesT]) Publish(
	_ context.Context,
	data DataT,
) (BytesT, error) {
	return data.MarshalSSZ()
}

// Request simply returns the reference it receives.
func (n NoopGossipHandler[DataT, BytesT]) Request(
	_ context.Context,
	ref BytesT,
) (DataT, error) {
	var (
		out DataT
	)

	// Use Empty() method to create a new instance of DataT
	out = out.Empty()
	return out, out.UnmarshalSSZ(ref)
}
