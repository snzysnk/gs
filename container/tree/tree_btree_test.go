package tree

import (
	"fmt"
	"testing"
)

func TestCopy(t *testing.T) {
	var (
		arr      = []string{"a", "b", "c", "d", "e", "f", "g"}
		position = 3
	)

	//copy(arr[position+1:], arr[position:])
	fmt.Println(arr[position+1:], arr[position:])
}
