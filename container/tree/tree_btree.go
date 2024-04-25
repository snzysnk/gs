package tree

import (
	"github.com/snzysnk/gs/v2/internal/srwmutex"
)

type IBTree interface {
}

type BTree struct {
	mu         srwmutex.IRWMutex
	root       *BTreeNode
	comparator func(v1, v2 interface{}) int
	size       int
	m          int
}

type BTreeNode struct {
	Parent   *BTreeNode
	Entries  []*BTreeEntry
	Children []*BTreeNode
}

type BTreeEntry struct {
	Key   interface{}
	Value interface{}
}

func NewBTree(m int, comparator func(v1, v2 interface{}) int, safe bool) *BTree {
	if m < 3 {
		panic("m must be >= 3")
	}
	if comparator == nil {
		panic("comparator must not be nil")
	}
	return &BTree{
		comparator: comparator,
		m:          m,
		mu:         srwmutex.New(safe),
	}
}

func NewBTreeFrom(m int, comparator func(v1, v2 interface{}) int, data map[interface{}]interface{}, safe bool) *BTree {
	node := NewBTree(m, comparator, safe)
	return node
}

func (tree *BTree) doSet(key, value interface{}) {
	entry := &BTreeEntry{Key: key, Value: value}
	if tree.root == nil {
		tree.root = &BTreeNode{Entries: []*BTreeEntry{entry}, Children: []*BTreeNode{}}
		tree.size++
		return
	}
}

func (tree *BTree) insert(key, value interface{}) {

}

func (tree *BTree) isLeaf(node *BTreeNode) bool {
	return len(node.Children) == 0
}

// 二分法查找位置
func (tree *BTree) search(node *BTreeNode, key interface{}) (index int, found bool) {
	low, mid, high := 0, 0, len(node.Entries)-1
	for low <= high {
		mid = low + (high-low)/2
		compare := tree.comparator(key, node.Entries[mid].Key)
		switch {
		case compare > 0:
			low = mid + 1
		case compare < 0:
			high = mid - 1
		default:
			return mid, true
		}
	}
	return low, false
}
