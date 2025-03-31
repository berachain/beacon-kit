package utils

import (
	"testing"

	"github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/stretchr/testify/require"
)

// GenerateValidBeaconBlock generates a valid beacon block for the Deneb.
func GenerateValidBeaconBlock(t *testing.T, forkVersion common.Version) *types.BeaconBlock {
	t.Helper()

	// Initialize your block here
	beaconBlock, err := types.NewBeaconBlockWithVersion(
		math.Slot(10),
		math.ValidatorIndex(5),
		common.Root{1, 2, 3, 4, 5}, // parent block root
		forkVersion,
	)
	require.NoError(t, err)

	versionable := types.NewVersionable(forkVersion)
	beaconBlock.StateRoot = common.Root{5, 4, 3, 2, 1}
	beaconBlock.Body = &types.BeaconBlockBody{
		Versionable: versionable,
		ExecutionPayload: &types.ExecutionPayload{
			Versionable: versionable,
			Timestamp:   10,
			ExtraData:   []byte("dummy extra data for testing"),
			Transactions: [][]byte{
				[]byte("0x"),
				[]byte("0x"),
				[]byte("0x"),
			},
			Withdrawals: engineprimitives.Withdrawals{
				{Index: 0, Amount: 100},
				{Index: 1, Amount: 200},
			},
			BaseFeePerGas: math.NewU256(0),
		},
		Eth1Data: &types.Eth1Data{},
		Deposits: []*types.Deposit{
			{
				Index: 1,
			},
		},
		BlobKzgCommitments: []eip4844.KZGCommitment{
			{1, 2, 3},
		},
	}
	body := beaconBlock.GetBody()
	body.SetProposerSlashings(types.ProposerSlashings{})
	body.SetAttesterSlashings(types.AttesterSlashings{})
	body.SetAttestations(types.Attestations{})
	body.SetSyncAggregate(&types.SyncAggregate{})
	body.SetVoluntaryExits(types.VoluntaryExits{})
	body.SetBlsToExecutionChanges(types.BlsToExecutionChanges{})
	if version.EqualsOrIsAfter(forkVersion, version.Electra()) {
		err = body.SetExecutionRequests(&types.ExecutionRequests{
			Deposits: []*types.DepositRequest{
				{
					Pubkey:      crypto.BLSPubkey{1, 2, 3},
					Credentials: types.WithdrawalCredentials(bytes.B32{4, 5, 6}),
					Amount:      100,
					Signature:   crypto.BLSSignature{1, 2, 3},
					Index:       1,
				},
			},
			Withdrawals: []*types.WithdrawalRequest{
				{
					SourceAddress:   common.ExecutionAddress{0, 1, 2, 3, 4, 5},
					ValidatorPubKey: crypto.BLSPubkey{4, 2, 0},
					Amount:          1000,
				},
			},
			Consolidations: []*types.ConsolidationRequest{},
		})
		require.NoError(t, err)
	}
	return beaconBlock
}
