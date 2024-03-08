// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package root_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/cosmos/cosmos-sdk/x/genutil/client/cli"

	cmdconfig "github.com/berachain/beacon-kit/config/cmd"
	beaconflags "github.com/berachain/beacon-kit/config/flags"
	"github.com/berachain/beacon-kit/examples/beacond/cmd/root"
)

func TestInitCommand(t *testing.T) {
	rootCmd := root.NewRootCmd()
	rootCmd.SetOut(os.NewFile(0, os.DevNull))
	rootCmd.SetArgs([]string{
		"init",           // Test the init cmd
		"BeaconApp-test", // Moniker
		fmt.Sprintf(
			"--%s=%s",
			cli.FlagOverwrite,
			"true",
		), // Overwrite genesis.json
		fmt.Sprintf("--%s", beaconflags.BeaconKitAcceptTos),
	})

	err := svrcmd.Execute(rootCmd, "", cmdconfig.DefaultNodeHome)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestHomeFlagRegistration(t *testing.T) {
	stdout := os.Stdout
	defer func() { os.Stdout = stdout }()
	os.Stdout = os.NewFile(0, os.DevNull)
	homeDir := os.TempDir()

	rootCmd := root.NewRootCmd()
	rootCmd.SetArgs([]string{
		"query",
		fmt.Sprintf("--%s", flags.FlagHome),
		homeDir,
		fmt.Sprintf("--%s", beaconflags.BeaconKitAcceptTos),
	})

	err := svrcmd.Execute(rootCmd, "", cmdconfig.DefaultNodeHome)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	result, err := rootCmd.Flags().GetString(flags.FlagHome)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != homeDir {
		t.Errorf("Expected homeDir to be %s, got %s", homeDir, result)
	}
}
