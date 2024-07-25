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
	"errors"
	"reflect"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
)

// NoopGossipHandlerDyn is a gossip handler that simply returns the
// ssz marshalled data as a "reference" to the object it receives.
type NoopGossipHandlerDyn[
	DataT constraints.SSZMarshallable, BytesT ~[]byte,
] struct{}

// Publish creates a new NoopGossipHandlerDyn.
func (n NoopGossipHandlerDyn[DataT, BytesT]) Publish(
	_ context.Context,
	data DataT,
) (BytesT, error) {
	return data.MarshalSSZ()
}

// Request simply returns the reference it receives.
func (n NoopGossipHandlerDyn[DataT, BytesT]) Request(
	_ context.Context,
	ref BytesT,
) (DataT, error) {
	var (
		out DataT
		ok  bool
	)

	// Alloc memory if DataT is a pointer.
	if reflect.ValueOf(&out).Elem().Kind() == reflect.Ptr {
		newInstance := reflect.New(reflect.TypeOf(out).Elem())
		out, ok = newInstance.Interface().(DataT)
		if !ok {
			return out, errors.New("failed to create new instance")
		}
	}

	return out, out.UnmarshalSSZ(ref)
}
