package bls12381

import (
	"cosmossdk.io/depinject"
	"fmt"
	"github.com/cometbft/cometbft/p2p"
	"github.com/cosmos/cosmos-sdk/client/flags"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cast"
	"os"
)

// DepInjectInput is the input for the dep inject framework.
type DepInjectInput struct {
	depinject.In

	AppOpts servertypes.AppOptions
}

// DepInjectOutput is the output for the dep inject framework.
type DepInjectOutput struct {
	depinject.Out

	BlsSigner BlsSigner
}

func ProvideBlsSigner(in DepInjectInput) DepInjectOutput {
	homeDir := cast.ToString(in.AppOpts.Get(flags.FlagHome))

	key, err := p2p.LoadNodeKey(fmt.Sprintf("%s/config/priv_validator_key.json", homeDir))
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	var pk [32]byte
	copy(pk[:], key.PrivKey.Bytes())

	return DepInjectOutput{
		BlsSigner: NewBlsSigner(pk),
	}
}
