package main

import (
	"os"

	"cosmossdk.io/log"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/itsdevbear/bolaris/cosmos/config"
	testapp "github.com/itsdevbear/bolaris/e2e/testapp"
	"github.com/itsdevbear/bolaris/e2e/testapp/polard/cmd"
)

func main() {
	config.SetupCosmosConfig()
	rootCmd := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, "", testapp.DefaultNodeHome); err != nil {
		log.NewLogger(rootCmd.OutOrStderr()).Error("failure when running app", "err", err)
		os.Exit(1)
	}
}
