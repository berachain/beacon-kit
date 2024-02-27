package randao

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/itsdevbear/bolaris/beacon/core/state"
	"github.com/itsdevbear/bolaris/types/consensus"
	"github.com/itsdevbear/bolaris/types/consensus/primitives"
)

// ProcessRandao checks the block proposer's
// randao commitment and generates a new randao mix to update
// in the beacon state's latest randao mixes slice.
//
// Spec pseudocode definition:
//
//	def process_randao(state: BeaconState, body: ReadOnlyBeaconBlockBody) -> None:
//	 epoch = get_current_epoch(state)
//	 # Verify RANDAO reveal
//	 proposer = state.validators[get_beacon_proposer_index(state)]
//	 signing_root = compute_signing_root(epoch, get_domain(state, DOMAIN_RANDAO))
//	 assert bls.Verify(proposer.pubkey, signing_root, body.randao_reveal)
//	 # Mix in RANDAO reveal
//	 mix = xor(get_randao_mix(state, epoch), hash(body.randao_reveal))
//	 state.randao_mixes[epoch % EPOCHS_PER_HISTORICAL_VECTOR] = mix
func ProcessRandao(
	ctx context.Context,
	beaconState state.BeaconState,
	b consensus.ReadOnlyBeaconKitBlock,
) (state.BeaconState, error) {
	buf, proposerPub, domain, err := randaoSigningData(ctx, beaconState)
	if err != nil {
		return nil, err
	}

	// randaoReveal := b.Body().RandaoReveal()
	randaoReveal := Reveal(b.RandaoReveal())
	if err := verifySignature(buf, proposerPub, randaoReveal[:], domain); err != nil {
		return nil, errors.Wrap(err, "could not verify block randao")
	}

	beaconState, err = ProcessRandaoNoVerify(beaconState, b.GetSlot(), randaoReveal)
	if err != nil {
		return nil, errors.Wrap(err, "could not process randao")
	}
	return beaconState, nil
}

// ProcessRandaoNoVerify generates a new randao mix to update
// in the beacon state's latest randao mixes slice.
//
// Spec pseudocode definition:
//
//	# Mix it in
//	state.latest_randao_mixes[get_current_epoch(state) % LATEST_RANDAO_MIXES_LENGTH] = (
//	    xor(get_randao_mix(state, get_current_epoch(state)),
//	        hash(body.randao_reveal))
//	)
func ProcessRandaoNoVerify(
	beaconState state.BeaconState,
	slot primitives.Slot,
	randaoReveal Reveal,
) (state.BeaconState, error) {
	// currentEpoch := slots.ToEpoch(slot)
	// currentEpoch := 0
	// If block randao passed verification, we XOR the state's latest randao mix with the block's
	// randao and update the state's corresponding latest randao mix value.
	// latestMixesLength := params.BeaconConfig().EpochsPerHistoricalVector
	// latestMixSlice, err := beaconState.RandaoMixAtIndex(uint64(currentEpoch % latestMixesLength))
	// if err != nil {
	// 	return nil, err
	// }

	latestMix := Mix([]byte{})
	if err := latestMix.MixInRandao(randaoReveal); err != nil {
		return nil, err
	}

	// if err := beaconState.UpdateRandaoMixesAtIndex(uint64(currentEpoch%latestMixesLength), [32]byte(latestMixSlice)); err != nil {
	// 	return nil, err
	// }
	return beaconState, nil
}

// retrieves the randao related signing data from the state.
func randaoSigningData(ctx context.Context, beaconState state.ReadOnlyBeaconState) ([]byte, []byte, []byte, error) {
	// proposerIdx, err := helpers.BeaconProposerIndex(ctx, beaconState)
	// if err != nil {
	// 	return nil, nil, nil, errors.Wrap(err, "could not get beacon proposer index")
	// }
	// proposerPub := beaconState.PubkeyAtIndex(proposerIdx)

	// currentEpoch := slots.ToEpoch(beaconState.Slot())
	// buf := make([]byte, 32)
	// binary.LittleEndian.PutUint64(buf, uint64(currentEpoch))

	// domain, err := signing.Domain(beaconState.Fork(), currentEpoch, params.BeaconConfig().DomainRandao, beaconState.GenesisValidatorsRoot())
	// if err != nil {
	// 	return nil, nil, nil, err
	// }
	// return buf, proposerPub[:], domain, nil
	return nil, nil, nil, nil
}

// verifies the signature from the raw data, public key and domain provided.
func verifySignature(signedData, pub, signature, domain []byte) error {
	// set, err := signatureBatch(signedData, pub, signature, domain, signing.UnknownSignature)
	// if err != nil {
	// 	return err
	// }
	// if len(set.Signatures) != 1 {
	// 	return errors.Errorf("signature set contains %d signatures instead of 1", len(set.Signatures))
	// }
	// // We assume only one signature set is returned here.
	// sig := set.Signatures[0]
	// publicKey := set.PublicKeys[0]
	// root := set.Messages[0]
	// rSig, err := bls.SignatureFromBytes(sig)
	// if err != nil {
	// 	return err
	// }
	// if !rSig.Verify(publicKey, root[:]) {
	// 	return signing.ErrSigFailedToVerify
	// }
	// return nil
	return nil
}
