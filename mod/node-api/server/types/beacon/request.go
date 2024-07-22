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

import "github.com/berachain/beacon-kit/mod/node-api/server/types"

type StateValidatorsGetRequest struct {
	types.StateIDRequest
	ID     []string `query:"id"     validate:"dive,validator_id"`
	Status []string `query:"status" validate:"dive,validator_status"`
}

type StateValidatorsPostRequest struct {
	types.StateIDRequest
	IDs      []string `json:"IDs"      validate:"dive,validator_id"`
	Statuses []string `json:"statuses" validate:"dive,validator_status"`
}

type StateValidatorRequest struct {
	types.StateIDRequest
	ValidatorID string `query:"validator_id" validate:"required,validator_id"`
}

type ValidatorBalancesGetRequest struct {
	types.StateIDRequest
	ID []string `query:"id" validate:"dive,validator_id"`
}

type ValidatorBalancesPostRequest struct {
	types.StateIDRequest
	IDs []string `validate:"dive,validator_id"`
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

type CommitteesRequest struct {
	types.StateIDRequest
	EpochOptionalRequest
	CommitteeIndexRequest
	SlotRequest
}

type SyncCommitteesRequest struct {
	types.StateIDRequest
	EpochOptionalRequest
}

type BeaconHeadersRequest struct {
	SlotRequest
	ParentRoot string `query:"parent_root" validate:"hex"`
}

type BlobSidecarRequest struct {
	types.BlockIDRequest
	Indices []string `query:"indices" validate:"dive,uint64"`
}
