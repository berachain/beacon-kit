package main

import (
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/berachain/beacon-kit/mod/cli/pkg/utils/parser"
	genesiscmd "github.com/berachain/beacon-kit/mod/cli/pkg/v2/commands/genesis"
	"github.com/berachain/beacon-kit/mod/config"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/genesis"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/cometbft"
	"github.com/berachain/beacon-kit/mod/errors"
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/app/components"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/app/components/signer"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	cmtjson "github.com/cometbft/cometbft/libs/json"
	cmttypes "github.com/cometbft/cometbft/types"
)

func AddExecutionPayloadToGenesis(cmtGenesis *cometbft.Genesis, chainSpec common.ChainSpec, config *cometbft.Config, ethGenesisPath string) error {
	ethGenesisBz, err := os.ReadFile(ethGenesisPath)
	if err != nil {
		return errors.Wrap(err, "failed to read eth1 genesis file")
	}
	ethGenesis := &gethprimitives.Genesis{}
	if err := ethGenesis.UnmarshalJSON(ethGenesisBz); err != nil {
		return errors.Wrap(err, "failed to unmarshal eth1 genesis")
	}
	genesisBlock := ethGenesis.ToBlock()
	payload := gethprimitives.BlockToExecutableData(
		genesisBlock,
		nil,
		nil,
	).ExecutionPayload

	genesisInfo := &genesis.Genesis[
		*types.Deposit, *types.ExecutionPayloadHeader,
	]{}
	if err = json.Unmarshal(
		cmtGenesis.AppState, genesisInfo,
	); err != nil {
		return errors.Wrap(err, "failed to unmarshal beacon state")
	}

	header, err := genesiscmd.ExecutableDataToExecutionPayloadHeader(
		version.ToUint32(genesisInfo.ForkVersion),
		payload,
		chainSpec.MaxWithdrawalsPerPayload(),
	)
	if err != nil {
		return errors.Wrap(
			err,
			"failed to convert executable data to execution payload header",
		)
	}
	genesisInfo.ExecutionPayloadHeader = header
	cmtGenesis.AppState, err = json.MarshalIndent(genesisInfo, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal beacon state")
	}

	return cmtGenesis.Export(config.GenesisFile())
}

func AddDepositToGenesis(cmtGenesis *cometbft.Genesis, chainSpec common.ChainSpec, cfg *config.Config, appOpts *components.AppOptions, pubKey crypto.BLSPubkey) error {
	blsSigner, err := components.ProvideBlsSigner(components.BlsSignerInput{
		AppOpts: appOpts,
		Config:  cfg,
	})
	if err != nil {
		return errors.Wrap(err, "failed to provide bls signer")
	}

	depositAmount, err := parser.ConvertAmount("32000000000")
	if err != nil {
		return err
	}

	currentVersion := version.FromUint32[common.Version](
		version.Deneb,
	)

	depositMsg, signature, err := types.CreateAndSignDepositMessage(
		types.NewForkData(currentVersion, common.Root{}),
		chainSpec.DomainTypeDeposit(),
		blsSigner,
		// TODO: configurable.
		types.NewCredentialsFromExecutionAddress(
			gethprimitives.ExecutionAddress{},
		),
		depositAmount,
	)
	if err != nil {
		return err
	}

	// Verify the deposit message.
	if err = depositMsg.VerifyCreateValidator(
		types.NewForkData(currentVersion, common.Root{}),
		signature,
		chainSpec.DomainTypeDeposit(),
		signer.BLSSigner{}.VerifySignature,
	); err != nil {
		return err
	}

	deposit := types.Deposit{
		Pubkey:      depositMsg.Pubkey,
		Amount:      depositMsg.Amount,
		Signature:   signature,
		Credentials: depositMsg.Credentials,
	}

	outputDocument, err := genesiscmd.MakeOutputFilepath(
		appOpts.HomeDir,
		pubKey.String(),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create output file path")
	}
	if err = genesiscmd.WriteDepositToFile(outputDocument, &deposit); err != nil {
		return errors.Wrap(err, "failed to write signed gen tx")
	}

	genesisInfo := &genesis.Genesis[
		*types.Deposit, *types.ExecutionPayloadHeader,
	]{}
	if err = json.Unmarshal(
		cmtGenesis.AppState, genesisInfo,
	); err != nil {
		return errors.Wrap(err, "failed to unmarshal beacon state")
	}

	//#nosec:G701 // won't realistically overflow.
	deposit.Index = uint64(len(genesisInfo.Deposits)) + 1
	genesisInfo.Deposits = append(genesisInfo.Deposits, &deposit)

	if cmtGenesis.AppState, err = json.MarshalIndent(
		genesisInfo, "", "  ",
	); err != nil {
		return err
	}
	return cmtGenesis.Export(cfg.CometBFT.GenesisFile())
}

func ConvertBartioGenesis(cmtConfig *cometbft.Config, bartioGenesisPath string) error {
	// Replace the file with the Bartio genesis file
	// Force copy the file at path into configDir
	sourceFile, err := os.Open(bartioGenesisPath)
	if err != nil {
		return errors.Newf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(cmtConfig.GenesisFile())
	if err != nil {
		return errors.Newf("failed to create destination file: %w", err)
	}

	sourceContent, err := io.ReadAll(sourceFile)
	if err != nil {
		return errors.Newf("failed to read source file: %w", err)
	}
	_, err = destFile.Write(sourceContent)
	if err != nil {
		return errors.Newf("failed to copy file: %w", err)
	}
	destFile.Close()

	// Read the genesis file
	genesisBytes, err := os.ReadFile(cmtConfig.GenesisFile())
	if err != nil {
		return errors.Newf("failed to read genesis file: %w", err)
	}

	// Unmarshal into map[string]json.RawMessage
	var genesisMap map[string]json.RawMessage
	if err := cmtjson.Unmarshal(genesisBytes, &genesisMap); err != nil {
		return errors.Newf("failed to unmarshal genesis file: %w", err)
	}

	var chainID string
	if err := json.Unmarshal(genesisMap["chain_id"], &chainID); err != nil {
		return errors.Newf("failed to unmarshal chain ID: %w", err)
	}

	var genesisState map[string]json.RawMessage
	if err := json.Unmarshal(genesisMap["app_state"], &genesisState); err != nil {
		return errors.Newf("failed to unmarshal genesis state: %w", err)
	}

	var consensus map[string]json.RawMessage
	if err := json.Unmarshal(genesisMap["consensus"], &consensus); err != nil {
		return errors.Newf("failed to unmarshal consensus: %w", err)
	}

	var consensusParams *cmttypes.ConsensusParams
	if err := cmtjson.Unmarshal(consensus["params"], &consensusParams); err != nil {
		return errors.Newf("failed to unmarshal consensus params: %w", err)
	}

	cmtGenesis := cometbft.NewGenesis(
		chainID,
		genesisState["beacon"],
		consensusParams,
	)
	t := &time.Time{}
	if err := t.UnmarshalJSON(genesisMap["genesis_time"]); err != nil {
		return errors.Newf("failed to unmarshal genesis time: %w", err)
	}
	cmtGenesis.GenesisTime = *t
	if err := cmtGenesis.Export(cmtConfig.GenesisFile()); err != nil {
		return err
	}
	return nil
}
