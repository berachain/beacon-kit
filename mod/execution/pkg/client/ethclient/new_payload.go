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

package ethclient

import (
	"context"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
)

// NewPayload calls the engine_newPayloadV3 method via JSON-RPC.
func (s *Eth1Client[ExecutionPayloadT]) NewPayload(
	ctx context.Context,
	payload ExecutionPayloadT,
	versionedHashes []common.ExecutionHash,
	parentBlockRoot *primitives.Root,
) (*engineprimitives.PayloadStatusV1, error) {
	switch payload.Version() {
	case version.Deneb:
		return s.NewPayloadV3(
			ctx, payload, versionedHashes, parentBlockRoot,
		)
	default:
		return nil, ErrInvalidVersion
	}
}

// newPayload is used to call the underlying JSON-RPC method for newPayload.
func (s *Eth1Client[ExecutionPayloadT]) NewPayloadV3(
	ctx context.Context,
	payload ExecutionPayloadT,
	versionedHashes []common.ExecutionHash,
	parentBlockRoot *primitives.Root,
) (*engineprimitives.PayloadStatusV1, error) {
	result := &engineprimitives.PayloadStatusV1{}
	if err := s.Client.Client().CallContext(
		ctx, result, NewPayloadMethodV3, payload, versionedHashes,
		(*common.ExecutionHash)(parentBlockRoot),
	); err != nil {
		return nil, err
	}
	return result, nil
}
