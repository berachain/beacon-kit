package components

import (
	"cosmossdk.io/core/address"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
)

func ProvideAddressCodec() (address.Codec, address.ValidatorAddressCodec, address.ConsensusAddressCodec) {
	addrCdc := addresscodec.NewBech32Codec("bera")
	return addrCdc, addrCdc, addrCdc
}
