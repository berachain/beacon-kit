package sszdb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/storage/pkg/sszdb/tree"
	"github.com/cockroachdb/pebble"
	ssz "github.com/ferranbt/fastssz"
)

const devDBPath = "./.tmp/sszdb.db"

type DB struct {
	db *pebble.DB
}

type Config struct {
	Path string
}

func New(cfg Config) (*DB, error) {
	if cfg.Path == "" {
		cfg.Path = devDBPath
	}
	db, err := pebble.Open(cfg.Path, &pebble.Options{})
	if err != nil {
		return nil, err
	}
	return &DB{
		db: db,
	}, nil
}

func (d *DB) Close() error {
	return d.db.Close()
}

func (d *DB) Get(key []byte) ([]byte, error) {
	if len(key) == 0 {
		return nil, errors.New("key cannot be empty")
	}

	res, closer, err := d.db.Get(key)
	if err != nil {
		if errors.Is(err, pebble.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	defer closer.Close()
	ret := make([]byte, len(res))
	copy(ret, res)
	return ret, nil
}

func (d *DB) Set(key []byte, value []byte) error {
	if len(key) == 0 {
		return errors.New("key cannot be empty")
	}

	wopts := pebble.NoSync
	err := d.db.Set(key, value, wopts)
	if err != nil {
		return err
	}
	return nil
}

func keyBytes(gindex uint64) []byte {
	key := make([]byte, 8)
	binary.LittleEndian.PutUint64(key, gindex)
	return key
}

func (d *DB) SaveMonolith(mono ssz.HashRoot) error {
	treeRoot, err := tree.NewTreeFromFastSSZ(mono)
	if err != nil {
		return err
	}
	treeRoot.Hash()
	return d.save(treeRoot, 1)
}

func (d *DB) save(node *tree.Node, gindex uint64) error {
	// Save the node
	key := keyBytes(gindex)
	if err := d.Set(key, node.Encode()); err != nil {
		return err
	}

	switch {
	case node.Left == nil && node.Right == nil:
		return nil
	case node.Left != nil && node.Right != nil:
		if err := d.save(node.Left, 2*gindex); err != nil {
			return err
		}
		if err := d.save(node.Right, 2*gindex+1); err != nil {
			return err
		}
	default:
		return errors.New("node has only one child")
	}
	return nil
}

func (d *DB) getNode(gindex uint64) (*tree.Node, error) {
	key := keyBytes(gindex)
	bz, err := d.Get(key)
	if err != nil {
		return nil, err
	}
	if bz == nil {
		return nil, nil
	}
	return tree.DecodeNode(bz)
}

func (d *DB) mustGetNode(gindex uint64) (*tree.Node, error) {
	key := keyBytes(gindex)
	bz, err := d.Get(key)
	if err != nil {
		return nil, err
	}
	if bz == nil {
		return nil, fmt.Errorf("node not found at gindex %d", gindex)
	}
	return tree.DecodeNode(bz)
}

// todo: refactor to use offset
func (d *DB) getNodeBytes(gindex uint64, lenBz uint64) ([]byte, error) {
	const chunksize = 32

	numNodes := int(math.Ceil(float64(lenBz) / chunksize))
	rem := lenBz % chunksize
	var (
		buf bytes.Buffer
	)
	for i := 0; i < numNodes; i++ {
		n, err := d.mustGetNode(gindex + uint64(i))
		if err != nil {
			return nil, err
		}
		// last node
		if i == numNodes-1 && rem != 0 {
			buf.Write(n.Value[:rem])
		} else {
			buf.Write(n.Value)
		}
	}

	return buf.Bytes(), nil
}

func getLeftNode(node *ssz.Node) *ssz.Node {
	left, err := node.Get(2)
	if err != nil {
		return nil
	}
	return left
}

func getRightNode(node *ssz.Node) *ssz.Node {
	right, err := node.Get(3)
	if err != nil {
		return nil
	}
	return right
}
