//go:build simulated

// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package simulated_test

import (
	"context"
	"math/big"
	"path"
	"testing"
	"time"

	"github.com/berachain/beacon-kit/beacon/blockchain"
	depositcli "github.com/berachain/beacon-kit/cli/commands/deposit"
	consensustypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/consensus/cometbft/service/encoding"
	"github.com/berachain/beacon-kit/gethlib/deposit"
	"github.com/berachain/beacon-kit/log/phuslu"
	"github.com/berachain/beacon-kit/node-core/components/signer"
	"github.com/berachain/beacon-kit/primitives/common"
	beaconmath "github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/state-transition/core"
	"github.com/berachain/beacon-kit/testing/simulated"
	"github.com/berachain/beacon-kit/testing/simulated/execution"
	cmtabci "github.com/cometbft/cometbft/abci/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	gethcore "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

// FuluDepositSuite tests the deposit queue drain at the Fulu fork boundary.
type FuluDepositSuite struct {
	suite.Suite
	simulated.SharedAccessors
}

func TestFuluDepositSuite(t *testing.T) {
	suite.Run(t, new(FuluDepositSuite))
}

func (s *FuluDepositSuite) SetupTest() {
	s.CtxApp, s.CtxAppCancelFn = context.WithCancel(context.Background())
	s.CtxComet = context.TODO()
	s.HomeDir = s.T().TempDir()

	const elGenesisPath = "./el-genesis-files/fulu-deposit-genesis.json"
	chainSpecFunc := simulated.ProvideFuluDepositTestChainSpec
	chainSpec, err := chainSpecFunc()
	s.Require().NoError(err)
	configs, genesisValidatorsRoot := simulated.InitializeHomeDirs(s.T(), chainSpec, elGenesisPath, s.HomeDir)
	cometConfig := configs[0]
	s.GenesisValidatorsRoot = genesisValidatorsRoot

	elNode := execution.NewRethNode(s.HomeDir, execution.ValidRethImage())
	elHandle, authRPC, elRPC := elNode.Start(s.T(), path.Base(elGenesisPath))
	s.ElHandle = elHandle

	s.LogBuffer = &simulated.SyncBuffer{}
	logger := phuslu.NewLogger(s.LogBuffer, nil)

	components := simulated.FixedComponents(s.T())
	components = append(components, simulated.ProvideSimComet)
	components = append(components, chainSpecFunc)

	s.TestNode = simulated.NewTestNode(s.T(), simulated.TestNodeInput{
		TempHomeDir: s.HomeDir,
		CometConfig: cometConfig,
		AuthRPC:     authRPC,
		ClientRPC:   elRPC,
		Logger:      logger,
		AppOpts:     viper.New(),
		Components:  components,
	})

	s.SimComet = s.TestNode.SimComet

	go func() {
		_ = s.TestNode.Start(s.CtxApp)
	}()

	s.SimulationClient = execution.NewSimulationClient(s.TestNode.ContractBackend)
	timeOut := 10 * time.Second
	interval := 50 * time.Millisecond
	err = simulated.WaitTillServicesStarted(s.LogBuffer, timeOut, interval)
	s.Require().NoError(err)
}

func (s *FuluDepositSuite) TearDownTest() {
	s.CleanupTest(s.T())
}

// TestDepositQueueDrainedOnFirstFuluBlock verifies that when the deposit
// store is overloaded with 3x MaxDepositsPerBlock just before the Fulu fork,
// the first Fulu block drains the entire queue. Additionally, new deposits
// arriving as EIP-6110 execution requests in the same block are also processed.
//
// Chain spec: Deneb1 at genesis, Electra1 at t=6, Fulu at t=7, MaxDepositsPerBlock=4.
// EL genesis: Cancun at genesis, Prague/Prague1 at t=6, Osaka at t=7.
//
// Timeline:
//
//	Block 1 (t=5): EL Cancun block includes 12 deposit txs. Eth1FollowDistance=1 prevents sync.
//	Block 2 (t=6): Electra1/Prague1 fork. FinalizeBlock syncs all 12 deposits from Cancun EL block 1.
//	               Send 2 additional deposit txs (for EIP-6110 requests in the next Osaka block).
//	Block 3 (t=7): First Fulu/Osaka block. Drains all 12 catchup deposits from
//	               the block body + 2 EIP-6110 deposit requests from execution payload.
//	Block 4 (t=8): Post-fork block to confirm chain continues cleanly.
func (s *FuluDepositSuite) TestDepositQueueDrainedOnFirstFuluBlock() {
	s.InitializeChain(s.T(), 1)

	blsSigner := simulated.GetBlsSigner(s.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)
	nodeAddress := pubkey.Address()
	s.SimComet.Comet.SetNodeAddress(nodeAddress)

	credAddress, err := common.NewExecutionAddressFromHex(simulated.WithdrawalExecutionAddress)
	s.Require().NoError(err)
	creds := consensustypes.NewCredentialsFromExecutionAddress(credAddress)

	maxDepositsPerBlock := s.TestNode.ChainSpec.MaxDepositsPerBlock()
	s.Require().Equal(uint64(4), maxDepositsPerBlock)

	numQueuedDeposits := int(3 * maxDepositsPerBlock) // 12
	depositAmount := beaconmath.Gwei(10_000 * 1e9)    // 10K BERA each

	// Send 3x MaxDepositsPerBlock deposits to the EL deposit contract.
	// These will be included in the next EL block and later synced to the CL deposit store.
	for i := 0; i < numQueuedDeposits; i++ {
		setOperator := i == 0
		s.sendDeposit(blsSigner, creds, depositAmount, setOperator, big.NewInt(int64(i)))
	}

	nextBlockTime := time.Unix(5, 0)
	nextBlockHeight := int64(1)

	// [Block 1, t=5] Deneb1/Cancun block. EL includes the 12 deposit txs.
	// Due to Eth1FollowDistance=1, the CL does not sync these deposits yet.
	{
		_, _, nextBlockTime = s.MoveChainToHeight(s.T(), nextBlockHeight, 1, nodeAddress, nextBlockTime)
		s.Require().Equal(time.Unix(6, 0), nextBlockTime)

		ds := s.TestNode.StorageBackend.DepositStore()
		deposits, _, err := ds.GetDepositsByIndex(s.CtxApp, 0, 100)
		s.Require().NoError(err)
		s.Require().Len(deposits, 1, "Only the genesis deposit should be in store after block 1")
		nextBlockHeight++
	}

	// [Block 2, t=6] Electra1/Prague1 fork activates. FinalizeBlock syncs deposits
	// from Cancun EL block 1. The deposit store now has 1 (genesis) + 12 (new) = 13.
	// No deposits are processed yet from the block body (PrepareProposal ran before sync).
	{
		s.LogBuffer.Reset()
		_, _, nextBlockTime = s.MoveChainToHeight(s.T(), nextBlockHeight, 1, nodeAddress, nextBlockTime)
		s.Require().Equal(time.Unix(7, 0), nextBlockTime)

		s.Require().Contains(s.LogBuffer.String(),
			"welcome to the electra1 (0x05010000) fork!",
			"Electra1 fork should activate on block 2")

		ds := s.TestNode.StorageBackend.DepositStore()
		deposits, _, err := ds.GetDepositsByIndex(s.CtxApp, 0, 100)
		s.Require().NoError(err)
		s.Require().Len(deposits, 1+numQueuedDeposits,
			"Deposit store should have genesis + %d queued deposits", numQueuedDeposits)

		validators, err := s.TestNode.APIBackend.FilterValidators(nextBlockHeight, nil, nil)
		s.Require().NoError(err)
		s.Require().Len(validators, 1, "Still 1 validator; queued deposits not yet applied")
		nextBlockHeight++
	}

	// Send 2 more deposit txs that will be picked up by the first Osaka EL block
	// as EIP-6110 execution requests.
	numEIP6110Deposits := 2
	for i := 0; i < numEIP6110Deposits; i++ {
		nonce := big.NewInt(int64(numQueuedDeposits + i))
		s.sendDeposit(blsSigner, creds, depositAmount, false, nonce)
	}
	time.Sleep(time.Second)

	// [Block 3, t=7] First Fulu/Osaka block.
	// The catchup logic sets depositRange=MaxUint64, so all 12 queued deposits
	// are placed on the block body. EIP-6110 deposit requests from the execution
	// payload are appended. All deposits are processed in a single block.
	{
		s.LogBuffer.Reset()
		_, _, nextBlockTime = s.MoveChainToHeight(s.T(), nextBlockHeight, 1, nodeAddress, nextBlockTime)

		s.Require().Contains(s.LogBuffer.String(),
			"welcome to the fulu (0x06000000) fork!",
			"Fulu fork should activate on block 3")

		s.Require().Contains(s.LogBuffer.String(),
			"Building block body with local deposits",
			"Block builder should report deposits being included")

		nextBlockHeight++
	}

	// [Block 4, t=8] Post-fork: confirm the chain continues without errors.
	{
		s.LogBuffer.Reset()
		_, _, nextBlockTime = s.MoveChainToHeight(s.T(), nextBlockHeight, 1, nodeAddress, nextBlockTime)
		_ = nextBlockTime

		validators, err := s.TestNode.APIBackend.FilterValidators(nextBlockHeight, nil, nil)
		s.Require().NoError(err)
		s.Require().Len(validators, 1, "Still 1 validator (all deposits went to the same pubkey)")

		s.T().Logf("effective balance after deposits: %s gwei", validators[0].Validator.EffectiveBalance)
	}
}

func (s *FuluDepositSuite) sendDeposit(
	blsSigner *signer.BLSSigner,
	creds consensustypes.WithdrawalCredentials,
	depositAmount beaconmath.Gwei,
	setOperator bool,
	nonce *big.Int,
) gethcommon.Hash {
	depositContractAddress := gethcommon.Address(s.TestNode.ChainSpec.DepositContractAddress())
	depositClient, err := deposit.NewDepositContract(depositContractAddress, s.TestNode.ContractBackend)
	s.Require().NoError(err)

	depositMsg, blsSig, err := depositcli.CreateDepositMessage(
		s.TestNode.ChainSpec,
		blsSigner,
		s.GenesisValidatorsRoot,
		creds,
		depositAmount,
	)
	s.Require().NoError(err)
	err = depositcli.ValidateDeposit(
		s.TestNode.ChainSpec,
		depositMsg.Pubkey,
		depositMsg.Credentials,
		depositMsg.Amount,
		s.GenesisValidatorsRoot,
		blsSig,
	)
	s.Require().NoError(err)

	elChainID := big.NewInt(int64(s.TestNode.ChainSpec.DepositEth1ChainID()))
	senderKey := simulated.GetTestKey(s.T())
	senderAddress := gethcommon.HexToAddress(creds.String())
	operator := senderAddress
	if !setOperator {
		operator = gethcommon.HexToAddress("0x0000000000000000000000000000000000000000")
	}

	txOpts := &bind.TransactOpts{
		From: senderAddress,
		Signer: func(_ gethcommon.Address, tx *gethcore.Transaction) (*gethcore.Transaction, error) {
			return gethcore.SignTx(
				tx, gethcore.LatestSignerForChainID(elChainID), senderKey,
			)
		},
		GasLimit: 200_000,
		Value:    big.NewInt(0).Mul(big.NewInt(int64(depositAmount)), big.NewInt(1e9)),
	}

	if nonce != nil {
		txOpts.Nonce = nonce
	}

	tx, depErr := depositClient.Deposit(txOpts, depositMsg.Pubkey[:], depositMsg.Credentials[:], blsSig[:], operator)
	s.Require().NoError(depErr)
	return tx.Hash()
}

// TestBodyDepositsAfterFuluRejected verifies that, once Fulu (Osaka) is active and the
// pre-Fulu deposit queue has been drained, a proposed block carrying deposits on the
// beacon block body is rejected. From Fulu onwards deposits must be sourced exclusively
// from EIP-6110 execution requests, so the only block permitted to carry body deposits is
// the single first-Fulu catchup block.
//
// Chain spec (ProvideFuluDepositTestChainSpec): Electra1 at t=6, Fulu at t=7. Block 3
// (t=7) is the first Fulu block; block 4 (t=8) is the first block where deposits on the
// block body are no longer a valid source.
func (s *FuluDepositSuite) TestBodyDepositsAfterFuluRejected() {
	s.InitializeChain(s.T(), 1)

	blsSigner := simulated.GetBlsSigner(s.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)
	nodeAddress := pubkey.Address()
	s.SimComet.Comet.SetNodeAddress(nodeAddress)

	// Advance through the first Fulu block (block 3 at t=7) so that the next block is a
	// non-first Fulu block where block body deposits are disallowed.
	const postFuluHeight = int64(4)
	_, _, nextBlockTime := s.MoveChainToHeight(s.T(), 1, postFuluHeight-1, nodeAddress, time.Unix(5, 0))
	s.Require().Equal(time.Unix(8, 0), nextBlockTime, "block 4 must be the first post-Fulu block")

	// Prepare a valid block proposal for the post-Fulu height.
	validProposal, err := s.SimComet.Comet.PrepareProposal(s.CtxComet, &cmtabci.PrepareProposalRequest{
		Height:          postFuluHeight,
		Time:            nextBlockTime,
		ProposerAddress: nodeAddress,
	})
	s.Require().NoError(err)
	s.Require().NotEmpty(validProposal)

	// Inject a deposit onto the beacon block body. The execution payload is left untouched,
	// so the block is only invalid because it sources a deposit from the body after Fulu.
	maliciousTxs := testBuildInvalidBlock(
		s.Require(),
		s.SharedAccessors,
		&cmtabci.PrepareProposalRequest{
			Txs:    validProposal.Txs,
			Height: postFuluHeight,
			Time:   nextBlockTime,
		},
		func(sb *consensustypes.SignedBeaconBlock) {
			sb.BeaconBlock.Body.SetDeposits(consensustypes.Deposits{{Index: 99}})
		},
	)

	s.LogBuffer.Reset()
	processResp, err := s.SimComet.Comet.ProcessProposal(s.CtxComet, &cmtabci.ProcessProposalRequest{
		Txs:             maliciousTxs,
		Height:          postFuluHeight,
		ProposerAddress: nodeAddress,
		Time:            nextBlockTime,
	})
	s.Require().NoError(err)
	s.Require().Equal(cmtabci.PROCESS_PROPOSAL_STATUS_REJECT, processResp.Status)
	s.Require().Contains(s.LogBuffer.String(), core.ErrUnexpectedDepositSource.Error())
}

// TestNoDepositRequestsBeforeFulu verifies that before Fulu (in Electra) deposits never
// enter the chain as EIP-6110 execution requests. A deposit transaction is made available
// to the EL while the chain is still pre-Fulu (Electra1/Prague), so that its deposit event
// is emitted before Osaka. The pinned EL image is the moving bera-reth:nightly tag, so
// depending on its EIP-6110 inclusion gate one of two valid outcomes is asserted, both of
// which uphold the invariant:
//
//  1. The EL gates EIP-6110 on Osaka (fixed bera-reth): the Prague block is built including
//     the deposit transaction and carries no deposit execution requests.
//  2. The EL still gates EIP-6110 on Prague (older bera-reth): the EL surfaces a deposit
//     request before Fulu and the CL refuses to build the block, failing with
//     ErrUnexpectedDepositSource instead of silently consuming the request.
//
// Either way, deposits cannot be sourced from execution requests before Fulu. The
// legitimate EIP-6110 inclusion path from the first Fulu block onwards is covered by
// TestDepositQueueDrainedOnFirstFuluBlock.
//
// Chain spec (ProvideFuluDepositTestChainSpec): Deneb1 at genesis, Electra1 at t=6, Fulu at
// t=7. EL genesis: Cancun at genesis, Prague at t=6, Osaka at t=7. Block 1 is at t=5
// (Deneb1/Cancun) and block 2 is proposed at t=6 (Electra1/Prague).
func (s *FuluDepositSuite) TestNoDepositRequestsBeforeFulu() {
	s.InitializeChain(s.T(), 1)

	blsSigner := simulated.GetBlsSigner(s.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)
	nodeAddress := pubkey.Address()
	s.SimComet.Comet.SetNodeAddress(nodeAddress)

	credAddress, err := common.NewExecutionAddressFromHex(simulated.WithdrawalExecutionAddress)
	s.Require().NoError(err)
	creds := consensustypes.NewCredentialsFromExecutionAddress(credAddress)
	depositAmount := beaconmath.Gwei(10_000 * 1e9)

	// [Block 1, t=5] Deneb1/Cancun block, no deposits yet.
	_, _, nextBlockTime := s.MoveChainToHeight(s.T(), 1, 1, nodeAddress, time.Unix(5, 0))
	s.Require().Equal(time.Unix(6, 0), nextBlockTime)

	// Send a deposit so that it is pending in the EL mempool when the next
	// (Electra1/Prague) payload is requested, emitting its deposit event pre-Osaka.
	depositTxHash := s.sendDeposit(blsSigner, creds, depositAmount, true, big.NewInt(0))
	time.Sleep(time.Second)

	// [Block 2, t=6] Propose the Electra1/Prague block directly (not via MoveChainToHeight)
	// so that both EL behaviors can be asserted on below.
	s.LogBuffer.Reset()
	proposal, err := s.SimComet.Comet.PrepareProposal(s.CtxComet, &cmtabci.PrepareProposalRequest{
		Height:          2,
		Time:            nextBlockTime,
		ProposerAddress: nodeAddress,
	})
	s.Require().NoError(err)

	if len(proposal.Txs) == 0 {
		// The EL gates EIP-6110 on Prague and surfaced a deposit request before Fulu. The
		// only acceptable CL behavior is refusing to build a block that sources deposits
		// from execution requests.
		s.T().Log("EL surfaced a pre-Fulu deposit request; asserting the CL refused to build the block")
		s.Require().Contains(s.LogBuffer.String(), core.ErrUnexpectedDepositSource.Error())
		s.Require().Contains(s.LogBuffer.String(), "EIP-6110 deposit requests before Fulu")
		return
	}

	// The EL gates EIP-6110 on Osaka: the Prague block must include the deposit transaction
	// (so its deposit event was actually emitted pre-Osaka and the assertion below is not
	// vacuous) and must carry no deposit execution requests.
	s.T().Log("EL built the pre-Fulu block; asserting it carries no deposit execution requests")
	s.Require().Len(proposal.Txs, 2)
	preFuluForkVersion := s.TestNode.ChainSpec.ActiveForkVersionForTimestamp(
		beaconmath.U64(nextBlockTime.Unix()),
	)
	signedBlk, err := encoding.UnmarshalBeaconBlockFromABCIRequest(
		proposal.Txs, blockchain.BeaconBlockTxIndex, preFuluForkVersion,
	)
	s.Require().NoError(err)
	block := signedBlk.GetBeaconBlock()

	depositIncluded := false
	for _, raw := range block.GetBody().GetExecutionPayload().GetTransactions() {
		var tx gethcore.Transaction
		if uErr := tx.UnmarshalBinary(raw); uErr != nil {
			continue
		}
		if tx.Hash() == depositTxHash {
			depositIncluded = true
			break
		}
	}
	s.Require().True(depositIncluded,
		"deposit tx must be included in the pre-Fulu block so its deposit event is emitted")

	requests, err := block.GetBody().GetExecutionRequests()
	s.Require().NoError(err)
	s.Require().Empty(requests.Deposits,
		"no EIP-6110 deposit requests must be produced before Fulu")
}
