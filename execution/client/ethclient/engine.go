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

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/ethereum/go-ethereum/beacon/engine"
)

/* -------------------------------------------------------------------------- */
/*                                 NewPayload                                 */
/* -------------------------------------------------------------------------- */

// NewPayload calls the engine_newPayloadV3 method via JSON-RPC.
func (s *Client) NewPayload(
	ctx context.Context,
	payload *ctypes.ExecutionPayload,
	versionedHashes []common.ExecutionHash,
	parentBlockRoot *common.Root,
) (*engineprimitives.PayloadStatusV1, error) {
	if payload.Version() < version.Deneb {
		return nil, ErrInvalidVersion
	}

	return s.NewPayloadV3(
		ctx, payload, versionedHashes, parentBlockRoot,
	)
}

// NewPayloadV3 is used to call the underlying JSON-RPC method for newPayload.
func (s *Client) NewPayloadV3(
	ctx context.Context,
	payload *ctypes.ExecutionPayload,
	versionedHashes []common.ExecutionHash,
	parentBlockRoot *common.Root,
) (*engineprimitives.PayloadStatusV1, error) {
	result := &engineprimitives.PayloadStatusV1{}
	if err := s.Call(
		ctx, result, NewPayloadMethodV3, payload, versionedHashes, parentBlockRoot,
	); err != nil {
		return nil, err
	}
	return result, nil
}

/* -------------------------------------------------------------------------- */
/*                              ForkchoiceUpdated                             */
/* -------------------------------------------------------------------------- */

// ForkchoiceUpdated is a helper function to call the appropriate version of
// the.
func (s *Client) ForkchoiceUpdated(
	ctx context.Context,
	state *engineprimitives.ForkchoiceStateV1,
	attrs any,
	forkVersion uint32,
) (*engineprimitives.ForkchoiceResponseV1, error) {
	if forkVersion < version.Deneb {
		return nil, ErrInvalidVersion
	}

	return s.ForkchoiceUpdatedV3(ctx, state, attrs)
}

// ForkchoiceUpdatedV3 calls the engine_forkchoiceUpdatedV3 method via JSON-RPC.
func (s *Client) ForkchoiceUpdatedV3(
	ctx context.Context,
	state *engineprimitives.ForkchoiceStateV1,
	attrs any,
) (*engineprimitives.ForkchoiceResponseV1, error) {
	return s.forkchoiceUpdated(ctx, ForkchoiceUpdatedMethodV3, state, attrs)
}

// forkchoiceUpdateCall is a helper function to call to any version
// of the forkchoiceUpdates method.
func (s *Client) forkchoiceUpdated(
	ctx context.Context,
	method string,
	state *engineprimitives.ForkchoiceStateV1,
	attrs any,
) (*engineprimitives.ForkchoiceResponseV1, error) {
	result := &engineprimitives.ForkchoiceResponseV1{}

	if err := s.Call(
		ctx, result, method, state, attrs,
	); err != nil {
		return nil, err
	}

	if (result.PayloadStatus == engineprimitives.PayloadStatusV1{}) {
		return nil, ErrNilResponse
	}

	return result, nil
}

/* -------------------------------------------------------------------------- */
/*                                 GetPayload                                 */
/* -------------------------------------------------------------------------- */

// GetPayload is a helper function to call the appropriate version of the
// engine_getPayload method.
func (s *Client) GetPayload(
	ctx context.Context,
	payloadID engineprimitives.PayloadID,
	forkVersion uint32,
) (ctypes.BuiltExecutionPayloadEnv, error) {
	if forkVersion < version.Deneb {
		return nil, ErrInvalidVersion
	}

	return s.GetPayloadV3(ctx, payloadID)
}

// GetPayloadV3 calls the engine_getPayloadV3 method via JSON-RPC.
func (s *Client) GetPayloadV3(
	ctx context.Context, payloadID engineprimitives.PayloadID,
) (ctypes.BuiltExecutionPayloadEnv, error) {
	var t *ctypes.ExecutionPayload
	result := &ctypes.ExecutionPayloadEnvelope[*engineprimitives.BlobsBundleV1[
		eip4844.KZGCommitment, eip4844.KZGProof, eip4844.Blob,
	]]{
		ExecutionPayload: t.Empty(version.Deneb),
	}

	if err := s.Call(
		ctx, result, GetPayloadMethodV3, payloadID,
	); err != nil {
		return nil, err
	}
	return result, nil
}

/* -------------------------------------------------------------------------- */
/*                                    Other                                   */
/* -------------------------------------------------------------------------- */

// ExchangeCapabilities calls the engine_exchangeCapabilities method via
// JSON-RPC.
func (s *Client) ExchangeCapabilities(
	ctx context.Context,
	capabilities []string,
) ([]string, error) {
	result := make([]string, 0)
	if err := s.Call(
		ctx, &result, ExchangeCapabilities, &capabilities,
	); err != nil {
		return nil, err
	}
	return result, nil
}

// GetClientVersionV1 calls the engine_getClientVersionV1 method via JSON-RPC.
func (s *Client) GetClientVersionV1(
	ctx context.Context,
) ([]engineprimitives.ClientVersionV1, error) {
	result := make([]engineprimitives.ClientVersionV1, 0)

	// NOTE: although the ethereum spec does not require us passing a
	// clientversion as param, it seems some clients require it and even
	// enforce a valid Code.
	if err := s.Call(
		ctx, &result, GetClientVersionV1, engine.ClientVersionV1{Code: "GE"},
	); err != nil {
		return nil, err
	}
	return result, nil
}
