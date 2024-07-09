// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package deposit

import (
	"crypto/ecdsa"
	"math/big"
	"net/url"
	"os"

	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/cli/pkg/utils/parser"
	"github.com/berachain/beacon-kit/mod/config"
	"github.com/berachain/beacon-kit/mod/config/pkg/spec"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/errors"
	engineclient "github.com/berachain/beacon-kit/mod/execution/pkg/client"
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/geth-primitives/pkg/deposit"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/signer"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/net/jwt"
	myUrl "github.com/berachain/beacon-kit/mod/primitives/pkg/net/url"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethCommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
)

// NewValidateDeposit creates a new command for validating a deposit message.
//

func NewCreateValidator(chainSpec common.ChainSpec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-validator",
		Short: "Creates a validator deposit",
		Long: `Creates a validator deposit with the necessary credentials. The 
		arguments are expected in the order of withdrawal credentials, deposit
		amount, current version, and genesis validator root. If the broadcast
		flag is set to true, a private key must be provided to sign the transaction.`,
		Args: cobra.ExactArgs(4), //nolint:mnd // The number of arguments.
		RunE: createValidatorCmd(chainSpec),
	}

	cmd.Flags().BoolP(
		broadcastDeposit, broadcastDepositShorthand,
		defaultBroadcastDeposit, broadcastDepositMsg,
	)
	cmd.Flags().String(privateKey, defaultPrivateKey, privateKeyMsg)
	cmd.Flags().BoolP(
		overrideNodeKey, overrideNodeKeyShorthand,
		defaultOverrideNodeKey, overrideNodeKeyMsg,
	)
	cmd.Flags().
		String(valPrivateKey, defaultValidatorPrivateKey, valPrivateKeyMsg)
	cmd.Flags().String(jwtSecretPath, defaultJWTSecretPath, jwtSecretPathMsg)
	cmd.Flags().String(engineRPCURL, defaultEngineRPCURL, engineRPCURLMsg)

	return cmd
}

// createValidatorCmd returns a command that builds a create validator request.
//
// TODO: Implement broadcast functionality. Currently, the implementation works
// for the geth client but something about the Deposit binding is not handling
// other execution layers correctly. Peep the commit history for what we had.
// 🤷‍♂️.
func createValidatorCmd(
	chainSpec common.ChainSpec,
) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		var (
			logger  = log.NewLogger(os.Stdout)
			privKey *ecdsa.PrivateKey
		)

		broadcast, err := cmd.Flags().GetBool(broadcastDeposit)
		if err != nil {
			return err
		}

		// If the broadcast flag is set, a private key must be provided.
		if broadcast {
			var fundingPrivKey string
			fundingPrivKey, err = cmd.Flags().GetString(privateKey)
			if err != nil {
				return err
			}
			if fundingPrivKey == "" {
				return parser.ErrPrivateKeyRequired
			}

			privKey, err = ethCrypto.HexToECDSA(fundingPrivKey)
			if err != nil {
				return err
			}
		}

		// Get the BLS signer.
		blsSigner, err := getBLSSigner(cmd)
		if err != nil {
			return err
		}

		credentials, err := parser.ConvertWithdrawalCredentials(args[0])
		if err != nil {
			return err
		}

		amount, err := parser.ConvertAmount(args[1])
		if err != nil {
			return err
		}

		currentVersion, err := parser.ConvertVersion(args[2])
		if err != nil {
			return err
		}

		genesisValidatorRoot, err := parser.ConvertGenesisValidatorRoot(args[3])
		if err != nil {
			return err
		}

		// Create and sign the deposit message.
		depositMsg, signature, err := types.CreateAndSignDepositMessage(
			types.NewForkData(currentVersion, genesisValidatorRoot),
			chainSpec.DomainTypeDeposit(),
			blsSigner,
			credentials,
			amount,
		)
		if err != nil {
			return err
		}

		// Verify the deposit message.
		if err = depositMsg.VerifyCreateValidator(
			types.NewForkData(currentVersion, genesisValidatorRoot),
			signature,
			chainSpec.DomainTypeDeposit(),
			signer.BLSSigner{}.VerifySignature,
		); err != nil {
			return err
		}

		// If the broadcast flag is not set, output the deposit message and
		// signature and return early.
		logger.Info(
			"Deposit Message CallData",
			"pubkey", depositMsg.Pubkey.String(),
			"withdrawal credentials", depositMsg.Credentials.String(),
			"amount", depositMsg.Amount,
			"signature", signature.String(),
		)

		if broadcast {
			var txHash gethCommon.Hash
			txHash, err = broadcastDepositTx(
				cmd, depositMsg, signature, privKey, logger, spec.DevnetChainSpec(),
			)
			if err != nil {
				return err
			}

			logger.Info(
				"Deposit transaction successful",
				"txHash", txHash.Hex(),
			)
		}

		// TODO: once broadcast is fixed, remove this.
		logger.Info("Send the above calldata to the deposit contract 🫡")

		return nil
	}
}

func broadcastDepositTx(
	cmd *cobra.Command,
	depositMsg *types.DepositMessage,
	signature crypto.BLSSignature,
	privKey *ecdsa.PrivateKey,
	logger log.Logger,
	chainSpec common.ChainSpec,
) (gethCommon.Hash, error) {
	// Spin up an engine client to broadcast the deposit transaction.
	// TODO: This should read in the actual config file. I'm going to rope
	// if I keep trying this right now so it's a flag lol! 🥲
	cfg := config.DefaultConfig()

	// Parse the engine RPC URL.
	engineRPCURL, err := cmd.Flags().GetString(engineRPCURL)
	if err != nil {
		return gethCommon.Hash{}, err
	}

	parsedURL, err := url.Parse(engineRPCURL)
	if err != nil {
		return gethCommon.Hash{}, err
	}

	cfg.Engine.RPCDialURL = convertURLToConnectionURL(parsedURL)
	logger.Info("RPCDialURL", "url", cfg.Engine.RPCDialURL)

	// Load the JWT secret.
	cfg.Engine.JWTSecretPath, err = cmd.Flags().GetString(jwtSecretPath)
	if err != nil {
		return gethCommon.Hash{}, err
	}

	jwtSecret, err := loadFromFile(cfg.Engine.JWTSecretPath)
	logger.Info("jwtSecret", "jwtSecret", jwtSecret)

	if err != nil {
		return gethCommon.Hash{}, errors.Wrapf(err, "error in loading jwt secret")
	}

	// TODO: This is a WIP. I'm not sure how to get the engine client to work
	// config := &beaconClient.Config{}
	//
	// var telemetrySink beaconClient.TelemetrySink

	// Spin up the engine client.
	// engineClient := engineclient.New[ExecutionPayload, PayloadAttributes](
	// config,
	// logger,
	// jwtSecret,
	// telemetrySink,
	// new(big.Int).SetUint64(chainSpec.DepositEth1ChainID()
	// ))
	// engineClient := engineclient.New(
	//	engineclient.WithEngineConfig(&cfg.Engine),
	//	engineclient.WithJWTSecret(jwtSecret),
	//	engineclient.WithLogger(logger),
	//)

	// engineClient, err := setupEngineClient(
	// cfg.Engine.RPCDialURL,
	// jwtSecret,
	// chainSpec,
	// logger)
	// fmt.Println("engineClient", engineClient)
	// err = engine ̰Client.Start(cmd.Context())
	// if err != nil {
	//	fmt.Println("err in starting engine client", err)
	//}

	engineClient, err := ethclient.Dial("http://localhost:8545")
	if err != nil || engineClient == nil {
		return gethCommon.Hash{}, errors.New("failed to create Ethereum client")
	}

	depositContractAddress := chainSpec.DepositContractAddress()

	chainID, err := engineClient.ChainID(cmd.Context())
	if err != nil {
		return gethCommon.Hash{}, err
	}

	// one way
	// contractAbi, err := deposit.BeaconDepositContractMetaData.GetAbi()
	// if err != nil {
	//	panic(err)
	//}
	// fmt.Println("CONTRACT ABI", contractAbi)
	//
	// callData, err := contractAbi.Pack(
	//	"deposit",
	//	depositMsg.Pubkey[:],
	//	depositMsg.Credentials[:],
	//	uint64(0),
	//	signature[:],
	//)
	// fmt.Println("CALL DATA", callData)
	//
	// if err != nil {
	//	fmt.Println("PANIC AT PACK")
	//	panic(err)
	//}
	//
	// tx := ethTypes.NewTx(
	//	&ethTypes.DynamicFeeTx{
	//		Nonce:     latestNonce,
	//		ChainID:   chainID,
	//		To:        &depositContractAddress,
	//		Value:     depositMsg.Amount.ToWei(),
	//		Data:      callData,
	//		GasTipCap: big.NewInt(1000000000),
	//		GasFeeCap: big.NewInt(1000000000),
	//		Gas:       500000,
	//	},
	//)
	//
	// fmt.Println("TX", tx)
	//
	// signedTx, err := ethTypes.SignTx(
	// tx, ethTypes.LatestSignerForChainID(chainID),
	// privKey
	// )
	// fmt.Println("SIGNED TX", signedTx)
	//
	// if err != nil {
	//	fmt.Errorf("error in signing tx: %v", err)
	//}

	////Now send this raw transaction through your RPC client
	// _, err = engineClient.CallContract(
	//	cmd.Context(),
	//	ethereum.CallMsg{
	//		From:  ethCrypto.PubkeyToAddress(privKey.PublicKey),
	//		To:    &depositContractAddress,
	//		Value: depositMsg.Amount.ToWei(),
	//		Data:  signedTx.Data(),
	//	},
	//	big.NewInt(0),
	//)
	//
	// fmt.Errorf("error in calling contract: %v", err)
	//
	// if err = engineClient.SendTransaction(
	//	cmd.Context(),
	//	signedTx,
	// ); err != nil {
	//	fmt.Println("PANIC AT SEND TRANSACTION")
	//	panic(err)
	//}
	//
	// fmt.Println("CONTRACT CALLED")
	// Getting same error using both approaches

	// One approach to send deposit txn - is through go bindings
	// Send the deposit to the deposit contract.

	depositContract, err := deposit.NewBeaconDepositContract(
		depositContractAddress,
		engineClient,
	)
	if err != nil {
		return gethCommon.Hash{}, err
	}

	logger.Info("from", "from", ethCrypto.PubkeyToAddress(privKey.PublicKey))
	fromAddress := ethCrypto.PubkeyToAddress(privKey.PublicKey)

	latestNonceForDeposit, errInNonce := engineClient.NonceAt(
		cmd.Context(),
		ethCrypto.PubkeyToAddress(privKey.PublicKey),
		nil,
	)
	if errInNonce != nil {
		return gethCommon.Hash{}, errInNonce
	}
	logger.Info("LATEST NONCE", "nonce", latestNonceForDeposit)

	depositTx, err := depositContract.Deposit(
		&bind.TransactOpts{
			From: fromAddress,
			Signer: func(
				_ common.ExecutionAddress, tx *ethTypes.Transaction,
			) (*ethTypes.Transaction, error) {
				return ethTypes.SignTx(
					tx, ethTypes.LatestSignerForChainID(chainID),
					privKey,
				)
			},
			Nonce: new(big.Int).SetUint64(latestNonceForDeposit),
			//nolint:mnd // The gas tip cap
			GasTipCap: big.NewInt(1000000000),
			//nolint:mnd // The gas fee cap
			GasFeeCap: big.NewInt(1000000000),
			//nolint:mnd // The gas limit
			GasLimit: 600000,
		},
		depositMsg.Pubkey[:],
		depositMsg.Credentials[:],
		uint64(depositMsg.Amount), // 32 eth is minimum deposit amount.
		signature[:],
	)
	if err != nil {
		return gethCommon.Hash{}, errors.Wrapf(err, "error in depositing")
	}

	// Wait for the transaction to be mined and check the status.
	depositReceipt, err := bind.WaitMined(cmd.Context(), engineClient, depositTx)
	if err != nil {
		return gethCommon.Hash{}, errors.Wrapf(
			err,
			"waiting for transaction to be mined",
		)
	}

	if depositReceipt.Status != 1 {
		return gethCommon.Hash{}, parser.ErrDepositTransactionFailed
	}

	return depositReceipt.TxHash, nil
}

func loadFromFile(path string) (*jwt.Secret, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return jwt.NewFromHex(string(data))
}

type ExecutionPayload struct {
	types.InnerExecutionPayload
}

func (ep ExecutionPayload) Empty(_ uint32) ExecutionPayload {
	return ExecutionPayload{}
}

type PayloadAttributes struct {
}

func (p PayloadAttributes) IsNil() bool {
	// TODO implement me
	panic("implement me")
}

//nolint:lll
func (p PayloadAttributes) GetSuggestedFeeRecipient() gethprimitives.ExecutionAddress {
	// TODO implement me
	panic("implement me")
}

func setupEngineClient(rpcURL *myUrl.ConnectionURL,
	secret *jwt.Secret,
	chainSpec common.ChainSpec,
	logger log.Logger,
) (*engineclient.EngineClient[ExecutionPayload, PayloadAttributes], error) {
	cfg := &engineclient.Config{RPCDialURL: rpcURL}
	var telemetrySink engineclient.TelemetrySink
	eth1ChainID := new(big.Int).SetUint64(chainSpec.DepositEth1ChainID())
	engineClient := engineclient.New[
		ExecutionPayload,
		PayloadAttributes](
		cfg, logger, secret, telemetrySink, eth1ChainID)
	return engineClient, nil
}

// getBLSSigner returns a BLS signer based on the override commands key flag.
func getBLSSigner(
	cmd *cobra.Command,
) (crypto.BLSSigner, error) {
	var blsSigner crypto.BLSSigner
	supplies := []interface{}{client.GetViperFromCmd(cmd)}
	overrideFlag, err := cmd.Flags().GetBool(overrideNodeKey)
	if err != nil {
		return nil, err
	}

	// Build the BLS signer.
	if overrideFlag {
		var (
			validatorPrivKey string
			legacyInput      components.LegacyKey
		)
		validatorPrivKey, err = cmd.Flags().GetString(valPrivateKey)
		if err != nil {
			return nil, err
		}
		if validatorPrivKey == "" {
			return nil, ErrValidatorPrivateKeyRequired
		}
		legacyInput, err = signer.LegacyKeyFromString(validatorPrivKey)
		if err != nil {
			return nil, err
		}
		supplies = append(supplies, legacyInput)
	}

	if err = depinject.Inject(
		depinject.Configs(
			depinject.Supply(supplies...),
			depinject.Provide(
				components.ProvideBlsSigner,
			),
		),
		&blsSigner,
	); err != nil {
		return nil, err
	}

	return blsSigner, nil
}

func convertURLToConnectionURL(u *url.URL) *myUrl.ConnectionURL {
	return &myUrl.ConnectionURL{
		URL: u,
	}
}
