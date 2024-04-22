package deposit

import (
	"crypto/ecdsa"
	"fmt"
	"os"

	"cosmossdk.io/depinject"
	"cosmossdk.io/log"

	engineclient "github.com/berachain/beacon-kit/mod/execution/client"
	"github.com/berachain/beacon-kit/mod/execution/client/ethclient"
	"github.com/berachain/beacon-kit/mod/node-builder/components"
	"github.com/berachain/beacon-kit/mod/node-builder/components/signer"
	"github.com/berachain/beacon-kit/mod/node-builder/config"
	"github.com/berachain/beacon-kit/mod/node-builder/config/spec"
	"github.com/berachain/beacon-kit/mod/node-builder/utils/jwt"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/runtime/services/staking/abi"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	gethclient "github.com/ethereum/go-ethereum/ethclient"
	"github.com/itsdevbear/comet-bls12-381/bls/blst"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewValidateDeposit creates a new command for validating a deposit message.
//
//nolint:gomnd // lots of magic numbers
func NewCreateValidator(clientCtx client.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-validator",
		Short: "Creates a validator deposit",
		Long: `Creates a validator deposit with the necessary credentials. The
		deposit message must include the public key, withdrawal credentials,
		and deposit amount. The arguments are expected in the order of withdrawal 
		credentials, deposit amount, current version, and genesis validator root.
		If the broadcast flag is set to true, a private key must be provided to
		sign the transaction.`,
		Args: cobra.ExactArgs(4),
		RunE: createValidatorCmd(clientCtx),
	}

	cmd.Flags().BoolP(
		broadcastDeposit, broadcastDepositShorthand,
		defaultBroadcastDeposit, broadcastDepositMsg,
	)
	cmd.Flags().String(privateKey, defaultPrivateKey, privateKeyMsg)

	return cmd
}

// validateDepositMessage validates a deposit message for creating a new
// validator.
func createValidatorCmd(clientCtx client.Context) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		var (
			blsSigner *signer.BLSSigner
			jwtSecret *jwt.Secret
			privKey   *ecdsa.PrivateKey

			logger = log.NewLogger(os.Stdout)
		)

		broadcastFlag, err := cmd.Flags().GetBool(broadcastDeposit)
		if err != nil {
			return err
		}

		// If the broadcast flag is set, a private key must be provided.
		if broadcastFlag {
			var fundingPrivKey string
			fundingPrivKey, err = cmd.Flags().GetString(privateKey)
			if err != nil {
				return err
			}
			if fundingPrivKey == "" {
				return ErrPrivateKeyRequired
			}

			privKey, err = crypto.HexToECDSA(fundingPrivKey)
			if err != nil {
				return err
			}
		}

		credentials, err := convertWithdrawalCredentials(args[0])
		if err != nil {
			return err
		}

		amount, err := convertAmountFromWei(args[1])
		if err != nil {
			return err
		}
		// amountBigInt, ok := new(big.Int).SetString(args[1], 10)
		// if !ok {
		// 	return ErrInvalidAmount
		// }

		currentVersion, err := convertVersion(args[2])
		if err != nil {
			return err
		}

		genesisValidatorRoot, err := convertGenesisValidatorRoot(args[3])
		if err != nil {
			return err
		}

		if err := depinject.Inject(
			depinject.Configs(
				depinject.Supply(
					logger,
					viper.GetViper(),
					spec.LocalnetChainSpec(),
				),
				depinject.Provide(
					components.ProvideBlsSigner,
				),
			),
			&blsSigner,
		); err != nil {
			panic(err)
		}

		// credentials = primitives.NewCredentialsFromExecutionAddress(
		// 	crypto.PubkeyToAddress(blsSigner.PublicKey()),
		// )

		// Create and sign the deposit message.
		depositMsg, signature, err := primitives.CreateAndSignDepositMessage(
			primitives.NewForkData(currentVersion, genesisValidatorRoot),
			spec.LocalnetChainSpec().DomainTypeDeposit(),
			blsSigner,
			credentials,
			amount,
		)
		if err != nil {
			return err
		}

		// Verify the deposit message.
		if err := depositMsg.VerifyCreateValidator(
			primitives.NewForkData(currentVersion, genesisValidatorRoot),
			signature,
			blst.VerifySignaturePubkeyBytes,
			spec.LocalnetChainSpec().DomainTypeDeposit(),
		); err != nil {
			return err
		}

		// If the broadcast flag is not set, output the deposit message and
		// signature and return early.
		logger.Info(
			"Deposit message created",
			"\nmessage", depositMsg,
			"\nsignature", signature,
		)
		if !broadcastFlag {
			return nil
		}

		// Spin up an engine client to broadcast the deposit transaction.
		// if err := depinject.Inject(
		// 	depinject.Configs(
		// 		depinject.Supply(
		// 			viper.GetViper(),
		// 		),
		// 		depinject.Provide(
		// 			components.ProvideJWTSecret,
		// 		),
		// 	),
		// 	&jwtSecret,
		// ); err != nil {
		// 	panic(err)
		// }
		jwtSecret, err = jwt.LoadFromFile("beacond/jwt.hex")
		if err != nil {
			panic(err)
		}

		cfg := config.MustReadConfigFromAppOpts(viper.GetViper())
		fmt.Println("CONFIG DUMP", cfg)

		cfg = config.DefaultConfig()

		ethClient, err := gethclient.Dial(cfg.Engine.RPCDialURL.String())
		if err != nil {
			return err
		}

		eth1client, err := ethclient.NewEth1Client(ethClient)
		if err != nil {
			return err
		}

		engineClient := engineclient.New(
			engineclient.WithEngineConfig(&cfg.Engine),
			engineclient.WithEth1Client(eth1client),
			engineclient.WithJWTSecret(jwtSecret),
			engineclient.WithLogger(logger),
		)
		engineClient.Start(cmd.Context())

		depositContract, err := abi.NewBeaconDepositContract(
			spec.LocalnetChainSpec().DepositContractAddress(),
			engineClient,
		)
		if err != nil {
			return err
		}

		chainID, err := engineClient.ChainID(cmd.Context())
		if err != nil {
			return err
		}

		// Send the deposit to the deposit contract.
		tx, err := depositContract.Deposit(
			&bind.TransactOpts{
				From: crypto.PubkeyToAddress(privKey.PublicKey),
				Signer: func(
					_ common.Address, tx *types.Transaction,
				) (*types.Transaction, error) {
					return types.SignTx(
						tx, types.LatestSignerForChainID(chainID),
						privKey,
					)
				},
				Value: depositMsg.Amount.ToWei(),
			},
			depositMsg.Pubkey[:],
			depositMsg.Credentials[:],
			0,
			signature[:],
		)
		if err != nil {
			return err
		}

		//
		receipt, err := bind.WaitMined(cmd.Context(), engineClient, tx)
		if err != nil {
			return err
		}

		if receipt.Status != 1 {
			return ErrDepositTransactionFailed
		}

		return nil
	}
}
