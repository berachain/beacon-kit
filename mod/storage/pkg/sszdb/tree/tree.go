package tree

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"reflect"
	"unsafe"

	"github.com/emicklei/dot"
	ssz "github.com/ferranbt/fastssz"
)

type Node struct {
	Left    *Node
	Right   *Node
	IsEmpty bool
	Value   []byte
}

func NewTreeFromFastSSZ(r ssz.HashRoot) (*Node, error) {
	root, err := ssz.ProofTree(r)
	if err != nil {
		return nil, err
	}
	return copyTree(root), nil
}

// TODO this is a big hack to speed up development
// to be replaced with either a custom walker or simply ssz/v2
// It can also be used for regression testing against the fastssz implementation
func copyTree(node *ssz.Node) *Node {
	if node == nil {
		return nil
	}
	reflectNode := reflect.Indirect(reflect.ValueOf(node))

	f := reflectNode.FieldByIndex([]int{0})
	left := reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem().Interface().(*ssz.Node)

	f = reflectNode.FieldByIndex([]int{1})
	right := reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem().Interface().(*ssz.Node)

	f = reflectNode.FieldByIndex([]int{2})
	isEmpty := f.Bool()

	f = reflectNode.FieldByIndex([]int{3})
	value := f.Bytes()

	return &Node{
		Left:    copyTree(left),
		Right:   copyTree(right),
		IsEmpty: isEmpty,
		Value:   value,
	}
}

func (n *Node) CachedHash() []byte {
	if (n.Left == nil && n.Right == nil) || n.Value != nil {
		return n.Value
	}
	h := sha256.Sum256(append(n.Left.CachedHash(), n.Right.CachedHash()...))
	n.Value = h[:]
	return n.Value
}

func (n *Node) Hash() []byte {
	if n.Left == nil && n.Right == nil {
		return n.Value
	}
	h := sha256.Sum256(append(n.Left.Hash(), n.Right.Hash()...))
	return h[:]
}

func (n *Node) DrawTree(w io.Writer) {
	n.CachedHash()
	g := dot.NewGraph(dot.Directed)
	drawNode(n, 1, g)
	g.Write(w)
}

func drawNode(n *Node, levelOrder int, g *dot.Graph) dot.Node {
	h := hex.EncodeToString(n.Value)
	dn := g.Node(fmt.Sprintf("n%d", levelOrder)).
		Label(fmt.Sprintf("%d\n%s..%s", levelOrder, h[:3], h[len(h)-3:]))

	if n.Left != nil {
		ln := drawNode(n.Left, 2*levelOrder, g)
		g.Edge(dn, ln).Label("0")
	}
	if n.Right != nil {
		rn := drawNode(n.Right, 2*levelOrder+1, g)
		g.Edge(dn, rn).Label("1")
	}
	return dn
}

func (n *Node) Encode() []byte {
	var buf bytes.Buffer
	if n.IsEmpty {
		buf.Write([]byte{0})
	} else {
		buf.Write([]byte{1})
	}
	buf.Write(n.Value)
	return buf.Bytes()
}

func DecodeNode(b []byte) (*Node, error) {
	if len(b) == 0 {
		return nil, errors.New("empty node")
	}
	isEmpty := b[0] == 0
	return &Node{
		IsEmpty: isEmpty,
		Value:   b[1:],
	}, nil
}
