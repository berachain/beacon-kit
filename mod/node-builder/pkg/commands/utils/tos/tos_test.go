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

package tos_test

import (
	"os"
	"testing"

	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/commands/utils/prompt/mocks"
	"github.com/stretchr/testify/mock"
)

// const (
// acceptTosFilename = "tosaccepted"
// )

func TestAcceptTosFlag(t *testing.T) {
	homeDir := makeTempDir(t)
	defer os.RemoveAll(homeDir)

	// rootCmd := root.NewRootCmd()
	// rootCmd.SetOut(os.NewFile(0, os.DevNull))
	// rootCmd.SetArgs([]string{
	// 	"query",
	// 	"--" + flags.FlagHome,
	// 	homeDir,
	// 	"--" + beaconflags.BeaconKitAcceptTos,
	// })

	// err := svrcmd.Execute(rootCmd, "", cmdconfig.DefaultNodeHome)
	// if err != nil {
	// 	t.Errorf("Expected no error, got %v", err)
	// }

	// expectTosAcceptSuccess(t, homeDir)
}

func TestAcceptWithCLI(t *testing.T) {
	homeDir := makeTempDir(t)
	defer os.RemoveAll(homeDir)

	// inputBuffer := bytes.NewReader([]byte("accept\n"))
	// rootCmd := root.NewRootCmd()
	// rootCmd.SetOut(os.NewFile(0, os.DevNull))
	// rootCmd.SetIn(inputBuffer)
	// rootCmd.SetArgs([]string{
	// 	"query",
	// 	"--" + flags.FlagHome,
	// 	homeDir,
	// })

	// if err := svrcmd.Execute(rootCmd, "", cmdconfig.DefaultNodeHome); err !=
	// nil {
	// 	t.Errorf("Expected no error, got %v", err)
	// }

	// expectTosAcceptSuccess(t, homeDir)
}

func TestDeclineWithCLI(t *testing.T) {
	homeDir := makeTempDir(t)
	defer os.RemoveAll(homeDir)

	// inputBuffer := bytes.NewReader([]byte("decline\n"))
	// rootCmd := root.NewRootCmd()
	// rootCmd.SetOut(os.NewFile(0, os.DevNull))
	// rootCmd.SetIn(inputBuffer)
	// rootCmd.SetArgs([]string{
	// 	"query",
	// 	"--" + flags.FlagHome,
	// 	homeDir,
	// })

	// err := svrcmd.Execute(rootCmd, "", cmdconfig.DefaultNodeHome)
	// if err == nil {
	// 	t.Errorf("Expected error, got nil")
	// } else if err.Error() != tos.DeclinedErrorString {
	// 	t.Errorf("Expected %v, got %v", tos.DeclinedErrorString, err)
	// }
}

func TestDeclineWithNonInteractiveCLI(t *testing.T) {
	homeDir := makeTempDir(t)
	defer os.RemoveAll(homeDir)

	// Setup non-interactive reader
	errReader := &mocks.Reader{}
	errReader.On("Read", mock.Anything).Return(0, errors.New("error"))

	// rootCmd := root.NewRootCmd()
	// rootCmd.SetIn(errReader)
	// rootCmd.SetOut(os.NewFile(0, os.DevNull))
	// rootCmd.SetArgs([]string{
	// 	"query",
	// 	"--" + flags.FlagHome,
	// 	homeDir,
	// })

	// err := svrcmd.Execute(rootCmd, "", cmdconfig.DefaultNodeHome)
	// if err == nil {
	// 	t.Errorf("Expected error, got nil")
	// } else if !strings.Contains(err.Error(), tos.BuildErrorPromptText("")) {
	// 	t.Errorf("Expected %v, got %v", tos.BuildErrorPromptText(""), err)
	// }
}

// func expectTosAcceptSuccess(t *testing.T, homeDir string) {
// 	if ok := file.Exists(filepath.Join(homeDir, acceptTosFilename)); !ok {
// 		t.Errorf("Expected tosaccepted file to exist in %s", homeDir)
// 	}
// }

func makeTempDir(t *testing.T) string {
	homeDir, err := os.MkdirTemp("", "beacond-test-*")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	return homeDir
}
