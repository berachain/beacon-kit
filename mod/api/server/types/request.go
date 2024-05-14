package types

// custom validator? <slot> <hex>
type StateIdRequest struct {
	StateId string `param:"state_id" validate:"required,oneof=head genesis finalized justified"`
}

// custom validator? <slot> <hex>
type BlockIdRequest struct {
	BlockId string `query:"block_id" validate:"required,oneof=genesis finalized justified"`
}

type StateValidatorsGetRequest struct {
	StateIdRequest
	Id     []string `query:"id" validate:"dive,hexadecimal,lte=98"`
	Status []string `query:"status" validate:"dive,oneof=pending_initialized pending_queued active_ongoing active_exiting active_slashed exited_unslashed exited_slashed withdrawal_possible withdrawal_done active pending exited withdrawal"`
}

type StateValidatorsPostRequest struct {
	StateIdRequest
	Ids      []string `json:"ids" validate:"required,dive,hexadecimal,lte=98"`
	Statuses []string `json:"statuses" validate:"required,dive,oneof=pending_initialized pending_queued active_ongoing active_exiting active_slashed exited_unslashed exited_slashed withdrawal_possible withdrawal_done active pending exited withdrawal"`
}

type StateValidatorRequest struct {
	StateIdRequest
	ValidatorId string `query:"validator_id" validate:"required,hexadecimal,lte=98"`
}

type ValidatorBalancesGetRequest struct {
	StateIdRequest
	Id []string `query:"id" validate:"dive,hexadecimal,lte=98"`
}

// test.. id is a top level array
type ValidatorBalancesPostRequest struct {
	StateIdRequest
	Ids []string `query:"" validate:"dive,hexadecimal,lte=98"`
}

type EpochOptionalRequest struct {
	Epoch string `query:"epoch" validate:"numeric"`
}

type EpochRequest struct {
	Epoch string `param:"epoch" validate:"required,numeric"`
}

type ComitteeIndexRequest struct {
	ComitteeIndex string `query:"index" validate:"numeric"`
}

type SlotRequest struct {
	Slot string `query:"slot" validate:"numeric"`
}

type ComitteesRequest struct {
	StateIdRequest
	EpochOptionalRequest
	ComitteeIndexRequest
	SlotRequest
}

type SyncComitteesRequest struct {
	StateIdRequest
	EpochOptionalRequest
}

type RandaoRequest struct {
	StateIdRequest
	EpochOptionalRequest
}

type BeaconHeadersRequest struct {
	SlotRequest
	ParentRoot string `query:"parent_root" validate:"hexadecimal"`
}

type BlobSidecarRequest struct {
	BlockIdRequest
	Indices []string `query:"indices" validate:"dive,numeric"`
}

type SyncComitteeAwardsRequest struct {
	BlockIdRequest
	Ids []string `json:"ids" validate:"required,dive,hexadecimal"`
}

type GetPoolAttestationRequest struct {
	SlotRequest
	CommitteeIndex string `query:"committee_index" validate:"numeric"`
}

type EventsRequest struct {
	Topics []string `query:"topics" validate:"required,dive,oneof=head block block_gossip attestation voluntary_exit bls_to_execution_change proposer_slashing attester_slashing finalized_checkpoint chain_reorg contribution_and_proof light_client_finality_update light_client_optimistic_update payload_attributes blob_sidecar"`
}

type GetAttestationRewardsRequest struct {
	EpochRequest
	Ids []string `query:"ids" validate:"dive,hexadecimal,lte=98"`
}
