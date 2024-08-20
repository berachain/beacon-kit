package sszdb

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/cockroachdb/pebble"
	"github.com/emicklei/dot"
)

const devDBPath = "./.tmp/sszdb.db"

var hashFn = func(b []byte) []byte {
	h := sha256.Sum256(b)
	return h[:]
}

type Backend struct {
	db             *pebble.DB
	stages         map[uint64][]byte
	zeroHashes     [65][]byte
	zeroHashLevels map[string]int
}

type BackendConfig struct {
	Path string
}

func NewBackend(cfg BackendConfig) (*Backend, error) {
	if cfg.Path == "" {
		panic("path is required")
	}
	db, err := pebble.Open(cfg.Path, &pebble.Options{})
	if err != nil {
		return nil, err
	}

	b := &Backend{
		db:     db,
		stages: make(map[uint64][]byte),
	}

	// init zero hashes
	zero := make([]byte, 32)
	b.zeroHashLevels = make(map[string]int)
	b.zeroHashLevels[string(zero)] = 0
	buf := make([]byte, 64)
	for i := 0; i < 64; i++ {
		copy(buf[:32], b.zeroHashes[i][:])
		copy(buf[32:], b.zeroHashes[i][:])
		b.zeroHashes[i+1] = hashFn(buf)
		b.zeroHashLevels[string(b.zeroHashes[i+1])] = i + 1
	}

	return b, nil
}

func (d *Backend) Close() error {
	return d.db.Close()
}

// Get retrieves a value from the database by (gindex, version).
// The version is ignored in this implementation.
// Returns nil if the key does not exist.
func (d *Backend) Get(gindex uint64, _ int64) ([]byte, error) {
	keyBz := keyBytes(gindex)
	return d.get(keyBz)
}

// get retrieves a value from the database or nil if the key does not exist.
func (d *Backend) get(key []byte) ([]byte, error) {
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

func (d *Backend) SaveMonolith(mono Treeable) error {
	treeRoot, err := NewTreeFromFastSSZ(mono)
	if err != nil {
		return err
	}
	treeRoot.CachedHash()
	return d.save(treeRoot, 1)
}

func (d *Backend) save(node *Node, gindex uint64) error {
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

// stage queues and caches a node, its descendants, and its ancestors for later
// commitment to the database.
func (d *Backend) stage(
	_ context.Context, node *Node, gindex uint64,
) error {
	err := d.deepStage(node, gindex)
	if err != nil {
		return err
	}
	d.stageToRoot(gindex)
	return nil
}

//nolint:mnd // there is nothing magic about the number 2
func (d *Backend) deepStage(
	node *Node,
	gindex uint64,
) error {
	d.stages[gindex] = node.Value
	switch {
	case node.Left == nil && node.Right == nil:
		return nil
	case node.Left != nil && node.Right != nil:
		if err := d.deepStage(node.Left, 2*gindex); err != nil {
			return err
		}
		if err := d.deepStage(node.Right, 2*gindex+1); err != nil {
			return err
		}
	default:
		return errors.New("node has only one child")
	}
	return nil
}

func (d *Backend) Commit(ctx context.Context) error {
	_, err := d.Hash(ctx)
	if err != nil {
		return err
	}
	for gindex, value := range d.stages {
		key := keyBytes(gindex)
		if err := d.Set(key, value); err != nil {
			return err
		}
	}
	d.stages = make(map[uint64][]byte)
	return nil
}

func (d *Backend) Hash(_ context.Context) ([]byte, error) {
	return d.hash(1)
}

func (d *Backend) hash(gindex uint64) ([]byte, error) {
	n, found := d.stages[gindex]
	if n != nil {
		return n, nil
	}
	if !found {
		n, err := d.get(keyBytes(gindex))
		if err != nil {
			return nil, err
		}
		if n != nil {
			return n, nil
		}
	}

	left, err := d.hash(2 * gindex)
	if err != nil {
		return nil, err
	}
	if left == nil {
		return nil, fmt.Errorf("left node not found at gindex %d", 2*gindex)
	}
	right, err := d.hash(2*gindex + 1)
	if err != nil {
		return nil, err
	}
	if right == nil {
		return nil, fmt.Errorf("right node not found at gindex %d", 2*gindex+1)
	}
	n = hashFn(append(left, right...))
	d.stages[gindex] = n
	return n, nil
}

// getNode first checks the stage for the node, then the database.
// return nil if the node does not exist.
func (d *Backend) getNode(ctx context.Context, gindex uint64) ([]byte, error) {
	nodeBz, ok := d.stages[gindex]
	if ok {
		return nodeBz, nil
	}
	key := keyBytes(gindex)
	return d.get(key)
}

// mustGetNode first checks the stage for the node, then the database.
// returns an error if the node does not exist.
func (d *Backend) mustGetNode(
	ctx context.Context,
	gindex uint64,
) (*Node, error) {
	nodeBz, ok := d.stages[gindex]
	if ok {
		return &Node{Value: nodeBz}, nil
	}

	key := keyBytes(gindex)
	bz, err := d.get(key)
	if err != nil {
		return nil, err
	}
	if bz == nil {
		return nil, fmt.Errorf("node not found at gindex %d", gindex)
	}
	return &Node{Value: bz}, nil
}

func (d *Backend) getNodeBytes(
	ctx context.Context,
	gindex uint64,
	length uint64,
	offset uint8,
) ([]byte, error) {
	const chunkSize = 32
	var (
		buf bytes.Buffer
		i   int
		l   = int(length)
		o   = int(offset)
	)
	for ; l > 0; i++ {
		node, err := d.mustGetNode(ctx, gindex+uint64(i))
		if err != nil {
			return nil, err
		}
		end := l + o
		if end > chunkSize {
			end = chunkSize
		}
		n, err := buf.Write(node.Value[o:end])
		if err != nil {
			return nil, err
		}
		l -= n + o
		o = 0
	}

	return buf.Bytes(), nil
}

func (d *Backend) stageToRoot(gindex uint64) {
	for gindex > 1 {
		gindex /= 2
		if _, ok := d.stages[gindex]; ok {
			break
		}
		d.stages[gindex] = nil
	}
}

func (d *Backend) drawNode(
	val []byte,
	levelOrder int,
	g *dot.Graph,
) (dot.Node, error) {
	h := hex.EncodeToString(val)
	dn := g.Node(fmt.Sprintf("n%d", levelOrder)).
		Label(fmt.Sprintf("%d\n%s..%s", levelOrder, h[:3], h[len(h)-3:]))

	left, err := d.Get(uint64(levelOrder)*2, 0)
	if err != nil {
		return dot.Node{}, err
	}
	if left != nil {
		ln, err := d.drawNode(left, 2*levelOrder, g)
		if err != nil {
			return dot.Node{}, err
		}
		g.Edge(dn, ln).Label("0")
	}

	right, err := d.Get(uint64(levelOrder)*2+1, 0)
	if err != nil {
		return dot.Node{}, err
	}
	if right != nil {
		rn, err := d.drawNode(right, 2*levelOrder+1, g)
		if err != nil {
			return dot.Node{}, err
		}
		g.Edge(dn, rn).Label("1")
	}
	return dn, nil
}

func (d *Backend) DrawTree(ctx context.Context, f io.Writer) error {
	root, err := d.mustGetNode(ctx, 1)
	if err != nil {
		return err
	}
	g := dot.NewGraph(dot.Directed)
	_, err = d.drawNode(root.Value, 1, g)
	if err != nil {
		return err
	}
	g.Write(f)
	return nil
}
