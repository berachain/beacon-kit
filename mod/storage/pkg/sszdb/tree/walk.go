package tree

import (
	"fmt"

	ssz "github.com/ferranbt/fastssz"
)

var _ ssz.HashWalker = (*TreeView)(nil)

type TreeView struct {
	nodes []*ssz.Node
	buf   []byte
}

func (t *TreeView) Index() int {
	return len(t.nodes)
}

func (t *TreeView) Append(i []byte) {
	t.buf = append(t.buf, i...)
}

func (t *TreeView) AppendUint64(i uint64) {
	t.buf = ssz.MarshalUint64(t.buf, i)
}

func (t *TreeView) AppendUint32(i uint32) {
	t.buf = ssz.MarshalUint32(t.buf, i)
}

func (t *TreeView) AppendUint8(i uint8) {
	t.buf = ssz.MarshalUint8(t.buf, i)
}

func (t *TreeView) AppendBytes32(b []byte) {
	t.buf = append(t.buf, b...)
	t.FillUpTo32()
}

func (t *TreeView) FillUpTo32() {
	// pad zero bytes to the left
	if rest := len(t.buf) % 32; rest != 0 {
		t.buf = append(t.buf, zeroBytes[:32-rest]...)
	}
}

func (t *TreeView) Merkleize(indx int) {
	if len(t.buf) != 0 {
		t.appendBytesAsNodes(t.buf)
		t.buf = t.buf[:0]
	}
	t.Commit(indx)
}

func (t *TreeView) MerkleizeWithMixin(indx int, num, limit uint64) {
	if len(t.buf) != 0 {
		t.appendBytesAsNodes(t.buf)
		t.buf = t.buf[:0]
	}
	t.CommitWithMixin(indx, int(num), int(limit))
}

func (t *TreeView) PutBitlist(bb []byte, maxSize uint64) {
	b, size := parseBitlist(nil, bb)

	indx := t.Index()
	t.appendBytesAsNodes(b)
	t.CommitWithMixin(indx, int(size), int((maxSize+255)/256))
}

func (t *TreeView) appendBytesAsNodes(b []byte) {
	// if byte list is empty, fill with zeros
	if len(b) == 0 {
		b = append(b, zeroBytes[:32]...)
	}
	// if byte list isn't filled with 32-bytes padded, pad
	if rest := len(b) % 32; rest != 0 {
		b = append(b, zeroBytes[:32-rest]...)
	}
	for i := 0; i < len(b); i += 32 {
		val := append([]byte{}, b[i:min(len(b), i+32)]...)
		t.nodes = append(t.nodes, ssz.LeafFromBytes(val))
	}
}

func (t *TreeView) PutBool(b bool) {
	t.AddNode(ssz.LeafFromBool(b))
}

func (t *TreeView) PutBytes(b []byte) {
	t.AddBytes(b)
}

func (t *TreeView) PutUint16(i uint16) {
	t.AddUint16(i)
}

func (t *TreeView) PutUint64(i uint64) {
	t.AddUint64(i)
}

func (t *TreeView) PutUint8(i uint8) {
	t.AddUint8(i)
}

func (t *TreeView) PutUint32(i uint32) {
	t.AddUint32(i)
}

func (t *TreeView) PutUint64Array(b []uint64, maxCapacity ...uint64) {
	indx := t.Index()
	for _, i := range b {
		t.AppendUint64(i)
	}

	// pad zero bytes to the left
	t.FillUpTo32()

	if len(maxCapacity) == 0 {
		// Array with fixed size
		t.Merkleize(indx)
	} else {
		numItems := uint64(len(b))
		limit := ssz.CalculateLimit(maxCapacity[0], numItems, 8)

		t.MerkleizeWithMixin(indx, numItems, limit)
	}
}

/// --- legacy ones ---

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}

func (t *TreeView) AddBytes(b []byte) {
	if len(b) <= 32 {
		t.AddNode(ssz.LeafFromBytes(b))
	} else {
		indx := t.Index()
		t.appendBytesAsNodes(b)
		t.Commit(indx)
	}
}

func (t *TreeView) AddUint64(i uint64) {
	t.AddNode(ssz.LeafFromUint64(i))
}

func (t *TreeView) AddUint32(i uint32) {
	t.AddNode(ssz.LeafFromUint32(i))
}

func (t *TreeView) AddUint16(i uint16) {
	t.AddNode(ssz.LeafFromUint16(i))
}

func (t *TreeView) AddUint8(i uint8) {
	t.AddNode(ssz.LeafFromUint8(i))
}

func (t *TreeView) AddNode(n *ssz.Node) {
	if t.nodes == nil {
		t.nodes = []*ssz.Node{}
	}
	t.nodes = append(t.nodes, n)
}

func (t *TreeView) Node() *ssz.Node {
	if len(t.nodes) != 1 {
		panic("BAD")
	}
	return t.nodes[0]
}

func (t *TreeView) Hash() []byte {
	return t.nodes[len(t.nodes)-1].Hash()
}

func (t *TreeView) Commit(i int) {
	fmt.Printf("Commit; i=%d TreeFromNodes=%v\n", i, len(t.nodes))
	// create tree from nodes
	res, err := ssz.TreeFromNodes(t.nodes[i:], t.getLimit(i))
	if err != nil {
		panic(err)
	}
	// remove the old nodes
	t.nodes = t.nodes[:i]
	// add the new node
	t.AddNode(res)
}

func (t *TreeView) CommitWithMixin(i, num, limit int) {
	fmt.Printf("CommitWithMixin; i=%d TreeFromNodes=%v\n", i, len(t.nodes))
	// create tree from nodes
	res, err := ssz.TreeFromNodesWithMixin(t.nodes[i:], num, limit)
	if err != nil {
		panic(err)
	}
	// remove the old nodes
	t.nodes = t.nodes[:i]

	// add the new node
	t.AddNode(res)
}

func (t *TreeView) AddEmpty() {
	t.AddNode(ssz.EmptyLeaf())
}

func (t *TreeView) getLimit(i int) int {
	size := len(t.nodes[i:])
	return int(nextPowerOfTwo(uint64(size)))
}
