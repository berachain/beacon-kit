// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2023, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
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

package execution

import (
	"context"
	"math/big"

	"github.com/prysmaticlabs/prysm/v4/beacon-chain/execution/types"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/interfaces"
	payloadattribute "github.com/prysmaticlabs/prysm/v4/consensus-types/payload-attribute"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	pb "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"

	"github.com/ethereum/go-ethereum/common"
	gethRPC "github.com/ethereum/go-ethereum/rpc"
)

// RPCClient defines the rpc methods required to interact with the eth1 node.
type RPCClient interface {
	Close()
	BatchCall(b []gethRPC.BatchElem) error
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
}

// ExecutionPayloadReconstructor defines a service that can reconstruct a full beacon
// block with an execution payload from a signed beacon block and a connection
// to an execution client's engine API.
type ExecutionPayloadReconstructor interface {
	ReconstructFullBlock(
		ctx context.Context, blindedBlock interfaces.ReadOnlySignedBeaconBlock,
	) (interfaces.SignedBeaconBlock, error)
	ReconstructFullBellatrixBlockBatch(
		ctx context.Context, blindedBlocks []interfaces.ReadOnlySignedBeaconBlock,
	) ([]interfaces.SignedBeaconBlock, error)
}

// EngineCaller defines a client that can interact with an Ethereum
// execution node's engine service via JSON-RPC.
type EngineCaller interface {
	NewPayload(ctx context.Context, payload interfaces.ExecutionData,
		versionedHashes []common.Hash, parentBlockRoot *common.Hash) ([]byte, error)
	ForkchoiceUpdated(
		ctx context.Context, state *pb.ForkchoiceState, attrs payloadattribute.Attributer,
	) (*pb.PayloadIDBytes, []byte, error)
	GetPayload(ctx context.Context, payloadID [8]byte,
		slot primitives.Slot) (interfaces.ExecutionData, *pb.BlobsBundle, bool, error)
	ExecutionBlockByHash(ctx context.Context, hash common.Hash,
		withTxs bool) (*pb.ExecutionBlock, error)
}

type ExecutionBlockCaller interface {
	LatestExecutionBlock(context.Context) (*pb.ExecutionBlock, error)
	LatestFinalizedBlock(ctx context.Context) (*pb.ExecutionBlock, error)
	LatestSafeBlock(ctx context.Context) (*pb.ExecutionBlock, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.HeaderInfo, error)
	EarliestBlock(ctx context.Context) (*pb.ExecutionBlock, error)
}

type (
	EngineAPI interface {
		EngineCaller
		ExecutionBlockCaller
	}
)

// ExecutionClient represents the execution layer client.
type ExecutionClient struct {
	EngineAPI
}
