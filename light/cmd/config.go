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

package cmd

import (
	"net/url"
	"strings"

	"github.com/berachain/beacon-kit/light/app"
	"github.com/berachain/beacon-kit/light/mod/provider"
	"github.com/berachain/beacon-kit/light/mod/provider/comet"
	engineclient "github.com/berachain/beacon-kit/mod/execution/client"
	"github.com/berachain/beacon-kit/mod/node-builder/commands/utils/prompt"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/spf13/cobra"
)

// ConfigFromCmd returns a new light node configuration from the given command.
func ConfigFromCmd(
	logger log.Logger,
	chainID string,
	cmd *cobra.Command,
) (*app.Config, error) {
	tl, err := cmd.Flags().GetString(trustLevel)
	if err != nil {
		return nil, err
	}
	directory, err := cmd.Flags().GetString(dir)
	if err != nil {
		return nil, err
	}
	listeningAddr, err := cmd.Flags().GetString(listenAddr)
	if err != nil {
		return nil, err
	}
	sequential, err := cmd.Flags().GetBool(sequential)
	if err != nil {
		return nil, err
	}
	trustedHeight, err := cmd.Flags().GetInt64(trustedHeight)
	if err != nil {
		return nil, err
	}
	trustedHash, err := cmd.Flags().GetBytesHex(trustedHash)
	if err != nil {
		return nil, err
	}
	trustingPeriod, err := cmd.Flags().GetDuration(trustingPeriod)
	if err != nil {
		return nil, err
	}
	maxOpenConnections, err := cmd.Flags().GetInt(maxOpenConnections)
	if err != nil {
		return nil, err
	}
	witnesses, err := cmd.Flags().GetString(witnessAddrsJoined)
	if err != nil {
		return nil, err
	}
	pAddr, err := cmd.Flags().GetString(primaryAddr)
	if err != nil {
		return nil, err
	}

	var witnessesAddrs []string
	if witnessAddrsJoined != "" {
		witnessesAddrs = strings.Split(witnesses, ",")
	}

	engine, err := cmd.Flags().GetString(engineURL)
	if err != nil {
		return nil, err
	}

	engineCfg := engineclient.DefaultConfig()
	engineCfg.RPCDialURL, err = url.Parse(engine)
	if err != nil {
		return nil, err
	}
	engineCfg.JWTSecretPath, err = cmd.Flags().GetString(jwtSecretPath)
	if err != nil {
		return nil, err
	}

	return app.NewConfig(
		comet.NewConfig(
			logger, chainID, trustingPeriod,
			trustedHeight, trustedHash, tl,
			listeningAddr, sequential,
			pAddr, witnessesAddrs,
			directory, maxOpenConnections,
			NewConfirmationFunc(cmd),
		),
		provider.NewConfig(chainID, listeningAddr, "/websocket"),
		&engineCfg,
	), nil
}

// NewConfirmationFunc returns a function that prompts the user for
// confirmation.
func NewConfirmationFunc(cmd *cobra.Command) func(string) bool {
	p := &prompt.Prompt{
		Cmd:        cmd,
		Default:    "n",
		ValidateFn: prompt.ValidateYesOrNo,
	}

	return func(action string) bool {
		p.Text = action
		for {
			input, err := p.AskAndValidate()
			if err != nil {
				p.Cmd.Println(err)
				continue
			}
			if input == "y" || input == "Y" {
				return true
			}
			return false
		}
	}
}
