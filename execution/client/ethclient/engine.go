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
	"github.com/berachain/beacon-kit/primitives/crypto"
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
	forkVersion := req.GetForkVersion()
	switch {
	case version.IsBefore(forkVersion, version.Deneb()):
		// Versions before Deneb are not supported for calling NewPayload.
		return nil, ErrInvalidVersion

	case version.Equals(forkVersion, version.Deneb()), version.Equals(forkVersion, version.Deneb1()):
		// Use V3 for Deneb versions (Deneb and Deneb1).
		return s.NewPayloadV3(
			ctx,
			req.GetExecutionPayload(),
			req.GetVersionedHashes(),
			req.GetParentBeaconBlockRoot(),
		)

	case version.Equals(forkVersion, version.Electra()):
		// Use V4 for Electra versions.
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

	case version.Equals(forkVersion, version.Electra1()):
		// Use V4P11 for Electra1 versions.
		executionRequests, err := req.GetEncodedExecutionRequests()
		if err != nil {
			return nil, err
		}
		return s.NewPayloadV4P11(
			ctx,
			req.GetExecutionPayload(),
			req.GetVersionedHashes(),
			req.GetParentBeaconBlockRoot(),
			executionRequests,
			req.GetParentProposerPubkey(),
		)

	default:
		return nil, ErrInvalidVersion
	}
}

// NewPayloadV3 calls the engine_newPayloadV3 via JSON-RPC.
func (s *Client) NewPayloadV3(
	ctx context.Context,
	payload *ctypes.ExecutionPayload,
	versionedHashes []common.ExecutionHash,
	parentBlockRoot common.Root,
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
	parentBlockRoot common.Root,
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

// NewPayloadV4P11 calls the engine_newPayloadV4P11 via JSON-RPC.
func (s *Client) NewPayloadV4P11(
	ctx context.Context,
	payload *ctypes.ExecutionPayload,
	versionedHashes []common.ExecutionHash,
	parentBlockRoot common.Root,
	executionRequests []ctypes.EncodedExecutionRequest,
	parentProposerPubKey *crypto.BLSPubkey,
) (*engineprimitives.PayloadStatusV1, error) {
	result := &engineprimitives.PayloadStatusV1{}
	if err := s.Call(
		ctx, result, NewPayloadMethodV4P11, payload, versionedHashes, parentBlockRoot, executionRequests, parentProposerPubKey,
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
	// V3 is used for beacon versions Deneb and onwards.
	switch {
	case version.IsBefore(forkVersion, version.Deneb()):
		// Versions before Deneb are not supported for calling ForkchoiceUpdated.
		return nil, ErrInvalidVersion

	case version.Equals(forkVersion, version.Deneb()),
		version.Equals(forkVersion, version.Deneb1()),
		version.Equals(forkVersion, version.Electra()):
		// Deneb versions and Electra use ForkchoiceUpdatedV3.
		return s.ForkchoiceUpdatedV3(ctx, state, attrs)

	case version.Equals(forkVersion, version.Electra1()):
		// Electra1 uses ForkchoiceUpdatedV3P11.
		return s.ForkchoiceUpdatedV3P11(ctx, state, attrs)

	default:
		return nil, ErrInvalidVersion
	}
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

// ForkchoiceUpdatedV3P11 calls the engine_forkchoiceUpdatedV3P11 method via JSON-RPC.
func (s *Client) ForkchoiceUpdatedV3P11(
	ctx context.Context,
	state *engineprimitives.ForkchoiceStateV1,
	attrs any,
) (*engineprimitives.ForkchoiceResponseV1, error) {
	result := &engineprimitives.ForkchoiceResponseV1{}
	if err := s.Call(
		ctx, result, ForkchoiceUpdatedMethodV3P11, state, attrs,
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
	switch {
	case version.IsBefore(forkVersion, version.Deneb()):
		// Versions before Deneb are not supported for calling GetPayload.
		return nil, ErrInvalidVersion

	case version.Equals(forkVersion, version.Deneb()), version.Equals(forkVersion, version.Deneb1()):
		return s.GetPayloadV3(ctx, payloadID, forkVersion)

	case version.Equals(forkVersion, version.Electra()):
		return s.GetPayloadV4(ctx, payloadID, forkVersion)

	case version.Equals(forkVersion, version.Electra1()):
		return s.GetPayloadV4P11(ctx, payloadID, forkVersion)

	default:
		return nil, ErrInvalidVersion
	}
}

// GetPayloadV3 calls the engine_getPayloadV3 method via JSON-RPC.
func (s *Client) GetPayloadV3(
	ctx context.Context,
	payloadID engineprimitives.PayloadID,
	forkVersion common.Version,
) (ctypes.BuiltExecutionPayloadEnv, error) {
	result := ctypes.NewEmptyExecutionPayloadEnvelope[*engineprimitives.BlobsBundleV1](forkVersion)
	if err := s.Call(ctx, result, GetPayloadMethodV3, payloadID); err != nil {
		return nil, fmt.Errorf("failed GetPayloadV3 call: %w", err)
	}
	return result, nil
}

// GetPayloadV4 calls the engine_getPayloadV4 method via JSON-RPC.
func (s *Client) GetPayloadV4(
	ctx context.Context,
	payloadID engineprimitives.PayloadID,
	forkVersion common.Version,
) (ctypes.BuiltExecutionPayloadEnv, error) {
	result := ctypes.NewEmptyExecutionPayloadEnvelope[*engineprimitives.BlobsBundleV1](forkVersion)
	if err := s.Call(ctx, result, GetPayloadMethodV4, payloadID); err != nil {
		return nil, fmt.Errorf("failed GetPayloadV4 call: %w", err)
	}
	return result, nil
}

// GetPayloadV4P11 calls the engine_getPayloadV4P11 method via JSON-RPC.
func (s *Client) GetPayloadV4P11(
	ctx context.Context,
	payloadID engineprimitives.PayloadID,
	forkVersion common.Version,
) (ctypes.BuiltExecutionPayloadEnv, error) {
	result := ctypes.NewEmptyExecutionPayloadEnvelope[*engineprimitives.BlobsBundleV1](forkVersion)
	if err := s.Call(ctx, result, GetPayloadMethodV4P11, payloadID); err != nil {
		return nil, fmt.Errorf("failed GetPayloadV4P11 call: %w", err)
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
