package smap

import (
	"fmt"
	"testing"
)

func TestClone(t *testing.T) {
	formMap := NewStrAnyFormMap(map[string]interface{}{
		"k1": "v1",
	}, false)
	formMap.Clone()
}

func TestSet(t *testing.T) {
	anyMap := NewStrAnyMap(true)
	anyMap.Set("name", "xie")
	anyMap.Set("age", func() interface{} {
		return 18
	})
	fmt.Println(anyMap.ToNewMap())
}
