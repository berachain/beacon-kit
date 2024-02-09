// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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

// import (
// 	"testing"
// 	"github.com/cosmos/cosmos-sdk/client"
// 	"github.com/itsdevbear/bolaris/config/tos"
// 	"github.com/itsdevbear/bolaris/examples/beacond/cmd/root"
// 	"github.com/prysmaticlabs/prysm/v4/testing/require".
// )

// func TestVerifyTosAcceptedOrPrompt(t *testing.T) {
// 	rootCmd := root.NewRootCmd()
// 	rootCmd.SetOut(os.NewFile(0, os.DevNull))
// 	appName := "TestApp"
// 	tosLink := "https://example.com/tos"
// 	clientCtx := client.Context{}.WithInput(os.Stdin)

// 	// replacing stdin
// tmpfile, err := os.CreateTemp("", "tmp")
// require.NoError(t, err)
// origStdin := os.Stdin
// os.Stdin = tmpfile
// defer func() { os.Stdin = origStdin }()

// // userprompt decline
// _, err = tmpfile.Write([]byte("decline"))
// require.NoError(t, err)
// _, err = tmpfile.Seek(0, 0)
// 	require.NoError(t, err)
// 	require.ErrorContains(t, "you have to
// 	accept Terms and Conditions", tos.VerifyTosAccep
// 	tedOrPrompt(appName, tosLink, clientCtx, rootCmd))

// 	// userprompt accept
// 	err = tmpfile.Truncate(0)
// 	require.NoError(t, err)
// 	_, err = tmpfile.Seek(0, 0)
// 	require.NoError(t, err)
// 	_, err = tmpfile.Write([]byte("accept"))
// 	require.NoError(t, err)
// 	_, err = tmpfile.Seek(0, 0)
// 	require.NoError(t, err)
// 	require.NoError(t, VerifyTosAcceptedOrPrompt(context))
// 	require.NoError(t, os.Remove(file
// 	path.Join(context.String(cmd.DataDirFlag.Name), acceptTosFilename)))

// 	require.NoError(t, tmpfile.Close())
// 	require.NoError(t, os.Remove(tmpfile.Name()))

// 	// saved in file
// 	require.NoError(t, os.WriteFile(filepath.Jo
// 	in(context.String(cmd.DataDirFlag.Name), acceptTosFilename), []byte(""), 0666))
// 	require.NoError(t, VerifyTosAcceptedOrPrompt(context))
// 	require.NoError(t, os.RemoveAll(context.String(cmd.DataDirFlag.Name)))

// 	// flag is set
// 	set.Bool(cmd.AcceptTosFlag.Name, true, "")
// 	require.NoError(t, VerifyTosAcceptedOrPrompt(context))
// 	require.NoError(t, os.RemoveAll(context.String(cmd.DataDirFlag.Name)))
// }
