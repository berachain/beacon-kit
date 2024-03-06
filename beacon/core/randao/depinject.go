package randao

import (
	"cosmossdk.io/depinject"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	bls12381 "github.com/itsdevbear/bolaris/crypto/bls12_381"
)

// DepInjectInput is the input for the dep inject framework.
type DepInjectInput struct {
	depinject.In

	BeaconState BeaconStateProvider
	AppOpts     servertypes.AppOptions
	Signer      bls12381.BlsSigner
}

// DepInjectOutput is the output for the dep inject framework.
type DepInjectOutput struct {
	depinject.Out

	RandaoProcessor *Processor
}

func ProvideRandaoProcessor(in DepInjectInput) DepInjectOutput {
	processor := NewProcessor(in.BeaconState, in.Signer, &Config{
		EpochsPerHistoricalVector: 0,
		ConfiguredPubKeyLength:    0,
	})

	return DepInjectOutput{
		RandaoProcessor: processor,
	}
}
