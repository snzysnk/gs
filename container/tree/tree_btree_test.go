package tree

import "testing"

func TestNewBTree(t *testing.T) {
	tree := NewBTree(3, func(v1, v2 interface{}) int {
		v1i := v1.(int)
		v2i := v2.(int)
		return v1i - v2i
	}, true)
	tree.Set(1, 1)
	tree.Set(2, 2)
	tree.Set(3, 3)
	tree.Dump()
}
