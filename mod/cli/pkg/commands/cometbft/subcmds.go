package cometbft

import (
	"github.com/berachain/beacon-kit/mod/cli/pkg/context"
	"github.com/cometbft/cometbft/p2p"
	pvm "github.com/cometbft/cometbft/privval"
	cmtversion "github.com/cometbft/cometbft/version"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"github.com/cosmos/cosmos-sdk/client"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ShowNodeIDCmd - ported from CometBFT, dump node ID to stdout
func ShowNodeIDCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show-node-id",
		Short: "Show this node's ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := context.GetServerContextFromCmd(cmd)
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
			serverCtx := context.GetServerContextFromCmd(cmd)
			cfg := serverCtx.Config

			privValidator := pvm.LoadFilePV(cfg.PrivValidatorKeyFile(),
				cfg.PrivValidatorStateFile())
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
			serverCtx := context.GetServerContextFromCmd(cmd)
			cfg := serverCtx.Config

			privValidator := pvm.LoadFilePV(cfg.PrivValidatorKeyFile(),
				cfg.PrivValidatorStateFile())

			valConsAddr := (sdk.ConsAddress)(privValidator.GetAddress())

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
