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

package tos

import (
	"os"
	"path/filepath"
	"strings"

	beaconprompt "github.com/berachain/beacon-kit/mod/node-builder/commands/utils/prompt"
	"github.com/berachain/beacon-kit/mod/node-builder/config/flags"
	"github.com/cockroachdb/errors"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

const (
	// acceptTosFilename is the name of the file that stores the accepted terms
	// of use.
	acceptTosFilename = "tosaccepted"
	// acceptTosPromptTextFormat is the format for the prompt text for accepting
	// the terms of use.
	//nolint:lll
	acceptTosPromptTextFormat = `
%s Terms of Use

By downloading, accessing or using the %s implementation (“%s”), you (referenced herein
as “you” or the “user”) certify that you have read and agreed to the terms and conditions below.

TERMS AND CONDITIONS: %s


Type "accept" to accept these terms and conditions [accept/decline]:`
	// acceptTosPromptErrTextFormat is the error prompt text for accepting the
	// terms of use.
	//nolint:lll
	AcceptTosPromptErrTextFormat = `could not scan text input, if you are trying to run in 
non-interactive environment, you can use the --accept-terms-of-use flag after reading the 
terms and conditions here: 
%s`
	//nolint:lll
	DeclinedErrorString = "you have to accept Terms and Conditions in order to continue"
)

// BuildTosPromptText builds the prompt text for accepting the terms of use.
func BuildTosPromptText(appName, tosLink string) string {
	return aurora.NewAurora(true).
		Sprintf(acceptTosPromptTextFormat, appName, appName, appName, tosLink)
}

// BuildErrorPromptText builds the prompt text for accepting the terms of use.
func BuildErrorPromptText(tosLink string) string {
	return aurora.NewAurora(true).
		Sprintf(AcceptTosPromptErrTextFormat, tosLink)
}

// VerifyTosAcceptedOrPrompt checks if Tos was accepted before or asks to
// accept.
func VerifyTosAcceptedOrPrompt(
	appName, tosLink string,
	clientCtx client.Context,
	cmd *cobra.Command,
) error {
	homedir := clientCtx.HomeDir
	tosFilePath := filepath.Join(homedir, acceptTosFilename)

	if exists, err := afero.Exists(
		afero.NewOsFs(), tosFilePath,
	); err != nil {
		return err
	} else if exists {
		return nil
	}

	if ok, err := cmd.Flags().
		GetBool(flags.BeaconKitAcceptTos); ok && err == nil {
		saveTosAccepted(homedir, cmd)
		return nil
	}

	prompt := &beaconprompt.Prompt{
		Cmd: cmd,
		Text: aurora.NewAurora(true).Bold(BuildTosPromptText(
			appName, tosLink,
		)).String(),
		Default: "decline",
		ValidateFn: func(input string) error {
			if !strings.EqualFold(input, "accept") {
				return errors.New(DeclinedErrorString)
			}
			return nil
		},
	}

	if _, err := prompt.AskAndValidate(); err != nil {
		if err.Error() == DeclinedErrorString {
			return err
		}
		return errors.New(BuildErrorPromptText(tosLink))
	}

	saveTosAccepted(homedir, cmd)
	return nil
}

// saveTosAccepted creates a file when Tos accepted.
func saveTosAccepted(dataDir string, cmd *cobra.Command) {
	fs := afero.NewOsFs()
	dataDirExists, err := afero.DirExists(fs, dataDir)
	if err != nil {
		cmd.PrintErrf("error checking directory: %s\n", dataDir)
	}

	if !dataDirExists {
		if err = fs.MkdirAll(dataDir, os.ModePerm); err != nil {
			cmd.PrintErrf("error creating directory: %s\n", dataDir)
		}
	}

	if err = afero.WriteFile(
		fs, filepath.Join(dataDir, acceptTosFilename),
		[]byte(""), os.ModePerm,
	); err != nil {
		cmd.PrintErrf(
			"error writing %s to file: %s\n",
			flags.BeaconKitAcceptTos,
			filepath.Join(dataDir, acceptTosFilename),
		)
	}
}
