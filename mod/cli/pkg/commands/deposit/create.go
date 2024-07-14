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
	"os"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/cli/pkg/utils/parser"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/execution/pkg/client/ethclient"
	"github.com/berachain/beacon-kit/mod/geth-primitives/pkg/bind"
	"github.com/berachain/beacon-kit/mod/geth-primitives/pkg/deposit"
	"github.com/berachain/beacon-kit/mod/geth-primitives/pkg/rpc"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/signer"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/cosmos/cosmos-sdk/client"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
)

// NewCreateValidator creates a new command to create a validator deposit.
func NewCreateValidator[
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
](
	chainSpec common.ChainSpec,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-validator",
		Short: "Creates a validator deposit",
		Long: `Creates a validator deposit with the necessary credentials. The 
		arguments are expected in the order of withdrawal credentials, deposit
		amount, current version, and genesis validator root. If the broadcast
		flag is set to true, a private key must be provided to sign the transaction.`,
		Args: cobra.ExactArgs(4), //nolint:mnd // The number of arguments.
		RunE: createValidatorCmd[ExecutionPayloadT](chainSpec),
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
	cmd.Flags().String(rpcURL, defaultRPCURL, rpcURLMsg)

	return cmd
}

// createValidatorCmd returns a command that builds a create validator request.
//
//nolint:gocognit // The function is not complex.
func createValidatorCmd[
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
](
	chainSpec common.ChainSpec,
) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		logger := log.NewLogger(os.Stdout)

		broadcast, err := cmd.Flags().GetBool(broadcastDeposit)
		if err != nil {
			return err
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
			var txHash common.ExecutionHash
			txHash, err = broadcastDepositTx[ExecutionPayloadT](
				cmd, depositMsg, signature, chainSpec,
			)
			if err != nil {
				return err
			}

			logger.Info(
				"Deposit transaction successful",
				"txHash", txHash.Hex(),
			)
		}

		return nil
	}
}

func broadcastDepositTx[
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
](
	cmd *cobra.Command,
	depositMsg *types.DepositMessage,
	signature crypto.BLSSignature,
	chainSpec common.ChainSpec,
) (common.ExecutionHash, error) {
	// Parse the private key.
	var (
		fundingPrivKey string
		privKey        *ecdsa.PrivateKey
	)
	fundingPrivKey, err := cmd.Flags().GetString(privateKey)
	if err != nil {
		return common.ExecutionHash{}, err
	}
	if fundingPrivKey == "" {
		return common.ExecutionHash{}, ErrPrivateKeyRequired
	}
	privKey, err = ethcrypto.HexToECDSA(fundingPrivKey)
	if err != nil {
		return common.ExecutionHash{}, err
	}

	// Parse the execution client RPC URL.
	rpcURL, err := cmd.Flags().GetString(rpcURL)
	if err != nil {
		return common.ExecutionHash{}, err
	}
	eth1ChainID := new(big.Int).SetUint64(chainSpec.DepositEth1ChainID())
	depositContractAddress := chainSpec.DepositContractAddress()

	// Dial the execution client.
	rpcClient, err := rpc.DialContext(
		cmd.Context(), rpcURL,
	)
	if err != nil {
		return common.ExecutionHash{}, err
	}
	ethClient, err := ethclient.NewFromRPCClient[ExecutionPayloadT](
		rpcClient,
	)
	if err != nil {
		return common.ExecutionHash{}, err
	}

	// Send the deposit to the deposit contract through abi bindings.
	depositContract, err := deposit.NewBeaconDepositContract(
		depositContractAddress,
		ethClient,
	)
	if err != nil {
		return common.ExecutionHash{}, err
	}
	depositTx, err := depositContract.Deposit(
		&bind.TransactOpts{
			From: ethcrypto.PubkeyToAddress(privKey.PublicKey),
			Signer: func(
				_ common.ExecutionAddress, tx *ethtypes.Transaction,
			) (*ethtypes.Transaction, error) {
				return ethtypes.SignTx(
					tx,
					ethtypes.LatestSignerForChainID(eth1ChainID),
					privKey,
				)
			},
			//nolint:mnd // The gas tip cap.
			// It is necessary for besu to work, not sure why though.
			GasTipCap: big.NewInt(1000000000),
			//nolint:mnd // The gas fee cap.
			// It is necessary for ethereumjs to work.
			GasFeeCap: big.NewInt(1000000000),
			//nolint:mnd // The gas limit.
			// It is necessary for ethereumjs to work.
			GasLimit: 600000,
		},
		depositMsg.Pubkey[:],
		depositMsg.Credentials[:],
		depositMsg.Amount.Unwrap(),
		signature[:],
	)
	if err != nil {
		return common.ExecutionHash{}, errors.Wrapf(err, "error in depositing")
	}

	// Wait for the transaction to be mined and check the status.
	depositReceipt, err := bind.WaitMined(cmd.Context(), ethClient, depositTx)
	if err != nil {
		return common.ExecutionHash{}, errors.Wrapf(
			err,
			"waiting for transaction to be mined",
		)
	}
	if depositReceipt.Status == ethtypes.ReceiptStatusFailed {
		return common.ExecutionHash{}, ErrDepositTransactionFailed
	}

	return depositReceipt.TxHash, nil
}

// getBLSSigner returns a BLS signer based on the override commands key flag.
func getBLSSigner(
	cmd *cobra.Command,
) (crypto.BLSSigner, error) {
	var legacyKey components.LegacyKey
	overrideFlag, err := cmd.Flags().GetBool(overrideNodeKey)
	if err != nil {
		return nil, err
	}

	// Build the BLS signer.
	if overrideFlag {
		var validatorPrivKey string
		validatorPrivKey, err = cmd.Flags().GetString(valPrivateKey)
		if err != nil {
			return nil, err
		}
		if validatorPrivKey == "" {
			return nil, ErrValidatorPrivateKeyRequired
		}
		legacyKey, err = signer.LegacyKeyFromString(validatorPrivKey)
		if err != nil {
			return nil, err
		}
	}

	return components.ProvideBlsSigner(
		components.BlsSignerInput{
			AppOpts: client.GetViperFromCmd(cmd),
			PrivKey: legacyKey,
		},
	)
}
