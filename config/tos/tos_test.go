package tos_test

import (
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/itsdevbear/bolaris/config/tos"
	"github.com/itsdevbear/bolaris/examples/beacond/cmd/root"
	"github.com/prysmaticlabs/prysm/v4/testing/require"
)

func TestVerifyTosAcceptedOrPrompt(t *testing.T) {
	rootCmd := root.NewRootCmd()
	rootCmd.SetOut(os.NewFile(0, os.DevNull))
	appName := "TestApp"
	tosLink := "https://example.com/tos"
	clientCtx := client.Context{}.WithInput(os.Stdin)

	// replacing stdin
	tmpfile, err := os.CreateTemp("", "tmp")
	require.NoError(t, err)
	origStdin := os.Stdin
	os.Stdin = tmpfile
	defer func() { os.Stdin = origStdin }()

	// userprompt decline
	_, err = tmpfile.Write([]byte("decline"))
	require.NoError(t, err)
	_, err = tmpfile.Seek(0, 0)
	require.NoError(t, err)
	require.ErrorContains(t, "you have to accept Terms and Conditions", tos.VerifyTosAcceptedOrPrompt(appName, tosLink, clientCtx, rootCmd))

	// userprompt accept
	err = tmpfile.Truncate(0)
	require.NoError(t, err)
	_, err = tmpfile.Seek(0, 0)
	require.NoError(t, err)
	_, err = tmpfile.Write([]byte("accept"))
	require.NoError(t, err)
	_, err = tmpfile.Seek(0, 0)
	require.NoError(t, err)
	// require.NoError(t, VerifyTosAcceptedOrPrompt(context))
	// require.NoError(t, os.Remove(filepath.Join(context.String(cmd.DataDirFlag.Name), acceptTosFilename)))

	// require.NoError(t, tmpfile.Close())
	// require.NoError(t, os.Remove(tmpfile.Name()))

	// // saved in file
	// require.NoError(t, os.WriteFile(filepath.Join(context.String(cmd.DataDirFlag.Name), acceptTosFilename), []byte(""), 0666))
	// require.NoError(t, VerifyTosAcceptedOrPrompt(context))
	// require.NoError(t, os.RemoveAll(context.String(cmd.DataDirFlag.Name)))

	// // flag is set
	// set.Bool(cmd.AcceptTosFlag.Name, true, "")
	// require.NoError(t, VerifyTosAcceptedOrPrompt(context))
	// require.NoError(t, os.RemoveAll(context.String(cmd.DataDirFlag.Name)))
}
