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

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	dastore "github.com/berachain/beacon-kit/da/store"
	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/stretchr/testify/mock"
)

// SimpleBlobProcessor is a testify mock for BlobProcessor interface.
// Used primarily for unit testing blob executor with controlled verification outcomes.
type SimpleBlobProcessor struct {
	mock.Mock
}

func (m *SimpleBlobProcessor) VerifySidecars(
	ctx context.Context,
	sidecars datypes.BlobSidecars,
	blkHeader *ctypes.BeaconBlockHeader,
	kzgCommitments eip4844.KZGCommitments[common.ExecutionHash],
) error {
	return m.Called(ctx, sidecars, blkHeader, kzgCommitments).Error(0)
}

func (m *SimpleBlobProcessor) ProcessSidecars(avs *dastore.Store, sidecars datypes.BlobSidecars) error {
	return m.Called(avs, sidecars).Error(0)
}
