package tree

import (
	"encoding/json"
	"fmt"
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

func (tree *BTree) insert(node *BTreeNode, entry *BTreeEntry) bool {
	if tree.isLeaf(node) {
		return tree.insertIntoLeaf(node, entry)
	}
	return tree.insertIntoInternal(node, entry)
}

func (tree *BTree) insertIntoInternal(node *BTreeNode, entry *BTreeEntry) (inserted bool) {
	index, found := tree.search(node, entry.Key)
	if found {
		node.Entries[index] = entry
		return false
	}
	return tree.insert(node, entry)
}

func (tree *BTree) dilatation(node *BTreeNode) {
	if !tree.shouldDilatation(node) {
		return
	}
	if tree.root == node {
		tree.rootDilatation()
		return
	}

	tree.noRootDilatation(node)
}

func (tree *BTree) noRootDilatation(node *BTreeNode) {
	var (
		middle = tree.middle()
		parent = node.Parent
	)
	left := &BTreeNode{Entries: append([]*BTreeEntry(nil), node.Entries[:middle]...), Parent: parent}
	right := &BTreeNode{Entries: append([]*BTreeEntry(nil), node.Entries[middle+1:]...), Parent: parent}

	if !tree.isLeaf(node) {
		left.Children = append([]*BTreeNode(nil), left.Children[:middle+1]...)
		right.Children = append([]*BTreeNode(nil), left.Children[middle+1:]...)
		setParent(left.Children, left)
		setParent(right.Children, right)
	}

	insertPosition, _ := tree.search(parent, node.Entries[middle].Key)
	parent.Entries = append(parent.Entries, nil)
	copy(parent.Entries[insertPosition+1:], parent.Entries[insertPosition:])
	parent.Entries[insertPosition] = node.Entries[middle]

	parent.Children[insertPosition] = left
	parent.Children = append(parent.Children, nil)
	copy(parent.Children[insertPosition+2:], parent.Children[insertPosition+1:])
	parent.Children[insertPosition+1] = right
	tree.dilatation(parent)
}

// 设置子节点的父节点引用
func setParent(nodes []*BTreeNode, parent *BTreeNode) {
	for _, node := range nodes {
		node.Parent = parent
	}
}

// 主节点的变动
func (tree *BTree) rootDilatation() {
	var (
		middle = tree.middle()
	)
	left := &BTreeNode{Entries: append([]*BTreeEntry(nil), tree.root.Entries[:middle]...)}
	right := &BTreeNode{Entries: append([]*BTreeEntry(nil), tree.root.Entries[middle+1:]...)}

	//从middle分左右节点
	if !tree.isLeaf(tree.root) {
		left.Children = append([]*BTreeNode(nil), tree.root.Children[:middle+1]...)
		right.Children = append([]*BTreeNode(nil), tree.root.Children[middle+1:]...)
		setParent(left.Children, left)
		setParent(right.Children, right)
	}

	//root变动
	newRoot := &BTreeNode{
		Entries:  []*BTreeEntry{tree.root.Entries[middle]},
		Children: []*BTreeNode{left, right},
	}

	//left.Parent = newRoot
	//right.Parent = newRoot
	tree.root = newRoot
	marshal, _ := json.Marshal(newRoot)
	fmt.Println(string(marshal))
}

// 当条目数大于最大的子节点数，可以移动条目形成新的node进行扩容
func (tree *BTree) shouldDilatation(node *BTreeNode) bool {
	return len(node.Entries) > tree.maxEntryNum()
}

// 最大子节点数
func (tree *BTree) maxChildrenNum() int {
	return tree.m
}

// 扩容条件
func (tree *BTree) maxEntryNum() int {
	return tree.m - 1
}

// 二分法启始位置
func (tree *BTree) middle() int {
	return tree.m / 2
}

// 判断是否是叶子节点(没有children的为页子节点)
func (tree *BTree) isLeaf(node *BTreeNode) bool {
	return len(node.Children) == 0
}

// 插入页子节点
func (tree *BTree) insertIntoLeaf(node *BTreeNode, entry *BTreeEntry) (inserted bool) {
	index, found := tree.search(node, entry.Key)
	if found {
		node.Entries[index] = entry
		return false
	}
	node.Entries = append(node.Entries, nil)
	copy(node.Entries[index+1:], node.Entries[index:])
	node.Entries[index] = entry
	tree.dilatation(node)
	return true
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

func (tree *BTree) Set(key, value interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	entry := &BTreeEntry{Key: key, Value: value}
	if tree.root == nil {
		tree.root = &BTreeNode{Entries: []*BTreeEntry{entry}, Children: []*BTreeNode{}}
		tree.size++
		return
	}

	if tree.insert(tree.root, entry) {
		tree.size++
	}
}

func (tree *BTree) Dump() {
	marshal, _ := json.Marshal(tree.root)
	fmt.Printf("result:%s\n", string(marshal))
}
