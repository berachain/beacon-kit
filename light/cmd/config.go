package cmd

// import (
// 	"strings"

// 	"github.com/cometbft/cometbft/libs/log"

// 	"github.com/spf13/cobra"

// 	"github.com/berachain/beacon-kit/io/cli/prompt"
// 	"github.com/berachain/beacon-kit/light/app"
// 	"github.com/berachain/beacon-kit/light/provider"
// 	"github.com/berachain/beacon-kit/light/provider/comet"
// )

// func ConfigFromCmd(logger log.Logger, chainID string, cmd *cobra.Command) (*app.Config, error) {
// 	tl, err := cmd.Flags().GetString(trustLevel)
// 	if err != nil {
// 		return nil, err
// 	}
// 	directory, err := cmd.Flags().GetString(dir)
// 	if err != nil {
// 		return nil, err
// 	}
// 	listeningAddr, err := cmd.Flags().GetString(listenAddr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// engineURL, err := cmd.Flags().GetString(engineURL)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	sequential, err := cmd.Flags().GetBool(sequential)
// 	if err != nil {
// 		return nil, err
// 	}
// 	trustedHeight, err := cmd.Flags().GetInt64(trustedHeight)
// 	if err != nil {
// 		return nil, err
// 	}
// 	trustedHash, err := cmd.Flags().GetBytesHex(trustedHash)
// 	if err != nil {
// 		return nil, err
// 	}
// 	trustingPeriod, err := cmd.Flags().GetDuration(trustingPeriod)
// 	if err != nil {
// 		return nil, err
// 	}
// 	maxOpenConnections, err := cmd.Flags().GetInt(maxOpenConnections)
// 	if err != nil {
// 		return nil, err
// 	}
// 	witnesses, err := cmd.Flags().GetString(witnessAddrsJoined)
// 	if err != nil {
// 		return nil, err
// 	}
// 	pAddr, err := cmd.Flags().GetString(primaryAddr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// pEthAddr, err := cmd.Flags().GetString(primaryEthAddr)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	var witnessesAddrs []string
// 	if witnessAddrsJoined != "" {
// 		witnessesAddrs = strings.Split(witnesses, ",")
// 	}

// 	return app.NewConfig(
// 		comet.NewConfig(
// 			logger, chainID, trustingPeriod,
// 			trustedHeight, trustedHash, tl,
// 			listeningAddr, sequential,
// 			pAddr, witnessesAddrs,
// 			directory, maxOpenConnections,
// 			NewConfirmationFunc(cmd),
// 		),
// 		provider.NewConfig(chainID, listeningAddr, "/websocket"),
// 	), nil
// }

// // NewConfirmationFunc returns a function that prompts the user for confirmation.
// func NewConfirmationFunc(cmd *cobra.Command) func(string) bool {
// 	p := &prompt.Prompt{
// 		Cmd:        cmd,
// 		Default:    "n",
// 		ValidateFn: prompt.ValidateYesOrNo,
// 	}

// 	return func(action string) bool {
// 		p.Text = action
// 		for {
// 			input, err := p.AskAndValidate()
// 			if err != nil {
// 				p.Cmd.Println(err)
// 				continue
// 			}
// 			if input == "y" || input == "Y" {
// 				return true
// 			}
// 			return false
// 		}
// 	}
// }
