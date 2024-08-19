package cometbft

import (
	"context"

	"github.com/cometbft/cometbft/p2p"
	pvm "github.com/cometbft/cometbft/privval"
	cmtversion "github.com/cometbft/cometbft/version"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	cmtjson "github.com/cometbft/cometbft/libs/json"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/server"
)

// StatusCommand returns the command to return the status of the network.
func StatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Query remote node for status",
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			status, err := cmtservice.GetNodeStatus(
				context.Background(), clientCtx,
			)
			if err != nil {
				return err
			}

			output, err := cmtjson.Marshal(status)
			if err != nil {
				return err
			}

			outputFormat, _ := cmd.Flags().GetString(flags.FlagOutput)
			if outputFormat == flags.OutputFormatJSON {
				clientCtx = clientCtx.WithOutputFormat(flags.OutputFormatJSON)
			}

			return clientCtx.PrintRaw(output)
		},
	}

	cmd.Flags().StringP(flags.FlagNode, "n",
		"tcp://localhost:26657", "Node to connect to")
	cmd.Flags().StringP(flags.FlagOutput, "o", "json",
		"Output format (text|json)")

	return cmd
}

// ShowNodeIDCmd - ported from CometBFT, dump node ID to stdout
func ShowNodeIDCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show-node-id",
		Short: "Show this node's ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			cfg := serverCtx.Config

			nodeKey, err := p2p.LoadNodeKey(cfg.NodeKeyFile())
			if err != nil {
				return err
			}

			cmd.Println(nodeKey.ID())
			return nil
		},
	}
}

// ShowValidatorCmd - ported from CometBFT, show this node's validator info
func ShowValidatorCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "show-validator",
		Short: "Show this node's CometBFT validator info",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			cfg := serverCtx.Config

			privValidator := pvm.LoadFilePV(
				cfg.PrivValidatorKeyFile(),
				cfg.PrivValidatorStateFile(),
			)
			pk, err := privValidator.GetPubKey()
			if err != nil {
				return err
			}

			sdkPK, err := cryptocodec.FromCmtPubKeyInterface(pk)
			if err != nil {
				return err
			}

			clientCtx := client.GetClientContextFromCmd(cmd)
			bz, err := clientCtx.Codec.MarshalInterfaceJSON(sdkPK)
			if err != nil {
				return err
			}

			cmd.Println(string(bz))
			return nil
		},
	}

	return &cmd
}

// ShowAddressCmd - show this node's validator address
func ShowAddressCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-address",
		Short: "Shows this node's CometBFT validator consensus address",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			cfg := serverCtx.Config

			privValidator := pvm.LoadFilePV(
				cfg.PrivValidatorKeyFile(),
				cfg.PrivValidatorStateFile(),
			)

			valConsAddr := privValidator.GetAddress()

			cmd.Println(valConsAddr.String())
			return nil
		},
	}

	return cmd
}

// VersionCmd prints CometBFT and ABCI version numbers.
//
//nolint:lll // long description
func VersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print CometBFT libraries' version",
		Long:  "Print protocols' and libraries' version numbers against which this app has been compiled.",
		RunE: func(cmd *cobra.Command, args []string) error {
			bs, err := yaml.Marshal(&struct {
				CometBFT      string
				ABCI          string
				BlockProtocol uint64
				P2PProtocol   uint64
			}{
				CometBFT:      cmtversion.CMTSemVer,
				ABCI:          cmtversion.ABCIVersion,
				BlockProtocol: cmtversion.BlockProtocol,
				P2PProtocol:   cmtversion.P2PProtocol,
			})
			if err != nil {
				return err
			}

			cmd.Println(string(bs))
			return nil
		},
	}
}

// TODO: implement
// BootstrapStateCmd - ported from CometBFT, bootstrap CometBFT state at an
// arbitrary block height using a light client
