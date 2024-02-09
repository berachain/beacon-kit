// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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

package engine

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/log"
	"github.com/itsdevbear/bolaris/config"
	"github.com/itsdevbear/bolaris/third_party/go-ethereum/common"
	enginev1 "github.com/itsdevbear/bolaris/third_party/prysm/proto/engine/v1"
	"github.com/itsdevbear/bolaris/types/consensus/blocks/blocks"
	"github.com/itsdevbear/bolaris/types/consensus/interfaces"
	"github.com/itsdevbear/bolaris/types/consensus/primitives"
	"github.com/itsdevbear/bolaris/types/consensus/version"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/execution"
	payloadattribute "github.com/prysmaticlabs/prysm/v4/consensus-types/payload-attribute"

	eth "github.com/itsdevbear/bolaris/engine/ethclient"
)

// Caller is implemented by engineClient.
var _ Caller = (*engineClient)(nil)

// engineClient is a struct that holds a pointer to an Eth1Client.
type engineClient struct {
	*eth.Eth1Client
	engineTimeout time.Duration
	beaconCfg     *config.Beacon
	logger        log.Logger
}

// NewClient creates a new engine client engineClient.
// It takes an Eth1Client as an argument and returns a pointer to an engineClient.
func NewClient(opts ...Option) Caller {
	ec := &engineClient{}
	for _, opt := range opts {
		if err := opt(ec); err != nil {
			panic(err)
		}
	}

	return ec
}

// NewPayload calls the engine_newPayloadVX method via JSON-RPC.
func (s *engineClient) NewPayload(
	ctx context.Context, payload interfaces.ExecutionData,
	versionedHashes []common.Hash, parentBlockRoot *common.Hash,
) ([]byte, error) {
	var (
		d            = time.Now().Add(s.engineTimeout)
		dctx, cancel = context.WithDeadline(ctx, d)
		err          error
		result       *enginev1.PayloadStatus
	)
	defer cancel()
	switch payload.Proto().(type) {
	case *enginev1.ExecutionPayloadCapella:
		payloadPb, ok := payload.Proto().(*enginev1.ExecutionPayloadCapella)
		if !ok {
			return nil, errors.New("execution data must be a Capella execution payload")
		}
		result, err = s.NewPayloadV2(dctx, payloadPb)
	case *enginev1.ExecutionPayloadDeneb:
		payloadPb, ok := payload.Proto().(*enginev1.ExecutionPayloadDeneb)
		if !ok {
			return nil, errors.New("execution data must be a Deneb execution payload")
		}
		result, err = s.NewPayloadV3(dctx, payloadPb, versionedHashes, parentBlockRoot)
	default:
		return nil, errors.New("unknown execution data type")
	}
	if err != nil {
		return nil, s.handleRPCError(err)
	}

	if result.GetValidationError() != "" {
		s.logger.Error("Got a validation error in newPayload", "err",
			errors.New(result.GetValidationError()))
	}

	return processPayloadStatusResult(result)
}

// ForkchoiceUpdated calls the engine_forkchoiceUpdatedV1 method via JSON-RPC.
func (s *engineClient) ForkchoiceUpdated(
	ctx context.Context, state *enginev1.ForkchoiceState, attrs payloadattribute.Attributer,
) (*enginev1.PayloadIDBytes, []byte, error) {
	d := time.Now().Add(s.engineTimeout)
	ctx, cancel := context.WithDeadline(ctx, d)
	defer cancel()
	result := &execution.ForkchoiceUpdatedResponse{}
	if attrs == nil {
		return nil, nil, errors.New("nil payload attributer")
	}
	switch attrs.Version() {
	case version.Deneb:
		a, err := attrs.PbV3()
		if err != nil {
			return nil, nil, err
		}
		err = s.RawClient().CallContext(ctx, result,
			execution.ForkchoiceUpdatedMethodV3, state, a)
		if err != nil {
			return nil, nil, s.handleRPCError(err)
		}
	case version.Capella:
		a, err := attrs.PbV2()
		if err != nil {
			return nil, nil, err
		}
		err = s.RawClient().CallContext(ctx, result,
			execution.ForkchoiceUpdatedMethodV2, state, a)
		if err != nil {
			return nil, nil, s.handleRPCError(err)
		}
	default:
		return nil, nil, fmt.Errorf("unknown payload attribute version: %v", attrs.Version())
	}

	if result.Status == nil {
		return nil, nil, execution.ErrNilResponse
	}
	if result.ValidationError != "" {
		s.logger.Error("Got validation error in forkChoiceUpdated", "err",
			errors.New(result.ValidationError))
	}
	lastestValidHash, err := processPayloadStatusResult(result.Status)
	if err != nil {
		return nil, lastestValidHash, err
	}
	return result.PayloadId, lastestValidHash, nil
}

// GetPayload calls the engine_getPayloadVX method via JSON-RPC.
// It returns the execution data as well as the blobs bundle.
func (s *engineClient) GetPayload(
	ctx context.Context, payloadID [8]byte, slot primitives.Slot,
) (interfaces.ExecutionData, *enginev1.BlobsBundle, bool, error) {
	d := time.Now().Add(s.engineTimeout)
	ctx, cancel := context.WithDeadline(ctx, d)
	defer cancel()
	if primitives.Epoch(slot) >= s.beaconCfg.Forks.DenebForkEpoch {
		result := &enginev1.ExecutionPayloadDenebWithValueAndBlobsBundle{}

		if err := s.RawClient().CallContext(ctx,
			result, execution.GetPayloadMethodV3, enginev1.PayloadIDBytes(payloadID),
		); err != nil {
			return nil, nil, false, s.handleRPCError(err)
		}

		ed, err := blocks.WrappedExecutionPayloadDeneb(result.GetPayload(),
			blocks.PayloadValueToWei(result.GetValue()))
		if err != nil {
			return nil, nil, false, err
		}

		return ed, result.GetBlobsBundle(), result.GetShouldOverrideBuilder(), nil
	}

	result := &enginev1.ExecutionPayloadCapellaWithValue{}
	if err := s.RawClient().CallContext(ctx,
		result, execution.GetPayloadMethodV2, enginev1.PayloadIDBytes(payloadID),
	); err != nil {
		return nil, nil, false, s.handleRPCError(err)
	}

	ed, err := blocks.WrappedExecutionPayloadCapella(result.GetPayload(),
		blocks.PayloadValueToWei(result.GetValue()))

	if err != nil {
		return nil, nil, false, err
	}
	return ed, nil, false, nil
}

// ExecutionBlockByHash fetches an execution engine block by hash by calling
// eth_blockByHash via JSON-RPC.
func (s *engineClient) ExecutionBlockByHash(ctx context.Context, hash common.Hash, withTxs bool,
) (*enginev1.ExecutionBlock, error) {
	result := &enginev1.ExecutionBlock{}
	err := s.Eth1Client.Client.Client().CallContext(
		ctx, result, "eth_getBlockByHash", hash, withTxs)
	return result, s.handleRPCError(err)
}
