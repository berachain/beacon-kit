package config_test

import (
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sgconfig "github.com/itsdevbear/bolaris/cosmos/config"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Configuration", func() {
	It("should set CoinType", func() {
		config := sdk.GetConfig()

		Expect(int(config.GetCoinType())).To(Equal(sdk.CoinType))
		Expect(config.GetFullBIP44Path()).To(Equal(sdk.FullFundraiserPath))

		sgconfig.SetupCosmosConfig()

		Expect(int(config.GetCoinType())).To(Equal(int(60)))
		Expect(config.GetCoinType()).To(Equal(sdk.GetConfig().GetCoinType()))
		Expect(config.GetFullBIP44Path()).To(Equal(sdk.GetConfig().GetFullBIP44Path()))
	})

	It("should generate HD path", func() {
		params := *hd.NewFundraiserParams(0, 60, 0)
		hdPath := params.String()

		Expect(hdPath).To(Equal("m/44'/60'/0'/0/0"))
		Expect(hdPath).To(Equal(60))
	})
})
