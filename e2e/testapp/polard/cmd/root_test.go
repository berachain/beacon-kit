package cmd_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/cosmos/cosmos-sdk/x/genutil/client/cli"

	testapp "github.com/itsdevbear/bolaris/e2e/testapp"
	"github.com/itsdevbear/bolaris/e2e/testapp/polard/cmd"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "e2e/testapp/polard/cmd")
}

var _ = Describe("Init command", func() {
	It("should initialize the app with given options", func() {
		stdout := os.Stdout
		defer func() { os.Stdout = stdout }()
		os.Stdout = os.NewFile(0, os.DevNull)
		rootCmd := cmd.NewRootCmd()
		rootCmd.SetArgs([]string{
			"init",        // Test the init cmd
			"simapp-test", // Moniker
			fmt.Sprintf("--%s=%s", cli.FlagOverwrite, "true"), // Overwrite genesis.json
		})

		err := svrcmd.Execute(rootCmd, "", testapp.DefaultNodeHome)
		Expect(err).ToNot(HaveOccurred())
	})
})

var _ = Describe("Home flag registration", func() {
	It("should set home directory correctly", func() {
		// Redirect standard out to null
		stdout := os.Stdout
		defer func() { os.Stdout = stdout }()
		os.Stdout = os.NewFile(0, os.DevNull)
		homeDir := os.TempDir()

		rootCmd := cmd.NewRootCmd()
		rootCmd.SetArgs([]string{
			"query",
			fmt.Sprintf("--%s", flags.FlagHome),
			homeDir,
		})

		err := svrcmd.Execute(rootCmd, "", testapp.DefaultNodeHome)
		Expect(err).ToNot(HaveOccurred())

		result, err := rootCmd.Flags().GetString(flags.FlagHome)
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(Equal(homeDir))
	})
})
