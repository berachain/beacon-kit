package types

type StateIdRequest struct {
	StateId string `param:"state_id" validate:"required,state_id"`
}

type BlockIdRequest struct {
	BlockId string `param:"block_id" validate:"required,block_id"`
}

type StateValidatorsGetRequest struct {
	StateIdRequest
	Id     []string `query:"id" validate:"dive,validator_id"`
	Status []string `query:"status" validate:"dive,validator_status"`
}

type StateValidatorsPostRequest struct {
	StateIdRequest
	Ids      []string `json:"ids" validate:"dive,validator_id"`
	Statuses []string `json:"statuses" validate:"dive,validator_status"`
}

type StateValidatorRequest struct {
	StateIdRequest
	ValidatorId string `query:"validator_id" validate:"required,validator_id"`
}

type ValidatorBalancesGetRequest struct {
	StateIdRequest
	Id []string `query:"id" validate:"dive,validator_id"`
}

type ValidatorBalancesPostRequest struct {
	StateIdRequest
	Ids []string `validate:"dive,validator_id"`
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
	StateIdRequest
	EpochOptionalRequest
	CommitteeIndexRequest
	SlotRequest
}

type SyncCommitteesRequest struct {
	StateIdRequest
	EpochOptionalRequest
}

type BeaconHeadersRequest struct {
	SlotRequest
	ParentRoot string `query:"parent_root" validate:"hex"`
}

type BlobSidecarRequest struct {
	BlockIdRequest
	Indices []string `query:"indices" validate:"dive,uint64"`
}
