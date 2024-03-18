package e2e_test

import (
	"encoding/hex"
	"log"
	"os/exec"

	"github.com/ethereum/go-ethereum/crypto"
)

// TestForgeScriptExecution tests the execution of a forge script
// against the beacon-kit network.
func (s *BeaconKitE2ESuite) TestForgeScriptExecution() {
	url := s.KurtosisE2ESuite.JSONRPCBalancer().URL()
	pk := hex.EncodeToString(
		crypto.FromECDSA(s.GenesisAccount().PrivateKey()),
	)

	// Change directory to /contracts/ before executing the command
	cmdStr := "cd ../contracts && " +
		"forge build && " +
		"forge script ./script/DeployAndCallERC20.s.sol " +
		"--broadcast --rpc-url=" + url + " " +
		"--private-key=" + pk

	// Execute the command
	cmd := exec.Command("bash", "-c", cmdStr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf(
			"Failed to execute command: %s, with the error: %s",
			cmdStr,
			err,
		)
	}

	s.Logger().Info("Output: %s\n", string(output))
}
