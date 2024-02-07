package tos

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/itsdevbear/bolaris/config/flags"
	"github.com/logrusorgru/aurora"
	"github.com/prysmaticlabs/prysm/v4/io/file"
	"github.com/prysmaticlabs/prysm/v4/io/prompt"
	"github.com/spf13/cobra"
)

const (
	acceptTosFilename         = "tosaccepted"
	acceptTosPromptTextFormat = `
%s Terms of Use

By downloading, accessing or using the %s implementation (“%s”), you (referenced herein
as “you” or the “user”) certify that you have read and agreed to the terms and conditions below.

TERMS AND CONDITIONS: %s


Type "accept" to accept this terms and conditions [accept/decline]:`

	acceptTosPromptErrTextFormat = `could not scan text input, if you are trying to run in non-interactive environment, you
can use the --accept-terms-of-use flag after reading the terms and conditions here: 
%s`
)

var (
	au = aurora.NewAurora(true)
)

// BuildTosPromptText builds the prompt text for accepting the terms of use.
func BuildTosPromptText(appName, tosLink string) string {
	return au.Sprintf(acceptTosPromptTextFormat, appName, appName, tosLink, tosLink)
}

// BuildErrorPromptText builds the prompt text for accepting the terms of use.
func BuildErrorPromptText(tosLink string) string {
	return au.Sprintf(acceptTosPromptErrTextFormat, tosLink)
}

// VerifyTosAcceptedOrPrompt checks if Tos was accepted before or asks to accept.
func VerifyTosAcceptedOrPrompt(
	appName, tosLink string,
	clientCtx client.Context,
	cmd *cobra.Command,
) error {
	homedir := clientCtx.HomeDir
	tosFilePath := filepath.Join(homedir, acceptTosFilename)
	if file.Exists(tosFilePath) {
		return nil
	}

	if ok, err := cmd.Flags().GetBool(flags.BeaconKitAcceptTos); ok && err == nil {
		saveTosAccepted(homedir, cmd)
		return nil
	}

	input, err := prompt.DefaultPrompt(au.Bold(BuildTosPromptText(
		appName, tosLink,
	)).String(), "decline")
	if err != nil {
		return errors.New(BuildErrorPromptText(tosLink))
	}

	if !strings.EqualFold(input, "accept") {
		return errors.New("you have to accept Terms and Conditions in order to continue")
	}

	saveTosAccepted(homedir, cmd)
	return nil
}

// saveTosAccepted creates a file when Tos accepted.
func saveTosAccepted(dataDir string, cmd *cobra.Command) {
	dataDirExists, err := file.HasDir(dataDir)
	if err != nil {
		cmd.PrintErrf("error checking directory: %s\n", dataDir)
	}
	if !dataDirExists {
		if err := file.MkdirAll(dataDir); err != nil {
			cmd.PrintErrf("error creating directory: %s\n", dataDir)
		}
	}
	if err := file.WriteFile(filepath.Join(dataDir, acceptTosFilename), []byte("")); err != nil {
		cmd.PrintErrf("error writing %s to file: %s\n", flags.BeaconKitAcceptTos,
			filepath.Join(dataDir, acceptTosFilename))
	}
}
