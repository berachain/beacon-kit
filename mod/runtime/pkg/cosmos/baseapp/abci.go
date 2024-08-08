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

package baseapp

import (
	"context"
	"errors"
	"fmt"
	"sort"

	corecomet "cosmossdk.io/core/comet"
	coreheader "cosmossdk.io/core/header"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/store/rootmulti"
	storetypes "cosmossdk.io/store/types"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	abci "github.com/cometbft/cometbft/api/cometbft/abci/v1"
	cmtproto "github.com/cometbft/cometbft/api/cometbft/types/v1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/gogoproto/proto"
)

func (app *BaseApp) InitChain(
	req *abci.InitChainRequest,
) (*abci.InitChainResponse, error) {
	if req.ChainId != app.chainID {
		return nil, fmt.Errorf(
			"invalid chain-id on InitChain; expected: %s, got: %s",
			app.chainID,
			req.ChainId,
		)
	}

	// On a new chain, we consider the init chain block height as 0, even though
	// req.InitialHeight is 1 by default.
	initHeader := cmtproto.Header{ChainID: req.ChainId, Time: req.Time}
	app.logger.Info(
		"InitChain",
		"initialHeight",
		req.InitialHeight,
		"chainID",
		req.ChainId,
	)

	// Set the initial height, which will be used to determine if we are
	// proposing
	// or processing the first block or not.
	app.initialHeight = req.InitialHeight
	if app.initialHeight == 0 { // If initial height is 0, set it to 1
		app.initialHeight = 1
	}

	// if req.InitialHeight is > 1, then we set the initial version on all
	// stores
	if req.InitialHeight > 1 {
		initHeader.Height = req.InitialHeight
		if err := app.cms.SetInitialVersion(req.InitialHeight); err != nil {
			return nil, err
		}
	}

	// initialize states with a correct header
	app.setState(execModeFinalize, initHeader)
	app.setState(execModeCheck, initHeader)

	// Store the consensus params in the BaseApp's param store. Note, this must
	// be done after the finalizeBlockState and context have been set as it's
	// persisted
	// to state.
	if req.ConsensusParams != nil {
		err := app.StoreConsensusParams(
			app.finalizeBlockState.Context(),
			*req.ConsensusParams,
		)
		if err != nil {
			return nil, err
		}
	}

	defer func() {
		// InitChain represents the state of the application BEFORE the first
		// block, i.e. the genesis block. This means that when processing the
		// app's InitChain handler, the block height is zero by default.
		// However, after Commit is called
		// the height needs to reflect the true block height.
		initHeader.Height = req.InitialHeight
		app.checkState.SetContext(
			app.checkState.Context().WithBlockHeader(initHeader).
				WithHeaderInfo(coreheader.Info{
					ChainID: req.ChainId,
					Height:  req.InitialHeight,
					Time:    req.Time,
				}),
		)
		app.finalizeBlockState.SetContext(
			app.finalizeBlockState.Context().WithBlockHeader(initHeader).
				WithHeaderInfo(coreheader.Info{
					ChainID: req.ChainId,
					Height:  req.InitialHeight,
					Time:    req.Time,
				}),
		)
	}()

	if app.initChainer == nil {
		return &abci.InitChainResponse{}, nil
	}

	// add block gas meter for any genesis transactions (allow infinite gas)
	app.finalizeBlockState.SetContext(
		app.finalizeBlockState.Context(),
	)

	res, err := app.initChainer(app.finalizeBlockState.Context(), req)
	if err != nil {
		return nil, err
	}

	if len(req.Validators) > 0 {
		if len(req.Validators) != len(res.Validators) {
			return nil, fmt.Errorf(
				"len(RequestInitChain.Validators) != len(GenesisValidators) (%d != %d)",
				len(req.Validators),
				len(res.Validators),
			)
		}

		sort.Sort(abcitypes.ValidatorUpdates(req.Validators))

		for i := range res.Validators {
			if !proto.Equal(&res.Validators[i], &req.Validators[i]) {
				return nil, fmt.Errorf(
					"genesisValidators[%d] != req.Validators[%d] ",
					i,
					i,
				)
			}
		}
	}

	// NOTE: We don't commit, but FinalizeBlock for block InitialHeight starts
	// from
	// this FinalizeBlockState.
	return &abci.InitChainResponse{
		ConsensusParams: res.ConsensusParams,
		Validators:      res.Validators,
		AppHash:         app.LastCommitID().Hash,
	}, nil
}

func (app *BaseApp) Info(_ *abci.InfoRequest) (*abci.InfoResponse, error) {
	lastCommitID := app.cms.LastCommitID()
	appVersion := InitialAppVersion
	if lastCommitID.Version > 0 {
		ctx, err := app.CreateQueryContext(lastCommitID.Version, false)
		if err != nil {
			return nil, fmt.Errorf("failed creating query context: %w", err)
		}
		appVersion, err = app.AppVersion(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed getting app version: %w", err)
		}
	}

	return &abci.InfoResponse{
		Data:             app.name,
		Version:          app.version,
		AppVersion:       appVersion,
		LastBlockHeight:  lastCommitID.Version,
		LastBlockAppHash: lastCommitID.Hash,
	}, nil
}

// PrepareProposal implements the PrepareProposal ABCI method and returns a
// ResponsePrepareProposal object to the client. The PrepareProposal method is
// responsible for allowing the block proposer to perform application-dependent
// work in a block before proposing it.
//
// Transactions can be modified, removed, or added by the application. Since the
// application maintains its own local mempool, it will ignore the transactions
// provided to it in RequestPrepareProposal. Instead, it will determine which
// transactions to return based on the mempool's semantics and the MaxTxBytes
// provided by the client's request.
//
// Ref:
// https://github.com/cosmos/cosmos-sdk/blob/main/docs/architecture/adr-060-abci-1.0.md
// Ref:
// https://github.com/cometbft/cometbft/blob/main/spec/abci/abci%2B%2B_basic_concepts.md
func (app *BaseApp) PrepareProposal(
	req *abci.PrepareProposalRequest,
) (resp *abci.PrepareProposalResponse, err error) {
	if app.prepareProposal == nil {
		return nil, errors.New("PrepareProposal handler not set")
	}

	// Always reset state given that PrepareProposal can timeout and be called
	// again in a subsequent round.
	header := cmtproto.Header{
		ChainID:            app.chainID,
		Height:             req.Height,
		Time:               req.Time,
		ProposerAddress:    req.ProposerAddress,
		NextValidatorsHash: req.NextValidatorsHash,
		AppHash:            app.LastCommitID().Hash,
	}
	app.setState(execModePrepareProposal, header)

	// CometBFT must never call PrepareProposal with a height of 0.
	//
	// Ref:
	// https://github.com/cometbft/cometbft/blob/059798a4f5b0c9f52aa8655fa619054a0154088c/spec/core/state.md?plain=1#L37-L38
	if req.Height < 1 {
		return nil, errors.New("PrepareProposal called with invalid height")
	}

	app.prepareProposalState.SetContext(
		app.getContextForProposal(app.prepareProposalState.Context(), req.Height).
			WithVoteInfos(toVoteInfo(req.LocalLastCommit.Votes)).

			// this is a set of votes that are not finalized yet, wait for
			// commit
			WithBlockHeight(req.Height).
			WithProposer(req.ProposerAddress).
			WithExecMode(sdk.ExecModePrepareProposal).
			WithCometInfo(corecomet.Info{
				Evidence:        sdk.ToSDKEvidence(req.Misbehavior),
				ValidatorsHash:  req.NextValidatorsHash,
				ProposerAddress: req.ProposerAddress,
				LastCommit:      sdk.ToSDKExtendedCommitInfo(req.LocalLastCommit),
			}).
			WithHeaderInfo(coreheader.Info{
				ChainID: app.chainID,
				Height:  req.Height,
				Time:    req.Time,
			}),
	)

	app.prepareProposalState.SetContext(app.prepareProposalState.Context())

	resp, err = app.prepareProposal(app.prepareProposalState.Context(), req)
	if err != nil {
		app.logger.Error(
			"failed to prepare proposal",
			"height",
			req.Height,
			"time",
			req.Time,
			"err",
			err,
		)
		return &abci.PrepareProposalResponse{Txs: req.Txs}, nil
	}

	return resp, nil
}

// ProcessProposal implements the ProcessProposal ABCI method and returns a
// ResponseProcessProposal object to the client. The ProcessProposal method is
// responsible for allowing execution of application-dependent work in a
// proposed
// block. Note, the application defines the exact implementation details of
// ProcessProposal. In general, the application must at the very least ensure
// that all transactions are valid. If all transactions are valid, then we
// inform
// CometBFT that the Status is ACCEPT. However, the application is also able
// to implement optimizations such as executing the entire proposed block
// immediately.
//
// If a panic is detected during execution of an application's ProcessProposal
// handler, it will be recovered and we will reject the proposal.
//
// Ref:
// https://github.com/cosmos/cosmos-sdk/blob/main/docs/architecture/adr-060-abci-1.0.md
// Ref:
// https://github.com/cometbft/cometbft/blob/main/spec/abci/abci%2B%2B_basic_concepts.md
func (app *BaseApp) ProcessProposal(
	req *abci.ProcessProposalRequest,
) (resp *abci.ProcessProposalResponse, err error) {
	if app.processProposal == nil {
		return nil, errors.New("ProcessProposal handler not set")
	}

	// CometBFT must never call ProcessProposal with a height of 0.
	// Ref:
	// https://github.com/cometbft/cometbft/blob/059798a4f5b0c9f52aa8655fa619054a0154088c/spec/core/state.md?plain=1#L37-L38
	if req.Height < 1 {
		return nil, errors.New("ProcessProposal called with invalid height")
	}

	// Always reset state given that ProcessProposal can timeout and be called
	// again in a subsequent round.
	header := cmtproto.Header{
		ChainID:            app.chainID,
		Height:             req.Height,
		Time:               req.Time,
		ProposerAddress:    req.ProposerAddress,
		NextValidatorsHash: req.NextValidatorsHash,
		AppHash:            app.LastCommitID().Hash,
	}
	app.setState(execModeProcessProposal, header)

	// Since the application can get access to FinalizeBlock state and write to
	// it, we must be sure to reset it in case ProcessProposal timeouts and is
	// called
	// again in a subsequent round. However, we only want to do this after we've
	// processed the first block, as we want to avoid overwriting the
	// finalizeState
	// after state changes during InitChain.
	if req.Height > app.initialHeight {
		app.setState(execModeFinalize, header)
	}

	app.processProposalState.SetContext(
		app.getContextForProposal(app.processProposalState.Context(), req.Height).
			WithVoteInfos(req.ProposedLastCommit.Votes).

			// this is a set of votes that are not finalized yet, wait for
			// commit
			WithBlockHeight(req.Height).
			WithHeaderHash(req.Hash).
			WithProposer(req.ProposerAddress).
			WithCometInfo(corecomet.Info{
				ProposerAddress: req.ProposerAddress,
				ValidatorsHash:  req.NextValidatorsHash,
				Evidence:        sdk.ToSDKEvidence(req.Misbehavior),
				LastCommit:      sdk.ToSDKCommitInfo(req.ProposedLastCommit),
			},
			).
			WithExecMode(sdk.ExecModeProcessProposal).
			WithHeaderInfo(coreheader.Info{
				ChainID: app.chainID,
				Height:  req.Height,
				Time:    req.Time,
			}),
	)

	resp, err = app.processProposal(app.processProposalState.Context(), req)
	if err != nil {
		app.logger.Error(
			"failed to process proposal",
			"height",
			req.Height,
			"time",
			req.Time,
			"hash",
			fmt.Sprintf("%X", req.Hash),
			"err",
			err,
		)
		return &abci.ProcessProposalResponse{
			Status: abci.PROCESS_PROPOSAL_STATUS_REJECT,
		}, nil
	}

	return resp, nil
}

// internalFinalizeBlock executes the block, called by the Optimistic
// Execution flow or by the FinalizeBlock ABCI method. The context received is
// only used to handle early cancellation, for anything related to state
// app.finalizeBlockState.Context()
// must be used.
func (app *BaseApp) internalFinalizeBlock(
	ctx context.Context,
	req *abci.FinalizeBlockRequest,
) (*abci.FinalizeBlockResponse, error) {
	if err := app.validateFinalizeBlockHeight(req); err != nil {
		return nil, err
	}

	header := cmtproto.Header{
		ChainID:            app.chainID,
		Height:             req.Height,
		Time:               req.Time,
		ProposerAddress:    req.ProposerAddress,
		NextValidatorsHash: req.NextValidatorsHash,
		AppHash:            app.LastCommitID().Hash,
	}

	// finalizeBlockState should be set on InitChain or ProcessProposal. If it
	// is nil, it means we are replaying this block and we need to set the state
	// here given that during block replay ProcessProposal is not executed by
	// CometBFT.
	if app.finalizeBlockState == nil {
		app.setState(execModeFinalize, header)
	}

	// Context is now updated with Header information.
	app.finalizeBlockState.SetContext(app.finalizeBlockState.Context().
		WithBlockHeader(header).
		WithHeaderHash(req.Hash).
		WithHeaderInfo(coreheader.Info{
			ChainID: app.chainID,
			Height:  req.Height,
			Time:    req.Time,
			Hash:    req.Hash,
			AppHash: app.LastCommitID().Hash,
		}).
		WithVoteInfos(req.DecidedLastCommit.Votes).
		WithExecMode(sdk.ExecModeFinalize).
		WithCometInfo(corecomet.Info{
			Evidence:        sdk.ToSDKEvidence(req.Misbehavior),
			ValidatorsHash:  req.NextValidatorsHash,
			ProposerAddress: req.ProposerAddress,
			LastCommit:      sdk.ToSDKCommitInfo(req.DecidedLastCommit),
		}))

	app.finalizeBlockState.SetContext(
		app.finalizeBlockState.Context(),
	)

	if app.checkState != nil {
		app.checkState.SetContext(app.checkState.Context().
			WithHeaderHash(req.Hash))
	}

	if err := app.preBlock(req); err != nil {
		return nil, err
	}

	if _, err := app.beginBlock(req); err != nil {
		return nil, err
	}

	// First check for an abort signal after beginBlock, as it's the first place
	// we spend any significant amount of time.
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		// continue
	}

	app.finalizeBlockState.SetContext(
		app.finalizeBlockState.Context(),
	)

	// Iterate over all raw transactions in the proposal and attempt to execute
	// them, gathering the execution results.
	//
	// NOTE: Not all raw transactions may adhere to the sdk.Tx interface, e.g.
	// vote extensions, so skip those.
	txResults := make([]*abci.ExecTxResult, 0, len(req.Txs))
	for _, rawTx := range req.Txs {
		response := app.deliverTx(rawTx)

		// check after every tx if we should abort
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			// continue
		}

		txResults = append(txResults, response)
	}

	endBlock, err := app.endBlock(app.finalizeBlockState.Context())
	if err != nil {
		return nil, err
	}

	// check after endBlock if we should abort, to avoid propagating the result
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		// continue
	}

	cp := app.GetConsensusParams(app.finalizeBlockState.Context())
	return &abci.FinalizeBlockResponse{
		TxResults:             txResults,
		ValidatorUpdates:      endBlock.ValidatorUpdates,
		ConsensusParamUpdates: &cp,
	}, nil
}

// FinalizeBlock will execute the block proposal provided by
// RequestFinalizeBlock. Specifically, it will execute an application's
// BeginBlock (if defined), followed
// by the transactions in the proposal, finally followed by the application's
// EndBlock (if defined).
//
// For each raw transaction, i.e. a byte slice, BaseApp will only execute it if
// it adheres to the sdk.Tx interface. Otherwise, the raw transaction will be
// skipped. This is to support compatibility with proposers injecting vote
// extensions into the proposal, which should not themselves be executed in
// cases
// where they adhere to the sdk.Tx interface.
func (app *BaseApp) FinalizeBlock(
	req *abci.FinalizeBlockRequest,
) (res *abci.FinalizeBlockResponse, err error) {

	res, err = app.internalFinalizeBlock(context.Background(), req)
	if res != nil {
		res.AppHash = app.workingHash()
	}

	return res, err
}

// Commit implements the ABCI interface. It will commit all state that exists in
// the deliver state's multi-store and includes the resulting commit ID in the
// returned abci.ResponseCommit. Commit will set the check state based on the
// latest header and reset the deliver state. Also, if a non-zero halt height is
// defined in config, Commit will execute a deferred function call to check
// against that height and gracefully halt if it matches the latest committed
// height.
func (app *BaseApp) Commit() (*abci.CommitResponse, error) {
	header := app.finalizeBlockState.Context().BlockHeader()
	retainHeight := app.GetBlockRetentionHeight(header.Height)

	rms, ok := app.cms.(*rootmulti.Store)
	if ok {
		rms.SetCommitHeader(header)
	}

	app.cms.Commit()

	resp := &abci.CommitResponse{
		RetainHeight: retainHeight,
	}

	// Reset the CheckTx state to the latest committed.
	//
	// NOTE: This is safe because CometBFT holds a lock on the mempool for
	// Commit. Use the header from this latest block.
	app.setState(execModeCheck, header)

	app.finalizeBlockState = nil

	return resp, nil
}

// workingHash gets the apphash that will be finalized in commit.
// These writes will be persisted to the root multi-store (app.cms) and flushed
// to disk in the Commit phase. This means when the ABCI client requests
// Commit(), the application state transitions will be flushed to disk and as a
// result, but we already have
// an application Merkle root.
func (app *BaseApp) workingHash() []byte {
	// Write the FinalizeBlock state into branched storage and commit the
	// MultiStore. The write to the FinalizeBlock state writes all state
	// transitions to the root
	// MultiStore (app.cms) so when Commit() is called it persists those values.
	app.finalizeBlockState.ms.Write()

	// Get the hash of all writes in order to return the apphash to the comet in
	// finalizeBlock.
	commitHash := app.cms.WorkingHash()
	app.logger.Debug(
		"hash of all writes",
		"workingHash",
		fmt.Sprintf("%X", commitHash),
	)

	return commitHash
}

// getContextForProposal returns the correct Context for PrepareProposal and
// ProcessProposal. We use finalizeBlockState on the first block to be able to
// access any state changes made in InitChain.
func (app *BaseApp) getContextForProposal(
	ctx sdk.Context,
	height int64,
) sdk.Context {
	if height == app.initialHeight {
		ctx, _ = app.finalizeBlockState.Context().CacheContext()

		// clear all context data set during InitChain to avoid inconsistent
		// behavior
		ctx = ctx.WithHeaderInfo(coreheader.Info{}).
			WithBlockHeader(cmtproto.Header{})
		return ctx
	}

	return ctx
}

// CreateQueryContext creates a new sdk.Context for a query, taking as args
// the block height and whether the query needs a proof or not.
func (app *BaseApp) CreateQueryContext(
	height int64,
	prove bool,
) (sdk.Context, error) {
	// use custom query multi-store if provided
	qms := app.cms.(storetypes.MultiStore)
	lastBlockHeight := qms.LatestVersion()
	if lastBlockHeight == 0 {
		return sdk.Context{}, errorsmod.Wrapf(
			sdkerrors.ErrInvalidHeight,
			"%s is not ready; please wait for first block",
			app.name,
		)
	}

	if height > lastBlockHeight {
		return sdk.Context{},
			errorsmod.Wrap(
				sdkerrors.ErrInvalidHeight,
				"cannot query with height in the future; please provide a valid height",
			)
	}

	// when a client did not provide a query height, manually inject the latest
	if height == 0 {
		height = lastBlockHeight
	}

	if height <= 1 && prove {
		return sdk.Context{},
			errorsmod.Wrap(
				sdkerrors.ErrInvalidRequest,
				"cannot query with proof when height <= 1; please provide a valid height",
			)
	}

	cacheMS, err := qms.CacheMultiStoreWithVersion(height)
	if err != nil {
		return sdk.Context{},
			errorsmod.Wrapf(
				sdkerrors.ErrNotFound,
				"failed to load state at height %d; %s (latest height: %d)",
				height,
				err,
				lastBlockHeight,
			)
	}

	// branch the commit multi-store for safety
	ctx := sdk.NewContext(cacheMS, true, app.logger).
		WithHeaderInfo(coreheader.Info{
			ChainID: app.chainID,
			Height:  height,
		}).
		WithBlockHeader(app.checkState.Context().BlockHeader()).
		WithBlockHeight(height)

	if height != lastBlockHeight {
		rms, ok := app.cms.(*rootmulti.Store)
		if ok {
			cInfo, err := rms.GetCommitInfo(height)
			if cInfo != nil && err == nil {
				ctx = ctx.WithHeaderInfo(
					coreheader.Info{Height: height, Time: cInfo.Timestamp},
				)
			}
		}
	}

	return ctx, nil
}

// GetBlockRetentionHeight returns the height for which all blocks below this
// height
// are pruned from CometBFT. Given a commitment height and a non-zero local
// minRetainBlocks configuration, the retentionHeight is the smallest height
// that
// satisfies:
//
// - Unbonding (safety threshold) time: The block interval in which validators
// can be economically punished for misbehavior. Blocks in this interval must be
// auditable e.g. by the light client.
//
// - Logical store snapshot interval: The block interval at which the underlying
// logical store database is persisted to disk, e.g. every 10000 heights. Blocks
// since the last IAVL snapshot must be available for replay on application
// restart.
//
// - State sync snapshots: Blocks since the oldest available snapshot must be
// available for state sync nodes to catch up (oldest because a node may be
// restoring an old snapshot while a new snapshot was taken).
//
// - Local (minRetainBlocks) config: Archive nodes may want to retain more or
// all blocks, e.g. via a local config option min-retain-blocks. There may also
// be a need to vary retention for other nodes, e.g. sentry nodes which do not
// need historical blocks.
func (app *BaseApp) GetBlockRetentionHeight(commitHeight int64) int64 {
	// pruning is disabled if minRetainBlocks is zero
	if app.minRetainBlocks == 0 {
		return 0
	}

	minNonZero := func(x, y int64) int64 {
		switch {
		case x == 0:
			return y

		case y == 0:
			return x

		case x < y:
			return x

		default:
			return y
		}
	}

	// Define retentionHeight as the minimum value that satisfies all non-zero
	// constraints. All blocks below (commitHeight-retentionHeight) are pruned
	// from CometBFT.
	var retentionHeight int64

	// Define the number of blocks needed to protect against misbehaving
	// validators
	// which allows light clients to operate safely. Note, we piggy back of the
	// evidence parameters instead of computing an estimated number of blocks
	// based
	// on the unbonding period and block commitment time as the two should be
	// equivalent.
	cp := app.GetConsensusParams(app.finalizeBlockState.Context())
	if cp.Evidence != nil && cp.Evidence.MaxAgeNumBlocks > 0 {
		retentionHeight = commitHeight - cp.Evidence.MaxAgeNumBlocks
	}

	//#nosec:G701 // bet.
	v := commitHeight - int64(app.minRetainBlocks)
	retentionHeight = minNonZero(retentionHeight, v)

	if retentionHeight <= 0 {
		// prune nothing in the case of a non-positive height
		return 0
	}

	return retentionHeight
}

// toVoteInfo converts the new ExtendedVoteInfo to VoteInfo.
func toVoteInfo(votes []abci.ExtendedVoteInfo) []abci.VoteInfo {
	legacyVotes := make([]abci.VoteInfo, len(votes))
	for i, vote := range votes {
		legacyVotes[i] = abci.VoteInfo{
			Validator: abci.Validator{
				Address: vote.Validator.Address,
				Power:   vote.Validator.Power,
			},
			BlockIdFlag: vote.BlockIdFlag,
		}
	}

	return legacyVotes
}

// LEgacy Helpers

// NewContextLegacy returns a new sdk.Context with the provided header.
func (app *BaseApp) NewContextLegacy(
	isCheckTx bool,
	header cmtproto.Header,
) sdk.Context {
	if isCheckTx {
		return sdk.NewContext(app.checkState.ms, true, app.logger).
			WithBlockHeader(header)
	}

	return sdk.NewContext(app.finalizeBlockState.ms, false, app.logger).
		WithBlockHeader(header)
}
