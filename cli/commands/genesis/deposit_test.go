package genesis_test

import (
	"path"
	"testing"

	"github.com/berachain/beacon-kit/cli/commands/genesis"
	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/node-core/components/signer"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/testing/utils"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cometbft/cometbft/crypto/bls12381"
	"github.com/cometbft/cometbft/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisDeposit(t *testing.T) {
	homeDir := utils.MakeTempHomeDir(t)
	t.Log("Home folder:", homeDir)
	defer utils.DeleteTempHomeDir(t, homeDir)

	chainSpec, err := spec.DevnetChainSpec()
	require.NoError(t, err)

	cometConfig := cmtcfg.DefaultConfig()

	cometConfig.SetRoot(homeDir)

	// Forces Comet to Create it
	cometConfig.NodeKey = "nodekey.json"

	depositAmount := math.Gwei(32000000000)
	withdrawalAdress := common.NewExecutionAddressFromHex("0x981114102592310C347E61368342DDA67017bf84")
	outputDocument := ""

	blsSigner := signer.BLSSigner{PrivValidator: types.NewMockPVWithKeyType(bls12381.KeyType)}

	err = genesis.AddGenesisDeposit(chainSpec, cometConfig, blsSigner, depositAmount, withdrawalAdress, outputDocument)
	require.NoError(t, err)

	require.FileExists(t, path.Join(homeDir, "nodekey.json"))
	require.FileExists(t, path.Join(homeDir, "data", "priv_validator_state.json"))
	require.FileExists(t, path.Join(homeDir, "config", "priv_validator_key.json"))
	require.DirExists(t, path.Join(homeDir, "config", "premined-deposits"))
}
