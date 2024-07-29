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

package types

import "github.com/berachain/beacon-kit/mod/node-api/handlers/types"

type GetGenesisRequest struct{}

type GetStateRootRequest struct {
	types.StateIDRequest
}

type GetStateForkRequest struct {
	types.StateIDRequest
}

type GetFinalityCheckpointsRequest struct {
	types.StateIDRequest
}

type GetStateValidatorsRequest struct {
	types.StateIDRequest
	IDs      []string `query:"id"     validate:"dive,validator_id"`
	Statuses []string `query:"status" validate:"dive,validator_status"`
}

type PostStateValidatorsRequest struct {
	types.StateIDRequest
	IDs      []string `json:"ids"      validate:"dive,validator_id"`
	Statuses []string `json:"statuses" validate:"dive,validator_status"`
}

type GetStateValidatorRequest struct {
	types.StateIDRequest
	ValidatorID string `query:"validator_id" validate:"required,validator_id"`
}

type GetValidatorBalancesRequest struct {
	types.StateIDRequest
	IDs []string `query:"id" validate:"dive,validator_id"`
}

type PostValidatorBalancesRequest struct {
	types.StateIDRequest
	IDs []string `validate:"dive,validator_id"`
}

type GetStateCommitteesRequest struct {
	types.StateIDRequest
	EpochOptionalRequest
	CommitteeIndexRequest
	SlotRequest
}

type GetSyncCommitteesRequest struct {
	types.StateIDRequest
	EpochOptionalRequest
}

type GetRandaoRequest struct {
	types.StateIDRequest
	EpochOptionalRequest
}

type GetBlockHeadersRequest struct {
	SlotRequest
	ParentRoot string `query:"parent_root" validate:"hex"`
}

type GetBlockHeaderRequest struct {
	types.BlockIDRequest
}

// TODO: body is big
//
//nolint:lll // tags get long
type PostBlindedBlocksV1Request struct {
	EthConsensusVersion string `json:"eth_consensus_version" validate:"required,eth_consensus_version"`
}

//nolint:lll // tags get long
type PostBlindedBlocksV2Request struct {
	PostBlindedBlocksV1Request
	BroadcastValidation string `json:"broadcast_validation" validate:"required,broadcast_validation"`
}

//nolint:lll // tags get long
type PostBlocksV1Request[BeaconBlockT any] struct {
	EthConsensusVersion string       `json:"eth_consensus_version" validate:"required,eth_consensus_version"`
	BeaconBlock         BeaconBlockT `json:"beacon_block"`
}

//nolint:lll // tags get long
type PostBlocksV2Request[BeaconBlockT any] struct {
	PostBlocksV1Request[BeaconBlockT]
	BroadcastValidation string `json:"broadcast_validation" validate:"required,broadcast_validation"`
}

type GetBlocksRequest struct {
	types.BlockIDRequest
}

type GetBlockRootRequest struct {
	types.BlockIDRequest
}

type GetBlockAttestationsRequest struct {
	types.BlockIDRequest
}

type GetBlobSidecarsRequest struct {
	types.BlockIDRequest
	Indices []string `query:"indices" validate:"dive,uint64"`
}

type PostRewardsSyncCommitteeRequest struct {
	types.BlockIDRequest
	IDs []string `validate:"dive,validator_id"`
}

type GetDepositTreeSnapshotRequest struct{}

type GetBlockRewardsRequest struct {
	types.BlockIDRequest
}

type PostAttestationsRewardsRequest struct {
	EpochRequest
	IDs []string `validate:"dive,validator_id"`
}

type GetBlindedBlockRequest struct {
	types.BlockIDRequest
}

type EpochOptionalRequest struct {
	Epoch string `query:"epoch" validate:"epoch"`
}

type EpochRequest struct {
	Epoch string `param:"epoch" validate:"required,epoch"`
}

type CommitteeIndexRequest struct {
	CommitteeIndex string `query:"committee_index" validate:"committee_index"`
}

type SlotRequest struct {
	Slot string `query:"slot" validate:"slot"`
}

type HeadersRequest struct {
	SlotRequest
	ParentRoot string `query:"parent_root" validate:"hex"`
}

type BlobSidecarRequest struct {
	types.BlockIDRequest
	Indices []string `query:"indices" validate:"dive,uint64"`
}
