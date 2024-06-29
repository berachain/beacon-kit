package ssz

import (
	"crypto/sha256"
	"reflect"
	"unsafe"

	fastssz "github.com/ferranbt/fastssz"
	ssz "github.com/ferranbt/fastssz"
)

type Node struct {
	GIndex uint64
	Left   *Node
	Right  *Node
	Value  []byte
}

func NewTreeFromFastSSZ(r fastssz.HashRoot) (*Node, error) {
	root, err := fastssz.ProofTree(r)
	if err != nil {
		return nil, err
	}
	return copyTree(root), nil
}

// TODO re-usable hashing, merge with other ssz work once completed

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

// TODO this is a big hack to speed up development
// to be replaced with either a custom walker or simply ssz/v2
// It can also be used for regression testing against the fastssz implementation
func copyTree(node *fastssz.Node) *Node {
	if node == nil {
		return nil
	}
	reflectNode := reflect.Indirect(reflect.ValueOf(node))

	f := reflectNode.FieldByIndex([]int{0})
	left := reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem().Interface().(*ssz.Node)

	f = reflectNode.FieldByIndex([]int{1})
	right := reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem().Interface().(*ssz.Node)

	f = reflectNode.FieldByIndex([]int{3})
	value := f.Bytes()

	return &Node{
		Left:  copyTree(left),
		Right: copyTree(right),
		Value: value,
	}
}

type SSZTreeable interface {
	GetRootNode() (*Node, error)
}
