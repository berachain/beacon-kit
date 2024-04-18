package deposit

import (
	"errors"
	"math/big"
	"os"

	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/node-builder/components"
	"github.com/berachain/beacon-kit/mod/node-builder/components/signer"
	"github.com/berachain/beacon-kit/mod/node-builder/config/spec"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/runtime/services/staking/abi"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
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
		and deposit amount. The arguments are expected in the order of public key,
		withdrawal credentials, deposit amount, signature, current version,
		and genesis validator root.`,
		Args: cobra.ExactArgs(4),
		RunE: createValidatorCmd(clientCtx),
	}

	return cmd
}

// validateDepositMessage validates a deposit message for creating a new
// validator.
func createValidatorCmd(clientCtx client.Context) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		var (
			blsSigner *signer.BLSSigner
		)
		if err := depinject.Inject(
			depinject.Configs(
				depinject.Supply(
					log.NewLogger(os.Stdout),
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

		credentials, err := ConvertWithdrawalCredentials(args[0])
		if err != nil {
			return err
		}

		amount, err := ConvertAmount(args[1])
		if err != nil {
			return err
		}

		currentVersion, err := ConvertVersion(args[2])
		if err != nil {
			return err
		}

		genesisValidatorRoot, err := ConvertGenesisValidatorRoot(args[3])
		if err != nil {
			return err
		}

		// TODO: modularize
		fundingPrivKey := "0xfffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306"

		// Create and sign the deposit message.
		depositMessage, signature, err := primitives.CreateAndSignDepositMessage(
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
		if err := depositMessage.VerifyCreateValidator(
			primitives.NewForkData(currentVersion, genesisValidatorRoot),
			signature,
			blst.VerifySignaturePubkeyBytes,
			spec.LocalnetChainSpec().DomainTypeDeposit(),
		); err != nil {
			return err
		}

		// todo: spin up and use engine api.
		ethClient, err := ethclient.Dial("http://localhost:8545")
		if err != nil {
			return err
		}

		// TODO: read from config.
		depositAddr := common.HexToAddress("0x00000000219ab540356cbb839cbe05303d7705fa")
		depositContract, err := abi.NewBeaconDepositContract(depositAddr, ethClient)
		if err != nil {
			return err
		}

		privKey, err := crypto.HexToECDSA(fundingPrivKey)
		if err != nil {
			return err
		}

		chainId := big.NewInt(80087)
		// Send the deposit to the deposit contract.
		tx, err := depositContract.Deposit(
			&bind.TransactOpts{
				From: crypto.PubkeyToAddress(privKey.PublicKey),
				Signer: func(
					_ common.Address, tx *types.Transaction,
				) (*types.Transaction, error) {
					return types.SignTx(
						tx, types.LatestSignerForChainID(chainId),
						privKey,
					)
				},
			},
			depositMessage.Pubkey[:],
			depositMessage.Credentials[:],
			depositMessage.Amount.Unwrap(),
			signature[:],
		)
		if err != nil {
			return err
		}

		//
		receipt, err := bind.WaitMined(cmd.Context(), ethClient, tx)
		if err != nil {
			return err
		}

		if receipt.Status != 1 {
			return errors.New("deposit transaction failed")
		}

		return nil
	}
}
