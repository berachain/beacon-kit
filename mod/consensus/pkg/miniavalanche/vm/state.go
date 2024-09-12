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

package vm

import (
	"errors"
	"fmt"

	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/database/prefixdb"
	"github.com/ava-labs/avalanchego/database/versiondb"
	"github.com/ava-labs/avalanchego/ids"

	"github.com/berachain/beacon-kit/mod/consensus/pkg/miniavalanche/block"
)

var (
	_ chainState = (*state)(nil)

	BlockPrefix    = []byte("block")
	BlockIDPrefix  = []byte("blockID")
	MetadataPrefix = []byte("metadata")

	LastAcceptedKey = []byte("last accepted")
)

type chainState interface {
	// Invariant: [block] is an accepted block.
	AddStatelessBlock(blk *block.StatelessBlock)
	GetBlock(blkID ids.ID) (*block.StatelessBlock, error)
	GetBlockID(blkHeight uint64) (ids.ID, error)

	GetLastAccepted() ids.ID
	SetLastAccepted(blkID ids.ID)

	Commit() error
	Close() error
}

// state holds vm's persisted state
type state struct {
	baseDB *versiondb.Database // allows atomic commits of all changes that happens in a block

	// baseDB is key-partitioned into the following DB
	// each storing a relevant db aspect
	addedBlock *block.StatelessBlock
	blockDB    database.Database
	blockIDDB  database.Database

	lastAcceptedBlockID ids.ID
	metadataDB          database.Database
}

func newState(
	db database.Database,
	genBlk *block.StatelessBlock,
) (chainState, error) {
	baseDB := versiondb.New(db)

	state := &state{
		baseDB:     baseDB,
		blockDB:    prefixdb.New(BlockPrefix, baseDB),
		blockIDDB:  prefixdb.New(BlockIDPrefix, baseDB),
		metadataDB: prefixdb.New(MetadataPrefix, baseDB),
	}

	if err := state.handleGenesis(genBlk); err != nil {
		return nil, fmt.Errorf("failed handling genesis: %w", err)
	}

	return state, state.loadInMemoryData()
}

// handleGenesis check whether we need to handle genesis and does it if needed
func (s *state) handleGenesis(genBlk *block.StatelessBlock) error {
	heightKey := database.PackUInt64(0)
	_, err := s.blockIDDB.Get(heightKey)
	switch err {
	case nil:
		// genesis already stored, nothing to do
		return nil
	case database.ErrNotFound:
		// store genesis data
		s.AddStatelessBlock(genBlk)
		s.lastAcceptedBlockID = genBlk.ID()
		return s.Commit()
	default:
		return fmt.Errorf("could not check for gensis: %w", err)
	}
}

func (s *state) loadInMemoryData() error {
	lastAccepted, err := database.GetID(s.metadataDB, LastAcceptedKey)
	if err != nil {
		return err
	}
	s.lastAcceptedBlockID = lastAccepted
	return nil
}

func (s *state) AddStatelessBlock(blk *block.StatelessBlock) {
	s.addedBlock = blk
}

func (s *state) GetBlock(blkID ids.ID) (*block.StatelessBlock, error) {
	blkBytes, err := s.blockDB.Get(blkID[:])
	switch err {
	case nil:
		return block.ParseStatelessBlock(blkBytes)
	case database.ErrNotFound:
		return nil, database.ErrNotFound
	default:
		return nil, fmt.Errorf("GetBlock internal error: %w", err)
	}
}

func (s *state) GetBlockID(blkHeight uint64) (ids.ID, error) {
	heightKey := database.PackUInt64(blkHeight)
	blkID, err := database.GetID(s.blockIDDB, heightKey)
	switch err {
	case nil:
		return blkID, nil
	case database.ErrNotFound:
		return ids.Empty, database.ErrNotFound
	default:
		return ids.Empty, fmt.Errorf("GetBlockID internal error: %w", err)
	}
}

func (s *state) GetLastAccepted() ids.ID {
	return s.lastAcceptedBlockID
}

func (s *state) SetLastAccepted(blkID ids.ID) {
	s.lastAcceptedBlockID = blkID
}

func (s *state) Commit() error {
	defer s.abort()
	batch, err := s.commitBatch()
	if err != nil {
		return err
	}
	return batch.Write()
}

func (s *state) commitBatch() (database.Batch, error) {
	err := errors.Join(
		s.writeBlocks(),
		s.writeMetadata(),
	)
	if err != nil {
		return nil, err
	}
	return s.baseDB.CommitBatch()
}

func (s *state) writeBlocks() error {
	var (
		blkID     = s.addedBlock.ID()
		blkBytes  = s.addedBlock.Bytes()
		blkHeight = s.addedBlock.Height()

		heightKey = database.PackUInt64(blkHeight)
		blkIDKey  = blkID[:]
	)

	if err := database.PutID(s.blockIDDB, heightKey, blkID); err != nil {
		return fmt.Errorf("failed to add blockID: %w", err)
	}

	if err := s.blockDB.Put(blkIDKey, blkBytes); err != nil {
		return fmt.Errorf("failed to write block %s: %w", blkID, err)
	}

	s.addedBlock = nil
	return nil
}

func (s *state) writeMetadata() error {
	if err := database.PutID(s.metadataDB, LastAcceptedKey, s.lastAcceptedBlockID); err != nil {
		return fmt.Errorf("failed to write last accepted: %w", err)
	}
	return nil
}

func (s *state) abort() {
	s.baseDB.Abort()
}

func (s *state) Close() error {
	return errors.Join(
		s.blockDB.Close(),
		s.blockIDDB.Close(),
		s.metadataDB.Close(),
	)
}
