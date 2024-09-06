package vm

import (
	"errors"
	"fmt"

	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/database/prefixdb"
	"github.com/ava-labs/avalanchego/database/versiondb"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/snow/validators"

	"github.com/berachain/beacon-kit/mod/consensus/pkg/miniavalanche/block"
)

var (
	_ chainState = (*state)(nil)

	ValidatorsPrefix = []byte("validators")
	BlockPrefix      = []byte("block")
	BlockIDPrefix    = []byte("blockID")
	MetadataPrefix   = []byte("metadata")

	LastAcceptedKey = []byte("last accepted")
)

type chainState interface {
	SetValidator(val *Validator)
	GetValidator(valID ids.ID) (*Validator, error)

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
	chainCtx *snow.Context
	baseDB   *versiondb.Database // allows atomic commits of all changes that happens in a block

	// baseDB is key-partitioned into the following DB
	// each storing a relevant db aspect
	addedValidators map[ids.ID]*Validator // valID -> validator
	validatorsDB    database.Database

	addedBlock *block.StatelessBlock
	blockDB    database.Database
	blockIDDB  database.Database

	lastAcceptedBlockID ids.ID
	metadataDB          database.Database
}

func newState(
	chainCtx *snow.Context,
	db database.Database,
	validators validators.Manager,
	genesisBytes []byte,
) (chainState, error) {
	baseDB := versiondb.New(db)

	state := &state{
		chainCtx:        chainCtx,
		baseDB:          baseDB,
		addedValidators: make(map[ids.ID]*Validator),
		validatorsDB:    prefixdb.New(ValidatorsPrefix, baseDB),
		blockDB:         prefixdb.New(BlockPrefix, baseDB),
		blockIDDB:       prefixdb.New(BlockIDPrefix, baseDB),
		metadataDB:      prefixdb.New(MetadataPrefix, baseDB),
	}

	if err := state.handleGenesis(genesisBytes); err != nil {
		return nil, fmt.Errorf("failed handling genesis: %w", err)
	}

	// For the time being, validators are static
	return state, state.loadInMemoryData(validators)
}

func (s *state) handleGenesis(genesisBytes []byte) error {
	// check whether we need to handle genesis
	heightKey := database.PackUInt64(0)
	_, err := s.blockIDDB.Get(heightKey)
	if err == nil {
		// genesis already stored, nothing to do
		return nil
	}
	if err != database.ErrNotFound {
		return fmt.Errorf("could not check for gensis: %w", err)
	}

	// parse genesis
	genBlk, genVals, err := parseGenesis(genesisBytes)
	if err != nil {
		return fmt.Errorf("failed initializing VM: %w", err)
	}

	// store genesis data
	s.AddStatelessBlock(genBlk)
	s.lastAcceptedBlockID = genBlk.ID()
	for _, val := range genVals {
		s.SetValidator(val)
	}

	return s.Commit()
}

func (s *state) loadInMemoryData(validators validators.Manager) error {
	lastAccepted, err := database.GetID(s.metadataDB, LastAcceptedKey)
	if err != nil {
		return err
	}
	s.lastAcceptedBlockID = lastAccepted

	it := s.validatorsDB.NewIterator()
	defer it.Release()
	for it.Next() {
		valBytes := it.Value()
		v, err := ParseValidator(valBytes)
		if err != nil {
			return fmt.Errorf("failed parsing validator, bytes %s: %w", valBytes, err)
		}
		err = validators.AddStaker(s.chainCtx.SubnetID, v.NodeID, nil, v.id, v.Weight)
		if err != nil {
			return fmt.Errorf("failed registration of validator %v: %w", v.id, err)
		}
	}
	return nil
}

func (s *state) SetValidator(val *Validator) {
	s.addedValidators[val.id] = val
}

func (s *state) GetValidator(valID ids.ID) (*Validator, error) {
	valBytes, err := s.validatorsDB.Get(valID[:])
	switch err {
	case nil:
		return ParseValidator(valBytes)
	case database.ErrNotFound:
		return nil, database.ErrNotFound
	default:
		return nil, fmt.Errorf("GetValidator internal error: %w", err)
	}
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
		s.writeValidators(),
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

func (s *state) writeValidators() error {
	for _, val := range s.addedValidators {
		if err := s.validatorsDB.Put(val.id[:], val.bytes); err != nil {
			return fmt.Errorf("failed to write validator %s: %w", val.id, err)
		}
	}
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
		s.validatorsDB.Close(),
		s.blockDB.Close(),
		s.blockIDDB.Close(),
		s.metadataDB.Close(),
	)
}
