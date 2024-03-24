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

package cli

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"cosmossdk.io/core/address"
	authclient "cosmossdk.io/x/auth/client"
	"cosmossdk.io/x/staking/client/cli"
	beacontypes "github.com/berachain/beacon-kit/runtime/modules/beacon/types"
	"github.com/cockroachdb/errors"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/cobra"
)

// GenTxCmd builds the application's gentx command.
//
//nolint:funlen,gocognit,maintidx // todo fix
func GenTxCmd(
	mm *module.Manager,
	_ client.TxEncodingConfig,
	_ types.GenesisBalancesIterator,
	valAdddressCodec address.Codec,
) *cobra.Command {
	ipDefault, err := server.ExternalIP()
	if err != nil {
		ipDefault = ""
	}
	fsCreateValidator, defaultsDesc := cli.CreateValidatorMsgFlagSet(ipDefault)

	cmd := &cobra.Command{
		Use:   "customgentx [key_name] [amount]",
		Short: "Generate a genesis tx carrying a self delegation",
		Args:  cobra.ExactArgs(2), //nolint:gomnd // there are two arguments.
		Long: fmt.Sprintf(
			`Generate a genesis transaction that creates a validator 
with a self-delegation, that is signed by the key in the Keyring referenced
by a given name. 
A node ID and consensus pubkey may optionally be provided. If they are omitted, 
they will be retrieved from the priv_validator.json
file. The following default parameters are included:
    %s

Example:
$ %s gentx my-key-name 1000000stake --home=/path/to/home --chain-id=test-1
`,
			defaultsDesc,
			version.AppName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				nodeID    string
				valPubKey cryptotypes.PubKey
				serverCtx = server.GetServerContextFromCmd(cmd)
				config    = serverCtx.Config
			)

			var clientCtx client.Context
			clientCtx, err = client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			nodeID, valPubKey, err = genutil.InitializeNodeValidatorFiles(
				serverCtx.Config,
			)
			if err != nil {
				return errors.Wrap(
					err,
					"failed to initialize node validator files",
				)
			}

			// read --nodeID, if empty take it from priv_validator.json
			var nodeIDString string
			if nodeIDString, err = cmd.Flags().GetString(
				cli.FlagNodeID,
			); err != nil {
				return errors.Wrap(err, "failed to get node ID")
			} else if nodeIDString != "" {
				nodeID = nodeIDString
			}

			// read --pubkey, if empty take it from priv_validator.json
			var pkStr string
			if pkStr, err = cmd.Flags().GetString(cli.FlagPubKey); err != nil ||
				pkStr != "" {
				if err = clientCtx.Codec.UnmarshalInterfaceJSON(
					[]byte(pkStr), &valPubKey); err != nil {
					return errors.Wrap(
						err,
						"failed to unmarshal validator public key",
					)
				}
			}

			// read the genesis doc from the file
			var appGenesis *types.AppGenesis
			appGenesis, err = types.AppGenesisFromFile(config.GenesisFile())
			if err != nil {
				return errors.Wrapf(
					err,
					"failed to read genesis doc file %s",
					config.GenesisFile(),
				)
			}

			var genesisState map[string]json.RawMessage
			if err = json.Unmarshal(appGenesis.AppState, &genesisState); err != nil {
				return errors.Wrap(err, "failed to unmarshal genesis state")
			}

			if err = mm.ValidateGenesis(genesisState); err != nil {
				return errors.Wrap(err, "failed to validate genesis state")
			}

			inBuf := bufio.NewReader(cmd.InOrStdin())

			name := args[0]
			var key *keyring.Record
			key, err = clientCtx.Keyring.Key(name)
			if err != nil {
				return errors.Wrapf(
					err,
					"failed to fetch '%s' from the keyring",
					name,
				)
			}

			moniker := config.Moniker

			var m string
			if m, err = cmd.Flags().GetString(cli.FlagMoniker); err != nil ||
				m != "" {
				moniker = m
			}

			// set flags for creating a gentx
			var createValCfg cli.TxCreateValidatorConfig
			createValCfg, err = cli.PrepareConfigForTxCreateValidator(
				cmd.Flags(),
				moniker,
				nodeID,
				appGenesis.ChainID,
				valPubKey,
			)
			if err != nil {
				return errors.Wrap(
					err,
					"error creating configuration to create validator msg",
				)
			}

			amount := args[1]
			// coins, err := sdk.ParseCoinsNormalized(amount)
			// if err != nil {
			// 	return errors.Wrap(err, "failed to parse coins")
			// }
			// addr, err := key.GetAddress()
			// if err != nil {
			// 	return err
			// }
			// err = genutil.ValidateAccountInGenesis(
			// 	genesisState,
			// 	genBalIterator,
			// 	addr,
			// 	coins,
			// 	cdc,
			// )
			// if err != nil {
			// 	return errors.Wrap(err, "failed to validate account in genesis")
			// }

			var txFactory tx.Factory
			txFactory, err = tx.NewFactoryCLI(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			var pub sdk.AccAddress
			pub, err = key.GetAddress()
			if err != nil {
				return err
			}
			clientCtx = clientCtx.WithInput(inBuf).WithFromAddress(pub)

			// The following line comes from a discrepancy between the `gentx`
			// and `create-validator` commands:
			// - `gentx` expects amount as an arg,
			// - `create-validator` expects amount as a required flag.
			// ref: https://github.com/cosmos/cosmos-sdk/issues/8251
			// Since gentx doesn't set the amount flag (which `create-validator`
			// reads from), we copy the amount arg into the valCfg directly.
			//
			// Ideally, the `create-validator` command should take a validator
			// config file instead of so many flags.
			// ref: https://github.com/cosmos/cosmos-sdk/issues/8177
			createValCfg.Amount = amount

			// create a 'create-validator' message
			var (
				txBldr tx.Factory
				msg    sdk.Msg
			)
			txBldr, msg, err = BuildCreateValidatorMsg(
				clientCtx,
				createValCfg,
				txFactory,
				true,
				valAdddressCodec,
			)
			if err != nil {
				return errors.Wrap(
					err,
					"failed to build create-validator message",
				)
			}

			// write the unsigned transaction to the buffer
			w := bytes.NewBuffer([]byte{})
			clientCtx = clientCtx.WithOutput(w)

			// if m, ok := msg.(sdk.HasValidateBasic); ok {
			// 	if err := m.ValidateBasic(); err != nil {
			// 		return err
			// 	}
			// }

			if err = txBldr.PrintUnsignedTx(clientCtx, msg); err != nil {
				return errors.Wrap(err, "failed to print unsigned std tx")
			}

			// read the transaction
			var stdTx sdk.Tx
			stdTx, err = readUnsignedGenTxFile(clientCtx, w)
			if err != nil {
				return errors.Wrap(err, "failed to read unsigned gen tx file")
			}

			// sign the transaction and write it to the output file
			var txBuilder client.TxBuilder
			txBuilder, err = clientCtx.TxConfig.WrapTxBuilder(stdTx)
			if err != nil {
				return fmt.Errorf("error creating tx builder: %w", err)
			}

			if err = authclient.SignTx(
				txFactory, clientCtx, name, txBuilder, true, true); err != nil {
				return errors.Wrap(err, "failed to sign std tx")
			}

			var outputDocument string
			outputDocument, err = cmd.Flags().
				GetString(flags.FlagOutputDocument)
			if err != nil || outputDocument == "" {
				outputDocument, err = makeOutputFilepath(config.RootDir, nodeID)
				if err != nil {
					return errors.Wrap(err, "failed to create output file path")
				}
			}

			if err = writeSignedGenTx(
				clientCtx, outputDocument, txBuilder.GetTx()); err != nil {
				return errors.Wrap(err, "failed to write signed gen tx")
			}

			cmd.PrintErrf("Genesis transaction written to %q\n", outputDocument)
			return nil
		},
	}

	cmd.Flags().
		String(flags.FlagOutputDocument, "",
			"Write the genesis transaction JSON document "+
				"to the given file instead of the default location")
	cmd.Flags().AddFlagSet(fsCreateValidator)
	flags.AddTxFlagsToCmd(cmd)
	err = cmd.Flags().
		MarkHidden(flags.FlagOutput)
	// signing makes sense to output only json
	if err != nil {
		panic(err)
	}

	return cmd
}

func makeOutputFilepath(rootDir, nodeID string) (string, error) {
	writePath := filepath.Join(rootDir, "config", "gentx")
	if err := os.MkdirAll(writePath, 0o700); err != nil {
		return "", fmt.Errorf(
			"could not create directory %q: %w",
			writePath,
			err,
		)
	}

	return filepath.Join(writePath, fmt.Sprintf("gentx-%v.json", nodeID)), nil
}

func readUnsignedGenTxFile(
	clientCtx client.Context,
	r io.Reader,
) (sdk.Tx, error) {
	bz, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	aTx, err := clientCtx.TxConfig.TxJSONDecoder()(bz)
	if err != nil {
		return nil, err
	}

	return aTx, err
}

func writeSignedGenTx(
	clientCtx client.Context,
	outputDocument string,
	tx sdk.Tx,
) error {
	//#nosec:G302,G304 // ignore error about file permissions
	outputFile, err := os.OpenFile(
		outputDocument,
		os.O_CREATE|os.O_EXCL|os.O_WRONLY,
		0o644,
	)
	if err != nil {
		return err
	}
	//#nosec:G307 // ignore error about file permissions
	defer outputFile.Close()

	json, err := clientCtx.TxConfig.TxJSONEncoder()(tx)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(outputFile, "%s\n", json)

	return err
}

// BuildCreateValidatorMsg makes a new MsgCreateValidator.
func BuildCreateValidatorMsg(
	clientCtx client.Context,
	config cli.TxCreateValidatorConfig,
	txBldr tx.Factory,
	generateOnly bool,
	valCodec address.Codec,
) (tx.Factory, sdk.Msg, error) {
	valAddr := clientCtx.GetFromAddress()
	valStr, err := valCodec.BytesToString(sdk.ValAddress(valAddr))
	if err != nil {
		return txBldr, nil, err
	}

	// var pkAny *codectypes.Any
	// if config.PubKey != nil {
	// 	var err error
	// 	if pkAny, err = codectypes.NewAnyWithValue(config.PubKey); err != nil {
	// 		return txBldr, nil, err
	// 	}
	// }

	msg := &beacontypes.MsgCreateValidatorX{
		Credentials: valStr,
		Pubkey:      config.PubKey.Bytes(),
	}

	// _ = msg2
	// msg := &banktypes.MsgSend{}

	if generateOnly {
		ip := config.IP
		p2pPort := config.P2PPort
		nodeID := config.NodeID

		if nodeID != "" && ip != "" && p2pPort > 0 {
			txBldr = txBldr.WithMemo(
				fmt.Sprintf("%s@%s:%d", nodeID, ip, p2pPort),
			)
		}
	}

	return txBldr, msg, nil
}
