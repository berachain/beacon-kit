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
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package ethclient

import (
	"context"
	"fmt"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/ethereum/go-ethereum/beacon/engine"
)

/* -------------------------------------------------------------------------- */
/*                                 NewPayload                                 */
/* -------------------------------------------------------------------------- */

// NewPayload calls the appropriate version of the Engine API NewPayload method.
func (s *Client) NewPayload(
	ctx context.Context,
	req ctypes.NewPayloadRequest,
) (*engineprimitives.PayloadStatusV1, error) {
	// Versions before Deneb are not supported for calling NewPayload.
	if version.IsBefore(req.GetForkVersion(), version.Deneb()) {
		return nil, ErrInvalidVersion
	}
	forkVersion := req.GetForkVersion()
	if version.Equals(forkVersion, version.Deneb()) || version.Equals(forkVersion, version.Deneb1()) {
		return s.NewPayloadV3(ctx, req.GetExecutionPayload(), req.GetVersionedHashes(), req.GetParentBeaconBlockRoot())
	}
	if version.Equals(forkVersion, version.Electra()) {
		executionRequests, err := req.GetEncodedExecutionRequests()
		if err != nil {
			return nil, err
		}
		return s.NewPayloadV4(
			ctx,
			req.GetExecutionPayload(),
			req.GetVersionedHashes(),
			req.GetParentBeaconBlockRoot(),
			executionRequests,
		)
	}
	return nil, ErrInvalidVersion
}

// NewPayloadV3 calls the engine_newPayloadV3 via JSON-RPC.
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

// NewPayloadV4 calls the engine_newPayloadV4 via JSON-RPC.
func (s *Client) NewPayloadV4(
	ctx context.Context,
	payload *ctypes.ExecutionPayload,
	versionedHashes []common.ExecutionHash,
	parentBlockRoot *common.Root,
	executionRequests []ctypes.EncodedExecutionRequest,
) (*engineprimitives.PayloadStatusV1, error) {
	result := &engineprimitives.PayloadStatusV1{}
	if err := s.Call(
		ctx, result, NewPayloadMethodV4, payload, versionedHashes, parentBlockRoot, executionRequests,
	); err != nil {
		return nil, err
	}
	return result, nil
}

/* -------------------------------------------------------------------------- */
/*                              ForkchoiceUpdated                             */
/* -------------------------------------------------------------------------- */

// ForkchoiceUpdated calls the appropriate version of the Engine API ForkchoiceUpdated method.
func (s *Client) ForkchoiceUpdated(
	ctx context.Context,
	state *engineprimitives.ForkchoiceStateV1,
	attrs any,
	forkVersion common.Version,
) (*engineprimitives.ForkchoiceResponseV1, error) {
	// Versions before Deneb are not supported for calling ForkchoiceUpdated.
	if version.IsBefore(forkVersion, version.Deneb()) {
		return nil, ErrInvalidVersion
	}

	// V3 is used for beacon versions Deneb and onwards.
	return s.ForkchoiceUpdatedV3(ctx, state, attrs)
}

// ForkchoiceUpdatedV3 calls the engine_forkchoiceUpdatedV3 method via JSON-RPC.
func (s *Client) ForkchoiceUpdatedV3(
	ctx context.Context,
	state *engineprimitives.ForkchoiceStateV1,
	attrs any,
) (*engineprimitives.ForkchoiceResponseV1, error) {
	result := &engineprimitives.ForkchoiceResponseV1{}
	if err := s.Call(
		ctx, result, ForkchoiceUpdatedMethodV3, state, attrs,
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

// GetPayload calls the appropriate version of the Engine API GetPayload method.
func (s *Client) GetPayload(
	ctx context.Context,
	payloadID engineprimitives.PayloadID,
	forkVersion common.Version,
) (ctypes.BuiltExecutionPayloadEnv, error) {
	// Versions before Deneb are not supported for calling GetPayload.
	if version.IsBefore(forkVersion, version.Deneb()) {
		return nil, ErrInvalidVersion
	}

	// V3 is used for beacon versions Deneb and onwards.
	return s.GetPayloadV3(ctx, payloadID, forkVersion)
}

// GetPayloadV3 calls the engine_getPayloadV3 method via JSON-RPC.
func (s *Client) GetPayloadV3(
	ctx context.Context,
	payloadID engineprimitives.PayloadID,
	forkVersion common.Version,
) (*ctypes.BuiltExecutionPayloadEnv, error) { // <-- Изменяем возвращаемый тип
	result := ctypes.NewEmptyExecutionPayloadEnvelope[*engineprimitives.BlobsBundleV1](forkVersion)
	if err := s.Call(ctx, result, GetPayloadMethodV3, payloadID); err != nil {
		return nil, fmt.Errorf("failed GetPayloadV3 call: %w", err)
	}
	return result, nil
}

/* -------------------------------------------------------------------------- */
/*                                    Other                                   */
/* -------------------------------------------------------------------------- */

// ExchangeCapabilities calls the engine_exchangeCapabilities method via JSON-RPC.
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
