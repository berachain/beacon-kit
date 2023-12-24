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
	"bytes"
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/execution"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/execution/types"
	"github.com/prysmaticlabs/prysm/v4/config/features"
	fieldparams "github.com/prysmaticlabs/prysm/v4/config/fieldparams"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/blocks"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/interfaces"
	payloadattribute "github.com/prysmaticlabs/prysm/v4/consensus-types/payload-attribute"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	pb "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
	"github.com/prysmaticlabs/prysm/v4/runtime/version"
	"github.com/prysmaticlabs/prysm/v4/time/slots"
	"go.opencensus.io/trace"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"cosmossdk.io/log"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethcoretypes "github.com/ethereum/go-ethereum/core/types"
	gethRPC "github.com/ethereum/go-ethereum/rpc"

	eth "github.com/itsdevbear/bolaris/beacon/execution/engine/ethclient"
	"github.com/itsdevbear/bolaris/types/config"
)

const (
	// Defines the seconds before timing out engine endpoints with non-block execution semantics.
	defaultEngineTimeout = time.Second
)

// zeroHash32 is a zeroed 32-byte array.
var zeroHash32 [32]byte

// Caller defines a client that can interact with an Ethereum
// execution node's engine engineCaller via JSON-RPC.
type Caller interface {
	NewPayload(ctx context.Context, payload interfaces.ExecutionData,
		versionedHashes []common.Hash, parentBlockRoot *common.Hash) ([]byte, error)
	ForkchoiceUpdated(
		ctx context.Context, state *pb.ForkchoiceState, attrs payloadattribute.Attributer,
	) (*pb.PayloadIDBytes, []byte, error)
	GetPayload(ctx context.Context, payloadID [8]byte,
		slot primitives.Slot) (interfaces.ExecutionData, *pb.BlobsBundle, bool, error)
	ExecutionBlockByHash(ctx context.Context, hash common.Hash,
		withTxs bool) (*pb.ExecutionBlock, error)

	// TODO: THESE NEED OT BE REMOVED
	LatestSafeBlock(ctx context.Context) (*pb.ExecutionBlock, error)
	LatestFinalizedBlock(ctx context.Context) (*pb.ExecutionBlock, error)
	LatestExecutionBlock(ctx context.Context) (*pb.ExecutionBlock, error)
	SyncProgress(ctx context.Context) (*ethereum.SyncProgress, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*gethcoretypes.Block, error)
	BlockByHash(ctx context.Context, hash common.Hash) (*gethcoretypes.Block, error)
}

// Caller is implemented by engineCaller.
var _ Caller = (*engineCaller)(nil)

// engineCaller is a struct that holds a pointer to an Eth1Client.
type engineCaller struct {
	*eth.Eth1Client
	beaconCfg *config.Beacon
	logger    log.Logger
}

// NewCaller creates a new engine client engineCaller.
// It takes an Eth1Client as an argument and returns a pointer to an engineCaller.
func NewCaller(ethclient *eth.Eth1Client, opts ...Option) Caller {
	ec := &engineCaller{
		Eth1Client: ethclient,
	}

	for _, opt := range opts {
		if err := opt(ec); err != nil {
			panic(err)
		}
	}

	return ec
}

// NewPayload calls the engine_newPayloadVX method via JSON-RPC.
func (s *engineCaller) NewPayload(
	ctx context.Context, payload interfaces.ExecutionData,
	versionedHashes []common.Hash, parentBlockRoot *common.Hash,
) ([]byte, error) {
	// d := time.Now().Add	(
	// 	time.Duration(
	// 		s.beaconCfg.ExecutionEngineTimeoutValue,
	// 	) * time.Second)
	// ctx, cancel := context.WithDeadline(ctx, d)
	// defer cancel()
	result := &pb.PayloadStatus{}

	switch payload.Proto().(type) {
	case *pb.ExecutionPayload:
		payloadPb, ok := payload.Proto().(*pb.ExecutionPayload)
		if !ok {
			return nil, errors.New("execution data must be a Bellatrix or Capella execution payload")
		}
		err := s.Eth1Client.Client.Client().CallContext(ctx, result,
			execution.NewPayloadMethod, payloadPb)
		if err != nil {
			return nil, s.handleRPCError(err)
		}
	case *pb.ExecutionPayloadCapella:
		payloadPb, ok := payload.Proto().(*pb.ExecutionPayloadCapella)
		if !ok {
			return nil, errors.New("execution data must be a Capella execution payload")
		}
		err := s.Eth1Client.Client.Client().CallContext(ctx, result,
			execution.NewPayloadMethodV2, payloadPb)
		if err != nil {
			return nil, s.handleRPCError(err)
		}
	case *pb.ExecutionPayloadDeneb:
		payloadPb, ok := payload.Proto().(*pb.ExecutionPayloadDeneb)
		if !ok {
			return nil, errors.New("execution data must be a Deneb execution payload")
		}
		err := s.Eth1Client.Client.Client().CallContext(ctx,
			result, execution.NewPayloadMethodV3, payloadPb, versionedHashes, parentBlockRoot,
		)
		if err != nil {
			return nil, s.handleRPCError(err)
		}
	default:
		return nil, errors.New("unknown execution data type")
	}

	if result.ValidationError != "" {
		s.logger.Error("Got a validation error in newPayload", "err",
			errors.New(result.ValidationError))
	}

	switch result.Status {
	case pb.PayloadStatus_INVALID_BLOCK_HASH:
		return nil, execution.ErrInvalidBlockHashPayloadStatus
	case pb.PayloadStatus_ACCEPTED:
		return nil, ErrAcceptedPayloadStatus
	case pb.PayloadStatus_SYNCING:
		return nil, ErrSyncingPayloadStatus
	case pb.PayloadStatus_INVALID:
		return result.LatestValidHash, execution.ErrInvalidPayloadStatus
	case pb.PayloadStatus_VALID:
		return result.LatestValidHash, nil
	case pb.PayloadStatus_UNKNOWN:
		return nil, execution.ErrUnknownPayloadStatus
	default:
		return nil, execution.ErrUnknownPayloadStatus
	}
}

// ForkchoiceUpdated calls the engine_forkchoiceUpdatedV1 method via JSON-RPC.
func (s *engineCaller) ForkchoiceUpdated(
	ctx context.Context, state *pb.ForkchoiceState, attrs payloadattribute.Attributer,
) (*pb.PayloadIDBytes, []byte, error) {
	// d := time.Now().Add(time.Duration(s.beaconCfg.ExecutionEngineTimeoutValue) * time.Second)
	// ctx, cancel := context.WithDeadline(ctx, d)
	// defer cancel()
	result := &execution.ForkchoiceUpdatedResponse{}
	if attrs == nil {
		return nil, nil, errors.New("nil payload attributer")
	}
	switch attrs.Version() {
	case version.Bellatrix:
		a, err := attrs.PbV1()
		if err != nil {
			return nil, nil, err
		}
		err = s.Eth1Client.Client.Client().CallContext(ctx, result,
			execution.ForkchoiceUpdatedMethod, state, a)
		if err != nil {
			return nil, nil, s.handleRPCError(err)
		}
	case version.Capella:
		a, err := attrs.PbV2()
		if err != nil {
			return nil, nil, err
		}
		err = s.Eth1Client.Client.Client().CallContext(ctx, result,
			execution.ForkchoiceUpdatedMethodV2, state, a)
		if err != nil {
			return nil, nil, s.handleRPCError(err)
		}
	case version.Deneb:
		a, err := attrs.PbV3()
		if err != nil {
			return nil, nil, err
		}
		err = s.Eth1Client.Client.Client().CallContext(ctx, result,
			execution.ForkchoiceUpdatedMethodV3, state, a)
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
	resp := result.Status
	switch resp.Status {
	case pb.PayloadStatus_ACCEPTED:
		return nil, nil, ErrAcceptedPayloadStatus
	case pb.PayloadStatus_SYNCING:
		return nil, nil, ErrSyncingPayloadStatus
	case pb.PayloadStatus_INVALID:
		return nil, resp.LatestValidHash, execution.ErrInvalidPayloadStatus
	case pb.PayloadStatus_VALID:
		return result.PayloadId, resp.LatestValidHash, nil
	case pb.PayloadStatus_UNKNOWN:
		return nil, nil, execution.ErrUnknownPayloadStatus
	case pb.PayloadStatus_INVALID_BLOCK_HASH:
		return nil, nil, execution.ErrInvalidBlockHashPayloadStatus
	}
	return nil, nil, execution.ErrUnknownPayloadStatus
}

// GetPayload calls the engine_getPayloadVX method via JSON-RPC.
// It returns the execution data as well as the blobs bundle.
func (s *engineCaller) GetPayload(ctx context.Context, payloadID [8]byte,
	slot primitives.Slot) (interfaces.ExecutionData, *pb.BlobsBundle, bool, error) {
	d := time.Now().Add(defaultEngineTimeout)
	ctx, cancel := context.WithDeadline(ctx, d)
	defer cancel()
	if slots.ToEpoch(slot) >= s.beaconCfg.DenebForkEpoch {
		result := &pb.ExecutionPayloadDenebWithValueAndBlobsBundle{}
		err := s.Eth1Client.Client.Client().CallContext(ctx,
			result, execution.GetPayloadMethodV3, pb.PayloadIDBytes(payloadID))
		if err != nil {
			return nil, nil, false, s.handleRPCError(err)
		}
		ed, err := blocks.WrappedExecutionPayloadDeneb(result.Payload,
			blocks.PayloadValueToGwei(result.Value))
		if err != nil {
			return nil, nil, false, err
		}
		return ed, result.BlobsBundle, result.ShouldOverrideBuilder, nil
	}

	if slots.ToEpoch(slot) >= s.beaconCfg.CapellaForkEpoch {
		result := &pb.ExecutionPayloadCapellaWithValue{}
		err := s.Eth1Client.Client.Client().CallContext(ctx,
			result, execution.GetPayloadMethodV2, pb.PayloadIDBytes(payloadID))
		if err != nil {
			return nil, nil, false, s.handleRPCError(err)
		}
		ed, err := blocks.WrappedExecutionPayloadCapella(result.Payload,
			blocks.PayloadValueToGwei(result.Value))
		if err != nil {
			return nil, nil, false, err
		}
		return ed, nil, false, nil
	}

	result := &pb.ExecutionPayload{}
	err := s.Eth1Client.Client.Client().CallContext(ctx,
		result, execution.GetPayloadMethod, pb.PayloadIDBytes(payloadID))
	if err != nil {
		return nil, nil, false, s.handleRPCError(err)
	}
	ed, err := blocks.WrappedExecutionPayload(result)
	if err != nil {
		return nil, nil, false, err
	}
	return ed, nil, false, nil
}

// LatestExecutionBlock fetches the latest execution engine block by calling
// eth_blockByNumber via JSON-RPC.
func (s *engineCaller) LatestExecutionBlock(ctx context.Context) (*pb.ExecutionBlock, error) {
	ctx, span := trace.StartSpan(ctx, "powchain.engine-api-client.LatestExecutionBlock")
	defer span.End()

	result := &pb.ExecutionBlock{}
	err := s.Eth1Client.Client.Client().CallContext(
		ctx,
		result,
		execution.ExecutionBlockByNumberMethod,
		"latest",
		false, /* no full transaction objects */
	)
	return result, s.handleRPCError(err)
}

// LatestExecutionBlock fetches the latest execution engine block by calling
// eth_blockByNumber via JSON-RPC.
func (s *engineCaller) LatestSafeBlock(ctx context.Context) (*pb.ExecutionBlock, error) {
	ctx, span := trace.StartSpan(ctx, "powchain.engine-api-client.LatestExecutionBlock")
	defer span.End()

	result := &pb.ExecutionBlock{}
	err := s.Eth1Client.Client.Client().CallContext(
		ctx,
		result,
		execution.ExecutionBlockByNumberMethod,
		"safe",
		false, /* no full transaction objects */
	)
	return result, s.handleRPCError(err)
}

// LatestExecutionBlock fetches the latest execution engine block by calling
// eth_blockByNumber via JSON-RPC.
func (s *engineCaller) EarliestBlock(ctx context.Context) (*pb.ExecutionBlock, error) {
	// ctx, span := trace.StartSpan(ctx, "powchain.engine-api-client.LatestExecutionBlock")
	// defer span.End()

	result := &pb.ExecutionBlock{}
	err := s.Eth1Client.Client.Client().CallContext(
		ctx,
		result,
		execution.ExecutionBlockByNumberMethod,
		"earliest",
		false, /* no full transaction objects */
	)
	return result, s.handleRPCError(err)
}

// LatestExecutionBlock fetches the latest execution engine block by calling
// eth_blockByNumber via JSON-RPC.
func (s *engineCaller) LatestFinalizedBlock(ctx context.Context) (*pb.ExecutionBlock, error) {
	ctx, span := trace.StartSpan(ctx, "powchain.engine-api-client.LatestExecutionBlock")
	defer span.End()

	result := &pb.ExecutionBlock{}
	err := s.Eth1Client.Client.Client().CallContext(
		ctx,
		result,
		execution.ExecutionBlockByNumberMethod,
		"finalized",
		false, /* no full transaction objects */
	)
	return result, s.handleRPCError(err)
}

// ExecutionBlockByHash fetches an execution engine block by hash by calling
// eth_blockByHash via JSON-RPC.
func (s *engineCaller) ExecutionBlockByHash(ctx context.Context, hash common.Hash, withTxs bool,
) (*pb.ExecutionBlock, error) {
	ctx, span := trace.StartSpan(ctx, "powchain.engine-api-client.ExecutionBlockByHash")
	defer span.End()
	result := &pb.ExecutionBlock{}
	err := s.Eth1Client.Client.Client().CallContext(
		ctx, result, execution.ExecutionBlockByHashMethod, hash, withTxs)
	return result, s.handleRPCError(err)
}

// ExecutionBlocksByHashes fetches a batch of execution engine blocks by hash by calling
// eth_blockByHash via JSON-RPC.
func (s *engineCaller) ExecutionBlocksByHashes(ctx context.Context,
	hashes []common.Hash, withTxs bool,
) ([]*pb.ExecutionBlock, error) {
	_, span := trace.StartSpan(ctx, "powchain.engine-api-client.ExecutionBlocksByHashes")
	defer span.End()
	numOfHashes := len(hashes)
	elems := make([]gethRPC.BatchElem, 0, numOfHashes)
	execBlks := make([]*pb.ExecutionBlock, 0, numOfHashes)
	if numOfHashes == 0 {
		return execBlks, nil
	}
	for _, h := range hashes {
		blk := &pb.ExecutionBlock{}
		newH := h
		elems = append(elems, gethRPC.BatchElem{
			Method: execution.ExecutionBlockByHashMethod,
			Args:   []interface{}{newH, withTxs},
			Result: blk,
			Error:  error(nil),
		})
		execBlks = append(execBlks, blk)
	}
	ioErr := s.Eth1Client.Client.Client().BatchCall(elems)
	if ioErr != nil {
		return nil, ioErr
	}
	for _, e := range elems {
		if e.Error != nil {
			return nil, s.handleRPCError(e.Error)
		}
	}
	return execBlks, nil
}

// HeaderByHash returns the relevant header details for the provided block hash.
func (s *engineCaller) HeaderByHash(ctx context.Context, hash common.Hash,
) (*types.HeaderInfo, error) {
	var hdr *types.HeaderInfo
	err := s.Eth1Client.Client.Client().CallContext(ctx, &hdr,
		execution.ExecutionBlockByHashMethod, hash, false /* no transactions */)
	if err == nil && hdr == nil {
		err = ethereum.NotFound
	}
	return hdr, err
}

// HeaderByNumber returns the relevant header details for the provided block number.
func (s *engineCaller) HeaderByNumber(ctx context.Context, number *big.Int,
) (*types.HeaderInfo, error) {
	var hdr *types.HeaderInfo
	err := s.Eth1Client.Client.Client().CallContext(ctx, &hdr,
		execution.ExecutionBlockByNumberMethod, toBlockNumArg(number), false /* no transactions */)
	if err == nil && hdr == nil {
		err = ethereum.NotFound
	}
	return hdr, err
}

// GetPayloadBodiesByHash returns the relevant payload bodies for the provided block hash.
func (s *engineCaller) GetPayloadBodiesByHash(
	ctx context.Context, executionBlockHashes []common.Hash,
) ([]*pb.ExecutionPayloadBodyV1, error) {
	ctx, span := trace.StartSpan(ctx, "powchain.engine-api-client.GetPayloadBodiesByHashV1")
	defer span.End()

	result := make([]*pb.ExecutionPayloadBodyV1, 0)
	err := s.Eth1Client.Client.Client().CallContext(ctx, &result,
		execution.GetPayloadBodiesByHashV1, executionBlockHashes)

	for i, item := range result {
		if item == nil {
			result[i] = &pb.ExecutionPayloadBodyV1{
				Transactions: make([][]byte, 0),
				Withdrawals:  make([]*pb.Withdrawal, 0),
			}
		}
	}
	return result, s.handleRPCError(err)
}

// GetPayloadBodiesByRange returns the relevant payload bodies for the provided range.
func (s *engineCaller) GetPayloadBodiesByRange(
	ctx context.Context, start, count uint64,
) ([]*pb.ExecutionPayloadBodyV1, error) {
	ctx, span := trace.StartSpan(ctx, "powchain.engine-api-client.GetPayloadBodiesByRangeV1")
	defer span.End()

	result := make([]*pb.ExecutionPayloadBodyV1, 0)
	err := s.Eth1Client.Client.Client().CallContext(ctx, &result,
		execution.GetPayloadBodiesByRangeV1, start, count)

	for i, item := range result {
		if item == nil {
			result[i] = &pb.ExecutionPayloadBodyV1{
				Transactions: make([][]byte, 0),
				Withdrawals:  make([]*pb.Withdrawal, 0),
			}
		}
	}
	return result, s.handleRPCError(err)
}

// ReconstructFullBlock takes in a blinded beacon block and reconstructs
// a beacon block with a full execution payload via the engine API.
func (s *engineCaller) ReconstructFullBlock(
	ctx context.Context, blindedBlock interfaces.ReadOnlySignedBeaconBlock,
) (interfaces.SignedBeaconBlock, error) {
	if err := blocks.BeaconBlockIsNil(blindedBlock); err != nil {
		return nil, errors.Wrap(err, "cannot reconstruct bellatrix block from nil data")
	}
	if !blindedBlock.Block().IsBlinded() {
		return nil, errors.New("can only reconstruct block from blinded block format")
	}
	header, err := blindedBlock.Block().Body().Execution()
	if err != nil {
		return nil, err
	}
	if header.IsNil() {
		return nil, errors.New("execution payload header in blinded block was nil")
	}

	// If the payload header has a block hash of 0x0, it means we are pre-merge and should
	// simply return the block with an empty execution payload.
	if bytes.Equal(header.BlockHash(), zeroHash32[:]) {
		var payload protoreflect.ProtoMessage
		payload, err = buildEmptyExecutionPayload(blindedBlock.Version())
		if err != nil {
			return nil, err
		}
		return blocks.BuildSignedBeaconBlockFromExecutionPayload(blindedBlock, payload)
	}

	executionBlockHash := common.BytesToHash(header.BlockHash())
	payload, err := s.retrievePayloadFromExecutionHash(ctx,
		executionBlockHash, header, blindedBlock.Version())
	if err != nil {
		return nil, err
	}
	fullBlock, err := blocks.BuildSignedBeaconBlockFromExecutionPayload(blindedBlock,
		payload.Proto())
	if err != nil {
		return nil, err
	}
	// reconstructedExecutionPayloadCount.Add(1)
	return fullBlock, nil
}

// ReconstructFullBellatrixBlockBatch takes in a batch of blinded beacon blocks and reconstructs
// them with a full execution payload for each block via the engine API.
func (s *engineCaller) ReconstructFullBellatrixBlockBatch(
	ctx context.Context, blindedBlocks []interfaces.ReadOnlySignedBeaconBlock,
) ([]interfaces.SignedBeaconBlock, error) {
	if len(blindedBlocks) == 0 {
		return []interfaces.SignedBeaconBlock{}, nil
	}
	executionHashes := []common.Hash{}
	validExecPayloads := []int{}
	zeroExecPayloads := []int{}
	for i, b := range blindedBlocks {
		if err := blocks.BeaconBlockIsNil(b); err != nil {
			return nil, errors.Wrap(err, "cannot reconstruct bellatrix block from nil data")
		}
		if !b.Block().IsBlinded() {
			return nil, errors.New("can only reconstruct block from blinded block format")
		}
		header, err := b.Block().Body().Execution()
		if err != nil {
			return nil, err
		}
		if header.IsNil() {
			return nil, errors.New("execution payload header in blinded block was nil")
		}
		// Determine if the block is pre-merge or post-merge. Depending on the result,
		// we will ask the execution engine for the full payload.
		if bytes.Equal(header.BlockHash(), zeroHash32[:]) {
			zeroExecPayloads = append(zeroExecPayloads, i)
		} else {
			executionBlockHash := common.BytesToHash(header.BlockHash())
			validExecPayloads = append(validExecPayloads, i)
			executionHashes = append(executionHashes, executionBlockHash)
		}
	}
	fullBlocks, err := s.retrievePayloadsFromExecutionHashes(ctx,
		executionHashes, validExecPayloads, blindedBlocks)
	if err != nil {
		return nil, err
	}
	// For blocks that are pre-merge we simply reconstruct them via an empty
	// execution payload.
	for _, realIdx := range zeroExecPayloads {
		bblock := blindedBlocks[realIdx]
		var payload protoreflect.ProtoMessage
		payload, err = buildEmptyExecutionPayload(bblock.Version())
		if err != nil {
			return nil, err
		}
		var fullBlock interfaces.SignedBeaconBlock
		fullBlock, err = blocks.BuildSignedBeaconBlockFromExecutionPayload(
			blindedBlocks[realIdx], payload,
		)
		if err != nil {
			return nil, err
		}
		fullBlocks[realIdx] = fullBlock
	}
	// reconstructedExecutionPayloadCount.Add(float64(len(blindedBlocks)))
	return fullBlocks, nil
}

func (s *engineCaller) retrievePayloadFromExecutionHash(ctx context.Context,
	executionBlockHash common.Hash, header interfaces.ExecutionData,
	version int) (interfaces.ExecutionData, error) {
	if features.Get().EnableOptionalEngineMethods {
		pBodies, err := s.GetPayloadBodiesByHash(ctx, []common.Hash{executionBlockHash})
		if err != nil {
			return nil, fmt.Errorf("could not get payload body by hash %#x: %w", executionBlockHash, err)
		}
		if len(pBodies) != 1 {
			return nil, errors.Errorf(
				"could not retrieve the correct number of payload bodies: wanted 1 but got %d",
				len(pBodies),
			)
		}
		bdy := pBodies[0]
		return fullPayloadFromPayloadBody(header, bdy, version)
	}

	executionBlock, err := s.ExecutionBlockByHash(ctx, executionBlockHash, true /* with txs */)
	if err != nil {
		return nil, fmt.Errorf("could not fetch execution block with txs by hash %#x: %w",
			executionBlockHash, err)
	}
	if executionBlock == nil {
		return nil, fmt.Errorf("received nil execution block for request by hash %#x",
			executionBlockHash)
	}
	if bytes.Equal(executionBlock.Hash.Bytes(), []byte{}) {
		return nil, ErrEmptyBlockHash
	}

	executionBlock.Version = version
	return fullPayloadFromExecutionBlock(version, header, executionBlock)
}

//nolint:gocognit // from prysm.
func (s *engineCaller) retrievePayloadsFromExecutionHashes(
	ctx context.Context,
	executionHashes []common.Hash,
	validExecPayloads []int,
	blindedBlocks []interfaces.ReadOnlySignedBeaconBlock) ([]interfaces.SignedBeaconBlock, error) {
	fullBlocks := make([]interfaces.SignedBeaconBlock, len(blindedBlocks))
	var execBlocks []*pb.ExecutionBlock
	var payloadBodies []*pb.ExecutionPayloadBodyV1
	var err error
	if features.Get().EnableOptionalEngineMethods {
		payloadBodies, err = s.GetPayloadBodiesByHash(ctx, executionHashes)
		if err != nil {
			return nil, fmt.Errorf("could not fetch payload bodies by hash %#x: %w",
				executionHashes, err)
		}
	} else {
		execBlocks, err = s.ExecutionBlocksByHashes(ctx, executionHashes, true /* with txs*/)
		if err != nil {
			return nil, fmt.Errorf("could not fetch execution blocks with txs by hash %#x: %w",
				executionHashes, err)
		}
	}

	// For each valid payload, we reconstruct the full block from it with the
	// blinded block.
	for sliceIdx, realIdx := range validExecPayloads {
		var payload interfaces.ExecutionData
		bblock := blindedBlocks[realIdx]
		//nolint:nestif // from prysm.
		if features.Get().EnableOptionalEngineMethods {
			b := payloadBodies[sliceIdx]
			if b == nil {
				return nil, fmt.Errorf("received nil payload body for request by hash %#x",
					executionHashes[sliceIdx])
			}
			var header interfaces.ExecutionData
			header, err = bblock.Block().Body().Execution()
			if err != nil {
				return nil, err
			}
			payload, err = fullPayloadFromPayloadBody(header, b, bblock.Version())
			if err != nil {
				return nil, err
			}
		} else {
			b := execBlocks[sliceIdx]
			if b == nil {
				return nil, fmt.Errorf("received nil execution block for request by hash %#x",
					executionHashes[sliceIdx])
			}
			var header interfaces.ExecutionData
			header, err = bblock.Block().Body().Execution()
			if err != nil {
				return nil, err
			}
			payload, err = fullPayloadFromExecutionBlock(bblock.Version(), header, b)
			if err != nil {
				return nil, err
			}
		}
		var fullBlock interfaces.SignedBeaconBlock
		fullBlock, err = blocks.BuildSignedBeaconBlockFromExecutionPayload(bblock,
			payload.Proto())
		if err != nil {
			return nil, err
		}
		fullBlocks[realIdx] = fullBlock
	}
	return fullBlocks, nil
}

func fullPayloadFromExecutionBlock(
	blockVersion int, header interfaces.ExecutionData, block *pb.ExecutionBlock,
) (interfaces.ExecutionData, error) {
	if header.IsNil() || block == nil {
		return nil, errors.New("execution block and header cannot be nil")
	}
	blockHash := block.Hash
	if !bytes.Equal(header.BlockHash(), blockHash[:]) {
		return nil, fmt.Errorf(
			"block hash field in execution header %#x does not match execution block hash %#x",
			header.BlockHash(),
			blockHash,
		)
	}
	blockTransactions := block.Transactions
	txs := make([][]byte, len(blockTransactions))
	for i, tx := range blockTransactions {
		txBin, err := tx.MarshalBinary()
		if err != nil {
			return nil, err
		}
		txs[i] = txBin
	}

	switch blockVersion {
	case version.Bellatrix:
		return blocks.WrappedExecutionPayload(&pb.ExecutionPayload{
			ParentHash:    header.ParentHash(),
			FeeRecipient:  header.FeeRecipient(),
			StateRoot:     header.StateRoot(),
			ReceiptsRoot:  header.ReceiptsRoot(),
			LogsBloom:     header.LogsBloom(),
			PrevRandao:    header.PrevRandao(),
			BlockNumber:   header.BlockNumber(),
			GasLimit:      header.GasLimit(),
			GasUsed:       header.GasUsed(),
			Timestamp:     header.Timestamp(),
			ExtraData:     header.ExtraData(),
			BaseFeePerGas: header.BaseFeePerGas(),
			BlockHash:     blockHash[:],
			Transactions:  txs,
		})
	case version.Capella:
		return blocks.WrappedExecutionPayloadCapella(&pb.ExecutionPayloadCapella{
			ParentHash:    header.ParentHash(),
			FeeRecipient:  header.FeeRecipient(),
			StateRoot:     header.StateRoot(),
			ReceiptsRoot:  header.ReceiptsRoot(),
			LogsBloom:     header.LogsBloom(),
			PrevRandao:    header.PrevRandao(),
			BlockNumber:   header.BlockNumber(),
			GasLimit:      header.GasLimit(),
			GasUsed:       header.GasUsed(),
			Timestamp:     header.Timestamp(),
			ExtraData:     header.ExtraData(),
			BaseFeePerGas: header.BaseFeePerGas(),
			BlockHash:     blockHash[:],
			Transactions:  txs,
			Withdrawals:   block.Withdrawals,
		}, 0) // We can't get the block value and don't care about the block value for this instance
	case version.Deneb:
		ebg, err := header.ExcessBlobGas()
		if err != nil {
			return nil, errors.Wrap(err,
				"unable to extract ExcessBlobGas attribute from excution payload header")
		}
		bgu, err := header.BlobGasUsed()
		if err != nil {
			return nil, errors.Wrap(err,
				"unable to extract BlobGasUsed attribute from excution payload header")
		}
		return blocks.WrappedExecutionPayloadDeneb(
			&pb.ExecutionPayloadDeneb{
				ParentHash:    header.ParentHash(),
				FeeRecipient:  header.FeeRecipient(),
				StateRoot:     header.StateRoot(),
				ReceiptsRoot:  header.ReceiptsRoot(),
				LogsBloom:     header.LogsBloom(),
				PrevRandao:    header.PrevRandao(),
				BlockNumber:   header.BlockNumber(),
				GasLimit:      header.GasLimit(),
				GasUsed:       header.GasUsed(),
				Timestamp:     header.Timestamp(),
				ExtraData:     header.ExtraData(),
				BaseFeePerGas: header.BaseFeePerGas(),
				BlockHash:     blockHash[:],
				Transactions:  txs,
				Withdrawals:   block.Withdrawals,
				ExcessBlobGas: ebg,
				BlobGasUsed:   bgu,
			}, 0) // We can't get the block value and don't care about the block value for this instance
	default:
		return nil, fmt.Errorf("unknown execution block version %d", block.Version)
	}
}

func fullPayloadFromPayloadBody(
	header interfaces.ExecutionData, body *pb.ExecutionPayloadBodyV1, bVersion int,
) (interfaces.ExecutionData, error) {
	if header.IsNil() || body == nil {
		return nil, errors.New("execution block and header cannot be nil")
	}

	switch bVersion {
	case version.Bellatrix:
		return blocks.WrappedExecutionPayload(&pb.ExecutionPayload{
			ParentHash:    header.ParentHash(),
			FeeRecipient:  header.FeeRecipient(),
			StateRoot:     header.StateRoot(),
			ReceiptsRoot:  header.ReceiptsRoot(),
			LogsBloom:     header.LogsBloom(),
			PrevRandao:    header.PrevRandao(),
			BlockNumber:   header.BlockNumber(),
			GasLimit:      header.GasLimit(),
			GasUsed:       header.GasUsed(),
			Timestamp:     header.Timestamp(),
			ExtraData:     header.ExtraData(),
			BaseFeePerGas: header.BaseFeePerGas(),
			BlockHash:     header.BlockHash(),
			Transactions:  body.Transactions,
		})
	case version.Capella:
		return blocks.WrappedExecutionPayloadCapella(&pb.ExecutionPayloadCapella{
			ParentHash:    header.ParentHash(),
			FeeRecipient:  header.FeeRecipient(),
			StateRoot:     header.StateRoot(),
			ReceiptsRoot:  header.ReceiptsRoot(),
			LogsBloom:     header.LogsBloom(),
			PrevRandao:    header.PrevRandao(),
			BlockNumber:   header.BlockNumber(),
			GasLimit:      header.GasLimit(),
			GasUsed:       header.GasUsed(),
			Timestamp:     header.Timestamp(),
			ExtraData:     header.ExtraData(),
			BaseFeePerGas: header.BaseFeePerGas(),
			BlockHash:     header.BlockHash(),
			Transactions:  body.Transactions,
			Withdrawals:   body.Withdrawals,
		}, 0) // We can't get the block value and don't care about the
		// block value for this instance
	case version.Deneb:
		ebg, err := header.ExcessBlobGas()
		if err != nil {
			return nil, errors.Wrap(err,
				"unable to extract ExcessBlobGas attribute from excution payload header")
		}
		bgu, err := header.BlobGasUsed()
		if err != nil {
			return nil, errors.Wrap(err,
				"unable to extract BlobGasUsed attribute from excution payload header")
		}
		return blocks.WrappedExecutionPayloadDeneb(
			&pb.ExecutionPayloadDeneb{
				ParentHash:    header.ParentHash(),
				FeeRecipient:  header.FeeRecipient(),
				StateRoot:     header.StateRoot(),
				ReceiptsRoot:  header.ReceiptsRoot(),
				LogsBloom:     header.LogsBloom(),
				PrevRandao:    header.PrevRandao(),
				BlockNumber:   header.BlockNumber(),
				GasLimit:      header.GasLimit(),
				GasUsed:       header.GasUsed(),
				Timestamp:     header.Timestamp(),
				ExtraData:     header.ExtraData(),
				BaseFeePerGas: header.BaseFeePerGas(),
				BlockHash:     header.BlockHash(),
				Transactions:  body.Transactions,
				Withdrawals:   body.Withdrawals,
				ExcessBlobGas: ebg,
				BlobGasUsed:   bgu,
			}, 0) // We can't get the block value and don't care about the
		// block value for this instance
	default:
		return nil, fmt.Errorf("unknown execution block version for payload %d", bVersion)
	}
}

// Handles errors received from the RPC server according to the specification.
func (s *engineCaller) handleRPCError(err error) error {
	if err == nil {
		return nil
	}
	if isTimeout(err) {
		return execution.ErrHTTPTimeout
	}
	e, ok := err.(gethRPC.Error) //nolint:errorlint // from prysm.
	if !ok {
		if strings.Contains(err.Error(), "401 Unauthorized") {
			s.logger.Error("HTTP authentication to your execution client is not working. " +
				"Please ensure you are setting a correct value for the --jwt-secret flag in " +
				"Prysm, or use an IPC connection if on the same machine. Please see our" +
				"documentation for more information on authenticating connections " +
				"here https://docs.prylabs.network/docs/execution-node/authentication")
			return fmt.Errorf("could not authenticate connection to execution client: %w", err)
		}
		return errors.Wrapf(err, "got an unexpected error in JSON-RPC response")
	}
	switch e.ErrorCode() {
	case -32700:
		// errParseCount.Inc()
		return execution.ErrParse
	case -32600:
		// errInvalidRequestCount.Inc()
		return execution.ErrInvalidRequest
	case -32601:
		// errMethodNotFoundCount.Inc()
		return execution.ErrMethodNotFound
	case -32602:
		// errInvalidParamsCount.Inc()
		return execution.ErrInvalidParams
	case -32603:
		// errInternalCount.Inc()
		return execution.ErrInternal
	case -38001:
		// errUnknownPayloadCount.Inc()
		return execution.ErrUnknownPayload
	case -38002:
		// errInvalidForkchoiceStateCount.Inc()
		return execution.ErrInvalidForkchoiceState
	case -38003:
		// errInvalidPayloadAttributesCount.Inc()
		return execution.ErrInvalidPayloadAttributes
	case -38004:
		// errRequestTooLargeCount.Inc()
		return execution.ErrRequestTooLarge
	case -32000:
		// errServerErrorCount.Inc()
		// Only -32000 status codes are data errors in the RPC specification.
		var errWithData gethRPC.DataError
		errWithData, ok = err.(gethRPC.DataError) //nolint:errorlint // from prysm.
		if !ok {
			return errors.Wrapf(err, "got an unexpected error in JSON-RPC response")
		}
		return errors.Wrapf(execution.ErrServer, "%v", errWithData.Error())
	default:
		return err
	}
}

// ErrHTTPTimeout returns true if the error is a http.Client timeout error.
var ErrHTTPTimeout = errors.New("timeout from http.Client")

type httpTimeoutError interface {
	Error() string
	Timeout() bool
}

func isTimeout(e error) bool {
	t, ok := e.(httpTimeoutError) //nolint:errorlint // from prysm.
	return ok && t.Timeout()
}

func buildEmptyExecutionPayload(v int) (proto.Message, error) {
	switch v {
	case version.Bellatrix:
		return &pb.ExecutionPayload{
			ParentHash:    make([]byte, fieldparams.RootLength),
			FeeRecipient:  make([]byte, fieldparams.FeeRecipientLength),
			StateRoot:     make([]byte, fieldparams.RootLength),
			ReceiptsRoot:  make([]byte, fieldparams.RootLength),
			LogsBloom:     make([]byte, fieldparams.LogsBloomLength),
			PrevRandao:    make([]byte, fieldparams.RootLength),
			BaseFeePerGas: make([]byte, fieldparams.RootLength),
			BlockHash:     make([]byte, fieldparams.RootLength),
			Transactions:  make([][]byte, 0),
			ExtraData:     make([]byte, 0),
		}, nil
	case version.Capella:
		return &pb.ExecutionPayloadCapella{
			ParentHash:    make([]byte, fieldparams.RootLength),
			FeeRecipient:  make([]byte, fieldparams.FeeRecipientLength),
			StateRoot:     make([]byte, fieldparams.RootLength),
			ReceiptsRoot:  make([]byte, fieldparams.RootLength),
			LogsBloom:     make([]byte, fieldparams.LogsBloomLength),
			PrevRandao:    make([]byte, fieldparams.RootLength),
			BaseFeePerGas: make([]byte, fieldparams.RootLength),
			BlockHash:     make([]byte, fieldparams.RootLength),
			Transactions:  make([][]byte, 0),
			ExtraData:     make([]byte, 0),
			Withdrawals:   make([]*pb.Withdrawal, 0),
		}, nil
	case version.Deneb:
		return &pb.ExecutionPayloadDeneb{
			ParentHash:    make([]byte, fieldparams.RootLength),
			FeeRecipient:  make([]byte, fieldparams.FeeRecipientLength),
			StateRoot:     make([]byte, fieldparams.RootLength),
			ReceiptsRoot:  make([]byte, fieldparams.RootLength),
			LogsBloom:     make([]byte, fieldparams.LogsBloomLength),
			PrevRandao:    make([]byte, fieldparams.RootLength),
			BaseFeePerGas: make([]byte, fieldparams.RootLength),
			BlockHash:     make([]byte, fieldparams.RootLength),
			Transactions:  make([][]byte, 0),
			ExtraData:     make([]byte, 0),
			Withdrawals:   make([]*pb.Withdrawal, 0),
		}, nil
	default:
		return nil, errors.Wrapf(execution.ErrUnsupportedVersion, "version=%s", version.String(v))
	}
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	pending := big.NewInt(-1)
	if number.Cmp(pending) == 0 {
		return "pending"
	}
	finalized := big.NewInt(int64(gethRPC.FinalizedBlockNumber))
	if number.Cmp(finalized) == 0 {
		return "finalized"
	}
	safe := big.NewInt(int64(gethRPC.SafeBlockNumber))
	if number.Cmp(safe) == 0 {
		return "safe"
	}
	return hexutil.EncodeBig(number)
}
