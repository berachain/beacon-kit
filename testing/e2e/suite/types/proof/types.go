package proof

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// JSON response returned from the
// `bkit/v1/proof/execution_number/{timestamp_id}` endpoint
type executionNumberResponse struct {
	ExecutionNumber math.U64      `json:"execution_number"`
	Proof           []common.Root `json:"execution_number_proof"`

	// TODO: other fields...
}

// JSON response returned from the
// `bkit/v1/proof/block_proposer/{timestamp_id}` endpoint
type blockProposerProofResponse struct {
	// TODO: Switch to use the json tag once resolved on the beacon API
	BeaconBlockHeader struct {
		ProposerIndex math.U64 `json:"proposer_index"`
	} `json:"beacon_block_header"`
	Pubkey bytes.B48     `json:"validator_pubkey"`
	Proof  []common.Root `json:"validator_pubkey_proof"`
}
