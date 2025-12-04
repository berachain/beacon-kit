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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package testhelpers

import (
	"context"

	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/stretchr/testify/mock"
)

// SimpleBlobRequester is a testify mock for BlobRequester interface.
// Used primarily for unit testing blob executor with controlled peer responses.
type SimpleBlobRequester struct {
	mock.Mock
}

func (m *SimpleBlobRequester) RequestBlobs(
	ctx context.Context,
	slot math.Slot,
	verifier func(datypes.BlobSidecars) error,
) ([]*datypes.BlobSidecar, error) {
	args := m.Called(ctx, slot, verifier)
	if args.Get(0) != nil {
		sidecars, ok := args.Get(0).([]*datypes.BlobSidecar)
		if !ok {
			return nil, args.Error(1)
		}
		// Call the verifier if provided (simulates Byzantine verification)
		if verifier != nil {
			if err := verifier(sidecars); err != nil {
				return nil, err // Verifier rejected
			}
		}
		return sidecars, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *SimpleBlobRequester) SetHeadSlot(_ math.Slot) {}
