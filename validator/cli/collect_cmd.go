package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"
)

const flagGenTxDir = "gentx-dir"

// CollectGenTxsCmd - return the cobra command to collect genesis transactions
func CollectGenTxsCmd(genBalIterator types.GenesisBalancesIterator, validator types.MessageValidator, valAddrCodec runtime.ValidatorAddressCodec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "customcollect-gentxs",
		Short: "Collect genesis txs and output a genesis.json file",
		RunE: func(cmd *cobra.Command, _ []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			clientCtx := client.GetClientContextFromCmd(cmd)
			cdc := clientCtx.Codec

			nodeID, valPubKey, err := genutil.InitializeNodeValidatorFiles(config)
			if err != nil {
				return errors.Wrap(err, "failed to initialize node validator files")
			}

			appGenesis, err := types.AppGenesisFromFile(config.GenesisFile())
			if err != nil {
				return errors.Wrap(err, "failed to read genesis doc from file")
			}

			genTxDir, _ := cmd.Flags().GetString(flagGenTxDir)
			genTxsDir := genTxDir
			if genTxsDir == "" {
				genTxsDir = filepath.Join(config.RootDir, "config", "gentx")
			}

			toPrint := newPrintInfo(config.Moniker, appGenesis.ChainID, nodeID, genTxsDir, json.RawMessage(""))
			initCfg := types.NewInitConfig(appGenesis.ChainID, genTxsDir, nodeID, valPubKey)

			appMessage, err := GenAppStateFromConfig(cdc, clientCtx.TxConfig, config, initCfg, appGenesis, genBalIterator, validator, valAddrCodec)
			if err != nil {
				return errors.Wrap(err, "failed to get genesis app state from config")
			}

			toPrint.AppMessage = appMessage

			return displayInfo(toPrint)
		},
	}

	cmd.Flags().String(flagGenTxDir, "", "override default \"gentx\" directory from which collect and execute genesis transactions; default [--home]/config/gentx/")

	return cmd
}

type printInfo struct {
	Moniker    string          `json:"moniker" yaml:"moniker"`
	ChainID    string          `json:"chain_id" yaml:"chain_id"`
	NodeID     string          `json:"node_id" yaml:"node_id"`
	GenTxsDir  string          `json:"gentxs_dir" yaml:"gentxs_dir"`
	AppMessage json.RawMessage `json:"app_message" yaml:"app_message"`
}

func newPrintInfo(moniker, chainID, nodeID, genTxsDir string, appMessage json.RawMessage) printInfo {
	return printInfo{
		Moniker:    moniker,
		ChainID:    chainID,
		NodeID:     nodeID,
		GenTxsDir:  genTxsDir,
		AppMessage: appMessage,
	}
}

func displayInfo(info printInfo) error {
	out, err := json.MarshalIndent(info, "", " ")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(os.Stderr, "%s\n", out)

	return err
}
