package sszdb

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
	"github.com/cockroachdb/pebble"
)

const devDBPath = "./.tmp/sszdb.db"

type Backend struct {
	db *pebble.DB
}

type BackendConfig struct {
	Path string
}

func NewBackend(cfg BackendConfig) (*Backend, error) {
	if cfg.Path == "" {
		cfg.Path = devDBPath
	}
	db, err := pebble.Open(cfg.Path, &pebble.Options{})
	if err != nil {
		return nil, err
	}
	return &Backend{
		db: db,
	}, nil
}

func (d *Backend) Close() error {
	return d.db.Close()
}

func (d *Backend) Get(key []byte) ([]byte, error) {
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

func (d *Backend) Set(key []byte, value []byte) error {
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

func (d *Backend) SaveMonolith(mono ssz.SSZTreeable) error {
	treeRoot, err := mono.GetRootNode()
	if err != nil {
		return err
	}
	treeRoot.CachedHash()
	return d.save(treeRoot, 1)
}

func (d *Backend) save(node *ssz.Node, gindex uint64) error {
	// Save the node
	key := keyBytes(gindex)
	if err := d.Set(key, node.Value); err != nil {
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

func (d *Backend) getNode(gindex uint64) (*ssz.Node, error) {
	key := keyBytes(gindex)
	bz, err := d.Get(key)
	if err != nil {
		return nil, err
	}
	if bz == nil {
		return nil, nil
	}
	return &ssz.Node{Value: bz}, nil
}

func (d *Backend) mustGetNode(gindex uint64) (*ssz.Node, error) {
	key := keyBytes(gindex)
	bz, err := d.Get(key)
	if err != nil {
		return nil, err
	}
	if bz == nil {
		return nil, fmt.Errorf("node not found at gindex %d", gindex)
	}
	return &ssz.Node{Value: bz}, nil
}

func (d *Backend) getNodeBytes(
	gindex uint64, length uint64, offset uint8,
) ([]byte, error) {
	const chunkSize = 32
	var (
		buf bytes.Buffer
		i   int
		l   = int(length)
		o   = int(offset)
	)
	for ; l > 0; i++ {
		node, err := d.mustGetNode(gindex + uint64(i))
		if err != nil {
			return nil, err
		}
		end := l + o
		if end > chunkSize {
			end = chunkSize
		}
		n, err := buf.Write(node.Value[offset:end])
		if err != nil {
			return nil, err
		}
		l -= n + o
		o = 0
	}

	return buf.Bytes(), nil
}
