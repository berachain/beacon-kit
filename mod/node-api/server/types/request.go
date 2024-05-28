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

package types

type StateIDRequest struct {
	StateID string `param:"state_id" validate:"required,state_id"`
}

type BlockIDRequest struct {
	BlockID string `param:"block_id" validate:"required,block_id"`
}

type StateValidatorsGetRequest struct {
	StateIDRequest
	ID     []string `query:"id"     validate:"dive,validator_id"`
	Status []string `query:"status" validate:"dive,validator_status"`
}

type StateValidatorsPostRequest struct {
	StateIDRequest
	IDs      []string `json:"IDs"      validate:"dive,validator_id"`
	Statuses []string `json:"statuses" validate:"dive,validator_status"`
}

type StateValidatorRequest struct {
	StateIDRequest
	ValidatorID string `query:"validator_id" validate:"required,validator_id"`
}

type ValidatorBalancesGetRequest struct {
	StateIDRequest
	ID []string `query:"id" validate:"dive,validator_id"`
}

type ValidatorBalancesPostRequest struct {
	StateIDRequest
	IDs []string `validate:"dive,validator_id"`
}

type EpochOptionalRequest struct {
	Epoch string `query:"epoch" validate:"epoch"`
}

type EpochRequest struct {
	Epoch string `param:"epoch" validate:"required,epoch"`
}

type CommitteeIndexRequest struct {
	ComitteeIndex string `query:"committee_index" validate:"committee_index"`
}

type SlotRequest struct {
	Slot string `query:"slot" validate:"slot"`
}

type CommitteesRequest struct {
	StateIDRequest
	EpochOptionalRequest
	CommitteeIndexRequest
	SlotRequest
}

type SyncCommitteesRequest struct {
	StateIDRequest
	EpochOptionalRequest
}

type BeaconHeadersRequest struct {
	SlotRequest
	ParentRoot string `query:"parent_root" validate:"hex"`
}

type BlobSidecarRequest struct {
	BlockIDRequest
	Indices []string `query:"indices" validate:"dive,uint64"`
}
