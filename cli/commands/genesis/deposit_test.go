package genesis

import (
	"testing"

	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/node-core/components/signer"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/stretchr/testify/require"
)

func TestGenesisDeposit(t *testing.T) {
	chainSpec, err := spec.DevnetChainSpec()
	require.NoError(t, err)

	cometConfig := cmtcfg.DefaultConfig()
	depositAmount := math.Gwei(32000000000)
	withdrawalAdress := common.NewExecutionAddressFromHex("0x981114102592310C347E61368342DDA67017bf84")
	outputDocument := ""

	blsSigner := signer.NewBLSSigner(privValKeyFile, privValStateFile)
	err = AddGenesisDeposit(chainSpec, cometConfig, blsSigner, depositAmount, withdrawalAdress, outputDocument)
	require.NoError(t, err)
}
