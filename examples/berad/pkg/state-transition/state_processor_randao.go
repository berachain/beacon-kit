package state_transition

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto/sha256"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	"github.com/go-faster/xor"
)

// processRandaoReveal processes the randao reveal and
// ensures it matches the local state.
func (sp *StateProcessor[
	BeaconBlockT, _, _, BeaconStateT,
	_, _, _, _, _, ForkDataT, _, _, _, _, _,
]) processRandaoReveal(
	st BeaconStateT,
	blk BeaconBlockT,
	skipVerification bool,
) error {
	slot, err := st.GetSlot()
	if err != nil {
		return err
	}

	// Ensure the proposer index is valid.
	proposer, err := st.ValidatorByIndex(blk.GetProposerIndex())
	if err != nil {
		return err
	}

	genesisValidatorsRoot, err := st.GetGenesisValidatorsRoot()
	if err != nil {
		return err
	}

	epoch := sp.cs.SlotToEpoch(slot)
	body := blk.GetBody()

	var fd ForkDataT
	fd = fd.New(
		version.FromUint32[common.Version](
			sp.cs.ActiveForkVersionForEpoch(epoch),
		), genesisValidatorsRoot,
	)

	if !skipVerification {
		signingRoot := fd.ComputeRandaoSigningRoot(
			sp.cs.DomainTypeRandao(), epoch,
		)
		reveal := body.GetRandaoReveal()
		if err = sp.signer.VerifySignature(
			proposer.GetPubkey(),
			signingRoot[:],
			reveal,
		); err != nil {
			return err
		}
	}

	prevMix, err := st.GetRandaoMixAtIndex(
		uint64(epoch) % sp.cs.EpochsPerHistoricalVector(),
	)
	if err != nil {
		return err
	}

	return st.UpdateRandaoMixAtIndex(
		uint64(epoch)%sp.cs.EpochsPerHistoricalVector(),
		sp.buildRandaoMix(prevMix, body.GetRandaoReveal()),
	)
}

// processRandaoMixesReset as defined in the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#randao-mixes-updates
//
//nolint:lll
func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, _, _, _, _, _, _, _, _, _, _,
]) processRandaoMixesReset(
	st BeaconStateT,
) error {
	slot, err := st.GetSlot()
	if err != nil {
		return err
	}

	epoch := sp.cs.SlotToEpoch(slot)
	mix, err := st.GetRandaoMixAtIndex(
		uint64(epoch) % sp.cs.EpochsPerHistoricalVector(),
	)
	if err != nil {
		return err
	}
	return st.UpdateRandaoMixAtIndex(
		uint64(epoch+1)%sp.cs.EpochsPerHistoricalVector(),
		mix,
	)
}

// buildRandaoMix as defined in the Ethereum 2.0 specification.
func (sp *StateProcessor[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) buildRandaoMix(
	mix common.Bytes32,
	reveal crypto.BLSSignature,
) common.Bytes32 {
	newMix := make([]byte, constants.RootLength)
	revealHash := sha256.Hash(reveal[:])
	// Apparently this library giga fast? Good project? lmeow.
	// It is safe to ignore this error, since it is guaranteed that
	// mix[:] and revealHash[:] are both Bytes32.
	_ = xor.Bytes(
		newMix, mix[:], revealHash[:],
	)
	return common.Bytes32(newMix)
}
