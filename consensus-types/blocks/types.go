package blocks

import (
	"github.com/berachain/beacon-kit/consensus-types/deneb"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/math"
)

// SignedBeaconBlock is the main signed beacon block structure. It can represent any block type.
type SignedBeaconBlock struct {
	Message   *BeaconBlock
	Signature crypto.BLSSignature
}

// BeaconBlock is the main beacon block structure. It can represent any block type.

type BeaconBlock struct {
	// Slot represents the position of the block in the chain.
	Slot math.Slot
	// ProposerIndex is the index of the validator who proposed the block.
	ProposerIndex math.ValidatorIndex
	// ParentRoot is the hash of the parent block
	ParentRoot common.Root
	// StateRoot is the hash of the state at the block.
	StateRoot common.Root
	// Body is the body of the BeaconBlock, containing the block's operations.
	Body *BeaconBlockBody

	// BbVersion is the BbVersion of the beacon block.
	// BbVersion must be not serialized but it's exported
	// to allow unit tests using reflect on beacon block.
	BbVersion common.Version `json:"-"`
}

// BeaconBlockBody is the main beacon block body structure. It can represent any block type.
type BeaconBlockBody struct {
	// RandaoReveal is the reveal of the RANDAO.
	RandaoReveal crypto.BLSSignature
	// Eth1Data is the data from the Eth1 chain.
	Eth1Data *deneb.Eth1Data
	// Graffiti is for a fun message or meme.
	Graffiti [32]byte
	// proposerSlashings is unused but left for compatibility.
	proposerSlashings []*deneb.ProposerSlashing
	// attesterSlashings is unused but left for compatibility.
	attesterSlashings []*deneb.AttesterSlashing
	// attestations is unused but left for compatibility.
	attestations []*deneb.Attestation
	// Deposits is the list of deposits included in the body.
	Deposits []*deneb.Deposit
	// voluntaryExits is unused but left for compatibility.
	voluntaryExits []*deneb.VoluntaryExit
	// syncAggregate is unused but left for compatibility.
	syncAggregate *deneb.SyncAggregate
	// ExecutionPayload is the execution payload of the body.
	ExecutionPayload *deneb.ExecutionPayload
	// blsToExecutionChanges is unused but left for compatibility.
	blsToExecutionChanges []*deneb.BlsToExecutionChange
	// BlobKzgCommitments is the list of KZG commitments for the EIP-4844 blobs.
	BlobKzgCommitments []eip4844.KZGCommitment
}
