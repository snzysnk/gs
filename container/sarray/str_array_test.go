package sarray

import (
	"fmt"
	"strings"
	"testing"
)

func TestSortFuncWithKey(t *testing.T) {
	arr := NewStrArrayFrom([]string{"a", "d", "e", "c", "f"}, true)
	fmt.Println(arr.ToStrSlice())
	arr.SortFunc(func(i, k string) bool {
		return strings.Compare(i, k) > 0
	})
	fmt.Println(arr.ToStrSlice())
	arr.SortFunc(func(i, k string) bool {
		return strings.Compare(i, k) < 0
	})
	fmt.Println(arr.ToStrSlice())
}

func TestUnique(t *testing.T) {
	arr := NewStrArrayFrom([]string{"a", "d", "e", "c", "a", "d", "e", "c", "d", "f"}, true)
	arrUnique := arr.Unique()
	arrUnique.SortFunc(func(i, k string) bool {
		return strings.Compare(i, k) < 0
	})
	fmt.Println(arrUnique.ToStrSlice())
	fmt.Println(arr.ToStrSlice())
}
