// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package client

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/itsdevbear/bolaris/config/version"
	eth "github.com/itsdevbear/bolaris/engine/client/ethclient"
	enginetypes "github.com/itsdevbear/bolaris/engine/types"
	enginev1 "github.com/itsdevbear/bolaris/engine/types/v1"
	"github.com/itsdevbear/bolaris/types/consensus/primitives"
)

// NewPayload calls the engine_newPayloadVX method via JSON-RPC.
func (s *EngineClient) NewPayload(
	ctx context.Context, payload enginetypes.ExecutionPayload,
	versionedHashes []common.Hash, parentBlockRoot *common.Hash,
) ([]byte, error) {
	dctx, cancel := context.WithTimeout(ctx, s.cfg.RPCTimeout)
	defer cancel()

	// Call the appropriate RPC method based on the payload version.
	result, err := s.callNewPayloadRPC(
		dctx,
		payload,
		versionedHashes,
		parentBlockRoot,
	)
	if err != nil {
		return nil, err
	}

	// This case is only true when the payload is invalid, so
	// `processPayloadStatusResult` below will return an error.
	if validationErr := result.GetValidationError(); validationErr != "" {
		s.logger.Error(
			"Got a validation error in newPayload",
			"err",
			errors.New(validationErr),
		)
	}

	return processPayloadStatusResult(result)
}

// callNewPayloadRPC calls the engine_newPayloadVX method via JSON-RPC.
func (s *EngineClient) callNewPayloadRPC(
	ctx context.Context, payload enginetypes.ExecutionPayload,
	versionedHashes []common.Hash, parentBlockRoot *common.Hash,
) (*enginev1.PayloadStatus, error) {
	switch payloadPb := payload.ToProto().(type) {
	case *enginev1.ExecutionPayloadDeneb:
		return s.NewPayloadV3(ctx, payloadPb, versionedHashes, parentBlockRoot)
	default:
		return nil, ErrInvalidPayloadType
	}
}

// ForkchoiceUpdated calls the engine_forkchoiceUpdatedV1 method via JSON-RPC.
func (s *EngineClient) ForkchoiceUpdated(
	ctx context.Context,
	state *enginev1.ForkchoiceState,
	attrs enginetypes.PayloadAttributer,
) (*enginev1.PayloadIDBytes, []byte, error) {
	dctx, cancel := context.WithTimeout(ctx, s.cfg.RPCTimeout)
	defer cancel()

	if attrs == nil {
		return nil, nil, ErrNilAttributesPassedToClient
	}

	result, err := s.callUpdatedForkchoiceRPC(dctx, state, attrs)
	if err != nil {
		return nil, nil, s.handleRPCError(err)
	}

	lastestValidHash, err := processPayloadStatusResult(result.Status)
	if err != nil {
		return nil, lastestValidHash, err
	}
	return result.PayloadID, lastestValidHash, nil
}

// updateForkChoiceByVersion calls the engine_forkchoiceUpdatedVX method via
// JSON-RPC.
func (s *EngineClient) callUpdatedForkchoiceRPC(
	ctx context.Context,
	state *enginev1.ForkchoiceState,
	attrs enginetypes.PayloadAttributer,
) (*eth.ForkchoiceUpdatedResponse, error) {
	switch v := attrs.ToProto().(type) {
	case *enginev1.PayloadAttributesV3:
		return s.ForkchoiceUpdatedV3(ctx, state, v)
	default:
		return nil, ErrInvalidPayloadAttributes
	}
}

// GetPayload calls the engine_getPayloadVX method via JSON-RPC. It returns
// the execution data as well as the blobs bundle.
func (s *EngineClient) GetPayload(
	ctx context.Context, payloadID primitives.PayloadID, slot primitives.Slot,
) (enginetypes.ExecutionPayload, *enginev1.BlobsBundle, bool, error) {
	dctx, cancel := context.WithTimeout(ctx, s.cfg.RPCTimeout)
	defer cancel()

	var fn func(
		context.Context, enginev1.PayloadIDBytes,
	) (*enginev1.ExecutionPayloadContainer, error)
	switch s.beaconCfg.ActiveForkVersion(primitives.Epoch(slot)) {
	case version.Deneb:
		fn = s.GetPayloadV3
	default:
		return nil, nil, false, ErrInvalidGetPayloadVersion
	}

	result, err := fn(dctx, enginev1.PayloadIDBytes(payloadID))
	if err != nil {
		return nil, nil, false, s.handleRPCError(err)
	}

	return result, result.GetBlobsBundle(), result.GetShouldOverrideBuilder(), nil
}

// ExchangeCapabilities calls the engine_exchangeCapabilities method via
// JSON-RPC.
func (s *EngineClient) ExchangeCapabilities(
	ctx context.Context,
) ([]string, error) {
	result, err := s.Eth1Client.ExchangeCapabilities(
		ctx, eth.BeaconKitSupportedCapabilities(),
	)
	if err != nil {
		s.statusErrMu.Lock()
		defer s.statusErrMu.Unlock()
		//#nosec:G703 wtf is even this problem here.
		s.statusErr = s.handleRPCError(err)
		return nil, s.statusErr
	}

	// Capture and log the capabilities that the execution client has.
	for _, capability := range result {
		s.logger.Info("exchanged capability", "capability", capability)
		s.capabilities[capability] = struct{}{}
	}

	// Log the capabilities that the execution client does not have.
	for _, capability := range eth.BeaconKitSupportedCapabilities() {
		if _, exists := s.capabilities[capability]; !exists {
			s.logger.Warn(
				"your execution client may require an update 🚸",
				"unsupported_capability", capability,
			)
		}
	}

	s.statusErr = nil
	return result, nil
}
